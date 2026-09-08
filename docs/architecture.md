# Architecture Overview

## System Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    Client (Browser)                       │
│               Nuxt 3 SSR / Vue 3 SPA                     │
│                Pinia State Management                     │
│              WebSocket (Real-time Quiz)                   │
└────────────────────┬─────────────────┬───────────────────┘
                     │ HTTP            │ WebSocket
                     ▼                 ▼
┌──────────────────────────────────────────────────────────┐
│               Go Backend (GoFiber)                        │
│              /api/v1/* endpoints                          │
│                                                           │
│  Routes → Middleware → Controllers → Services → Models    │
│                                                           │
│  goqu (SQL query builder) + database/sql                  │
└───────┬──────────┬──────────────┬──────────────┬─────────┘
        │          │              │              │
        ▼          ▼              ▼              ▼
┌───────────┐ ┌─────────┐ ┌───────────┐ ┌─────────────┐
│PostgreSQL │ │ Valkey/ │ │Ory Kratos │ │SMTP (Mailpit│
│    15     │ │  Redis  │ │   (Auth)  │ │  in dev)    │
└───────────┘ └─────────┘ └───────────┘ └─────────────┘
```

## Component Breakdown

### Frontend (`app/`)

**Framework:** Nuxt 3 (SSR + SPA hybrid)

| Layer | Location | Purpose |
|-------|----------|---------|
| Pages | `pages/` | File-based routing |
| Components | `components/` | Reusable Vue components (auto-imported) |
| Composables | `composables/` | Reusable logic (auto-imported) |
| Store | `store/` | Pinia state management |
| Plugins | `plugins/` | Nuxt plugins (chart.js) |

**Key Patterns:**
- Auto-imports: components and composables need no manual imports
- File-based routing: `pages/admin/index.vue` → `/admin`
- Dynamic routes: `pages/join/[code].vue` → `/join/:code`
- SSR for SEO; client-side WebSocket for quiz playing

### Backend (`api/`)

**Framework:** GoFiber (fast HTTP framework)

```
api/
├── app.go              # Entry point, CLI commands
├── cli/                # CLI command definitions
├── config/             # Env parsing (envconfig)
├── constants/          # Application constants
├── controllers/        # HTTP handlers
│   └── api/v1/         # Versioned controllers
├── database/           # Connection + migrations
│   └── migrations/     # SQL migration files
├── helpers/            # Utility functions
├── middlewares/        # HTTP middleware
├── models/            # Data models (goqu queries)
├── routes/            # Route definitions
├── services/          # Business logic
└── utils/             # Helpers (CSV, scoring, etc.)
```

**Request Flow:**
1. Route matched in `routes/main.go`
2. Middleware chain executes (auth, logging, permissions)
3. Controller handles request
4. Service layer contains business logic
5. Models execute database queries via goqu

### Authentication

- **Ory Kratos** manages user registration, login, sessions
- JWT tokens for API authentication (guest users)
- Cookie-based sessions for Kratos-authenticated users
- WebSocket connections use custom auth middleware

### Real-Time Communication

- WebSocket connections for live quiz sessions
- Two socket endpoints:
  - `/api/v1/socket/admin/arrange/:session_id` — Admin quiz control
  - `/api/v1/socket/join/:invitation_code` — Participant quiz play
- Message format: `{component, event, data}`

### Database

- **PostgreSQL 15** with UUID primary keys
- **golang-migrate** for schema migrations
- **goqu** query builder for type-safe SQL
- **Valkey/Redis** for caching and session management

### Key Tables

| Table | Purpose |
|-------|---------|
| `users` | User accounts (Kratos-managed) |
| `quizzes` | Quiz metadata |
| `questions` | Quiz questions with options |
| `active_quizzes` | Running quiz sessions |
| `active_quiz_questions` | Questions in active session |
| `user_played_quizzes` | Quiz participation records |
| `user_quiz_responses` | Individual answers |
| `shared_quizzes` | Quiz sharing permissions |
| `quiz_categories` | Quiz categorization |
