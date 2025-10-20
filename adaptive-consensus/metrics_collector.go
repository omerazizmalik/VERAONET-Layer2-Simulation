// adaptive-consensus/metrics_collector.go
// Metrics helpers used by the adaptive switcher.
// You can plug in real Geth/Prometheus/JSON-RPC metrics here.

package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// NodeMetrics is the minimal set we need to decide consensus
type NodeMetrics struct {
	ActiveUsers      int     // concurrent users / load proxy
	ThroughputTPS    float64 // transactions per second
	LatencyMS        int     // average end-to-end latency (ms)
	EnergyNormalized float64 // 0..1 (or map from J/tx)
}

// MetricsSource is a simple interface to fetch time-series metrics
type MetricsSource interface {
	Next() (NodeMetrics, error)
	Close() error
}

/* -------------------- RandomSource (default demo) -------------------- */

type randomSource struct{}

func NewRandomSource() MetricsSource {
	rand.Seed(time.Now().UnixNano())
	return &randomSource{}
}

func (r *randomSource) Next() (NodeMetrics, error) {
	// Simulate realistic ranges; tweak as needed
	active := rand.Intn(12000) // 0..11999
	var tps float64
	switch {
	case active < 500:
		tps = float64(50 + rand.Intn(120)) // 50–170
	case active < 5000:
		tps = float64(200 + rand.Intn(400)) // 200–600
	default:
		tps = float64(500 + rand.Intn(900)) // 500–1400
	}
	lat := 80 + rand.Intn(1400) // 80..1480 ms
	energy := rand.Float64() * 1.0

	return NodeMetrics{
		ActiveUsers:      active,
		ThroughputTPS:    tps,
		LatencyMS:        lat,
		EnergyNormalized: energy,
	}, nil
}

func (r *randomSource) Close() error { return nil }

/* -------------------- CSVSource (for reproducible runs) -------------------- */
// CSV Header expected (order flexible, case-insensitive):
// ActiveUsers, ThroughputTPS, LatencyMS, EnergyNormalized
// Example row: 1200, 380.5, 640, 0.62

type csvSource struct {
	file    *os.File
	reader  *csv.Reader
	colIx   map[string]int
	pending []string
}

func NewCSVSource(path string) (MetricsSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	r.ReuseRecord = true

	header, err := r.Read()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("csv read header: %w", err)
	}
	colIx := map[string]int{}
	for i, h := range header {
		colIx[normalize(h)] = i
	}
	required := []string{"activeusers", "throughputtps", "latencyms", "energynormalized"}
	for _, need := range required {
		if _, ok := colIx[need]; !ok {
			f.Close()
			return nil, fmt.Errorf("missing CSV column: %s", need)
		}
	}
	return &csvSource{file: f, reader: r, colIx: colIx}, nil
}

func (c *csvSource) Next() (NodeMetrics, error) {
	rec, err := c.reader.Read()
	if err != nil {
		return NodeMetrics{}, err
	}
	get := func(key string) string { return rec[c.colIx[key]] }

	active, err := strconv.Atoi(get("activeusers"))
	if err != nil {
		return NodeMetrics{}, fmt.Errorf("parse ActiveUsers: %w", err)
	}
	tps, err := strconv.ParseFloat(get("throughputtps"), 64)
	if err != nil {
		return NodeMetrics{}, fmt.Errorf("parse ThroughputTPS: %w", err)
	}
	lat, err := strconv.Atoi(get("latencyms"))
	if err != nil {
		return NodeMetrics{}, fmt.Errorf("parse LatencyMS: %w", err)
	}
	energy, err := strconv.ParseFloat(get("energynormalized"), 64)
	if err != nil {
		return NodeMetrics{}, fmt.Errorf("parse EnergyNormalized: %w", err)
	}
	return NodeMetrics{
		ActiveUsers:      active,
		ThroughputTPS:    tps,
		LatencyMS:        lat,
		EnergyNormalized: energy,
	}, nil
}

func (c *csvSource) Close() error {
	if c.file != nil {
		return c.file.Close()
	}
	return nil
}

func normalize(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			out = append(out, r+('a'-'A'))
		} else if r != ' ' && r != '\t' {
			out = append(out, r)
		}
	}
	return string(out)
}
