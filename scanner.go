package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

func scan(hostname string, port int, wg *sync.WaitGroup, openPorts chan int) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout("tcp", address, 20*time.Second)

	if err != nil {
		openPorts <- 0
		return
	}

	defer conn.Close()
	openPorts <- port
}

func checkRange(start, end int) error {
	if start > end {
		return errors.New("start of range after end")
	}
	return nil
}

func main() {
	// Enter start of range
	fmt.Println("Start Port: ")
	var start int
	fmt.Scanln(&start)

	fmt.Println("End Port: ")
	var end int
	fmt.Scanln(&end)

	if checkRange(start, end) != nil {
		fmt.Println("RANGE ERROR! ")
	} else {
		openPorts := make(chan int, (end - start + 1))
		var wg sync.WaitGroup

		go func() {
			for port := range openPorts {
				if port != 0 {
					fmt.Println("Open port:", port)
				}
			}
		}()

		for i := start; i < end; i++ {
			wg.Add(1)
			go scan("localhost", i, &wg, openPorts)
		}

		wg.Wait()
		close(openPorts)
	}
}
