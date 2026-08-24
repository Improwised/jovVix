import { test as base, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';

// Cache test data in memory to avoid repeated file I/O
let cachedTestData = null;

function getTestData() {
  if (!cachedTestData) {
    try {
      const dataPath = path.join('tests/.test-data', 'shared.json');
      if (fs.existsSync(dataPath)) {
         cachedTestData = JSON.parse(fs.readFileSync(dataPath, 'utf8'));
      } else {
         // Fallback or empty if setup failed (should catch elsewhere)
         cachedTestData = {}; 
      }
    } catch (e) {
      console.error('Error reading test data:', e);
      cachedTestData = {};
    }
  }
  return cachedTestData;
}

export const test = base.extend({
  authRequest: async ({ request }, use) => {
    await use(request);
  },

  createdQuiz: async ({}, use) => {
    const data = getTestData();
    await use({ quizId: data.quizId, quizTitle: data.quizTitle });
  },

  activeSession: async ({}, use) => {
    const data = getTestData();
    await use({ sessionId: data.sessionId, code: data.code });
  },
});

export { expect };
