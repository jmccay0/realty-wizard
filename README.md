# RealtyWizard - Texas Real Estate Transaction Management

A TurboTax-style application for guiding Texas home sellers through residential resale transactions.

## ✨ New in v3.0: Buyer-Initiated Transactions & PostgreSQL

- **Buyer-initiated transactions**: Buyers can now start their own transaction workflows
- **Role-based project creation**: Choose buyer or seller role when creating projects
- **PostgreSQL support**: Production-ready database with full SQLite compatibility
- **Dual database mode**: Switch between SQLite and PostgreSQL via environment variables

## v2.0 Features: User Authentication & Multi-Participant Support

- Secure JWT-based authentication
- Role-based access (buyer/seller/admin)
- Multi-participant transactions
- Project-specific roles

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
- SQLite / PostgreSQL (dual database support)
- JWT authentication
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

**Option 1: SQLite (Quick Start - Default)**
```bash
cd backend
go mod download

# Optional: Copy environment file
cp .env.example .env
# Edit .env with your JWT_SECRET (DB_TYPE defaults to sqlite)

# Seed database with admin user
go run cmd/seed/main.go

# Start server
go run cmd/api/main.go
# Server runs on http://localhost:8080
```

**Option 2: PostgreSQL (Production-Ready)**
```bash
cd backend
go mod download

# Start PostgreSQL with Docker
cd ..
docker-compose up -d

# Use PostgreSQL environment configuration
cp backend/.env.postgres backend/.env
# Edit .env with your JWT_SECRET if needed

# Seed database with admin user
cd backend
DB_TYPE=postgres POSTGRES_URL="postgres://realwiz:realwiz_dev_password@localhost:5432/realty_wizard?sslmode=disable" go run cmd/seed/main.go

# Start server with PostgreSQL
DB_TYPE=postgres POSTGRES_URL="postgres://realwiz:realwiz_dev_password@localhost:5432/realty_wizard?sslmode=disable" go run cmd/api/main.go
# Server runs on http://localhost:8080
```

**Default Admin Credentials:**
- Email: `admin@realwiz.local`
- Password: `admin123`
- **⚠️ Change password after first login!**

**Environment Variables:**
- `DB_TYPE`: `sqlite` (default) or `postgres`
- `DB_PATH`: Path to SQLite database (default: `./data/realty-wizard.db`)
- `POSTGRES_URL`: PostgreSQL connection string (required when `DB_TYPE=postgres`)
- `AUTH_ENABLED`: `true` (default) or `false`
- `JWT_SECRET`: Secret key for JWT tokens (required for production)

#### Frontend Setup
```bash
cd frontend
npm install

# Optional: Copy environment file
cp .env.example .env

# Start dev server
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
3. Database: SQLite in `backend/data/realty-wizard.db` OR PostgreSQL via Docker on `:5432`

## Buyer vs Seller Workflows

**When creating a new project, users select their role:**

**Seller Workflow:**
- Seller provides their name(s), contact info
- System guides through property disclosure
- Generates listing documents
- Tracks deadlines from listing to close

**Buyer Workflow:**
- Buyer provides their name, contact info
- System guides through offer preparation
- Generates purchase documents
- Tracks deadlines from contract to close

## Phase 1 Features

- Property profile setup
- Seller's Disclosure guided workflow
- TREC form awareness and preparation
- Deadline calculation engine
- Calendar export (.ics)
- Document checklist generation

## License

Private project - All rights reserved
