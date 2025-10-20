// experiments/geth_tests/adaptive_switch.go
// CSV-driven demo for consensus selection.
//
// Build & run:
//   cd experiments/geth_tests
//   go run adaptive_switch.go --thresholds ../../adaptive-consensus/config_thresholds.json --metrics ../results/latency_results.csv
package main

import (
    "encoding/csv"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strconv"
)

type Consensus string
const (
    PoW  Consensus = "PoW"
    APoW Consensus = "APoW"
    PoS  Consensus = "PoS"
    DPoS Consensus = "DPoS"
)

type Thresholds struct {
    LowLoad          int     `json:"low_load"`
    MediumLoad       int     `json:"medium_load"`
    HighLoad         int     `json:"high_load"`
    LatencyThreshold int     `json:"latency_threshold"`
    EnergyThreshold  float64 `json:"energy_threshold"`
}

func loadThresholds(path string) (Thresholds, error) {
    f, err := os.Open(path)
    if err != nil { return Thresholds{}, err }
    defer f.Close()
    var t Thresholds
    if err := json.NewDecoder(f).Decode(&t); err != nil { return Thresholds{}, err }
    return t, nil
}

func decide(users, latency int, t Thresholds) Consensus {
    if latency > t.LatencyThreshold {
        if users >= t.MediumLoad { return DPoS }
        return PoS
    }
    switch {
    case users < t.LowLoad:
        return PoW
    case users < t.MediumLoad:
        return APoW
    case users < t.HighLoad:
        return PoS
    default:
        return DPoS
    }
}

func main() {
    thresholdsPath := flag.String("thresholds", "../../adaptive-consensus/config_thresholds.json", "Path to thresholds JSON")
    metricsCSV := flag.String("metrics", "../results/latency_results.csv", "Metrics CSV (Test,Consensus,Users,Latency_ms)")
    flag.Parse()

    t, err := loadThresholds(*thresholdsPath)
    if err != nil { log.Fatalf("load thresholds: %v", err) }

    f, err := os.Open(*metricsCSV)
    if err != nil { log.Fatalf("open metrics: %v", err) }
    defer f.Close()

    r := csv.NewReader(f)
    _, err = r.Read() // header
    if err != nil { log.Fatalf("read header: %v", err) }

    fmt.Printf("Using thresholds: low=%d med=%d high=%d latency_thr=%dms\n",
        t.LowLoad, t.MediumLoad, t.HighLoad, t.LatencyThreshold)

    outPath := filepath.Join(filepath.Dir(*metricsCSV), "decisions.csv")
    out, err := os.Create(outPath)
    if err != nil { log.Fatalf("create output: %v", err) }
    defer out.Close()
    w := csv.NewWriter(out)
    defer w.Flush()
    w.Write([]string{"Test","Users","Latency_ms","SelectedConsensus"})

    for {
        rec, err := r.Read()
        if err != nil {
            break
        }
        testID := rec[0]
        users, _ := strconv.Atoi(rec[2])
        lat, _ := strconv.Atoi(rec[3])

        c := decide(users, lat, t)
        w.Write([]string{testID, strconv.Itoa(users), strconv.Itoa(lat), string(c)})
        fmt.Printf("[Test %s] users=%d latency=%dms -> %s\n", testID, users, lat, c)
    }
    fmt.Printf("Wrote decisions to %s\n", outPath)
}