# Frontend Development Guide

## Project Structure

```
app/
├── pages/            # File-based routing (Nuxt)
├── components/       # Vue components (auto-imported)
│   ├── ui/           # shadcn-vue components
│   ├── common/       # Shared layout components
│   ├── Quiz/         # Quiz-specific components
│   └── landing/      # Homepage components
├── composables/      # Reusable logic (auto-imported)
├── store/            # Pinia state management
├── plugins/          # Nuxt plugins
├── layouts/          # Page layouts
├── assets/           # Static assets
├── config/           # App config
└── public/           # Public static files
```

## Adding a New Page

Create a file in `pages/` — Nuxt handles routing automatically:

```bash
# File: pages/admin/settings.vue
# Route: /admin/settings
```

### Dynamic Routes

```bash
# File: pages/quiz/[id].vue
# Route: /quiz/:id
```

### Nested Routes

```bash
# File: pages/admin/reports/index.vue  → /admin/reports
# File: pages/admin/reports/[id].vue   → /admin/reports/:id
```

## Adding a Component

Create a `.vue` file in `components/` — auto-imported everywhere:

```vue
<script setup lang="ts">
const props = defineProps<{
  title: string
  count: number
}>()

const emit = defineEmits<{
  select: [id: string]
}>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ title }}</CardTitle>
    </CardHeader>
    <CardContent>
      <p>Count: {{ count }}</p>
      <Button @click="emit('select', 'some-id')">
        Select
      </Button>
    </CardContent>
  </Card>
</template>
```

### Using shadcn-vue Components

Components in `components/ui/` are available globally:

```vue
<template>
  <Button variant="default">Click me</Button>
  <Badge variant="secondary">New</Badge>
  <Card>
    <CardContent>Content</CardContent>
  </Card>
  <Input placeholder="Type here" />
</template>
```

### Adding a New shadcn-vue Component

```bash
cd app
npx shadcn-vue@2.2.0 add <component-name>
```

## Adding a Composable

Create a `.js` or `.ts` file in `composables/` — auto-imported:

```typescript
// composables/useQuizTimer.ts
export function useQuizTimer(duration: number) {
  const timeRemaining = ref(duration)
  const isRunning = ref(false)

  function start() {
    isRunning.value = true
    const interval = setInterval(() => {
      if (timeRemaining.value <= 0) {
        clearInterval(interval)
        isRunning.value = false
        return
      }
      timeRemaining.value--
    }, 1000)
  }

  return { timeRemaining, isRunning, start }
}
```

Usage in any component (auto-imported):
```vue
<script setup>
const { timeRemaining, start } = useQuizTimer(60)
</script>
```

## Adding a Pinia Store

Create a file in `store/`:

```typescript
// store/quizList.ts
export const useQuizListStore = defineStore('quizList', () => {
  const quizzes = ref([])
  const loading = ref(false)

  async function fetchQuizzes() {
    loading.value = true
    try {
      const config = useRuntimeConfig()
      const data = await $fetch(`${config.public.apiUrl}/quizzes`, {
        credentials: 'include',
      })
      quizzes.value = data.data
    } finally {
      loading.value = false
    }
  }

  return { quizzes, loading, fetchQuizzes }
})
```

### Persisting State

```typescript
export const useUserStore = defineStore('users', () => {
  // ... state and actions
}, {
  persist: true, // Full state persistence
})
```

Or selectively:
```typescript
{
  persist: {
    paths: ['activeQuizTitle'], // Only persist specific fields
  },
}
```

## API Communication

### HTTP Requests

Use Nuxt's `$fetch` or `useFetch`:

```typescript
// Client-side only
const data = await $fetch('/api/v1/quizzes', {
  credentials: 'include', // Send cookies
})

// SSR-compatible
const { data } = await useFetch('/api/v1/quizzes', {
  credentials: 'include',
})
```

### WebSocket Connections

Use the `quiz_operation` composable pattern:

```typescript
const handler = new QuizHandler(
  socketUrl,
  identifier,
  (component, event, data) => {
    // Handle messages
  },
  { reconnect: true, maxRetries: 3 }
)
```

## Testing

### Unit Tests with Vitest

```typescript
// components/MyComponent.test.js
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MyComponent from './MyComponent.vue'

describe('MyComponent', () => {
  it('renders correctly', () => {
    const wrapper = mount(MyComponent, {
      props: { title: 'Test' }
    })
    expect(wrapper.text()).toContain('Test')
  })
})
```

### Running Tests

```bash
npm run test           # Run once
npm run test:watch     # Watch mode
npm run test:coverage  # With coverage
```

## Code Style

- ESLint with `@typescript-eslint`, `vue3-recommended`, `prettier`
- Zero warnings allowed (`--max-warnings 0`)
- Run `npm run lint-fix` to auto-fix
- Components: `<script setup lang="ts">` preferred
- Use Tailwind CSS v4 utility classes
