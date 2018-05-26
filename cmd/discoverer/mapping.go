package main

import (
	"fmt"
	"log"
	"sync"
)

type Mapping interface {
	Lookup(host string) (targetHost string, ok bool)
	Update(port uint32, pid int32, host string)
	Has(port uint32, pid int32) bool
	Clear()
}

func NewMapping() Mapping {
	m := &mappingImpl{}
	m.toAddr = func(port uint32) string {
		return fmt.Sprintf("localhost:%d", port)
	}
	m.Clear()
	return m
}

type mappingImpl struct {
	portByHost *sync.Map
	toAddr     func(port uint32) string
}

type process struct {
	pid  int32
	host string
}

func (m *mappingImpl) Lookup(host string) (string, bool) {
	if v, ok := m.portByHost.Load(host); ok {
		if proc, ok := v.(process); ok {
			return proc.host, true
		}
	}
	return "", false
}

func (m *mappingImpl) Update(port uint32, pid int32, host string) {
	dst := m.toAddr(port)
	m.portByHost.Store(host, process{pid: pid, host: dst})
	log.Printf("detect new maping: %s -> %s (pid = %d)\n", host, dst, pid)
}

func (m *mappingImpl) Has(port uint32, pid int32) bool {
	v, ok := m.portByHost.Load(m.toAddr(port))
	if ok {
		if proc, ok := v.(process); ok {
			return proc.pid == pid
		}
	}
	return false
}

func (m *mappingImpl) Clear() {
	m.portByHost = new(sync.Map)
}
