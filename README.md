# Realty Wizard - Texas Real Estate Transaction Guide

A TurboTax-style application for guiding Texas home sellers through residential resale transactions.

## Project Structure

```
realty-wizard/
├── backend/          # Go API server
│   ├── cmd/api/     # API server entry point
│   ├── internal/    # Private application code
│   │   ├── models/      # Domain models
│   │   ├── rules/       # Business rules & deadline engine
│   │   ├── storage/     # Data persistence layer
│   │   └── handlers/    # HTTP handlers
│   └── pkg/trec/    # Public TREC form utilities
├── cli/             # Go CLI tool (future)
├── frontend/        # React web UI
├── shared/          # Shared resources
│   ├── forms/       # TREC PDF/DOCX templates
│   └── docs/        # Reference documentation
└── docs/            # Project documentation
```

## Tech Stack

**Backend:**
- Go 1.21+
- Chi router (lightweight, idiomatic)
- SQLite (local-first storage)
- Testify (testing)

**Frontend:**
- React 18
- TypeScript
- Tailwind CSS (RealWiz design)
- React Router
- Axios

**Development:**
- Air (Go hot reload)
- Vite (frontend dev server)

## Getting Started

### Local Development

#### Prerequisites
- Go 1.21+
- Node.js 18+
- Make (optional)

#### Backend Setup
```bash
cd backend
go mod download
go run cmd/api/main.go
# Server runs on http://localhost:8080
```

#### Frontend Setup
```bash
cd frontend
npm install
npm run dev
# Dev server runs on http://localhost:5173
```

#### Running Tests
```bash
cd backend
go test ./...
```

## 🚀 AWS Deployment

Deploy to AWS serverless infrastructure (~$5/month):

**Quick Start (10 minutes):** See [AWS-QUICKSTART.md](AWS-QUICKSTART.md)

**Detailed Guide:** See [DEPLOYMENT.md](DEPLOYMENT.md)

**One-Command Deploy:**
```bash
./deploy-all.sh
```

This deploys:
- Backend: AWS Lambda + API Gateway
- Frontend: S3 + CloudFront CDN
- Database: SQLite on Lambda (ephemeral) or EFS (persistent)

## Development Workflow

1. Backend API runs on `:8080`
2. Frontend dev server on `:5173` with proxy to backend
3. SQLite database stored in `backend/data/realty-wizard.db`

## Phase 1 Features

- Property profile setup
- Seller's Disclosure guided workflow
- TREC form awareness and preparation
- Deadline calculation engine
- Calendar export (.ics)
- Document checklist generation

## License

Private project - All rights reserved
