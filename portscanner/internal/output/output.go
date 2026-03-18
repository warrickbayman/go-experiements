package output

import (
	"encoding/json"
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
	return writePlain(f, summary)
}

func writeJSON(f *os.File, summary scanner.ScanSummary) error {
	type portEntry struct {
		Port   int    `json:"port"`
		Banner string `json:"banner,omitempty"`
	}

	type hostEntry struct {
		Host      string      `json:"host"`
		OpenPorts []portEntry `json:"open_ports"`
	}

	type jsonOutput struct {
		Hosts       []string    `json:"hosts"`
		TotalProbes int         `json:"total_probes"`
		TotalOpen   int         `json:"total_open"`
		DurationMs  int64       `json:"duration_ms"`
		Results     []hostEntry `json:"results"`
	}

	results := make([]hostEntry, 0, len(summary.HostResults))
	for _, hr := range summary.HostResults {
		ports := make([]portEntry, 0, len(hr.OpenPorts))
		for _, r := range hr.OpenPorts {
			ports = append(ports, portEntry{Port: r.Port, Banner: r.Banner})
		}

		results = append(results, hostEntry{Host: hr.Host, OpenPorts: ports})
	}

	out := jsonOutput{
		Hosts:       summary.Hosts,
		TotalProbes: summary.TotalPorts,
		TotalOpen:   summary.TotalOpen,
		DurationMs:  summary.Duration.Milliseconds(),
		Results:     results,
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(out)
}

func writePlain(f *os.File, summary scanner.ScanSummary) error {
	fmt.Fprintf(f, "Scan Duration : %s\n", summary.Duration.Round(1000000))
	fmt.Fprintf(f, "Total Proves  : %d\n", summary.TotalPorts)
	fmt.Fprintf(f, "Open Ports.   : %d\n\n", summary.TotalOpen)

	for _, hr := range summary.HostResults {
		if len(hr.OpenPorts) == 0 {
			continue
		}

		fmt.Fprintf(f, "Host: %s\n", hr.Host)
		sorted := hr.OpenPorts
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Port < sorted[j].Port
		})

		for _, r := range sorted {
			if r.Banner != "" {
				fmt.Fprintf(f, "  %d\t%s\n", r.Port, r.Banner)
			} else {
				fmt.Fprintf(f, "%d\n", r.Port)
			}
		}
		fmt.Fprintln(f)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func PrintOpen(result scanner.ScanResult) {
	if result.Banner != "" {
		fmt.Printf("  %-6d open    %s\n", result.Port, result.Banner)
	} else {
		fmt.Printf("  %-6d open\n", result.Port)
	}
}
