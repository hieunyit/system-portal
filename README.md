# System Portal

This repository contains the Go backend API and a simple Next.js frontend to interact with it.

## Backend

The backend API is written in Go. To build and run it:

```bash
# build
make build
# or directly
# go build ./cmd/api
```

## Frontend

The `frontend` directory contains a minimal [Next.js](https://nextjs.org/) application.

### Install dependencies

```bash
cd frontend
npm install
```

### Development server

```bash
npm run dev
```

The app will be available at `http://localhost:3000` and expects the Go API to be running on the same host.
