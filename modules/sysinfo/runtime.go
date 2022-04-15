package sysinfo

import (
	"runtime"
	"time"
)

// RuntimeStatus runtime status
type RuntimeStatus struct {
	NumGoroutine int

	// General statistics.
	MemAllocated uint64 // bytes allocated and still in use
	MemTotal     uint64 // bytes allocated (even if freed)
	MemSys       uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    uint64 // bytes allocated and still in use
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	// Inuse is bytes used now.
	// Sys is bytes obtained from system.
	StackInuse  uint64 // bootstrap stacks
	StackSys    uint64
	MSpanInuse  uint64 // mspan structures
	MSpanSys    uint64
	MCacheInuse uint64 // mcache structures
	MCacheSys   uint64
	BuckHashSys uint64 // profiling bucket hash table
	GCSys       uint64 // GC metadata
	OtherSys    uint64 // other system allocations

	// Garbage collector statistics.
	NextGC         uint64 // next run in HeapAlloc time (bytes)
	LastGC         uint64 // last run in absolute time (ns)
	LastGCRelative uint64 // last run in relative time (ns)
	PauseTotalNs   uint64
	LastPauseNs    uint64 // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC          uint32
}

// Status runtime status
func newStatus() *RuntimeStatus {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return &RuntimeStatus{
		NumGoroutine:   runtime.NumGoroutine(),
		MemAllocated:   m.Alloc,
		MemTotal:       m.TotalAlloc,
		MemSys:         m.Sys,
		Lookups:        m.Lookups,
		MemMallocs:     m.Mallocs,
		MemFrees:       m.Frees,
		HeapAlloc:      m.HeapAlloc,
		HeapSys:        m.HeapSys,
		HeapIdle:       m.HeapIdle,
		HeapInuse:      m.HeapInuse,
		HeapReleased:   m.HeapReleased,
		HeapObjects:    m.HeapObjects,
		StackInuse:     m.StackInuse,
		StackSys:       m.StackSys,
		MSpanInuse:     m.MSpanInuse,
		MSpanSys:       m.MSpanSys,
		MCacheInuse:    m.MCacheInuse,
		MCacheSys:      m.MCacheSys,
		BuckHashSys:    m.BuckHashSys,
		GCSys:          m.GCSys,
		OtherSys:       m.OtherSys,
		NextGC:         m.NextGC,
		LastGC:         m.LastGC,
		LastGCRelative: uint64(time.Now().UnixNano()) - m.LastGC,
		PauseTotalNs:   m.PauseTotalNs,
		LastPauseNs:    m.PauseNs[(m.NumGC+255)%256],
		NumGC:          m.NumGC,
	}
}
