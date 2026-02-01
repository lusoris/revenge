// Package main is the entry point for the Revenge media server.
package main

import (
	"fmt"
	"os"

	"github.com/revenge/revenge/internal/config"
	"github.com/revenge/revenge/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("revenge", version.Info())
		return
	}

	cfg := config.Default()
	fmt.Printf("Starting Revenge server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	// TODO: Initialize and run application
}
