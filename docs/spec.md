# 🗝️ Valdock — Valkey Instance & ACL Management Dashboard

## 📖 Overview

**Valdock** is a self‑hosted dashboard and API service for managing multiple **Valkey** servers (single, clustered, or containerized).
It automates instance creation via Docker, handles ACL user management, and centralizes configuration control — all through a simple web interface and REST API.

Built with:

- **Backend:** Go + `valkey-go` + Docker SDK
- **Frontend:** Svelte SPA
- **Storage:** Internal persistent Valkey instance (for metadata, ACL data, settings)

---

## ⚙️ System Architecture

```
┌────────────────────────────────┐
│          Svelte SPA            │
│  ─ REST/JSON → Backend (Go) ── │
└────────────────────────────────┘
               │
               ▼
┌────────────────────────────────┐
│         Go Backend API         │
│  - Instance management         │
│  - ACL handling (regen + load) │
│  - Config templates            │
│  - Internal Valkey store       │
│  - Docker container control    │
└────────────────────────────────┘
               │
               ▼
┌────────────────────────────────┐
│        Internal Valkey DB      │
│  Keys:                         │
│  - instances:*                 │
│  - aclusers:instance:id:*      │
│  - config:templates:*          │
│  - users:* (dashboard users)   │
└────────────────────────────────┘
               │
               ▼
┌────────────────────────────────┐
│        Docker Engine           │
│  - One container per instance  │
│  - Mounted ACL + config files  │
└────────────────────────────────┘
```

---

## 🗓️ Roadmap

### **Phase 1 — MVP**

- ✅ Instance lifecycle management (create/start/stop/delete)
- ✅ Regenerate & reload ACLs automatically
- ✅ Internal persistent Valkey database
- ✅ REST API
- ✅ Basic Svelte dashboard: list instances & ACL users

### **Phase 2 — Enhancements**

- Add live metrics (memory, keys, uptime)
- Editable instance configuration (Valkey conf)
- Config template management
- User authentication for dashboard (JWT tokens)
- Embedded single binary deployment (`valdock`)

### **Phase 3 — Advanced Management**

- Centralized role management for dashboard users
- Scheduled instance backups & restores
- Instance logs & runtime monitoring
- Cluster management (multi-node Valkey clusters)
- Theme & settings UI

### **Phase 4 — Production Features**

- TLS support
- Integration with Vault (secret rotation)
- Remote Docker host support
- Automatic scaling & templates

---

## 🌐 REST API Design (v1 Draft)

**Base URL:**

```
/api/v1/
```

---

### 🧱 1. Instance Management

| Method   | Path                   | Description                        |
| -------- | ---------------------- | ---------------------------------- |
| `GET`    | `/instances`           | List all managed Valkey instances  |
| `POST`   | `/instances`           | Create new instance                |
| `GET`    | `/instances/:id`       | Get details for one instance       |
| `DELETE` | `/instances/:id`       | Delete instance (and container)    |
| `POST`   | `/instances/:id/start` | Start an existing container        |
| `POST`   | `/instances/:id/stop`  | Stop a running container           |
| `GET`    | `/instances/:id/stats` | Get metrics (memory, uptime, etc.) |

**▶ Request**

```json
POST /api/v1/instances
{
  "name": "valkey_app1",
  "port": 6380,
  "configTemplate": "default",
  "users": [
    { "username": "appuser", "password": "secret" }
  ]
}
```

**◀ Response**

```json
{
  "id": "valkey_app1",
  "status": "running",
  "port": 6380,
  "created_at": "2025-10-20T22:00:00Z"
}
```

---

### 🔐 2. ACL Management

| Method   | Path                            | Description                  |
| -------- | ------------------------------- | ---------------------------- |
| `GET`    | `/instances/:id/acls`           | List users for that instance |
| `POST`   | `/instances/:id/acls`           | Add new ACL user             |
| `PUT`    | `/instances/:id/acls/:username` | Update existing ACL user     |
| `DELETE` | `/instances/:id/acls/:username` | Delete user                  |

**▶ Example — Add user**

```json
POST /api/v1/instances/valkey_app1/acls
{
  "username": "readonly",
  "password": "abc123",
  "permissions": "+get",
  "keys": "~public:*"
}
```

**◀ Response**

```json
{
  "status": "ok",
  "message": "User 'readonly' added and ACL reloaded."
}
```

**▶ Example — Update password**

```json
PUT /api/v1/instances/valkey_app1/acls/readonly
{
  "password": "newpass123"
}
```

**◀ Response**

```json
{ "status": "ok" }
```

---

### ⚙️ 3. Configuration

| Method   | Path                      | Description                      |
| -------- | ------------------------- | -------------------------------- |
| `GET`    | `/instances/:id/config`   | Read current config              |
| `PUT`    | `/instances/:id/config`   | Update certain config parameters |
| `GET`    | `/config/templates`       | List available templates         |
| `POST`   | `/config/templates`       | Add new template                 |
| `GET`    | `/config/templates/:name` | Get template details             |
| `DELETE` | `/config/templates/:name` | Delete template                  |

**▶ Request**

```json
PUT /api/v1/instances/valkey_app1/config
{
  "maxmemory": "512mb",
  "appendonly": "yes"
}
```

**◀ Response**

```json
{
  "status": "ok",
  "message": "Config updated and instance reloaded."
}
```

---

### 🧠 4. System / Health / Settings

| Method | Path               | Description                         |
| ------ | ------------------ | ----------------------------------- |
| `GET`  | `/system/health`   | Basic health check (`200 OK`)       |
| `GET`  | `/system/info`     | Info about Valdock, Docker, Valkey  |
| `GET`  | `/system/metrics`  | Host resource metrics               |
| `GET`  | `/system/settings` | Dashboard settings                  |
| `PUT`  | `/system/settings` | Update dashboard/global preferences |

**◀ Example Response (`/system/info`)**

```json
{
  "valdock_version": "0.1.0",
  "valkey_version": "8.0.1",
  "docker_version": "27.1.0",
  "instances_count": 3
}
```

---

## 🧰 Example Data Model (Internal Valkey)

| Key Pattern                  | Description                    |
| ---------------------------- | ------------------------------ |
| `instance:<id>`              | JSON metadata for instance     |
| `aclusers:<instance>:<user>` | JSON for ACL user              |
| `config:template:<name>`     | Valkey config template content |
| `user:<dashboardUser>`       | Dashboard UI login (future)    |

---

## 🧱 Deployment Layout

```
valdock/
 ├── backend/
 │     ├── main.go
 │     ├── api/
 │     ├── storage/
 │     ├── internal/
 │     │      ├── docker/
 │     │      ├── acl/
 │     │      └── config/
 │     └── webui/ (built Svelte app)
 ├── frontend/
 │     └── src/
 │         ├── routes/
 │         ├── components/
 │         └── lib/
 ├── docker-compose.yml
 ├── valkey-data/
 └── README.md
```

---

## 🔒 Security Notes

- Internal Valkey user/password stored locally (protected).
- ACL file updates trigger `ACL LOAD` for immediate reapplication.
- Dashboard admin endpoints protected with API tokens (added in Phase 2).
- Sensitive data (passwords) always hashed before storing.

---

## 🧭 Future Extensions

- Cluster orchestration (`/clusters` API)
- Scheduled backups/export configs
- Metrics/Prometheus endpoint
- Role templates for ACLs
- Plugin interface for external secret stores (Vault, SSM, etc.)
- CLI client (`valdockctl`)

---

## 📜 License

Open-source (MIT or Apache‑2 suggested)
