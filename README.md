# Chirpy

A simple Go REST API server for experimenting with Go web development.

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- [Goose](https://github.com/pressly/goose) for database migrations
- [sqlc](https://sqlc.dev/) for generating type-safe database code 

### Setup

1. Clone the repository
2. Setup a .env file with the required environment variables (see below)
2. Start the PostgreSQL database:
   ```bash
   docker-compose up -d
   ```

3. Run database migrations:
   ```bash
   goose up
   ```

4. Start the server:
   ```bash
   go run .
   ```

## Environment Variables

Copy the `.env` file and configure:
- `DB_URL` - PostgreSQL connection string
- `SECRET` - JWT secret key
- `POLKA_KEY` - API key for external services

## API Endpoints

This is a learning project focused on building REST API fundamentals with Go.