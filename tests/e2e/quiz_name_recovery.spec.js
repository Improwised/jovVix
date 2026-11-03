import { test, expect } from '../fixtures';

test.use({ storageState: { cookies: [], origins: [] } });

test('a direct play link with a taken name returns to the join form', async ({ page }) => {
  await page.route('**/api/v1/user_played_quizes/123456', async (route) => {
    await route.fulfill({ status: 409, contentType: 'application/json', body: JSON.stringify({ status: 'fail', data: { message: 'Name taken', suggestions: ['Smit42'] } }) });
  });
  await page.route('**/api/v1/user/who', async (route) => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ status: 'success', data: { firstname: 'Smit', username: 'internal', role: 'admin-user' } }) });
  });
  await page.goto('/join/play/123456');
  await expect(page).toHaveURL(/\/join\?/);
  await expect(page.getByText(/name is already taken/i)).toBeVisible();
  await expect(page.getByLabel('Invitation Code')).toHaveValue('123 456');
  await expect(page.getByText('That quiz has ended')).toHaveCount(0);
});

test('a delayed profile response does not replace a name the player has typed', async ({ page }) => {
  let release;
  const waitForRelease = new Promise((resolve) => { release = resolve; });
  await page.route('**/api/v1/user/who', async (route) => {
    await waitForRelease;
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ status: 'success', data: { firstname: 'Smit', username: 'internal', role: 'admin-user' } }) });
  });
  const profileRequest = page.waitForRequest('**/api/v1/user/who');
  await page.goto('/join', { waitUntil: 'domcontentloaded' });
  await profileRequest;
  const input = page.getByLabel('Your quiz name');
  await input.fill('Custom');
  const profileResponse = page.waitForResponse((response) => response.url().includes('/api/v1/user/who'));
  release();
  await profileResponse;
  await expect(input).toHaveValue('Custom');
  await expect(page.getByRole('button', { name: 'Enter Lobby' })).toBeEnabled();
});
