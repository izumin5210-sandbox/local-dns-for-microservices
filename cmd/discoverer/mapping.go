package main

import (
	"fmt"
	"log"
	"sync"
)

type Mapping interface {
	Lookup(host string) (targetHost string, ok bool)
	Update(port uint32, pid int32, host string)
	IsChecked(port uint32, pid int32) bool
	CanMap(host string) bool
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
	pidByPort *sync.Map
	hostTable *sync.Map
	toAddr    func(port uint32) string
}

func (m *mappingImpl) Lookup(host string) (string, bool) {
	if v, ok := m.hostTable.Load(host); ok {
		if dst, ok := v.(string); ok {
			return dst, true
		}
	}
	return "", false
}

func (m *mappingImpl) Update(port uint32, pid int32, host string) {
	dst := m.toAddr(port)
	m.pidByPort.Store(port, pid)
	m.hostTable.Store(host, m.toAddr(port))
	log.Printf("detect new maping: %s -> %s (pid = %d)\n", host, dst, pid)
}

func (m *mappingImpl) IsChecked(port uint32, pid int32) bool {
	v, ok := m.pidByPort.Load(port)
	if ok {
		if oldPid, ok := v.(int32); ok {
			return oldPid == pid
		}
	}
	return false
}

func (m *mappingImpl) CanMap(host string) bool {
	_, ok := m.hostTable.Load(host)
	return ok
}

func (m *mappingImpl) Clear() {
	m.pidByPort = new(sync.Map)
	m.hostTable = new(sync.Map)
}
