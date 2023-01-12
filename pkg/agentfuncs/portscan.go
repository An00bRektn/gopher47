package agentfuncs

import (
	"fmt"
	"log"
	"net"
	"sort"
)

// Based off of Blackhat Go: https://github.com/blackhat-go/bhg/blob/master/ch-2/tcp-scanner-final/main.go

func worker(ports, results chan int, addr string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func PortScanTCP(addr string, targets []int, workers int) string {
	ports := make(chan int, workers)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, addr)
	}

	go func() {
		for _, i := range targets {
			ports <- i
		}
	}()

	for i := 0; i < len(targets); i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	output := ""
	for _, port := range openports {
		log.Printf("%d open\n", port)
		output = fmt.Sprintf(output + "%d,", port)
	}

	return output
}