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

## Quick start (backend + DB)

```bash
# Start PostgreSQL + Go backend
docker compose up --build

# The API is now available at http://localhost:8080
```

## Mobile app

```bash
cd mobile
npm install
# Edit src/api/client.ts and set EXPO_PUBLIC_API_URL to your machine's LAN IP
npx expo start
```

Scan the QR code with the Expo Go app on your device.

## API reference

| Method | Path                    | Description                            |
|--------|-------------------------|----------------------------------------|
| GET    | /api/health             | Health check                           |
| GET    | /api/logs               | List logs (`?limit=20&offset=0`)       |
| POST   | /api/logs               | Create log (`{odometer_km, fuel_level?, notes?}`) |
| GET    | /api/logs/:id           | Get single log                         |
| DELETE | /api/logs/:id           | Delete log                             |
| POST   | /api/logs/:id/image     | Upload dashboard image (OCR analysis)  |

## Environment variables (backend)

| Variable       | Default                                      | Description             |
|----------------|----------------------------------------------|-------------------------|
| `DATABASE_URL` | —                                            | PostgreSQL DSN (required)|
| `UPLOAD_DIR`   | `/tmp/hojlogg-uploads`                       | Image storage directory |
| `ADDR`         | `:8080`                                      | Listen address          |
| `GIN_MODE`     | `release`                                    | `debug` or `release`    |
