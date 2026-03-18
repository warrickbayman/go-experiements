package scanner

import (
	"fmt"
	"net"
	"time"
)

type ScanTarget struct {
	Host string
	Port int
}

type ScanResult struct {
	Host   string
	Port   int
	Open   bool
	Banner string
	Error  error
}

type HostResult struct {
	Host      string
	OpenPorts []ScanResult
}

type ScanSummary struct {
	Hosts       []string
	TotalPorts  int
	TotalOpen   int
	Duration    time.Duration
	HostResults []HostResult
}

func (r ScanResult) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func probe(host string, port int, timeout time.Duration) ScanResult {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	result := ScanResult{
		Host: host,
		Port: port,
	}

	if err != nil {
		result.Open = false
		result.Error = err
		return result
	}

	defer conn.Close()
	result.Open = true

	conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	if n > 0 {
		result.Banner = cleanBanner(buf[:n])
	}

	return result
}

func cleanBanner(raw []byte) string {
	out := make([]byte, 0, len(raw))
	for _, b := range raw {
		if b >= 32 && b < 127 {
			out = append(out, b)
		}
	}

	return string(out)
}
