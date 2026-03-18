package scanner

import (
	"context"
	"sync"
	"time"
)

func Run(ctx context.Context, hosts []string, ports []int, concurrency int, timeout int) ScanSummary {
	start := time.Now()

	timeoutDuration := time.Duration(timeout) * time.Millisecond
	
	targets := make(chan ScanTarget)
	results := make(chan ScanResult)

	// start workers...
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range targets {
				select {
				case <-ctx.Done():
					return
				default:
				}
				result := probe(target.Host, target.Port, timeoutDuration)
				results <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		defer close(targets)
		for _, host := range hosts {
			for _, port := range ports {
				select {
				case <-ctx.Done():
					return
				case targets <- ScanTarget{Host: host, Port: port}:
				}
			}
		}
	}()

	hostMap := make(map[string]*HostResult)
	for _, host := range hosts {
		hostMap[host] = &HostResult{Host: host}
	}

	totalOpen := 0
	for result := range results {
		if result.Open {
			totalOpen++
			hr := hostMap[result.Host]
			hr.OpenPorts = append(hr.OpenPorts, result)
		}
	}

	hostResults := make([]HostResult, 0, len(hosts))
	for _, host := range hosts {
		hostResults = append(hostResults, *hostMap[host])
	}

	return ScanSummary{
		Hosts:       hosts,
		TotalPorts:  len(hosts) * len(ports),
		TotalOpen:   totalOpen,
		Duration:    time.Since(start),
		HostResults: hostResults,
	}
}
