package monitor

import "time"

const (
	Version        = "v1.0.0"
	CheckInterval  = 500 * time.Millisecond
	MaxHistory     = 1000
	MaxLogs        = 6
	GraphHeight    = 5
	MaxLatencies   = 30
	UptimeWindow   = 10000 // Rolling window for uptime calculation
)

type LogEntry struct {
	Timestamp  time.Time
	Success    bool
	StatusCode int
	Latency    time.Duration
	ErrMsg     string
}

type CheckResult struct {
	Success    bool
	StatusCode int
	Latency    time.Duration
	Err        error
	Timestamp  time.Time
}

type TickMsg time.Time
