# Harkon-NFT (A Super App Portfolio Project)

Welcome to Harkon-NFT! This is a full-stack portfolio project demonstrating how to build a composable NFT "Super App" on the Flow blockchain.

This project includes Cadence smart contracts, a Go backend (API & Indexer using Echo + ent), and a React/Vite frontend (dApp).

### 🔮 Future Vision & Disclaimer
This project began as a technical portfolio to explore the advanced capabilities of the Flow blockchain (Composable NFTs, VRF, Storefronts).

However, due to the robust architecture and the positive feedback from the initial concept, this project is under active development and is intended to be improved and potentially launched as a serious, long-term product.

### ✨ Core Features

* **Composable NFTs:** `NFTMoment` (UGC photos) that can be physically "equipped" with `NFTAccessories`.
* **On-Chain Gacha:** An `AccessoryPack` that uses Flow's native VRF (on-chain randomness) via the "Receipt-as-NFT" commit-reveal pattern.
* **Proof of Attendance:** A "Soul-Bound Token" (`NFTEventPass`) system for online/offline events, used for "Earn-to-Mint."
* **Standard Marketplace:** Full integration with `NFTStorefrontV2` for listings and sales.
* **Secure UGC:** A secure backend flow for AI moderation and uploading assets to IPFS (Pinata).
* **API & Indexer:** A Go (Echo) backend connected to an `ent` database for fast off-chain queries.

---

### 🛠️ Tech Stack

* **Blockchain:** Flow, Cadence, Flow CLI
* **Backend:** Go (GoLang), Echo (API Framework), `ent` (ORM/Indexer), SQLite
* **Frontend:** React, TypeScript, Vite, TailwindCSS, FCL (Flow Client Library)
* **Asset Storage:** IPFS (via Pinata)

---

###  Prerequisites

Before you begin, ensure you have the following installed:
* [Go (v1.20+)](https://go.dev/doc/install)
* [Node.js (v18+)](https://nodejs.org/)
* [Flow CLI](https://developers.flow.com/tools/flow-cli/install)

---

### 🚀 How to Run the Project (Locally)

You will need **4 separate terminals** running to fully operate this project.

#### 1. Installation (Run Once)

First, install all dependencies for the frontend and backend.

```bash
# Install frontend (React/Vite) dependencies
cd frontend
npm install

# Install backend (Go/Echo/ent) dependencies
cd ../backend
go mod tidy
```

#### 2. Setup Backend .env
```bash
PRIVATE_KEY=
DATABASE_URL="postgres://username:pass@localhost:5432/db_name"
PINATA_JWT_KEY=
```

#### 3. Run FLOW CLI
in terminal 1
```bash
flow emulator start
```

in terminal 2
```bash
flow emulator start
```

in terminal 3
```bash
flow project deploy
```

#### 4. Run Backend and Frontend
in terminal 4
```bash
cd backend
go run ./api
```

in terminal 5
```bash
cd frontend
npm run dev
```
