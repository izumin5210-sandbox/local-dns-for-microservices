package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	lstvPat = regexp.MustCompile("\\s+")
)

type Scanner struct {
	pidByPort *sync.Map
}

func NewScanner() *Scanner {
	return &Scanner{
		pidByPort: new(sync.Map),
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
	fmt.Println(cwd)
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
