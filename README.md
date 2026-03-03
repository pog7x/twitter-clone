<div align="center">
  <img src="./backend/uploads/avatar_user_1.svg" width="80" alt="Twitter Clone Logo" />
  <h1>Twitter Clone</h1>
  <p>Full-stack social media application inspired by X (Twitter), built with Nuxt 3 and Go.</p>

  <p>
    <img src="https://img.shields.io/badge/Nuxt-3.x-00DC82?style=flat-square&logo=nuxt.js&logoColor=white" alt="Nuxt 3" />
    <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.24" />
    <img src="https://img.shields.io/badge/PostgreSQL-12-336791?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker" />
    <img src="https://img.shields.io/badge/TailwindCSS-via_@nuxt/ui-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white" alt="TailwindCSS" />
  </p>
</div>

---

## About

Twitter Clone replicates the core experience of X (Twitter): three-column layout, tweet composition with up to 4 image attachments, user profiles, likes, trending topics, and dark/light mode. The frontend is a Nuxt 3 SPA, the backend is a RESTful Go API, and everything runs in Docker containers with PostgreSQL.

**Key features:** authentication (JWT + bcrypt), tweet feed, image lightbox, inline tweet editing & deletion, user profiles with editable bio, trending sidebar, responsive mobile layout.

---

## Tech Stack

| Layer | Technologies |
|---|---|
| **Frontend** | Nuxt 3 (SPA), Vue 3 + Composition API, Pinia, @nuxt/ui, TailwindCSS, TypeScript |
| **Backend** | Go 1.24, Iris, GORM, JWT (HS256), Viper, Logrus, bcrypt |
| **Infrastructure** | PostgreSQL 12, Docker + Docker Compose, nginx |

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)

### Development

Hot reload enabled for both frontend (Nuxt dev server) and backend (`reflex`).

```bash
git clone https://github.com/your-username/twitter-clone.git
cd twitter-clone
docker compose -f docker-compose.dev.yml up --build
```

### Production

```bash
docker compose up --build -d
```

After startup:
- Frontend: [http://localhost:3000](http://localhost:3000)
- Backend API: [http://localhost:8080](http://localhost:8080)

### Seed credentials

| Username | Password |
|---|---|
| `pog7x` | `pass123` |
| `pog8x` | `pass456` |
