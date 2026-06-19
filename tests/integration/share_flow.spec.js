import { test, expect } from '../fixtures';

test.describe('Share Quiz (Integration)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin/quiz/list-quiz');
  });

  test('UI Share Flow: Modal -> Submit -> Verify API', async ({ page, authRequest }) => {
    const frontendUrl = process.env.FRONTEND_URL || 'http://127.0.0.1:5000';
    await page.goto(`${frontendUrl}/admin/quiz/create-quiz`);
    
    const quizTitle = `UI_Int_Quiz_${Date.now()}`;
    await page.locator('input#title').fill(quizTitle);
    await page.locator('input#description').fill('Integration Test Description');
    
    const csvContent = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"What is 2+2?",single answer,10,1,2,3,4,,4,text,text,`;
    const buffer = Buffer.from(csvContent);
    
    await page.setInputFiles('input[id="attachment"]', {
      name: 'test.csv',
      mimeType: 'text/csv',
      buffer: buffer
    });

    const responsePromise = page.waitForResponse(response => 
      response.url().includes('/upload') && response.status() === 202
    );

    await page.locator('button:has-text("Create Quiz")').click();

    const response = await responsePromise;
    const responseBody = await response.json();
    const quizId = responseBody.data;
    
    await expect(page.locator('.Vue-Toastification__toast--success')).toBeVisible();

    // Perform Client-Side Navigation to avoid SSR Docker networking issues
    await page.locator('text=Quizzes').click();
    
    await expect(page).toHaveURL(/\/admin\/quiz$/);
    
    await page.locator('text=My Quizzes').click();
    await expect(page).toHaveURL(/\/admin\/quiz\/list-quiz/);

    const card = page.locator('.card', { hasText: quizTitle }).first();
    await expect(card).toBeVisible();
    await card.locator('a:has-text("View Quiz")').click();

    // Now client-side fetch should succeed
    await expect(page.locator('text=Total Questions')).toBeVisible({ timeout: 15000 });

    const shareBtn = page.locator('button[title="Share Quiz"]');
    await expect(shareBtn).toBeVisible();
    await shareBtn.click();

    const modal = page.locator('#shareQuizModal');
    await expect(modal).toBeVisible();
    await modal.locator('button[title="Add People"]').click();

    const friendEmail = `int_friend_${Date.now()}@example.com`;
    await modal.locator('input#email').fill(friendEmail);
    await modal.locator('select#permission').selectOption('write');

    await modal.locator('button:has-text("Share Quiz")').click();

    // Use specific text to distinguish from previous "File uploaded" toast
    await expect(page.locator('.Vue-Toastification__toast--success', { hasText: 'Quiz shared successfully' })).toBeVisible();

    // Use polling to handle potential DB/API consistency latency
    await expect.poll(async () => {
        const apiRes = await authRequest.get(`/api/v1/shared_quizzes/${quizId}`);
        if (!apiRes.ok()) return false;
        const apiData = await apiRes.json();
        if (!apiData.data) return false;
        return apiData.data.find(u => u.shared_to === friendEmail);
    }, {
        timeout: 5000,
        intervals: [1000]
    }).toBeTruthy();
    
    const finalRes = await authRequest.get(`/api/v1/shared_quizzes/${quizId}`);
    const finalData = await finalRes.json();
    const sharedUser = finalData.data.find(u => u.shared_to === friendEmail);
    expect(sharedUser.permission).toBe('write');
  });
});
