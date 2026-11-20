# Architecture Overview

## System Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    Client (Browser)                       │
│               Nuxt 4 SSR / Vue 3 SPA                     │
│                Pinia State Management                     │
│              WebSocket (Real-time Quiz)                   │
│              HTTP Polling (Room State, 3s)                │
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

**Framework:** Nuxt 4 (SSR + SPA hybrid)

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
- Host connection lock via Redis `SETNX` prevents duplicate tabs

### Room System (Multi-Quiz Sessions)

- Rooms group multiple quiz sessions under one join code
- Room state is derived from database columns (no stored state column):
  - **lobby** — session active, no question started
  - **playing** — session active, question is open
  - **results** — session ended, room still open (between quizzes)
  - **closed** — room ended or expired
- Clients poll `GET /rooms/:id` every 3 seconds via `RoomPanel` component
- When the host starts the next quiz, `current_session_id` updates and clients auto-redirect
- Player names persist across quizzes via `quiz_room_members` table

### Solo Mode

- `active_quizzes.is_solo` marks a session as a solo attempt: same engine, same questions, same scoring, with the host as the only player
- Solo sessions create no room, accept no joins, and ignore host pacing events such as skip, pause and timer skip
- The quiz is fingerprinted into `active_quizzes.quiz_version` when the attempt starts; the solo leaderboard only compares attempts with the same fingerprint, so editing a quiz starts a new group
- `active_quizzes.solo_completed` separates finished attempts from abandoned ones
- Results are served by `GET /api/v1/solo/results/:attempt`
- Private quiz hosts cannot participate in their own session

### AI Quiz Generation

jovVix generates quiz questions through OpenAI-compatible chat-completions providers. Each user brings their own provider and API key, configured from the UI.

**Flow:**
1. User saves provider, model, and API key via `AiSettingsForm` (frontend)
2. Credentials are encrypted with AES-256-GCM (key derived via Argon2id) and stored in the `ai_settings` table; the Vault Password is never stored
3. The frontend unlocks the vault, then calls `/api/v1/ai/*` endpoints
4. `AIController` (`api/controllers/api/v1/ai_controller.go`) validates the request and resolves per-user credentials
5. `AIQuizService` (`api/services/ai_quiz.go`) builds the prompt, calls the provider, and parses the JSON response; it adapts the request if the provider rejects a parameter (JSON mode, `max_tokens`, temperature)
6. Generated questions are validated, normalized, and stored or returned to the frontend

**Endpoint group** (`api/routes/main.go`): `/api/v1/ai/{status,settings,models,test,questions/generate,quizzes}` plus `/api/v1/quizzes/:quiz_id/questions/ai{/generate}`.

See the [AI Quiz Generation Guide](./ai-quiz-generation.md) for providers, configuration, and troubleshooting.

### Database

- **PostgreSQL 15** with UUID primary keys
- **golang-migrate** for schema migrations
- **goqu** query builder for type-safe SQL
- **Valkey/Redis** for caching, session management, and host connection locking

### Key Tables

| Table | Purpose |
|-------|---------|
| `users` | User accounts (Kratos-managed) |
| `quizzes` | Quiz metadata |
| `questions` | Quiz questions with options |
| `active_quizzes` | Running quiz sessions (live and solo) |
| `active_quiz_questions` | Questions in active session |
| `user_played_quizzes` | Quiz participation records |
| `user_quiz_responses` | Individual answers |
| `shared_quizzes` | Quiz sharing permissions |
| `quiz_categories` | Quiz categorization |
| `ai_settings` | Encrypted AI provider credentials (AES-256-GCM + Argon2id) |
| `quiz_rooms` | Room grouping for multi-quiz sessions |
| `quiz_room_members` | Player name reservations across quizzes in a room |
