package ledger

import (
	"context"
	"log"
	"time"
)

const (
	BaseRate = 1_000_000 // bytes per Hash (1MB = 1 Hash)
)

var tierMultipliers = map[int]float64{
	1: 1.0,
	2: 2.0,
	3: 4.0,
}

// StartEpochCalculator runs a background goroutine that calculates
// Hash earnings at the end of each 24-hour epoch.
func (d *DB) StartEpochCalculator(ctx context.Context, tier int) {
	go func() {
		// Calculate time until next epoch boundary (midnight UTC)
		now := time.Now().UTC()
		nextEpoch := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		waitDuration := time.Until(nextEpoch)

		log.Printf("epoch calculator: next epoch in %s", waitDuration.Round(time.Minute))

		select {
		case <-time.After(waitDuration):
			d.calculateEpoch(tier)
		case <-ctx.Done():
			return
		}

		// After the first epoch, run every 24 hours
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				d.calculateEpoch(tier)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// calculateEpoch calculates Hash earnings for the previous epoch and stores them.
func (d *DB) calculateEpoch(tier int) {
	yesterday := epochID(time.Now()) - 1
	epochStart := epochStart(yesterday)
	epochEnd := epochStart + 86400

	log.Printf("epoch calculator: calculating epoch %d", yesterday)

	// Get bytes served during the epoch
	var bytesServed int64
	err := d.sql.QueryRow(`
		SELECT COALESCE(SUM(bytes), 0) FROM traffic_log
		WHERE timestamp >= ? AND timestamp < ?`,
		epochStart, epochEnd,
	).Scan(&bytesServed)
	if err != nil {
		log.Printf("epoch calculator: error querying bytes served: %v", err)
		return
	}

	// Get uptime during the epoch
	var uptimeMinutes int64
	now := epochEnd
	err = d.sql.QueryRow(`
		SELECT COALESCE(SUM(
			(MIN(COALESCE(stopped_at, ?), ?) - MAX(started_at, ?)) / 60
		), 0)
		FROM node_uptime
		WHERE started_at < ? AND (stopped_at IS NULL OR stopped_at > ?)`,
		now, epochEnd, epochStart, epochEnd, epochStart,
	).Scan(&uptimeMinutes)
	if err != nil {
		log.Printf("epoch calculator: error querying uptime: %v", err)
		return
	}

	// Calculate Hashes earned
	uptimeFactor := float64(uptimeMinutes) / 1440.0
	multiplier := tierMultipliers[tier]
	if multiplier == 0 {
		multiplier = 1.0
	}

	hashesEarned := (float64(bytesServed) / float64(BaseRate)) * uptimeFactor * multiplier

	// Minimum Hash reward for being online at all (1 Hash per 10% uptime)
	if uptimeFactor >= 0.1 && hashesEarned < uptimeFactor*float64(tier) {
		hashesEarned = uptimeFactor * float64(tier)
	}

	log.Printf("epoch calculator: epoch %d — %d bytes served, %d min uptime (%.1f%%), %.4f Hashes earned",
		yesterday, bytesServed, uptimeMinutes, uptimeFactor*100, hashesEarned)

	// Store the epoch summary
	_, err = d.sql.Exec(`
		INSERT OR REPLACE INTO epoch_summary 
		(epoch, hashes_earned, bytes_served, uptime_minutes, tier, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		yesterday, hashesEarned, bytesServed, uptimeMinutes, tier, time.Now().Unix(),
	)
	if err != nil {
		log.Printf("epoch calculator: error storing epoch summary: %v", err)
		return
	}

	log.Printf("epoch calculator: epoch %d complete — %.4f Hashes minted", yesterday, hashesEarned)
}

// EpochHistory returns the last N epoch summaries.
func (d *DB) EpochHistory(limit int) ([]EpochSummary, error) {
	rows, err := d.sql.Query(`
		SELECT epoch, hashes_earned, bytes_served, uptime_minutes, tier
		FROM epoch_summary ORDER BY epoch DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []EpochSummary
	for rows.Next() {
		var s EpochSummary
		if err := rows.Scan(&s.Epoch, &s.HashesEarned, &s.BytesServed, &s.UptimeMinutes, &s.Tier); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}
