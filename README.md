# System Portal

This repository contains the Go backend API and a Next.js based web frontend to interact with it.

## Backend

The backend API is written in Go. To build and run it:

```bash
# build
make build
# or directly
# go build ./cmd/api
```

## Frontend

The `frontend` directory contains a [Next.js](https://nextjs.org/) 14 application built with Tailwind CSS.

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
You can configure the backend URL via `NEXT_PUBLIC_API_BASE_URL` in a `.env.local` file.
