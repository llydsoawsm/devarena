# DevCode Execution & Setup Guide (`exe.md`)

This document provides step-by-step instructions to install all prerequisites, set up the database and services, scaffold the monorepo, and run the complete **DevCode** platform locally on Windows.

---

## 1. System Prerequisites & Installation

Open **PowerShell as Administrator** or run the following commands via terminal.

### 1.1 Install Go (1.22+)
```powershell
# Using winget
winget install GoLang.Go

# After installation, verify in a new terminal:
go version
```
*Manual download if preferred: [https://go.dev/dl/](https://go.dev/dl/)*

---

### 1.2 Install Docker Desktop (for Code Execution Sandbox & Services)
```powershell
# Install WSL2 (required for Docker Desktop on Windows)
wsl --install

# Install Docker Desktop
winget install Docker.DockerDesktop
```
* After installation, restart Windows if prompted.
* Start Docker Desktop and ensure the engine is running (`docker info`).

---

### 1.3 Install pnpm (Package Manager for Next.js)
```powershell
# Node.js v24 is already installed on your machine. Install pnpm globally:
npm install -g pnpm

# Verify:
pnpm -v
```

---

### 1.4 Database & Cache Options

You have two options to run PostgreSQL and Redis:

#### Option A: Docker Compose (Recommended)
Once Docker Desktop is running, launch both PostgreSQL and Redis with a single command from the project root:
```powershell
docker compose up -d
```

#### Option B: Standalone Services
- **Redis**: Already available in your Laragon directory at `C:\laragon\bin\redis\redis-x64-5.0.14.1\redis-server.exe`.
- **PostgreSQL**: Install via winget:
  ```powershell
  winget install PostgreSQL.PostgreSQL.16
  ```

---

## 2. Project Directory Structure

```
c:\Users\Lloyd Pili\projects\devcode\
├── backend/                     # Go 1.22+ API & WebSocket Server
│   ├── cmd/server/main.go       # Entry point
│   ├── internal/
│   │   ├── config/              # Environment config
│   │   ├── handlers/            # REST API controllers
│   │   ├── middleware/          # JWT, CORS, Rate limiting
│   │   ├── models/              # GORM database models
│   │   ├── services/            # Matchmaking, Match, Rating, Judge
│   │   ├── ws/                  # Gorilla WebSocket hub & rooms
│   │   └── judge/               # Docker Engine SDK sandbox runner
│   ├── docker/                  # Dockerfiles for the 8 language runners
│   │   ├── python/Dockerfile
│   │   ├── javascript/Dockerfile
│   │   ├── typescript/Dockerfile
│   │   ├── cpp/Dockerfile
│   │   ├── go/Dockerfile
│   │   ├── java/Dockerfile
│   │   ├── rust/Dockerfile
│   │   └── csharp/Dockerfile
│   ├── .env
│   ├── go.mod
│   └── go.sum
├── frontend/                    # Next.js 15 (React 19) Frontend
│   ├── src/
│   │   ├── app/                 # App Router (Dashboard, Match, Results, Admin)
│   │   ├── components/          # shadcn/ui & Monaco Editor components
│   │   ├── hooks/               # useWebSocket, useMatch, useAuth
│   │   ├── lib/                 # API client, types, utilities
│   │   └── stores/              # Zustand state store
│   ├── .env.local
│   ├── package.json
│   └── tsconfig.json
├── docker-compose.yml           # PostgreSQL + Redis local orchestration
├── plan-devcode.md              # Original specification
├── devcode-implementation-plan.md # Approved architecture plan
└── exe.md                       # This file
```

---

## 3. Environment Configuration Files

### 3.1 Root `docker-compose.yml`
Save this in `c:\Users\Lloyd Pili\projects\devcode\docker-compose.yml`:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: devcode-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: devcode
      POSTGRES_USER: devcode
      POSTGRES_PASSWORD: devcode_secret_password
    ports:
      - "5432:5432"
    volumes:
      - devcode_pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    container_name: devcode-redis
    restart: unless-stopped
    command: redis-server --appendonly yes
    ports:
      - "6379:6379"
    volumes:
      - devcode_redisdata:/data

volumes:
  devcode_pgdata:
  devcode_redisdata:
```

---

### 3.2 Backend Environment File (`backend/.env`)
Create `c:\Users\Lloyd Pili\projects\devcode\backend\.env`:

```ini
# Server
PORT=8080
APP_ENV=development
CLIENT_URL=http://localhost:3000

# PostgreSQL
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=devcode
DB_PASSWORD=devcode_secret_password
DB_NAME=devcode
DB_SSLMODE=disable

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=

# JWT Secrets
JWT_SECRET=super_secret_devcode_jwt_access_key_change_in_production_32b
JWT_REFRESH_SECRET=super_secret_devcode_jwt_refresh_key_change_in_prod
JWT_EXPIRATION_MINUTES=15
JWT_REFRESH_EXPIRATION_DAYS=7

# Judge Execution
DOCKER_HOST=npipe:////./pipe/docker_engine
JUDGE_MAX_WORKERS=4
JUDGE_TIMEOUT_SECONDS=5
```

---

### 3.3 Frontend Environment File (`frontend/.env.local`)
Create `c:\Users\Lloyd Pili\projects\devcode\frontend\.env.local`:

```ini
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
```

---

## 4. Step-by-Step Project Setup & Execution

### Step 1: Start Infrastructure (PostgreSQL & Redis)
In `c:\Users\Lloyd Pili\projects\devcode`:
```powershell
docker compose up -d
```
Verify containers are active:
```powershell
docker compose ps
```

---

### Step 2: Initialize & Run Backend

1. Navigate to the backend directory:
   ```powershell
   cd c:\Users\Lloyd Pili\projects\devcode\backend
   ```

2. Initialize Go module and install dependencies:
   ```powershell
   go mod init devcode/backend

   # Framework & Routing
   go get -u github.com/gin-gonic/gin
   go get -u github.com/gin-contrib/cors
   go get -u github.com/gorilla/websocket

   # Database (GORM + Postgres)
   go get -u gorm.io/gorm
   go get -u gorm.io/driver/postgres
   go get -u gorm.io/datatypes

   # Cache & Queues (Redis)
   go get -u github.com/redis/go-redis/v9

   # Auth & Utilities
   go get -u github.com/golang-jwt/jwt/v5
   go get -u golang.org/x/crypto/bcrypt
   go get -u github.com/oklog/ulid/v2
   go get -u github.com/joho/godotenv

   # Docker Engine SDK
   go get -u github.com/docker/docker/client
   ```

3. Run the Go server:
   ```powershell
   go run cmd/server/main.go
   ```
   *The server starts on `http://localhost:8080` and auto-runs database migrations.*

---

### Step 3: Build Judge Sandbox Docker Images

Build the lightweight runner images for the supported programming languages:

```powershell
cd c:\Users\Lloyd Pili\projects\devcode\backend\docker

# 1. Python
docker build -t devcode-runner-python ./python

# 2. JavaScript (Node.js)
docker build -t devcode-runner-javascript ./javascript

# 3. TypeScript
docker build -t devcode-runner-typescript ./typescript

# 4. C++
docker build -t devcode-runner-cpp ./cpp

# 5. Go
docker build -t devcode-runner-go ./go

# 6. Java
docker build -t devcode-runner-java ./java

# 7. Rust
docker build -t devcode-runner-rust ./rust

# 8. C#
docker build -t devcode-runner-csharp ./csharp
```

---

### Step 4: Initialize & Run Frontend

1. Navigate to the frontend directory:
   ```powershell
   cd c:\Users\Lloyd Pili\projects\devcode\frontend
   ```

2. Create Next.js 15 application (if starting fresh):
   ```powershell
   pnpm create next-app . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
   ```

3. Install required frontend packages:
   ```powershell
   pnpm add @monaco-editor/react zustand axios lucide-react sonner clsx tailwind-merge class-variance-authority
   pnpm add -D @types/node @types/react @types/react-dom
   ```

4. Initialize shadcn/ui components:
   ```powershell
   pnpm dlx shadcn@latest init
   pnpm dlx shadcn@latest add button card dialog dropdown-menu input label badge progress separator tabs
   ```

5. Start the Next.js dev server:
   ```powershell
   pnpm run dev
   ```
   *Frontend is accessible at `http://localhost:3000`.*

---

## 5. Daily Development Workflow

To start development every day:

```powershell
# Terminal 1: Database & Redis
cd c:\Users\Lloyd Pili\projects\devcode
docker compose up -d

# Terminal 2: Go Backend API & WebSockets
cd c:\Users\Lloyd Pili\projects\devcode\backend
go run cmd/server/main.go

# Terminal 3: Next.js Frontend
cd c:\Users\Lloyd Pili\projects\devcode\frontend
pnpm run dev
```

---

## 6. Verification Checklist

1. **Database Connection**: Verify PostgreSQL is listening on `5432` and Redis on `6379`.
2. **Backend Health Check**: Open `http://localhost:8080/api/health` — should return `{"status":"ok"}`.
3. **Frontend Accessible**: Open `http://localhost:3000` — displays DevCode landing/dashboard.
4. **WebSocket Test**: Connect to `ws://localhost:8080/ws?token=<jwt>` using a test client.
5. **Code Execution Sandbox**: Test a basic Python submission against `devcode-runner-python` inside Docker.
