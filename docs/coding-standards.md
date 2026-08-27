# Coding Standards

## Go Backend

### Linting

We use `golangci-lint` with the project configuration in `api/.golangci.yaml`:

```bash
cd api
golangci-lint run
```

### Style

- Follow [Effective Go](https://go.dev/doc/effective_go) conventions
- Use `gofmt` / `goimports` for formatting
- Exported functions: PascalCase, godoc comments required
- Unexported functions: camelCase
- Package names: lowercase, single-word, no underscores
- File names: lowercase with underscores (`quiz_controller.go`)

### Error Handling

```go
// Always check errors
result, err := someFunction()
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Use structured errors
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "message": "Invalid quiz ID",
    "error":   "invalid_id",
})
```

### Database Queries

- Use goqu query builder (not raw SQL strings)
- Always use parameterized queries
- Handle `sql.Null*` types explicitly

```go
query := db.Goqu.Select("id", "title").
    From("quizzes").
    Where(goqu.C("creator_id").Eq(userID))
```

### Testing

- Table-driven tests for multiple cases
- Use `testify/assert` for assertions
- Test file naming: `*_test.go`
- Aim for meaningful coverage, not 100%

```go
func TestCalculateScore(t *testing.T) {
    tests := []struct {
        name     string
        answers  []int
        correct  []int
        points   int16
        expected int16
    }{
        {"all correct", []int{1, 2}, []int{1, 2}, 10, 20},
        {"partial", []int{1, 0}, []int{1, 2}, 10, 10},
        {"none correct", []int{0, 0}, []int{1, 2}, 10, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CalculateScore(tt.answers, tt.correct, tt.points)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Vue/TypeScript Frontend

### Linting

ESLint with zero warnings tolerance:

```bash
cd app
npm run lint        # Check
npm run lint-fix    # Auto-fix
```

### Component Style

- Always use `<script setup lang="ts">`
- Props: use `defineProps<{}>()` with TypeScript types
- Emits: use `defineEmits<{}>()`
- Prefer composition API over options API

```vue
<script setup lang="ts">
interface Props {
  title: string
  count?: number
}

const props = withDefaults(defineProps<Props>(), {
  count: 0,
})

const emit = defineEmits<{
  update: [value: string]
}>()
</script>
```

### Composables

- File naming: `use*.ts` or descriptive `.js`
- Return refs, not raw values
- Use `useRuntimeConfig()` for API URLs

```typescript
export function useQuizData(quizId: Ref<string>) {
  const data = ref(null)
  const loading = ref(false)

  async function fetch() { ... }

  return { data, loading, fetch }
}
```

### Naming Conventions

| Type | Convention | Example |
|------|-----------|---------|
| Components | PascalCase | `QuizListCard.vue` |
| Composables | camelCase | `useQuizTimer.ts` |
| Stores | camelCase | `useQuizListStore` |
| Pages | kebab-case dirs | `admin/quiz/list-quiz/index.vue` |
| CSS classes | Tailwind utilities | `flex items-center gap-4` |

### File Organization

- One component per file
- Keep files under 300 lines
- Extract complex logic to composables
- Group related components in subdirectories (`Quiz/`, `common/`)

## General

### Git

- Conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `chore:`
- Branch naming: `feat/`, `fix/`, `docs/`, `refactor/`, `test/`
- No secrets in commits or environment files committed to repo

### Documentation

- Use Markdown for all docs
- Include code examples
- Keep docs updated when changing code
- Link between related docs
