# API Development Guide

This guide covers adding endpoints, middleware, and database migrations to the jovVix API.

## Project Structure

```
api/
├── routes/main.go          # Route registration
├── controllers/api/v1/     # Request handlers
├── middlewares/             # Auth, logging, permissions
├── models/                 # Database queries (goqu)
├── services/               # Business logic
├── database/migrations/    # SQL migrations
├── config/                 # Environment config
└── constants/              # App constants
```

## Adding a New Endpoint

### 1. Define the Route

In `api/routes/main.go`, add your route to the appropriate group:

```go
// Example: GET /api/v1/quizzes/:quiz_id/stats
quizGroup.Get("/:quiz_id/stats", controllers.GetQuizStats)
```

### 2. Create the Controller

Create or modify `api/controllers/api/v1/quiz_controller.go`:

```go
func GetQuizStats(c *fiber.Ctx) error {
    quizID := c.Params("quiz_id")

    stats, err := services.GetQuizStats(quizID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Failed to get quiz stats",
        })
    }

    return c.JSON(fiber.Map{
        "data":    stats,
        "message": "Success",
    })
}
```

### 3. Add the Service (Business Logic)

Create or modify `api/services/quiz_service.go`:

```go
func GetQuizStats(quizID string) (*models.QuizStats, error) {
    // Business logic here
    return models.GetQuizStats(quizID)
}
```

### 4. Add the Model (Database Query)

Create or modify `api/models/quiz.go`:

```go
func GetQuizStats(quizID string) (*QuizStats, error) {
    db := database.DB

    query := db.Goqu.Select(
        "q.id",
        goqu.COUNT("qq.question_id").As("total_questions"),
    ).From(
        goqu.T("quizzes").As("q"),
    ).Join(
        goqu.T("quiz_questions").As("qq"),
        goqu.On(goqu.C("q.id").Eq(goqu.C("qq.quiz_id"))),
    ).Where(
        goqu.C("q.id").Eq(quizID),
    ).GroupBy(goqu.C("q.id"))

    var stats QuizStats
    _, err := query.ScanStructs(&stats)
    return &stats, err
}
```

### 5. Add Swagger Annotations

Add Swagger comments to your controller for API documentation:

```go
// GetQuizStats godoc
// @Summary      Get quiz statistics
// @Description  Get statistics for a specific quiz
// @Tags         quizzes
// @Accept       json
// @Produce      json
// @Param        quiz_id path string true "Quiz ID"
// @Success      200 {object} models.QuizStats
// @Router       /quizzes/{quiz_id}/stats [get]
func GetQuizStats(c *fiber.Ctx) error { ... }
```

### 6. Regenerate Swagger

```bash
cd api
make swagger-gen
```

## Middleware

### Adding Authentication

Use existing middleware from `api/middlewares/`:

```go
// Kratos-authenticated (registered users only)
router.Use(middlewares.KratosAuthenticated)

// Custom JWT auth (guest users)
router.Use(middlewares.CustomAuthenticated)

// Quiz permission check
router.Use(middlewares.QuizPermission)

// Verify edit access
router.Use(middlewares.VerifyQuizEditAccess)
```

### Creating Custom Middleware

```go
func MyMiddleware(c *fiber.Ctx) error {
    // Pre-processing
    value := c.Get("X-Custom-Header")

    // Call next handler
    err := c.Next()
    if err != nil {
        return err
    }

    // Post-processing
    return nil
}
```

## Database Migrations

### Creating a Migration

```bash
cd api
make migrate file_name=add_new_column
```

This creates:
- `database/migrations/YYYYMMDDHHMMSS_add_new_column.up.sql`
- `database/migrations/YYYYMMDDHHMMSS_add_new_column.down.sql`

### Writing Migrations

**Up migration:**
```sql
-- +migrate Up
ALTER TABLE quizzes ADD COLUMN new_column TEXT NOT NULL DEFAULT '';
```

**Down migration:**
```sql
-- +migrate Down
ALTER TABLE quizzes DROP COLUMN new_column;
```

### Running Migrations

```bash
go run app.go migrate up     # Apply pending migrations
go run app.go migrate down   # Rollback last migration
```

## Testing

### Running Tests

```bash
go test ./...              # All tests
go test -v ./...           # Verbose
go test -cover ./...       # With coverage
make test-coverage         # Coverage report
```

### Writing Tests

```go
func TestGetQuizStats(t *testing.T) {
    // Setup test database
    db := testutils.SetupTestDB()
    defer db.Close()

    // Insert test data
    quizID := testutils.CreateTestQuiz(db)

    // Call function
    stats, err := GetQuizStats(quizID)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, quizID, stats.ID)
}
```

## Code Style

- Follow standard Go conventions
- Use `golangci-lint` for enforcement: `golangci-lint run`
- Error handling: always check and handle errors
- Naming: PascalCase for exported, camelCase for unexported
- Comments: godoc format for exported functions
