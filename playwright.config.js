// @ts-check
import { defineConfig, devices } from '@playwright/test';

/**
 * @see https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  testDir: './tests',
  /* Global Setup via Project Dependencies */
  
  /* Run tests in files in parallel */
  fullyParallel: true,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: 'html',
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.BASE_URL || 'http://127.0.0.1:3000',

    /* Use the authenticated state */
    storageState: 'tests/.auth/admin.json',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
    
    /* Ignore HTTP errors for testing negative paths */
    ignoreHTTPSErrors: true,
  },

  /* Configure projects for major browsers and API layers */
  projects: [
    {
      name: 'setup',
      testMatch: /auth\.setup\.js/,
    },
    {
      name: 'unit',
      testMatch: /tests\/unit\/.*\.spec\.js/,
      dependencies: ['setup'],
    },
    {
      name: 'integration',
      testMatch: /tests\/integration\/.*\.spec\.js/,
      use: { ...devices['Desktop Chrome'] },
      dependencies: ['setup'], // Depends on auth
      workers: 1, // Run sequentially to avoid session conflicts
    },
    {
      name: 'e2e',
      testMatch: /tests\/e2e\/.*\.spec\.js/,
      use: {
        ...devices['Desktop Chrome'],
        baseURL: process.env.FRONTEND_URL || 'http://127.0.0.1:5000'
      },
      dependencies: ['setup'], // Depends on auth
      workers: 1, // Run sequentially
    },
  ],
});
