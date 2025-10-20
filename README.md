VERAONET: Adaptive Layer-2 Blockchain for Digital Heritage

Author: Omer Aziz
Affiliation: University of Management and Technology, Lahore, Pakistan
Email: f2018288006@umt.edu.pk

Associated Publication:

Aziz, O., Farooq, M.S., Khelifi, A., & Omer, A. (2025). VERAONET: A Virtual Ecosystem for Rewards and Archaeological Operations Network. npj Heritage Science.

📖 Overview

This repository provides the experimental implementation and simulation scripts used in the paper
“VERAONET: A Virtual Ecosystem for Rewards and Archaeological Operations Network” —
a Layer-2 blockchain framework designed for archaeological data management and digital heritage preservation.

VERAONET introduces a pluggable consensus mechanism that dynamically switches between:

Proof of Work (PoW)

Adjustable Proof of Work (APoW)

Proof of Stake (PoS)

Delegated Proof of Stake (DPoS)

based on real-time network load, system resources, and operational context.
The goal is to achieve scalability, sustainability, and cost efficiency for cultural heritage applications such as artifact provenance, milestone-based rewards, and museum record verification.

🧩 Repository Structure
├── adaptive-consensus/
│   ├── adaptive_switch.go          # Core logic for consensus switching
│   ├── metrics_collector.go        # Geth node metrics for latency/gas tracking
│   ├── config_thresholds.json      # Switching thresholds for PoW/APoW/PoS/DPoS
│
├── experiments/
│   ├── ganache_tests/              # Local Layer-2 simulations (10–500 users)
│   ├── geth_tests/                 # Public node simulations (50–10,000 users)
│   ├── results/                    # CSVs with latency, gas, and energy metrics
│
├── visualization/
│   ├── plots.ipynb                 # Notebook to reproduce Figures 4–5
│   └── comparative_energy_chart.png
│
├── docs/
│   ├── Algorithm1_AdaptiveConsensus.pdf
│   ├── SystemArchitecture_Figure3.pdf
│   └── EvaluationReport.md
│
└── README.md                       # (this file)

⚙️ Experimental Setup
1. Dependencies

Go Ethereum (Geth): v1.14.2 (base fork)

Ganache CLI: v7.9+

Python 3.11 / Jupyter Notebook

Node.js: v20+ (for local simulation scripts)

GoLang: v1.21+

2. Configuration

Edit config_thresholds.json to define switching triggers:

{
  "low_load": 250,
  "medium_load": 2500,
  "high_load": 7000,
  "latency_threshold": 500,
  "energy_threshold": 0.7
}

3. Running Simulations

Local Simulation (Ganache):

cd experiments/ganache_tests
python3 run_simulation.py --users 500 --consensus adaptive


Public Geth Testbed:

cd experiments/geth_tests
go run adaptive_switch.go --config config_thresholds.json


Results are stored automatically in /results as CSVs.

📊 Results Summary

Gas Consumption: ↓ 82% compared to Ethereum PoW baseline

Latency: ↓ 67% under high-load DPoS scenarios

Energy Consumption: ↓ 58% overall

Transaction Throughput: ↑ 2.5× improvement on simulated 10k-node network

🧠 Integration with Virtual Museum Platform

This implementation links directly with the companion repository:
🔗 Virtual-Reality-Museum

That repository demonstrates how VERAONET integrates with ArchaeoMeta and Archaeological Workflows, including:

Artifact upload and provenance verification

Smart-contract reward distribution

Decentralized governance and reputation tracking

⚖️ License & Attribution

This work extends the official Go Ethereum (Geth) implementation:

https://github.com/ethereum/go-ethereum

Licensed under GNU LGPL v3.0

All derivative components and simulation scripts in this repository are released under the same license,
with attribution to the original Geth developers.

🧾 Citation

If you use this code or dataset, please cite:

@article{aziz2025verao,
  title={VERAONET: A Virtual Ecosystem for Rewards and Archaeological Operations Network},
  author={Aziz, Omer and Farooq, Muhammad Shoaib and Khelifi, Adel and Omer, Abdullah},
  journal={npj Heritage Science},
  year={2025}
}

💡 Acknowledgment

This work was supported by Abu Dhabi University and the University of Management and Technology, Lahore,
as part of the collaborative initiative on Digital Heritage and Blockchain Technologies.
