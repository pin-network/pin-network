// meshd — PiN node daemon
// Pi Integrated Network
// https://github.com/pin-network/pin-network

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"meshd/config"
	"meshd/ledger"
	"meshd/limits"
	"meshd/node"
	"meshd/scheduler"
	"meshd/server"
	"meshd/store"
)

var (
	configPath = flag.String("config", "", "Path to config file (default: ~/.pin/config.yaml)")
	devMode    = flag.Bool("dev", false, "Run in development mode with verbose logging")
	initMode   = flag.Bool("init", false, "Initialise a new node and exit")
	version    = flag.Bool("version", false, "Print version and exit")
)

const Version = "0.1.0-dev"

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("meshd version %s\n", Version)
		fmt.Println("PiN — Pi Integrated Network")
		fmt.Println("https://github.com/pin-network/pin-network")
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if *devMode {
		cfg.Dev = true
		log.Println("running in development mode")
	}

	// Set low OS priority — meshd runs as a background process
	limits.SetLowPriority()

	// Initialise resource limiter
	limiter := limits.New(cfg)

	if *initMode {
		if err := node.Init(cfg); err != nil {
			log.Fatalf("failed to initialise node: %v", err)
		}
		fmt.Println("node initialised successfully")
		fmt.Printf("data directory: %s\n", cfg.DataDir)
		fmt.Println("run meshd to start the node")
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialise ledger
	db, err := ledger.Open(cfg.LedgerPath())
	if err != nil {
		log.Fatalf("failed to open ledger: %v", err)
	}
	defer db.Close()

	// Start epoch Hash calculator
	db.StartEpochCalculator(ctx, cfg.Node.Tier)

	// Initialise content store
	st, err := store.New(cfg.StorePath())
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}

	// Start scheduler
	sched := scheduler.New(cfg)
	sched.Start(ctx)

	// Watch config file for changes — allows app to update settings at runtime
	configWatcher := config.NewWatcher(*configPath, func(newCfg *config.Config) {
		sched.UpdateConfig(newCfg)
		limiter.UpdateConfig(newCfg)
	})
	configWatcher.Start(ctx)

	// Initialise the PiN node (libp2p host + DHT)
	n, err := node.New(ctx, cfg, db)
	if err != nil {
		log.Fatalf("failed to create node: %v", err)
	}

	log.Printf("PiN node started")
	log.Printf("  NodeID:   %s", n.ID())
	log.Printf("  Tier:     %d", cfg.Node.Tier)
	log.Printf("  Storage:  %s (limit %dGB)", cfg.Node.StoragePath, cfg.Node.StorageLimitGB)
	log.Printf("  Listen:   %v", n.Addrs())
	log.Printf("  Schedule: always_on=%v", cfg.Schedule.AlwaysOn)
	log.Printf("  Config:   watching for changes every 30s")

	// Start the local API server
	api := server.NewAPI(cfg, n, db, st, sched, limiter)
	go func() {
		if err := api.ListenAndServe(); err != nil {
			log.Printf("API server error: %v", err)
		}
	}()
	log.Printf("  API:      http://0.0.0.0:%d", cfg.Network.APIPort)

	// Bootstrap into the network
	if err := n.Bootstrap(ctx); err != nil {
		log.Printf("warning: bootstrap incomplete: %v", err)
		log.Println("continuing — will retry peer discovery in background")
	}

	// Record uptime start
	uptimeID, err := db.RecordStart()
	if err != nil {
		log.Printf("warning: could not record uptime start: %v", err)
	}

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal %v, shutting down", sig)
	case <-ctx.Done():
		log.Println("context cancelled, shutting down")
	}

	if uptimeID > 0 {
		if err := db.RecordStop(uptimeID); err != nil {
			log.Printf("warning: could not record uptime stop: %v", err)
		}
	}

	log.Println("meshd stopped")
}
