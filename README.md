# GoTickitz 🎬

<br />
<br />

<div align="center">
  <img src="./assets/logo.svg" alt="Logo" width="200" />
</div>

<br />
<br />

[![Go Version](https://img.shields.io/badge/Go-v1.21+-blue.svg)](https://golang.org/doc/devel/release.html)
[![Gin](https://img.shields.io/badge/Gin-v1.9.0+-00ADD8.svg)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v15+-336791.svg)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-v7.0+-DC382D.svg)](https://redis.io/)

> Modern movie ticket booking platform API built with Go and Gin - Find movies, book seats, enjoy the show!

## Tech Stack

- **Go** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Primary database (via pgx driver)
- **Redis** - Caching
- **JWT** - Authentication

## Project Structure

```
gotickitz/
├── cmd/
│   └── main.go         # Application entry point
├── db/
│   └── migrations/     # Database migration files
├── internal/
│   ├── handlers/       # HTTP request handlers
│   ├── middlewares/    # Custom middlewares
│   ├── models/         # Data models
│   ├── repositories/   # Database interaction
│   ├── routes/         # API routes
│   └── utils/          # Helper functions
├── pkg/
│   ├── hash.go         # Password hashing
│   ├── jwt.go          # JWT authentication
│   ├── postgre.go      # PostgreSQL connection
│   └── redis.go        # Redis connection
├── public/             # Static files
│   └── img/            # Image storage
└── .env                # Environment variables
```

## API Endpoints

### Authentication

- `POST /api/v1/users/auth` - User login
- `POST /api/v1/users/auth/new` - User registration

### User Profile

- `GET /api/v1/users/profile` - Get user profile
- `PATCH /api/v1/users/profile` - Update user profile

### Movies

- `GET /api/v1/movies` - Get movies list (with filtering options)
- `GET /api/v1/movies/:id` - Get movie details

### Showings & Seats

- `GET /api/v1/showings` - Get movie showings
- `GET /api/v1/showings/:id/seat` - Get available seats for a showing

### Transactions

- `POST /api/v1/transactions` - Create a new transaction (book tickets)
- `GET /api/v1/transactions` - Get user's transaction history

### Admin Operations

- `POST /api/v1/admin/movies` - Add new movie
- `PATCH /api/v1/admin/movies/:id` - Update movie
- `DELETE /api/v1/admin/movies` - Delete movie

## Setup

1. Clone the repository
2. Create a [`.env`](.env ) file with the following variables:
   ```
   DBUSER=your_db_user
   DBPASS=your_db_password
   DBHOST=localhost
   DBPORT=5432
   DBNAME=tickitz
   RDSHOST=localhost
   RDSPORT=6379
   JWT_SECRET=your_jwt_secret
   JWT_ISSUER=your_issuer
   ```
3. Set up the PostgreSQL database and run the migrations
4. Start the Redis server
5. Run the application:
   ```
   go run cmd/main.go
   ```

## Features

- JWT-based authentication system
- Role-based access control (user/admin)
- Image upload for user profiles
- Movie filtering by name, genre, and status (upcoming, popular)
- Seat selection with availability checking
- Transaction history
- Points system for users
- Redis caching for frequently accessed data

## Development

To run in development mode with hot reload, install [Fresh](https://github.com/gravityblast/fresh) and run:

```
fresh
```

## API Documentation

For detailed API documentation, please refer to the API documentation (not included in this repository).
