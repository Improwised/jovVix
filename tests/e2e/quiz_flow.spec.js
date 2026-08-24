import { test, expect } from '../fixtures';

test.describe('Full Quiz Lifecycle - Creator to Player', () => {

  test('Complete Quiz Creation to Leaderboard', async ({ browser }) => {
    test.setTimeout(120000);
    const hostContext = await browser.newContext({ storageState: 'tests/.auth/admin.json' });
    const hostPage = await hostContext.newPage();

    const playerContext = await browser.newContext({ storageState: { cookies: [], origins: [] } });
    const playerPage = await playerContext.newPage();

    await hostPage.goto('/admin/quiz/create-quiz');
    await expect(hostPage).toHaveURL(/\/admin\/quiz\/create-quiz/);

    const quizTitle = `E2E_Quiz_${Date.now()}`;

    await expect(hostPage.locator('#title')).toBeVisible({ timeout: 10000 });
    await hostPage.locator('#title').fill(quizTitle);

    await hostPage.locator('#description').fill('E2E Test Quiz');

    const csvContent = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"What is 1+1?",single answer,10,2,3,4,5,,2,text,text,`;

    const fileBuffer = Buffer.from(csvContent);
    const fileInput = hostPage.locator('#attachment');
    await fileInput.setInputFiles({
      name: 'e2e_quiz.csv',
      mimeType: 'text/csv',
      buffer: fileBuffer
    });

    await hostPage.locator('button[type=submit]').click();

    await expect(hostPage.locator('text=Start Quiz')).toBeVisible({ timeout: 15000 });

    await hostPage.locator('text=Start Quiz').click();

    await expect(hostPage).toHaveURL(/\/admin\/arrange\//, { timeout: 15000 });

    const gameCode = await hostPage.locator('h2.code').textContent();
    expect(gameCode).toBeTruthy();
    expect(gameCode.length).toBe(6); // Assuming 6-digit code

    await playerPage.goto('/join');
    
    await expect(playerPage.getByRole('heading', { name: 'Join Quiz' })).toBeVisible({ timeout: 20000 });

    await expect(playerPage.getByLabel('User Name')).toBeVisible({ timeout: 30000 });
    await playerPage.getByLabel('User Name').fill('E2E_Player');

    const codeInputs = playerPage.locator('input[placeholder="0"]');
    if (await codeInputs.count() > 0) {
      // OTP input
      const digits = gameCode.split('');
      for (let i = 0; i < digits.length; i++) {
        await codeInputs.nth(i).fill(digits[i]);
      }
    } else {
      await playerPage.locator('input[name="code"]').fill(gameCode);
    }

    await playerPage.locator('button[type=submit]').click();

    await expect(playerPage.locator('text=Ready Steady Go')).toBeVisible();

    await expect(hostPage.locator('text=E2E_Player')).toBeVisible();

    await hostPage.locator('button:has-text("Start Quiz")').click();

    // Wait for question to appear (socket delay)
    await expect(playerPage.locator('h3.font-bold')).toContainText('What is 1+1?', { timeout: 30000 });

    await playerPage.locator('.option-box').first().click();

    await expect(playerPage.locator('text=Answer Submitted')).toBeVisible();

    // For 1-question quiz, 'Finish' appears directly.
    await hostPage.locator('button:has-text("Finish")').click();

    await expect(hostPage).toHaveURL(/\/admin\/scoreboard/, { timeout: 15000 });
    
    await hostPage.locator('button:has-text("Next")').click();

    await expect(hostPage.getByRole('heading', { name: 'Rankings' })).toBeVisible({ timeout: 15000 });

    await expect(hostPage.locator('text=E2E_Player')).toBeVisible();

    // Player should see completion or redirect
    await Promise.race([
        expect(playerPage).toHaveURL(/\/scoreboard/),
        expect(playerPage.locator('text=Rank')).toBeVisible(),
        expect(playerPage.locator('text=Game Over')).toBeVisible()
    ]);

    await hostContext.close();
    await playerContext.close();
  });
});