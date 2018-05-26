package main

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	psnet "github.com/shirou/gopsutil/net"
	"github.com/spf13/afero"
)

var (
	lstvPat = regexp.MustCompile("\\s+")
)

type Watcher struct {
	mapping Mapping
	fs      afero.Fs
}

func NewWatcher(mapping Mapping) *Watcher {
	return &Watcher{
		mapping: mapping,
		fs:      afero.NewOsFs(),
	}
}

func (w *Watcher) Watch(ctx context.Context, interval time.Duration) error {
	log.Println("start watcing ports")

	if err := w.do(); err != nil {
		log.Println(err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.do(); err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (w *Watcher) do() error {
	log.Println("start scanning ports")

	stats, err := psnet.Connections("tcp")
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	sem := make(chan struct{}, 4)

	for _, s := range stats {
		sem <- struct{}{}
		wg.Add(1)
		go func(s psnet.ConnectionStat) {
			defer wg.Done()
			defer func() { <-sem }()
			w.scan(s.Laddr.Port, s.Pid)
		}(s)
	}

	wg.Wait()
	log.Printf("finish scanning ports (found %d processes)\n", len(stats))

	return nil
}

func (w *Watcher) scan(port uint32, pid int32) error {
	if w.mapping.IsChecked(port, pid) {
		return nil
	}

	out, err := exec.Command("lsof", "-p", strconv.FormatInt(int64(pid), 10)).Output()
	if err != nil {
		return err
	}

	var cwd string
	for _, line := range strings.Split(string(out), "\n") {
		cols := lstvPat.Split(line, 9)
		if len(cols) == 9 && cols[3] == "cwd" {
			cwd = cols[8]
			break
		}
	}

	data, err := afero.ReadFile(w.fs, filepath.Join(cwd, "localhost"))
	if err != nil {
		return err
	}

	w.mapping.Update(port, pid, strings.Split(string(data), "\n")[0])
	return nil
}
