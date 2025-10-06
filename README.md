# 🐦 Chirpy — A Mini Twitter Clone Built in Go

**Chirpy** is a lightweight, local-first Twitter-style API server written in Go.  
It allows users to **create accounts, log in, post "chirps" (short messages of up to 140 characters)**, and manage their data—all backed by a **PostgreSQL** database.

If you’re learning how to build a production-style backend with Go or want a simple microservice to test authentication, database interactions, and RESTful API design, Chirpy is a perfect starting point.

---

## 🚀 Features

- User registration and authentication (login, refresh, revoke tokens)
- Create, retrieve, and delete **chirps**
- PostgreSQL database integration for persistent storage
- Admin endpoints for metrics and server reset
- Simple, RESTful API design
- Webhook integration example (`/api/polka/webhooks`)

---

## 🧩 Tech Stack

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Architecture:** RESTful API
- **Routing:** Go `net/http` with `ServeMux`

---

## ⚙️ Setup & Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

### 2. Before using
To use it you will need to install Postgres to manage the program's
database. As well as the programming language Go.

### 3. To run the server
```bash
go run .
```

### 4. Test it
```bash
curl -X POST http://localhost:8080/api/chirps \
-d '{"password": "1234", "email": "user@example.com"}'
```
Check out the API endpoints [here](./docs/API_DOCS.md).