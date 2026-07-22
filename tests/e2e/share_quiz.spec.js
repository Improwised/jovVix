import { test, expect } from '../fixtures';

test.describe('Share Quiz (Smoke)', () => {
  test('Share Quiz Flow: Create -> Share -> Success', async ({ page }) => {
    const quizTitle = `E2E_Share_${Date.now()}`;
    const friendEmail = `e2e_friend_${Date.now()}@example.com`;

    await page.locator('input#title').fill(quizTitle);
    await page.locator('input#description').fill('E2E Smoke Test');
    
    const csvContent = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"What is 2+2?",single answer,10,1,2,3,4,,4,text,text,`;
    const buffer = Buffer.from(csvContent);
    
    await page.setInputFiles('input[id="attachment"]', {
      name: 'test.csv',
      mimeType: 'text/csv',
      buffer: buffer
    });

    await page.locator('button:has-text("Create Quiz")').click();
    
    await expect(page.locator('.Vue-Toastification__toast--success')).toContainText('file uploaded successfully', { timeout: 15000 });

    await page.locator('text=Quizzes').click();
    await page.locator('text=My Quizzes').click();
    await expect(page).toHaveURL(/\/admin\/quiz\/list-quiz/);

    const card = page.locator('.card', { hasText: quizTitle }).first();
    await expect(card).toBeVisible();
    await card.locator('a:has-text("View Quiz")').click();

    const shareBtn = page.locator('button[title="Share Quiz"]');
    await expect(shareBtn).toBeVisible({ timeout: 10000 });
    await shareBtn.click();

    const modal = page.locator('#shareQuizModal');
    await expect(modal).toBeVisible();

    await modal.locator('button[title="Add People"]').click();

    await modal.locator('input#email').fill(friendEmail);
    await modal.locator('select#permission').selectOption('write');
    await modal.locator('button:has-text("Share Quiz")').click();

    await expect(page.locator('.Vue-Toastification__toast--success', { hasText: 'Quiz shared successfully' })).toBeVisible();
  });
});
