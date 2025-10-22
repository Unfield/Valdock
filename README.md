# 🗝️ Valdock

> **Spin. Secure. Scale. Your Valkeys — docked in one place.**
>
> **Valdock** is a lightweight, self‑hosted dashboard and API to **manage multiple Valkey (Redis‑compatible) servers** via Docker.
> It handles instance creation, ACL management, and configuration — all through a modern Svelte web interface and Go backend.

---

## 🚀 Features

- 🧱 **Multi‑instance orchestration**
  Create, start, stop, and remove Valkey containers easily — all from a single web UI.

- 🔐 **Built‑in ACL management**
  Add, remove, or update users and passwords for each instance. ACL files are generated automatically and applied live.

- ⚙️ **Dynamic configuration**
  Adjust Valkey settings (like `maxmemory` or persistence options) from the dashboard. Changes are immediately applied.

- 🧠 **Internal Valkey database**
  Valdock uses its own embedded Valkey for metadata and user storage — fast, simple, and persistent.

- 🖥️ **Svelte single‑page app frontend**
  A modern, reactive interface that feels native and dev‑friendly.

- 🧩 **Docker‑first design**
  Everything runs in containers for clean isolation and easy updates.

---

## 🧭 Architecture Overview

```
                    ┌──────────────────────────────┐
                    │       Valdock Dashboard      │
                    │   (Go REST API + Svelte UI)  │
                    └─────────────┬────────────────┘
                                  │
                       Internal Valkey DB (metadata)
                                  │
                ┌─────────────────┴──────────────────┐
                │                                    │
         [Valkey Instance 1]                 [Valkey Instance 2]
           - Unique port                        - Unique ACL
           - Own config                         - Own persistence
```

Valdock stores all instance metadata, ACL user info, and templates in its internal Valkey database, then **renders configuration and ACL files** for each instance on disk, applying them via the Docker Engine API.

---

## ⚙️ Tech Stack

| Layer        | Technology                                    |
| ------------ | --------------------------------------------- |
| **Backend**  | Go + valkey‑go SDK                            |
| **Frontend** | Svelte + Vite SPA                             |
| **Storage**  | Internal Valkey (persistent RDB/AOF)          |
| **Runtime**  | Docker (each Valkey instance = one container) |

---

## 📦 Installation (Dev Mode)

1. **Clone the repo**

   ```bash
   git clone https://github.com/yourname/valdock.git
   cd valdock
   ```

2. **Start the stack**

   ```bash
   docker compose up -d
   ```

   This runs:
   - `valdock-api` — Go backend with REST API
   - `valdock-ui` — Svelte frontend
   - `valdock-db` — internal Valkey instance

3. **Open your dashboard**

   ```
   http://localhost:8080
   ```

---

## 🔌 API Overview

All API endpoints are under `/api/v1/`.

| Category  | Example Endpoints                                            |
| --------- | ------------------------------------------------------------ |
| Instances | `GET /instances`, `POST /instances`, `DELETE /instances/:id` |
| ACLs      | `GET /instances/:id/acls`, `POST /instances/:id/acls`        |
| Config    | `GET /instances/:id/config`, `PUT /instances/:id/config`     |
| System    | `GET /system/health`, `GET /system/info`                     |

Valdock follows simple JSON request/response patterns for full Svelte integration.

---

## 🧱 Directory Structure

```
valdock/
 ├── backend/
 │    ├── api/              ← REST endpoints (instances, acl, system)
 │    ├── storage/          ← internal Valkey handling
 │    ├── docker/           ← container lifecycle logic
 │    ├── internal/acl/     ← ACL generator & sync
 │    └── main.go
 ├── frontend/              ← Svelte SPA (served by Go in production)
 ├── docker-compose.yml
 └── README.md
```

---

## 🔮 Roadmap

| Phase     | Focus                                               |
| --------- | --------------------------------------------------- |
| ✅ MVP    | Multiple Valkey instances, ACL management, core API |
| 🟡 v0.2   | Live metrics, config templates, persistent settings |
| 🟢 v0.3   | Dashboard user roles & authentication               |
| 🟣 v0.4   | Cluster support: create and manage Valkey clusters  |
| 🏁 Future | Backups, monitoring, TLS management, auto‑scaling   |

---

## 💡 Name Story

> **Valdock** = _Valkey + Dock(er)_
>
> The tool that docks, guards, and organizes all your Valkeys — securely and simply.

---

## 🧑💻 License

**Elastic License 2.0 (ELv2)**
See the `LICENSE` file for full details.

By using or reproducing the software, you agree to the terms of the Elastic License 2.0.
Commercial use is permitted under its conditions — attribution required.

---

## ❤️ Contributing

1. Fork the project
2. Create a feature branch
3. Open a pull request

Style, readability, and simplicity are the main project pillars.
