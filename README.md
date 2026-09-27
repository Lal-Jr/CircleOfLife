# 🌍 Circle of Life

Circle of Life is a full-stack, location-based social feed and discovery platform. Built with a focus on geographic proximity, it enables users to seamlessly broadcast "Help" requests or organize "Meetups" with people immediately around them. 

The platform is engineered for modern performance, utilizing **PostGIS** for geospatial querying, **Redis** for robust caching and rate-limiting, and **Server-Sent Events (SSE)** for real-time feed updates.

## 🔗 Live Demo

**[circleoflife-web.fly.dev](https://circleoflife-web.fly.dev)**. Click **"Try the demo account"** on the login page and you're in, with no sign-up.

| | |
|---|---|
| Email | `demo@circleoflife.app` |
| Password | `CircleDemo123!` |

The demo account always sees the sample neighbourhood in central Bangalore, wherever you are, and there's no location prompt. Other accounts use your real location.

![Demo: one-click login, the nearby feed, a comment thread, the map and creating a post](docs/screenshots/demo.gif)

### Screenshots

| Nearby feed | Activity map |
|---|---|
| ![Nearby feed](docs/screenshots/feed.png) | ![Activity map](docs/screenshots/map.png) |
| **Post and comment thread** | **Create a post** |
| ![Post detail](docs/screenshots/post-detail.png) | ![Create post](docs/screenshots/create-post.png) |
| **Login with one-click demo** | **Profile** |
| ![Login](docs/screenshots/login.png) | ![Profile](docs/screenshots/profile.png) |

<p align="center"><img src="docs/screenshots/mobile-feed.png" alt="Mobile feed" width="300" /><br/><em>Mobile layout with bottom navigation</em></p>

### Sample data

On first startup the backend seeds a small neighbourhood so the app never looks empty ([`backend/internal/db/seed.go`](backend/internal/db/seed.go)). Seeding is idempotent and runs in a single transaction:

- **The demo account** with two posts of its own: a help request and a meetup.
- **Five neighbours** (Priya, Arjun, Meera, Rahul and Ananya) with 10 posts spread across 0.5 to 2.5 km. The help requests cover a flat car battery, a medicine pickup, a lost dog and tech help for a grandparent. The meetups include a 5K run, a book swap, a garden day, a cricket match and a coding group, each with an upcoming time.
- **Comment threads and Helpful votes** between the neighbours and the demo user, so post pages show real conversations.

Neighbour accounts can't be logged into; only the demo account has a usable password.

## 🚀 Key Features

### 📍 Location-Based Feed
- **Geospatial Discovery:** Leverages PostgreSQL with PostGIS to accurately query, rank, and deliver posts within a user-defined radius (e.g., 5km).
- **Ranking Algorithm:** Feed priority is dynamically calculated based on a decay formula weighing distance, time since creation, and post type (e.g., Help vs Meetup). 
- **Nearby Highlight:** Hyper-local, urgent activity (under 500m) is automatically pinned to the top of the feed as an "Urgent Nearby Activity" card to encourage immediate community assistance.

### ⚡ Real-Time & Optimistic UI
- **Live Updates:** Integrated Server-Sent Events (SSE) backed by Redis Pub/Sub push real-time notifications to users when new posts are created nearby, ensuring the feed always feels alive.
- **Optimistic Rendering:** Leveraging React Query, creating posts or comments updates the UI instantly before the server confirms the request, providing a snappy, app-like experience.

### 🗺️ Interactive Map Visualization
- **Mapbox Integration:** A fully interactive map view built with `react-map-gl`, plotting nearby posts as distinct markers so users can visually explore their surroundings.

### 🛡️ Production-Ready Security & Polish
- **Authentication:** Stateless, secure JWT-based authentication flow with bcrypt password hashing.
- **Rate Limiting:** Granular, token-bucket user-level rate limiting (backed by Redis) prevents spam on post creation (5/min) and commenting (10/min).
- **Input Sanitization:** Strict HTML and XSS sanitization applied server-side using `bluemonday` to prevent malicious payloads, coupled with tuned character length constraints.

## 🛠️ Tech Stack

### Frontend (Next.js Application)
A responsive, mobile-first web app built for modern UX.
- **Framework:** Next.js (React 19) utilizing the App Router.
- **Styling:** Tailwind CSS integrated with `shadcn/ui`, Radix UI primitives, and Lucide Icons.
- **State & Data:** TanStack Query (React Query) for caching, optimistic updates, and infinite scrolling feeds.
- **Forms:** `react-hook-form` paired with `zod` for rigorous client-side schema validation.
- **Maps:** `mapbox-gl` for geospatial rendering.

### Backend (Golang HTTP API)
A high-performance, Domain-Driven Design (DDD) Go backend.
- **Framework:** Gin Web Framework for lightning-fast HTTP routing and middleware chaining.
- **Database:** PostgreSQL utilizing `jackc/pgx` for connection pooling, featuring PostGIS extensions for spatial data.
- **Caching & PubSub:** Redis for reducing database load (feed caching), rate limiting, and cross-node event broadcasting.
- **Security:** `golang-jwt/jwt`, `crypto/bcrypt`, `bluemonday` for XSS protection, and `go-playground/validator`.

## 📂 Project Structure

The repository maintains a clean separation of concerns:

- `/frontend`
  - `/src/app`: Application routes (`/feed`, `/map`, `/post`, `/create`, `/login`, `/signup`).
  - `/src/components`: UI components, including the interactive map, feed cards, and layout elements.
  - `/src/hooks`: Custom React hooks mapping complex logic (e.g., `useFeed`, `useRealtimeFeed`, `useLocation`).
- `/backend`
  - `/internal`: Core business logic, separated into `services`, `handlers`, `middleware`, `models`, and `repository` layers.
  - `/cmd`: The entry point for the application (`cmd/server/main.go`).

## 🏁 Getting Started

### Prerequisites
- Node.js (v20+)
- Go (v1.25.0+)
- PostgreSQL (with PostGIS extension installed)
- Redis (Optional for local dev, required for caching/rate-limiting/real-time events)

### Environment Configuration

You will need to set up `.env` files in both directories:

**Frontend (`/frontend/.env.local`)**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_MAPBOX_TOKEN=your_mapbox_public_token
```

**Backend (`/backend/.env`)**
```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/circleoflife?sslmode=disable
REDIS_URL=redis://localhost:6379/0
JWT_SECRET=your_super_secret_jwt_key
```

### 1. Running the Database
Ensure PostgreSQL is running and the database `circleoflife` is created with the `postgis` extension enabled:
```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

### 2. Running the Backend
Navigate to the `/backend` directory and start the Go server:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
The API will be available at `http://localhost:8080`.

### 3. Running the Frontend
In a separate terminal, navigate to the `/frontend` directory:
```bash
cd frontend
npm install
npm run dev
```
The client app will be accessible at `http://localhost:3000`.

On first run the backend creates the tables, enables PostGIS and seeds the demo account and sample neighbourhood. Sign in with `demo@circleoflife.app` / `CircleDemo123!`.

Alternatively, `docker compose up --build` starts PostGIS, Redis, the API and the frontend (on `http://localhost:3100`).

## ☁️ Deployment

The live demo runs on **Fly.io** (region `sin`) with **Supabase** Postgres, which has PostGIS built in:

```bash
# API
cd backend
fly secrets set -a circleoflife-api DATABASE_URL='postgresql://…:5432/postgres?sslmode=require' JWT_SECRET='…'
fly deploy

# Web: the API URL and Mapbox token are baked in at build time
cd ../frontend
fly deploy --build-arg NEXT_PUBLIC_MAPBOX_TOKEN=pk.…
```

`REDIS_URL` is optional. Without it, caching and rate limiting are skipped and live SSE updates are off, but everything else works.
