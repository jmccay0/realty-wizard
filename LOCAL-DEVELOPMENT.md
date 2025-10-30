# Local Development Guide

This guide will help you run the Realty Wizard application on your local machine for development and testing.

## Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** and npm - [Download](https://nodejs.org/)
- A terminal/command line

## Quick Start (5 minutes)

### 1. Start the Backend API

Open a terminal and run:

```bash
cd backend
go run cmd/api/main.go
```

You should see:
```
🚀 Server starting on :8080
```

The backend API is now running at `http://localhost:8080`

**Keep this terminal open!**

---

### 2. Start the Frontend (New Terminal)

Open a **new terminal window** and run:

```bash
cd frontend
npm install  # Only needed first time
npm run dev
```

You should see:
```
  VITE v7.1.10  ready in 234 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

The frontend is now running at `http://localhost:5173`

**Keep this terminal open too!**

---

### 3. Open the Application

Open your browser and go to:
```
http://localhost:5173
```

You should see the Realty Wizard home page!

---

## Using the Application

### Create Your First Transaction

1. Click **"Start New Transaction"**
2. Fill out the form:
   - Property Address: `123 Main St, Austin, TX 78701`
   - Seller Names: `John Doe, Jane Doe`
   - Email: `seller@example.com`
   - Phone: `512-555-0100`
   - Check "I have a real estate agent"
   - Agent Name: `Agent Smith`
3. Click **"Create Project"**

### Complete Property Details

1. You'll be redirected to the property setup wizard
2. Fill out property information:
   - Year Built: `2015`
   - Legal Description: `Lot 42, Block 7, Test Subdivision`
   - Tax ID: `12345-67890`
   - Check "This is the seller's homestead"
   - HOA: Select "Yes" or "No"
   - Survey: Select if you have one
3. Continue to Disclosure
4. Fill out seller's disclosure questions
5. Click **"Complete Setup"**

### Enter Contract Terms

1. From the dashboard, click **"Enter Contract Terms"**
2. Fill out the contract form:
   - **Effective Date:** Today's date or contract date
   - **Sales Price:** e.g., `450000`
   - **Closing Date:** 45 days from effective date
   - **Option Fee:** e.g., `500`
   - **Option Period:** `10` days (typical)
   - **Earnest Money:** e.g., `5000`
   - **Title Commitment Days:** `20` (default)
   - Check **"Buyer will use financing"** if applicable
   - Check **"Seller will remain in property"** if doing a lease-back
3. Click **"Create Contract & Generate Deadlines"**

### View Deadlines

After creating the contract, the dashboard will show:
- **Contract Summary** - Sales price, closing date, option details
- **Upcoming Deadlines** - Auto-generated deadlines sorted by date
  - Option fee due (3 business days)
  - Earnest money due (3 business days)
  - Inspection period ends
  - Title commitment due
  - Loan approval deadline (if financing)
  - And more!

---

## Development Workflow

### Backend Development

The backend uses **hot reload** - just save your Go files and the server will restart automatically.

**Edit backend code:**
```bash
# Edit files in backend/internal/
# Save changes
# Server auto-restarts
```

**Run tests:**
```bash
cd backend
go test ./... -v
```

**Check for errors:**
```bash
cd backend
go build ./...
```

---

### Frontend Development

Vite provides **instant hot reload** - save your files and see changes immediately in the browser!

**Edit frontend code:**
```bash
# Edit files in frontend/src/
# Save changes
# Browser updates automatically
```

**Run tests:**
```bash
cd frontend
npm test
```

**Build for production:**
```bash
cd frontend
npm run build
# Output in frontend/dist/
```

---

## API Endpoints

The backend exposes these REST endpoints:

### Projects
- `GET /api/projects` - List all projects
- `POST /api/projects` - Create new project
- `GET /api/projects/:id` - Get project details
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project
- `GET /api/projects/:id/summary` - Get full project summary

### Property & Disclosure
- `POST /api/projects/:id/property` - Create property details
- `GET /api/projects/:id/property` - Get property
- `PUT /api/projects/:id/property` - Update property
- `POST /api/projects/:id/disclosure` - Create disclosure
- `GET /api/projects/:id/disclosure` - Get disclosure
- `PUT /api/projects/:id/disclosure` - Update disclosure

### Contracts & Deadlines
- `POST /api/projects/:id/contract` - Create contract (auto-generates deadlines)
- `GET /api/projects/:id/contract` - Get contract
- `GET /api/projects/:id/deadlines` - List deadlines
- `PUT /api/deadlines/:id` - Update deadline
- `DELETE /api/deadlines/:id` - Delete deadline

---

## Testing the API with curl

```bash
# Create a project
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "property_address": "123 Main St, Austin, TX 78701",
    "seller_names": ["John Doe"],
    "seller_email": "john@example.com",
    "seller_phone": "512-555-0100"
  }'

# List all projects
curl http://localhost:8080/api/projects

# Get project summary
curl http://localhost:8080/api/projects/{project-id}/summary
```

---

## Database

The backend uses **SQLite** stored in:
```
backend/data/realty-wizard.db
```

**Reset database:**
```bash
# Stop the backend server (Ctrl+C)
rm backend/data/realty-wizard.db
# Restart the server - fresh database created automatically
```

**View database contents:**
```bash
# Install sqlite3 if needed
sqlite3 backend/data/realty-wizard.db

# Run queries
sqlite> .tables
sqlite> SELECT * FROM projects;
sqlite> .quit
```

---

## Environment Variables

### Backend (Optional)

Create `backend/.env`:
```bash
PORT=8080
DB_PATH=./data/realty-wizard.db
```

### Frontend (Development)

The frontend automatically uses `http://localhost:8080` for the API in development mode.

To change it, create `frontend/.env.local`:
```bash
VITE_API_URL=http://localhost:8080/api
```

---

## Troubleshooting

### Backend won't start

**Error: "port already in use"**
```bash
# Find and kill the process using port 8080
lsof -ti:8080 | xargs kill -9
```

**Error: "database is locked"**
```bash
# Another process is using the database
# Stop all backend instances and restart
```

### Frontend won't start

**Error: "EADDRINUSE"**
```bash
# Port 5173 is in use
# Kill the process or use a different port
npm run dev -- --port 3000
```

**Error: "Cannot find module"**
```bash
# Dependencies not installed
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Tests failing

**Backend tests:**
```bash
cd backend
go mod tidy  # Update dependencies
go test ./... -v
```

**Frontend tests:**
```bash
cd frontend
rm -rf node_modules
npm install
npm test
```

### Browser shows blank page

1. Check browser console (F12) for errors
2. Verify backend is running: `curl http://localhost:8080/api/projects`
3. Check frontend terminal for build errors
4. Clear browser cache (Cmd+Shift+R or Ctrl+Shift+R)

---

## VS Code Setup (Recommended)

Install these extensions:
- **Go** (golang.go) - Go language support
- **ES7+ React/Redux snippets** - React snippets
- **Prettier** - Code formatter
- **ESLint** - JavaScript linting

Create `.vscode/launch.json` for debugging:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/api/main.go"
    }
  ]
}
```

---

## Next Steps

- **Deploy to AWS:** See [AWS-QUICKSTART.md](AWS-QUICKSTART.md)
- **Run tests:** `go test ./...` and `npm test`
- **Add features:** Check out the codebase structure below

## Codebase Structure

```
realty-wizard/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go           # Local server entry point
│   │   └── lambda/main.go        # AWS Lambda entry point
│   ├── internal/
│   │   ├── handlers/             # HTTP request handlers
│   │   │   ├── handlers.go       # Main handlers
│   │   │   └── handlers_test.go  # Handler tests
│   │   ├── models/               # Data structures
│   │   │   └── models.go
│   │   ├── rules/                # Business logic
│   │   │   ├── deadlines.go      # Deadline calculation
│   │   │   └── deadlines_test.go # Deadline tests
│   │   └── storage/              # Database layer
│   │       ├── storage.go        # Interface
│   │       └── sqlite.go         # SQLite implementation
│   └── data/                     # SQLite database files
│
├── frontend/
│   ├── src/
│   │   ├── pages/                # Page components
│   │   │   ├── Home.tsx
│   │   │   ├── NewProject.tsx
│   │   │   ├── Wizard.tsx
│   │   │   ├── Dashboard.tsx
│   │   │   ├── EnterContract.tsx
│   │   │   └── EnterContract.test.tsx
│   │   ├── api.ts                # API client
│   │   ├── types.ts              # TypeScript types
│   │   └── index.css             # Global styles
│   └── dist/                     # Build output (production)
│
└── docs/
    ├── LOCAL-DEVELOPMENT.md      # This file
    ├── AWS-QUICKSTART.md         # AWS deployment
    └── DEPLOYMENT.md             # Detailed deployment guide
```

---

## Have Questions?

- Check the [AWS deployment guide](AWS-QUICKSTART.md) for production setup
- Review test files for usage examples
- Check GitHub issues for known problems

Happy coding! 🚀
