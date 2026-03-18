package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Hosts       []string
	Ports       []int
	Timeout     int
	Concurrency int
	OutputFile  string
	OutputJSON  bool
}

var commonPorts = []int{
	21, 22, 23, 25, 53, 80, 110, 111, 135, 139,
	143, 443, 445, 993, 995, 1723, 3306, 3389,
	5900, 8080, 8443, 8888, 27017,
}

func Parse() *Config {
	hosts := flag.String("hosts", "", "Comma-separated list of hosts to scan (e.g. 192.168.1.1, example.com)")
	ports := flag.String("ports", "", "Port range or list (e.g. 80, 1-1024, 22,80,443)")
	timeout := flag.Int("timeout", 500, "Connection timeout in milliseconds")
	concurrency := flag.Int("concurrency", 100, "Number of concurrent workers")
	outputFile := flag.String("output", "", "File path to write results to")
	outputJSON := flag.Bool("json", false, "Output results as JSON")

	flag.Parse()

	if *hosts == "" {
		fmt.Fprintln(os.Stderr, "Error: -hosts flag is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg := &Config{
		Hosts:       parseHosts(*hosts),
		Ports:       commonPorts,
		Timeout:     *timeout,
		Concurrency: *concurrency,
		OutputFile:  *outputFile,
		OutputJSON:  *outputJSON,
	}

	if *ports != "" {
		parsed, err := parsePorts(*ports)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing ports: %v\n", err)
			os.Exit(1)
		}
		cfg.Ports = parsed
	}

	return cfg

}

func parseHosts(raw string) []string {
	parts := strings.Split(raw, ",")
	hosts := make([]string, 0, len(parts))

	for _, h := range parts {
		h = strings.TrimSpace(h)
		if h != "" {
			hosts = append(hosts, h)
		}
	}

	return hosts
}

func parsePorts(raw string) ([]int, error) {
	ports := []int{}
	parts := strings.Split(raw, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			var start, end int
			_, err := fmt.Sscanf(part, "%d-%d", &start, &end)
			if err != nil {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			if start > end || start < 1 || end > 65535 {
				return nil, fmt.Errorf("invalid range: %s", part)
			}
			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			var p int
			_, err := fmt.Sscanf(part, "%d", &p)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("port out of range, %d", p)
			}
			ports = append(ports, p)
		}
	}

	return ports, nil
}
