# VERAONET Evaluation Report (Summary)

This document summarizes the evaluation settings and key results used in the paper.

## Environment
- Ganache CLI v7.9+, Geth v1.14.2
- Go 1.21+, Python 3.11, Node.js 20+
- Hardware: 8 vCPU, 32GB RAM, SSD

## Workloads
- Users: 100–10,000 concurrent
- Mix: provenance writes, artifact updates, reward distribution

## Metrics
- Latency (ms), Throughput (TPS), Gas (gwei), Energy (normalized)

## Results (Highlights)
- Gas ↓ ~82% vs PoW baseline (PoS/DPoS)
- Latency ↓ ~67% under high load (DPoS)
- Energy ↓ ~58% overall through adaptive switching

## Reproducibility
- Scripts: `experiments/` and `adaptive-consensus/`
- CSV outputs: `experiments/results/`
- Visualization: `visualization/plots.ipynb`
