// PiN Browser — Pi Integrated Network
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

	"pin-browser/browser"
)

var (
	apiAddr  = flag.String("api", "127.0.0.1:4002", "meshd API address")
	headless = flag.Bool("headless", false, "Run in headless proxy mode (no UI)")
	version  = flag.Bool("version", false, "Print version and exit")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("PiN Browser %s\n", browser.Version)
		fmt.Println("Pi Integrated Network")
		fmt.Println("https://github.com/pin-network/pin-network")
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("browser: shutting down")
		cancel()
	}()

	b := browser.New(browser.Config{
		APIAddr:  *apiAddr,
		Headless: *headless,
	})

	if err := b.Start(ctx); err != nil {
		log.Fatalf("browser: %v", err)
	}
}
