```markdown
# Go Microservice Template

### Web2 + Web3 Hybrid Chassis

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean_Architecture_%7C_DDD-orange)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Stars](https://img.shields.io/github/stars/zy99978455-otw/go-micro-template?style=social)

An **enterprise-grade**, **production-ready** microservice chassis designed for the transition from Web2 to Web3.  
This framework provides a unified infrastructure for building **hybrid applications** (MySQL/Redis + Blockchain), engineered with **Pure Dependency Injection (DI)** and **Domain-Driven Design (DDD)** principles.

---

## 🚀 Key Features

### 📐 Clean Architecture & DI
- **Pure Dependency Injection** — No global variables (`global.DB` removed). Components are explicitly initialized and wired in `main.go`.
- **Wire Ready** — Built-in `ProviderSet` definitions for Google Wire support, making future expansion effortless.
- **DDD Layering** — Strict separation of concerns: `Server` (Transport) → `Biz` (Domain) → `Data` (Repository).

### 🌐 Hybrid Infrastructure (Web2 + Web3)
- **Seamless Integration** — Run Web2 (User/Order) and Web3 (Wallet/Indexer) logic in the same process.
- **Unified Data Layer** — Centralized repository management for MySQL (GORM), Redis, and RPC Clients.

### 🔗 High-Availability RPC Manager (The Core)
- **Multi-Chain Support** — Config-driven multi-chain setup (Ethereum, BSC, Polygon, etc.).
- **Health Checks** — Background worker pool that periodically checks RPC node latency and block height.
- **Load Balancing** — Automatically routes requests to the healthiest and fastest RPC node.
- **Auto-Failover** — Smartly switches to backup nodes upon connection failure.

### 🛡️ Microservice Governance
- **Service Discovery** — Built-in **Consul** registration with Docker-friendly IP resolution (`register_ip`).
- **Graceful Shutdown** — Unified cleanup mechanism ensuring database connections and servers are closed safely.
- **Config-Driven** — Fully dynamic `config.yaml` to switch modes without code changes.

---

## 📂 Project Structure

Following standard Go layout and Clean Architecture:

```
├── cmd/                    # Main application entry points
├── configs/                # Configuration files (YAML)
├── internal/
│   ├── biz/                # Domain logic & Interfaces (Pure Go)
│   ├── data/               # Data access implementation (DB, RPC, Redis)
│   │   └── rpc_manager.go  # 🔥 Core RPC Load Balancer
│   └── server/             # Transport layer (HTTP/gRPC)
│       ├── http.go         # Router registration & DI Wiring
│       └── chain_handler.go# Web3 API Handlers
├── pkg/                    # Infrastructure libraries (Logger, DB Drivers)
└── README.md
```

---

## 🛠️ Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose (Optional)
- An Ethereum/BSC RPC URL (Infura/Alchemy/Ankr)

### 1. Installation
```bash
git clone https://github.com/zy99978455-otw/go-micro-template.git
cd go-micro-template
go mod tidy
```

### 2. Configuration
Copy the debug config and setup your environment:
```bash
cp configs/config-debug.yaml configs/config-local.yaml
```
Edit `configs/config-local.yaml` (example for Web3 mode):
```yaml
chains:
  - chain_id: 1
    chain_name: "eth_mainnet"
    rpc_url: "https://rpc.ankr.com/eth"
```

### 3. Run
```bash
go run cmd/server/main.go
```

Startup logs will show:
```
INFO ... ✅ [Web2] MySQL Connected
INFO ... ✅ [Web3] RPC Node Added: ChainID 1
INFO ... ✅ [验证成功] 通过 RPCManager 拿到了客户端! ChainID: 1
INFO ... ✅ Consul Service Registered
```

---

## 📡 API Reference

Built-in Web3 endpoints powered by the RPC Manager.

### Get Block Height
- **URL**: `/api/v1/web3/block`
- **Method**: `GET`
- **Query Params**:
  - `chain_id` (int, optional): e.g. `1` for ETH, `56` for BSC. Default: `1`

**Response Example:**
```json
{
  "code": 200,
  "data": {
    "chain_id": 1,
    "height": 24080901
  }
}
```

---

## 🧩 Architecture Overview

The `internal/data/rpc_manager.go` implements a robust multi-chain RPC load balancer:

1. **Init** — Loads chain configs via DI and dials all endpoints
2. **Monitor** — Background goroutine checks latency & block height every 30s
3. **Serve** — `GetClient(chainID)` returns the healthiest node

---

## 📝 Roadmap

- [x] Refactor to Pure DI & Clean Architecture
- [x] Web2 Infrastructure (MySQL/Redis/Zap)
- [x] Web3 Infrastructure (High-Availability RPC Manager)
- [x] Service Discovery (Consul)
- [ ] Code Generation for Repository Layer
- [ ] Prometheus Metrics Integration
- [ ] Distributed Tracing (OpenTelemetry)

---

## 🤝 Contribution

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📄 License

MIT License
```