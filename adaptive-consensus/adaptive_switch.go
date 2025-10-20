// adaptive-consensus/adaptive_switch.go
// VERAONET: Adaptive consensus switcher (PoW/APoW/PoS/DPoS)
// Usage examples:
//   go run . --config ./config_thresholds.json --oneshot
//   go run . --config ./config_thresholds.json --interval 5
//
// This file depends on metrics helpers in metrics_collector.go (same package).

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// Consensus enumerates supported mechanisms
type Consensus string

const (
	ConsensusPoW  Consensus = "PoW"
	ConsensusAPoW Consensus = "APoW"
	ConsensusPoS  Consensus = "PoS"
	ConsensusDPoS Consensus = "DPoS"
)

// Thresholds config structure (loaded from JSON)
type Thresholds struct {
	LowLoad          int     `json:"low_load"`            // e.g., 250 users/tps
	MediumLoad       int     `json:"medium_load"`         // e.g., 2500
	HighLoad         int     `json:"high_load"`           // e.g., 7000
	LatencyThreshold int     `json:"latency_threshold"`   // ms, e.g., 500
	EnergyThreshold  float64 `json:"energy_threshold"`    // normalized 0..1 (or J/tx if you prefer)
}

// DecisionInput bundles metrics + thresholds for the policy
type DecisionInput struct {
	Metrics    NodeMetrics
	Thresholds Thresholds
}

// decideConsensus implements a simple, explainable policy.
// You can replace this with a more advanced RL-based tuner later.
func decideConsensus(in DecisionInput) (Consensus, string) {
	m := in.Metrics
	t := in.Thresholds

	// 1) Guardrails: if latency or energy are critical, prefer stake-based modes
	if m.LatencyMS > t.LatencyThreshold || m.EnergyNormalized > t.EnergyThreshold {
		// If load is high, DPoS; else PoS
		if m.ActiveUsers >= t.MediumLoad {
			return ConsensusDPoS, "High latency/energy & medium/high load → DPoS"
		}
		return ConsensusPoS, "High latency/energy & low/medium load → PoS"
	}

	// 2) Load-based switching
	switch {
	case m.ActiveUsers < t.LowLoad:
		return ConsensusPoW, "Low load → PoW (max security; low cost impact)"
	case m.ActiveUsers < t.MediumLoad:
		return ConsensusAPoW, "Medium-ish load → APoW (difficulty adapts)"
	case m.ActiveUsers < t.HighLoad:
		return ConsensusPoS, "High load (tier 1) → PoS"
	default:
		return ConsensusDPoS, "Very high load → DPoS (maximize throughput)"
	}
}

func main() {
	var (
		configPath = flag.String("config", "./config_thresholds.json", "Path to thresholds JSON config")
		interval   = flag.Int("interval", 0, "Loop interval in seconds (0 = oneshot)")
		oneshot    = flag.Bool("oneshot", false, "Run a single decision and exit")
		inputCSV   = flag.String("metrics-csv", "", "Optional CSV file to stream metrics from (overrides random generator)")
	)
	flag.Parse()

	thr, err := loadThresholds(*configPath)
	if err != nil {
		log.Fatalf("failed to load thresholds: %v", err)
	}

	// If oneshot is set, force interval=0 behavior
	if *oneshot {
		*interval = 0
	}

	// Wire up the metrics source (CSV or random sim)
	var source MetricsSource
	if *inputCSV != "" {
		csvSrc, err := NewCSVSource(*inputCSV)
		if err != nil {
			log.Fatalf("failed to open CSV metrics source: %v", err)
		}
		defer csvSrc.Close()
		source = csvSrc
	} else {
		source = NewRandomSource()
	}

	// Run once or loop
	if *interval <= 0 {
		m, err := source.Next()
		if err != nil {
			log.Fatalf("failed to read metrics: %v", err)
		}
		cons, why := decideConsensus(DecisionInput{Metrics: m, Thresholds: thr})
		printDecision(m, thr, cons, why)
		return
	}

	// Graceful shutdown
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			m, err := source.Next()
			if err != nil {
				log.Printf("metrics read error: %v", err)
				continue
			}
			cons, why := decideConsensus(DecisionInput{Metrics: m, Thresholds: thr})
			printDecision(m, thr, cons, why)
		case <-done:
			fmt.Println("\nreceived interrupt, exiting.")
			return
		}
	}
}

func loadThresholds(path string) (Thresholds, error) {
	f, err := os.Open(path)
	if err != nil {
		return Thresholds{}, err
	}
	defer f.Close()
	var t Thresholds
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return Thresholds{}, err
	}
	// Minimal sanity checks
	if t.LowLoad <= 0 || t.MediumLoad <= 0 || t.HighLoad <= 0 {
		return Thresholds{}, fmt.Errorf("invalid load thresholds in %s", path)
	}
	return t, nil
}

func printDecision(m NodeMetrics, t Thresholds, c Consensus, why string) {
	fmt.Printf(
		"\n[VERAONET] decision @ %s\n"+
			"   users=%d  tps=%.1f  latency_ms=%d  energy_norm=%.2f\n"+
			"   thresholds: low=%d  med=%d  high=%d  L_thr=%dms  E_thr=%.2f\n"+
			"   → SELECTED CONSENSUS: %s\n   reason: %s\n",
		time.Now().Format(time.RFC3339),
		m.ActiveUsers, m.ThroughputTPS, m.LatencyMS, m.EnergyNormalized,
		t.LowLoad, t.MediumLoad, t.HighLoad, t.LatencyThreshold, t.EnergyThreshold,
		c, why,
	)
}
