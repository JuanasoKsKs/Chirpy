# CHIRPY: Practice HTTP Server fundamentals
REST API that Handles authentication, encrypt passwords, POST/GET/UPDATE/DELETE content from a PostgreSQL database

## Requirements:
- Go 1.22+
- PostgreSQL

## INSTALLATION:
1. clone the repo: `git clone <url>`
2. Install dependencies: `go mod tidy`

## CONFIG:
1. Run the Server in the root of the project: `go run main.go` or `go run .`
2. The server run in the port `8080` so you should use `localhost:8080/` in Browser or http Client.

## USAGE / ENDPOINTS:
- `GET /api/healthz` : Returns a 200 OK status to verify the server is healthy.
- `POST /api/users`: Allows to create a user profile, and store the connection data in the database .
- `POST /api/login`: Allows users to communicate with the server after authentication.
- `POST /api/chirps`: Allows authenticated users to create a new chirp.
- `GET /api/chirps/{UserID}`: Allows authenticated users to retrieve owned chirps.
Content must be as JSON in the `http.Request.Body`
Others would be added shortly.


## Environment Variables
Create a `.env` file with:
- `DB_URL`: PostgreSQL connection string in the shape of `protocol://username:password@host:port/database`
- `SECRET`: Secret key for authentication.