# Collaborative Task Manager

A real-time collaborative task management app built with Next.js, Go, PostgreSQL, and Redis.

## Tech Stack

- **Frontend**: Next.js 14 (TypeScript), Tailwind CSS, shadcn/ui
- **Backend**: Go, Chi router, WebSockets
- **Database**: PostgreSQL
- **Cache / Pub-Sub**: Redis
- **Auth**: JWT
- **Deployment**: Vercel (frontend), Fly.io (backend), Supabase (DB)

## Features

- Google OAuth + email/password authentication
- Kanban boards with drag-and-drop
- Real-time updates via WebSockets
- Role-based access (admin / member)
- REST API

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### Run locally

```bash
docker-compose up -d        # starts postgres + redis
cd backend && go run .      # starts Go API on :8080
cd frontend && npm run dev  # starts Next.js on :3000
```

## Project Structure

```
collaborative-tasks/
├── backend/     # Go REST API + WebSocket server
├── frontend/    # Next.js app
└── docker/      # Docker + docker-compose config
```
