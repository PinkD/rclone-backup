package stat

import (
	"sync"
	"time"
)

const (
	None = iota
	Syncing
	Fail
	Success
)

type Status int

type SyncStatus struct {
	Name   string
	Status Status
	Time   time.Time

	LastSuccessTime time.Time
}

type StatusMap struct {
	sync.Map
}

func (m *StatusMap) Store(name string, status *SyncStatus) {
	if status.LastSuccessTime.IsZero() {
		v, ok := m.Map.Load(name)
		if ok {
			s := v.(*SyncStatus)
			status.LastSuccessTime = s.LastSuccessTime
		}
	}
	m.Map.Store(name, status)
}

func (m *StatusMap) LoadStatus(name string) Status {
	v, ok := m.Map.Load(name)
	if ok {
		s := v.(*SyncStatus)
		return s.Status
	}
	return None
}

// map[name]SyncStatus
var Map StatusMap
