package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	psnet "github.com/shirou/gopsutil/net"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	scanner := NewScanner()

	for {
		select {
		case <-ticker.C:
			stats, err := psnet.Connections("tcp")
			if err == nil {
				wg := new(sync.WaitGroup)
				for _, s := range stats {
					wg.Add(1)
					go func(pid int32, port uint32) {
						defer wg.Done()
						scanner.Scan(pid, port)
					}(s.Pid, s.Laddr.Port)
				}
				wg.Wait()
			} else {
				log.Println(err)
			}
		case <-sigCh:
			return nil
		}
	}

	return nil
}
