package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"portscanner/internal/cli"
	"portscanner/internal/output"
	"portscanner/internal/scanner"
	"syscall"
)

func main() {
	cfg := cli.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nInterrupted - stopping scan...")
		cancel()
	}()

	fmt.Printf(
		"Scanning %d host(s) across %d port(s)\n\n",
		len(cfg.Hosts),
		len(cfg.Ports),
	)

	// Run the scanner...

	summary := scanner.Run(ctx, cfg.Hosts, cfg.Ports, cfg.Concurrency, cfg.Timeout)

	output.PrintSummary(summary)

	if cfg.OutputFile != "" {
		err := output.WriteToFile(cfg.OutputFile, summary, cfg.OutputJSON)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nResults written to %s\n", cfg.OutputFile)
	}
}
