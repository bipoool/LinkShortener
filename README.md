Got it!  
Here’s the updated **README.md** including the API usage and your deployment details 👇:

---

# LinkShortener

A simple and efficient URL shortener service built with **Go**, **PostgreSQL**, and **DiceDB** (in-memory cache).

## Features

- Shorten URLs easily.
- Persistent storage with **PostgreSQL**.
- High-speed caching with **DiceDB**.
- Ready-to-use **Docker** setup for easy deployment.

---

## Getting Started

### 1. Setup Configuration

- Copy the sample config file:

```bash
cp config.sample.json config.json
```

- Update `config.json` with your PostgreSQL and DiceDB connection details if needed.

---

### 2. Run the Application

Simply run:

```bash
go run main.go
```

This will start the URL shortener server on `localhost:8081`.

---

## API Usage

### Create a Short URL

- **Endpoint:** `POST http://localhost:8081/url`
- **Payload:**

```json
{
    "original-url": "https://example-url.com"
}
```

- **Response:**

```json
{
    "short-url": "http://localhost:8081/000003"
}
```

---

## Using Docker (Recommended)

A `Dockerfile` is included to containerize the application.

### To build and run:

```bash
docker build -t urlshortener .
docker run -p 8081:8081 urlshortener
```

---

## Using Docker Compose (Spin up everything)

If you want to run **PostgreSQL**, **DiceDB**, and **the URL shortener** together:

```bash
cd docker
docker-compose up -d
```

This will:

- Start a **PostgreSQL** container.
- Start a **DiceDB** container.
- Start the **URL Shortener** service.

All services will be ready and connected automatically.

---

## Deployment

This project is deployed on:

🔗 **[shtln.xyz](http://shtln.xyz)**

> **Note:** There is a **rate limit of 3 requests per minute** for creating a short URL.

---

## Project Structure

```
├── main.go
├── config.sample.json
├── config.json (after you create it)
├── Dockerfile
├── docker/
│   └── docker-compose.yml
├── internal/ (application code)
└── README.md
```

---

## Requirements

- Go 1.24+
- Docker (optional but recommended)
