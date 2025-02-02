package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// scanUDPPorts scans a range of UDP ports on a given host within a specified timeout
func scanUDPPorts(host string, minPort, maxPort int,
	timeout time.Duration) {
	var wg sync.WaitGroup

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			addr := fmt.Sprintf("%s:%d", host, p)
			conn, err := net.DialTimeout("udp", addr, timeout)
			if err != nil {
				fmt.Printf("Port %d is closed on %s\n", p, host)
				return
			}
			defer conn.Close()

			// Send a dummy packet
			_, err = conn.Write([]byte("ping"))
			if err != nil {
				fmt.Printf("Port %d is closed on %s\n", p, host)
				return
			}

			// Set read deadline
			buffer := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(timeout))
			_, err = conn.Read(buffer)
			if err != nil {
				fmt.Printf("Port %d is closed on %s\n", p, host)
				return
			}

			fmt.Printf("Port %d is open on %s\n", p, host)
		}(port)
	}

	wg.Wait()
}

// scanTCPPorts scans a range of UDP ports on a given host within a specified timeout
func scanTCPPorts(minPort int, maxPort int,
	hostname string, timeoutInSec time.Duration) error {
	if minPort > maxPort {
		return fmt.Errorf("minPort is higher than maxPort")
	}

	if minPort > 65535 {
		return fmt.Errorf("minPort is higher than 65535")
	}

	if maxPort > 65535 {
		return fmt.Errorf("maxPort is higher than 65535")
	}

	var wg sync.WaitGroup
	for i := minPort; i <= maxPort; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", hostname, j)
			conn, err := net.DialTimeout("tcp", address, timeoutInSec)
			if err != nil {
				log.Printf("hostname: %s, port: %d", hostname, j)
				return
			}
			conn.Close()
			log.Printf("%d open\n", j)
		}(i)
	}
	wg.Wait()

	return nil
}

func main() {

}
