package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScanner struct {
	Host      string
	StartPort int
	EndPort   int
}

type ScanResult struct {
	Port    int
	IsOpen  bool
	Message string
}

func NewPortScanner(host string, startPort, endPort int) *PortScanner {
	return &PortScanner{
		Host:      host,
		StartPort: startPort,
		EndPort:   endPort,
	}
}

func (ps *PortScanner) ScanPorts() []ScanResult {
	var wg sync.WaitGroup
	results := make(chan ScanResult, ps.EndPort-ps.StartPort+1)

	for port := ps.StartPort; port <= ps.EndPort; port++ {
		wg.Add(1)
		go ps.scanPort(port, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var scanResults []ScanResult
	for result := range results {
		if result.IsOpen {
			scanResults = append(scanResults, result)
		}
	}

	return scanResults
}

func (ps *PortScanner) scanPort(port int, wg *sync.WaitGroup, results chan<- ScanResult) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", ps.Host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	result := ScanResult{Port: port}
	if err == nil {
		conn.Close()
		result.IsOpen = true
		result.Message = fmt.Sprintf("Port %d is open", port)
	}

	results <- result
}
