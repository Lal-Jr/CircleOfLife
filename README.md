# Circle Of Life

Circle of Life is a full-stack social feed and location-based discovery platform. It enables users to create posts for "Help" or "Meetups", view them in a chronologically sorted feed, and discover them geolocated on an interactive map.

## Features

- **Authentication System**: Secure signup and login flow leveraging JWT and bcrypt.
- **Social Feed**: View, create, and interact with posts in a beautifully designed feed environment.
- **Interactive Map**: Discover nearby help requests and meetups using an integrated Mapbox-powered interface.
- **Post Management**: Create formatted posts categorized as "Help" or "Meetup" with rich details including geospatial coordinates (latitude/longitude), titles, descriptions, and meetup times.
- **Comment System**: Users can collaboratively engage with posts through an integrated real-time comment section.

## Tech Stack

### Frontend (Next.js Application)
Modern, responsive, and dynamic user interface built with:
- **Framework**: [Next.js](https://nextjs.org/) (React 19)
- **Styling & UI**: [Tailwind CSS](https://tailwindcss.com/) with Shadcn UI, Radix UI primitives, and Lucide React icons.
- **State Management & Data Fetching**: [TanStack/React Query](https://tanstack.com/query/latest) and Axios.
- **Forms & Validation**: `react-hook-form` organically combined with precise `zod` schema definitions.
- **Mapping**: `mapbox-gl` and `react-map-gl` for immersive geospatial capabilities.

### Backend (Golang HTTP Server)
Robust, high-performance Go API designed for seamless scalability:
- **Framework**: [Gin](https://gin-gonic.com/) web framework handling routing and core middleware.
- **Database**: Standardized on PostgreSQL (handling users, posts, and comments) queried utilizing `jackc/pgx`.
- **Caching & Real-time**: Redis natively structured for caching and event streaming.
- **Authentication**: JWT-based stateless authentication flow combined with Go's `crypto/bcrypt` algorithm.
- **Data Validation**: Strict inputs enforced using `go-playground/validator`.

## Project Structure

The repository maintains a clean separation of concerns:

- `/frontend`: Contains the complete Next.js application leveraging the `app` directory router.
  - `/src/app`: Application pages structured into `feed`, `map`, `post`, `create`, `login`, and `signup`.
  - `/src/components`: Reusable generic UI components and smart layout features.
- `/backend`: Contains the extensive Golang REST API structured around Domain-Driven Design (DDD) concepts.
  - `/internal`: Core business encapsulated logic, effectively separated into `services`, `handlers`, `models`, `repository`, and data `config`/`db` initialization.

## Getting Started

### Prerequisites
- Node.js (v20+)
- Go (v1.25.0+)
- PostgreSQL 
- Redis (Optional, dependent on local config)

### Working Locally

You will need to set up your `.env` files in both the frontend and backend directories:
- **Frontend** requires connection secrets, including a valid *Mapbox API* token.
- **Backend** requires a PostgreSQL connection string, a Redis connection layer, and a random JWT encryption secret.

#### 1. Running the Frontend
Navigate into the `/frontend` directory:
```bash
cd frontend
npm install
npm run dev
```
The client app should spin up and be accessible locally at `http://localhost:3000`.

#### 2. Running the Backend
In a new terminal, navigate into the `/backend` directory:
```bash
cd backend
go mod tidy
go run main.go
```
The Go HTTP server will listen natively at `http://localhost:8080`.
