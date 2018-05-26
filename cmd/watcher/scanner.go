package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/afero"
)

var (
	lstvPat = regexp.MustCompile("\\s+")
)

type Scanner struct {
	pidByPort *sync.Map
	fs        afero.Fs
}

func NewScanner() *Scanner {
	return &Scanner{
		pidByPort: new(sync.Map),
		fs:        afero.NewOsFs(),
	}
}

func (s *Scanner) Scan(pid int32, port uint32) error {
	if s.checked(pid, port) {
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
	data, err := afero.ReadFile(s.fs, filepath.Join(cwd, "localhost"))
	if err != nil {
		return err
	}
	fmt.Printf(":%d (pid = %d) => %s\n", port, pid, string(data))
	return nil
}

func (s *Scanner) checked(pid int32, port uint32) bool {
	v, ok := s.pidByPort.LoadOrStore(port, pid)
	if ok {
		if storedPid, ok := v.(int32); ok {
			return pid == storedPid
		}
	}
	return false
}
