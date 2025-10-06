## 📡 API Endpoints

### Health Check
| Method | Endpoint | Description |
|---------|-----------|-------------|
| `GET` | `/api/healthz` | Check if the server is up and running |

---

### Authentication & User Management
| Method | Endpoint | Description |
|---------|-----------|-------------|
| `POST` | `/api/users` | Create a new user |
| `PUT` | `/api/users` | Update existing user info |
| `POST` | `/api/login` | Log in a user and receive tokens |
| `POST` | `/api/refresh` | Refresh authentication token |
| `POST` | `/api/revoke` | Revoke user authentication token |

---

### Chirps (Posts)
| Method | Endpoint | Description |
|---------|-----------|-------------|
| `POST` | `/api/chirps` | Create a new chirp |
| `GET` | `/api/chirps/` | Retrieve all chirps |
| `GET` | `/api/chirps/{chirpID}` | Retrieve a specific chirp by ID |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp by ID |

---

### Webhooks
| Method | Endpoint | Description |
|---------|-----------|-------------|
| `POST` | `/api/polka/webhooks` | Handle incoming webhook events (example integration) |

---

### Admin Endpoints
| Method | Endpoint | Description |
|---------|-----------|-------------|
| `GET` | `/admin/metrics` | Retrieve server metrics |
| `POST` | `/admin/reset` | Reset server data (use with caution) |