package main

import (
	"flag"
	"log"

	"./daemon"
)

func processFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3000", "HTTP listen spec")

	flag.Parse()

	return cfg
}

func main() {
	// process command line
	cfg := processFlags()

	// start the daemon
	if err := daemon.Run(cfg); err != nil {
		log.Printf("Daemon error: %v\n", err)
	}
}
