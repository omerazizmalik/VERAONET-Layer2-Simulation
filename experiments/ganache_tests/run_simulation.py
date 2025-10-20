#!/usr/bin/env python3
"""
VERAONET – Ganache local simulation
Generates synthetic metrics and writes CSVs in ../results
Usage:
  python3 run_simulation.py --users 500 --steps 50 --out ../results
"""
import argparse, csv, random
from pathlib import Path

def simulate_row(users:int):
    if users < 500:
        tps = random.uniform(50, 170)
    elif users < 5000:
        tps = random.uniform(200, 600)
    else:
        tps = random.uniform(500, 1400)
    latency = int(random.uniform(80, 220))
    energy = random.uniform(0.05, 0.20)
    return {
        "ActiveUsers": users,
        "ThroughputTPS": round(tps, 2),
        "LatencyMS": latency,
        "EnergyNormalized": round(energy, 3),
    }

def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--users", type=int, default=500)
    ap.add_argument("--steps", type=int, default=50)
    ap.add_argument("--out", type=str, default="../results")
    args = ap.parse_args()

    outdir = Path(args.out)
    outdir.mkdir(parents=True, exist_ok=True)

    latency_path = outdir / "latency_results.csv"
    gas_path = outdir / "gas_usage.csv"
    energy_path = outdir / "energy_comparison.csv"

    with open(latency_path, "w", newline="") as f_lat, \
         open(gas_path, "w", newline="") as f_gas, \
         open(energy_path, "w", newline="") as f_energy:

        lat_writer = csv.writer(f_lat)
        gas_writer = csv.writer(f_gas)
        energy_writer = csv.writer(f_energy)

        lat_writer.writerow(["Test","Consensus","Users","Latency_ms"])
        gas_writer.writerow(["Test","Consensus","Users","Gas_gwei"])
        energy_writer.writerow(["Test","Consensus","Users","EnergyNormalized"])

        cycles = ["PoW","APoW","PoS","DPoS"]
        for i in range(1, args.steps+1):
            users = max(5, int(random.gauss(args.users, max(10, args.users*0.15))))
            row = simulate_row(users)
            consensus = cycles[(i-1) % len(cycles)]

            base_lat = row["LatencyMS"]
            if consensus == "PoW":
                lat = int(base_lat * random.uniform(2.0, 3.0))
            elif consensus == "APoW":
                lat = int(base_lat * random.uniform(1.4, 2.1))
            elif consensus == "PoS":
                lat = int(base_lat * random.uniform(0.6, 0.9))
            else:
                lat = int(base_lat * random.uniform(0.45, 0.75))
            lat_writer.writerow([i, consensus, users, lat])

            if consensus == "PoW":
                gas = int(random.uniform(30000, 50000))
            elif consensus == "APoW":
                gas = int(random.uniform(26000, 47000))
            elif consensus == "PoS":
                gas = int(random.uniform(2000, 5000))
            else:
                gas = int(random.uniform(5000, 12000))
            gas_writer.writerow([i, consensus, users, gas])

            if consensus == "PoW":
                en = round(random.uniform(0.75, 0.95), 3)
            elif consensus == "APoW":
                en = round(random.uniform(0.55, 0.75), 3)
            elif consensus == "PoS":
                en = round(random.uniform(0.10, 0.25), 3)
            else:
                en = round(random.uniform(0.15, 0.30), 3)
            energy_writer.writerow([i, consensus, users, en])

    print(f"[OK] Wrote:\n  - {latency_path}\n  - {gas_path}\n  - {energy_path}")

if __name__ == "__main__":
    main()