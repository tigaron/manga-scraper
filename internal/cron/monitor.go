package cron

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type cronMonitor struct {
	mu      sync.Mutex
	counter map[string]int
	time    map[string][]time.Duration
}

func newCronMonitor() *cronMonitor {
	return &cronMonitor{
		counter: make(map[string]int),
		time:    make(map[string][]time.Duration),
	}
}

func (t *cronMonitor) IncrementJob(_ uuid.UUID, name string, _ []string, _ gocron.JobStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.counter[name]
	if !ok {
		t.counter[name] = 0
	}
	t.counter[name]++
}

func (t *cronMonitor) RecordJobTiming(startTime, endTime time.Time, _ uuid.UUID, name string, _ []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.time[name]
	if !ok {
		t.time[name] = make([]time.Duration, 0)
	}
	t.time[name] = append(t.time[name], endTime.Sub(startTime))
}
