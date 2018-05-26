package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	m := NewMapping()
	w := NewWatcher(m)
	d := NewDNSServer(m)
	p := NewProxyServer(m)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go w.Watch(ctx, 5*time.Second)
	go d.Serve(os.Getenv("DNS_ADDR"))
	go p.Serve(os.Getenv("PROXY_ADDR"))

	<-sigCh

	cancel()

	sdCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		d.Shutdown(sdCtx)
	}()

	go func() {
		defer wg.Done()
		p.Shutdown(ctx)
	}()

	wg.Wait()
}
