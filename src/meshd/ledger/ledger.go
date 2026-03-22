// Package ledger manages the local proof-of-service SQLite database.
package ledger

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	sql *sql.DB
}

// TrafficRecord represents a single served request logged to the ledger.
type TrafficRecord struct {
	ID         int64
	Timestamp  time.Time
	Requester  string
	ContentCID string
	Bytes      int64
	DurationMS int64
	Verified   bool
}

// EpochSummary holds the Hash earnings for a 24-hour epoch.
type EpochSummary struct {
	Epoch         int64
	HashesEarned  float64
	BytesServed   int64
	UptimeMinutes int64
	Tier          int
}

// Open opens (or creates) the ledger database at the given path.
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("creating ledger directory: %w", err)
	}

	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("opening ledger database: %w", err)
	}

	l := &DB{sql: db}
	if err := l.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrating ledger schema: %w", err)
	}

	return l, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.sql.Close()
}

// LogTraffic records a served request to the ledger.
func (d *DB) LogTraffic(r TrafficRecord) error {
	_, err := d.sql.Exec(`
		INSERT INTO traffic_log (timestamp, requester, content_cid, bytes, duration_ms, verified)
		VALUES (?, ?, ?, ?, ?, ?)`,
		r.Timestamp.Unix(),
		r.Requester,
		r.ContentCID,
		r.Bytes,
		r.DurationMS,
		r.Verified,
	)
	if err != nil {
		return fmt.Errorf("logging traffic record: %w", err)
	}
	return nil
}

// Balance returns the current total Hash balance.
func (d *DB) Balance() (float64, error) {
	var balance float64
	err := d.sql.QueryRow(`
		SELECT COALESCE(SUM(hashes_earned), 0) FROM epoch_summary`).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("querying Hash balance: %w", err)
	}
	return balance, nil
}

// BytesServedToday returns total bytes served in the current epoch.
func (d *DB) BytesServedToday() (int64, error) {
	today := epochID(time.Now())
	var bytes int64
	err := d.sql.QueryRow(`
		SELECT COALESCE(SUM(bytes), 0) FROM traffic_log
		WHERE timestamp >= ?`, epochStart(today)).Scan(&bytes)
	if err != nil {
		return 0, fmt.Errorf("querying bytes served: %w", err)
	}
	return bytes, nil
}

// RecentTraffic returns the most recent traffic records.
func (d *DB) RecentTraffic(limit int) ([]TrafficRecord, error) {
	rows, err := d.sql.Query(`
		SELECT id, timestamp, requester, content_cid, bytes, duration_ms, verified
		FROM traffic_log ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("querying recent traffic: %w", err)
	}
	defer rows.Close()

	var records []TrafficRecord
	for rows.Next() {
		var r TrafficRecord
		var ts int64
		if err := rows.Scan(&r.ID, &ts, &r.Requester, &r.ContentCID,
			&r.Bytes, &r.DurationMS, &r.Verified); err != nil {
			return nil, err
		}
		r.Timestamp = time.Unix(ts, 0)
		records = append(records, r)
	}
	return records, nil
}

// migrate creates the database schema if it does not exist.
func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS traffic_log (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp   INTEGER NOT NULL,
		requester   TEXT NOT NULL,
		content_cid TEXT NOT NULL,
		bytes       INTEGER NOT NULL,
		duration_ms INTEGER NOT NULL,
		verified    BOOLEAN DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_traffic_timestamp ON traffic_log(timestamp);
	CREATE INDEX IF NOT EXISTS idx_traffic_cid ON traffic_log(content_cid);

	CREATE TABLE IF NOT EXISTS epoch_summary (
		epoch          INTEGER PRIMARY KEY,
		hashes_earned  REAL NOT NULL DEFAULT 0,
		bytes_served   INTEGER NOT NULL DEFAULT 0,
		uptime_minutes INTEGER NOT NULL DEFAULT 0,
		tier           INTEGER NOT NULL DEFAULT 1,
		created_at     INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS node_uptime (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		started_at INTEGER NOT NULL,
		stopped_at INTEGER
	);
	`

	if _, err := d.sql.Exec(schema); err != nil {
		return fmt.Errorf("creating schema: %w", err)
	}
	return nil
}

// epochID returns the epoch number for a given time (days since Unix epoch).
func epochID(t time.Time) int64 {
	return t.Unix() / 86400
}

// epochStart returns the Unix timestamp of the start of an epoch.
func epochStart(epoch int64) int64 {
	return epoch * 86400
}

// RecordStart records the node coming online and returns the uptime record ID.
func (d *DB) RecordStart() (int64, error) {
	result, err := d.sql.Exec(`
		INSERT INTO node_uptime (started_at) VALUES (?)`,
		time.Now().Unix(),
	)
	if err != nil {
		return 0, fmt.Errorf("recording uptime start: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RecordStop records the node going offline.
func (d *DB) RecordStop(id int64) error {
	_, err := d.sql.Exec(`
		UPDATE node_uptime SET stopped_at = ? WHERE id = ?`,
		time.Now().Unix(), id,
	)
	if err != nil {
		return fmt.Errorf("recording uptime stop: %w", err)
	}
	return nil
}

// UptimeToday returns the number of minutes the node has been online today.
func (d *DB) UptimeToday() (int64, error) {
	todayStart := epochStart(epochID(time.Now()))
	todayEnd := todayStart + 86400
	now := time.Now().Unix()

	// Cap end time at now (not future)
	if now < todayEnd {
		todayEnd = now
	}

	var minutes int64
	err := d.sql.QueryRow(`
		SELECT COALESCE(SUM(
			(MIN(COALESCE(stopped_at, ?), ?) - MAX(started_at, ?)) / 60
		), 0)
		FROM node_uptime
		WHERE started_at < ? AND (stopped_at IS NULL OR stopped_at > ?)
		AND (MIN(COALESCE(stopped_at, ?), ?) - MAX(started_at, ?)) > 0`,
		todayEnd, todayEnd, todayStart,
		todayEnd, todayStart,
		todayEnd, todayEnd, todayStart,
	).Scan(&minutes)
	if err != nil {
		return 0, fmt.Errorf("querying uptime: %w", err)
	}
	return minutes, nil
}
