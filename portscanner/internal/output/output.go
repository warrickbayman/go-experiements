package output

import (
	"fmt"
	"os"
	"portscanner/internal/scanner"
	"sort"
)

func PrintSummary(summary scanner.ScanSummary) {
	fmt.Println()
	fmt.Println("--------------------------------")
	fmt.Printf("Scan complete in %s\n", summary.Duration.Round(1000000))
	fmt.Printf("Hosts scanned : %d\n", len(summary.Hosts))
	fmt.Printf("Ports per host: %d\n", summary.TotalPorts/max(len(summary.Hosts), 1))
	fmt.Printf("Total probes  : %d\n", summary.TotalPorts)
	fmt.Printf("Open ports    : %d\n", summary.TotalOpen)
	fmt.Println("--------------------------------")

	for _, hr := range summary.HostResults {
		if len(hr.OpenPorts) == 0 {
			continue
		}
		fmt.Printf("\n%s\n", hr.Host)
		sorted := hr.OpenPorts
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[1].Port < sorted[j].Port
		})
		for _, r := range sorted {
			PrintOpen(r)
		}
	}
}

func WriteToFile(path string, summary scanner.ScanSummary, asJSON bool) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer f.Close()

	if asJSON {
		return writeJSON(f, summary)
	}
	return writePlain
}

func PrintOpen(result scanner.ScanResult) {
	if result.Banner != "" {
		fmt.Printf("  %-6d open    %s\n", result.Port, result.Banner)
	} else {
		fmt.Printf("  %-6d open\n", result.Port)
	}
}
