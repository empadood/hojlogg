# hojlogg 🏍️

Motorcycle dashboard logging app.

Scan your dashboard with your phone → the backend reads the odometer & fuel level via OCR → data is saved to PostgreSQL.

## Stack

| Layer    | Technology                          |
|----------|-------------------------------------|
| Mobile   | React Native (Expo) + TypeScript    |
| Backend  | Go (Gin) REST API                   |
| Database | PostgreSQL 16                       |

## Project layout

```
hojlogg/
├── backend/          Go REST API
│   ├── main.go
│   ├── Dockerfile
│   ├── db/           DB connection + SQL migrations
│   ├── handlers/     HTTP handlers (logs CRUD + image upload)
│   ├── middleware/   Request logging
│   └── models/       Shared data models
├── mobile/           Expo React Native app
│   ├── App.tsx       Navigation root
│   └── src/
│       ├── api/      Axios client (logsApi)
│       ├── components/
│       └── screens/  Home · NewLog · Camera · LogDetail
└── docker-compose.yml
```

---

## Prerequisites

| Tool | Minimum version | Install |
|------|-----------------|---------|
| Docker & Docker Compose | 24+ | https://docs.docker.com/get-docker/ |
| Go | 1.23 | https://go.dev/dl/ |
| Node.js | 20 | https://nodejs.org/ |
| Expo Go (phone) | latest | App Store / Google Play |

---

## Option A — Run everything with Docker Compose (recommended)

This is the easiest way to get the **database** and **backend** running together.

```bash
# 1. Clone the repo (if you haven't already)
git clone https://github.com/empadood/hojlogg.git
cd hojlogg

# 2. Build and start PostgreSQL + the Go backend
docker compose up --build
```

Docker Compose will:
- Start a PostgreSQL 16 container and wait until it is healthy.
- Build the Go backend image and start it.
- Apply the SQL migrations automatically on first run.

The API is now available at **http://localhost:8080**.

```bash
# Verify it is running
curl http://localhost:8080/api/health
# → {"status":"ok"}
```

To stop:
```bash
docker compose down          # stop containers, keep data
docker compose down -v       # stop containers AND delete all data
```

---

## Option B — Run each service manually

Use this approach if you prefer to iterate on the backend without rebuilding a Docker image.

### 1. Start the database

```bash
# Start only the PostgreSQL container
docker compose up -d db

# Wait a few seconds, then confirm it is ready
docker compose exec db pg_isready -U hojlogg
# → /var/run/postgresql:5432 - accepting connections
```

The database is reachable at `localhost:5432` with:

| Setting  | Value      |
|----------|------------|
| Host     | localhost  |
| Port     | 5432       |
| Database | hojlogg    |
| User     | hojlogg    |
| Password | hojlogg    |

### 2. Start the backend

```bash
cd backend

# (Optional) copy and edit the sample env file
cp .env.example .env   # if present, or create one manually — see below

# Install Go dependencies
go mod download

# Run the server
DATABASE_URL="postgres://hojlogg:hojlogg@localhost:5432/hojlogg" go run .
```

The server starts on **http://localhost:8080** and applies SQL migrations automatically.

You can also set variables in a `.env` file in the `backend/` directory:

```dotenv
DATABASE_URL=postgres://hojlogg:hojlogg@localhost:5432/hojlogg
ADDR=:8080
UPLOAD_DIR=/tmp/hojlogg-uploads
GIN_MODE=debug
```

Then simply run:
```bash
go run .
```

### 3. Run backend tests

```bash
cd backend
go test ./...
```

---

## Start the mobile app

The mobile app connects to the backend. You need to tell it where the backend is running.

### On a physical device (recommended)

Your phone and development machine must be on the **same Wi-Fi network**.

```bash
cd mobile

# 1. Install dependencies
npm install

# 2. Set the backend URL to your machine's local IP address
#    (find it with `ipconfig` on Windows or `ifconfig` / `ip addr` on Linux/macOS)
export EXPO_PUBLIC_API_URL=http://192.168.1.42:8080

# 3. Start the Expo development server
npx expo start
```

Scan the QR code that appears in the terminal with the **Expo Go** app on your phone.

### On an Android emulator

```bash
cd mobile
npm install
# The emulator can reach the host machine via 10.0.2.2
EXPO_PUBLIC_API_URL=http://10.0.2.2:8080 npx expo start --android
```

### On an iOS simulator (macOS only)

```bash
cd mobile
npm install
# The simulator shares the host network, so localhost works
EXPO_PUBLIC_API_URL=http://localhost:8080 npx expo start --ios
```

---

## API reference

| Method | Path                    | Description                                           |
|--------|-------------------------|-------------------------------------------------------|
| GET    | /api/health             | Health check                                          |
| GET    | /api/logs               | List logs (`?limit=20&offset=0`)                      |
| POST   | /api/logs               | Create log (`{odometer_km, fuel_level?, notes?}`)     |
| GET    | /api/logs/:id           | Get single log                                        |
| DELETE | /api/logs/:id           | Delete log                                            |
| POST   | /api/logs/:id/image     | Upload dashboard image (runs OCR, updates log values) |

### Example: create a log entry

```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{"odometer_km": 12345, "fuel_level": 75, "notes": "Morning ride"}'
```

### Example: upload a dashboard photo

```bash
curl -X POST http://localhost:8080/api/logs/<LOG_ID>/image \
  -F "image=@/path/to/dashboard.jpg"
```

---

## Environment variables (backend)

| Variable       | Default                    | Description                    |
|----------------|----------------------------|--------------------------------|
| `DATABASE_URL` | —                          | PostgreSQL DSN (**required**)  |
| `UPLOAD_DIR`   | `/tmp/hojlogg-uploads`     | Directory for uploaded images  |
| `ADDR`         | `:8080`                    | Listen address                 |
| `GIN_MODE`     | `release`                  | `debug` or `release`           |

