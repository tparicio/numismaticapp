# NumismaticApp

A Dockerized web application for managing numismatic coin collections with AI-powered analysis.

## Description
NumismaticApp allows users to upload photos of coins, automatically crops and rotates them, analyzes them using Google's Gemini AI to extract metadata (Country, Year, Value, etc.), and stores them in a PostgreSQL database. The frontend provides a dashboard and a gallery to view the collection.

## Prerequisites
- **Docker** & **Docker Compose**
- **Make**
- **Google Gemini API Key**
- **Go 1.21+** (for local development)
- **Node.js 20+** (for local frontend development)

## Installation & Execution

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd numismaticapp
   ```

2. **Environment Setup:**
   Create a `.env` file based on the example (or set variables in `docker-compose.yml`):
   ```bash
   GEMINI_API_KEY=your_api_key_here
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_DB=numismatic
   ```

3. **Run with Docker Compose:**
   ```bash
   make run
   ```
   This will start the Backend (API), Frontend, and Database.

4. **Access the App:**
   - Frontend: `http://localhost:5173` (or configured port)
   - API: `http://localhost:8080`

## Tech Stack
- **Backend:** Go (Fiber), SQLC, pgx
- **Database:** PostgreSQL
- **AI:** Google Gemini
- **Image Processing:** libvips
- **Frontend:** Vue 3, Vite, TailwindCSS, DaisyUI

## Development Commands
- `make run`: Start everything with Docker.
- `make lint`: Run linters (Go & Vue).
- `make test`: Run unit tests.
- `sqlc generate`: Regenerate DB code after schema changes.
