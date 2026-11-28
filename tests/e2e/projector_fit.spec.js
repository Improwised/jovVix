import { test, expect } from '@playwright/test';

// A host puts the lobby and the score page on a projector (#400). At common
// projector sizes the QR code, Start Quiz and the top 3 rankings must show
// without scrolling; phones must not scroll sideways.
// Runs on its own stack: it signs up a host and uploads its own quiz.
// Env: FRONTEND_URL (web), BASE_URL (API), as in playwright.config.js.

const WEB = process.env.FRONTEND_URL || 'http://127.0.0.1:5000';
const API = process.env.BASE_URL || 'http://127.0.0.1:3000';
const PROJECTORS = [
  { width: 1366, height: 768 },
  { width: 1280, height: 720 },
];
const PHONE = { width: 390, height: 844 };
const PLAYERS = ['Asha', 'Bilal', 'Chen'];
const blank = { cookies: [], origins: [] };

const fitsViewport = async (page, locator) => {
  const box = await locator.boundingBox();
  const { width, height } = page.viewportSize();
  expect(box, 'element is rendered').not.toBeNull();
  expect(box.y, 'top edge on screen').toBeGreaterThanOrEqual(0);
  expect(box.x, 'left edge on screen').toBeGreaterThanOrEqual(0);
  expect(box.y + box.height, 'bottom edge on screen').toBeLessThanOrEqual(height);
  expect(box.x + box.width, 'right edge on screen').toBeLessThanOrEqual(width);
};

const atTop = async (page, size) => {
  await page.setViewportSize(size);
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.waitForTimeout(300);
};

const noSideScroll = async (page) => {
  const overflow = await page.evaluate(
    () => document.documentElement.scrollWidth - window.innerWidth
  );
  expect(overflow, 'horizontal overflow in px').toBeLessThanOrEqual(0);
};

test('lobby and score page fit a projector without scrolling', async ({ browser }) => {
  test.setTimeout(180000);
  const host = await (
    await browser.newContext({ storageState: blank, baseURL: WEB, viewport: PROJECTORS[0] })
  ).newPage();

  // Sign up a host.
  await host.goto('/account/register');
  await host.fill('input[name="traits.name.first"]', 'Meera');
  await host.fill('input[name="traits.name.last"]', 'Host');
  await host.fill('input[name="traits.email"]', `projector${Date.now()}@example.test`);
  await host.fill('input[name="password"]', 'Quiz#Night2026!x');
  await host.click('button[type="submit"]');
  // Kratos signs the host in and sends them back to the app's home page.
  await host.waitForURL(
    (url) => url.origin === new URL(WEB).origin && !url.pathname.startsWith('/account'),
    { timeout: 30000 }
  );

  // Upload a one-question quiz and open its lobby. The calls run in the page,
  // as the app makes them: Kratos marks its cookies Secure, which the browser
  // sends to http://127.0.0.1 but Playwright's request client does not.
  const csv = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"Which planet is the largest?",single answer,10,Jupiter,Saturn,Mars,Venus,,1,text,text,`;
  const sessionId = await host.evaluate(
    async ({ api, csv, title }) => {
      const form = new FormData();
      form.append('description', 'Projector fit check');
      form.append('attachment', new Blob([csv], { type: 'text/csv' }), 'quiz.csv');
      const call = async (path, body) => {
        const res = await fetch(`${api}/api/v1/${path}`, {
          method: 'POST',
          credentials: 'include',
          body,
        });
        if (!res.ok) throw new Error(`${path}: ${res.status} ${await res.text()}`);
        return (await res.json()).data;
      };
      const quizId = await call(`quizzes/${title}/upload`, form);
      return call(`quizzes/${quizId}/demo_session`);
    },
    { api: API, csv, title: `Projector_${Date.now()}` }
  );
  await host.goto(`/admin/arrange/${sessionId}`);
  const codeText = host.locator('.code').first();
  await expect(codeText).toHaveText(/\d{6}/, { timeout: 30000 });
  const code = (await codeText.innerText()).replace(/\D/g, '');

  // (a) Lobby while the host is alone, as the room first sees it: the whole
  // QR code and Start Quiz are on the first screen.
  const qr = host.locator('.qr-card');
  const start = host.locator('button[type="submit"]', { hasText: 'Start Quiz' });
  for (const size of PROJECTORS) {
    await atTop(host, size);
    await fitsViewport(host, qr);
    await fitsViewport(host, start);
  }

  // Players join from phones.
  const players = [];
  for (const name of PLAYERS) {
    const player = await (
      await browser.newContext({ storageState: blank, baseURL: WEB, viewport: PHONE })
    ).newPage();
    await player.goto(`/join?code=${code}`);
    await player.fill('#username', name);
    await player.locator('button[type="submit"]').first().click();
    await expect(player.getByText('Ready, Steady, Go!')).toBeVisible({ timeout: 30000 });
    players.push(player);
  }
  for (const name of PLAYERS) {
    await expect(host.getByText(name, { exact: true })).toBeVisible({ timeout: 30000 });
  }

  // (a) Lobby with players: the QR code and Start Quiz still fit.
  for (const size of PROJECTORS) {
    await atTop(host, size);
    await fitsViewport(host, qr);
    await fitsViewport(host, start);
  }
  // (c) Lobby on a phone: no sideways scroll.
  await atTop(host, PHONE);
  await noSideScroll(host);
  await noSideScroll(players[0]);

  // One question: everyone answers, then the score page shows.
  await atTop(host, PROJECTORS[0]);
  await start.click();
  for (const player of players) {
    const answer = player.getByRole('button', { name: /Jupiter/ });
    await expect(answer).toBeEnabled({ timeout: 30000 });
    await answer.click();
  }
  await expect(host.getByRole('heading', { name: 'Leaderboard' })).toBeVisible({ timeout: 60000 });
  await expect(host.getByText('Rank board')).toHaveCount(0);

  // (b) Score page: the top 3 rankings are on the first screen.
  for (const size of PROJECTORS) {
    await atTop(host, size);
    for (const name of PLAYERS) {
      await fitsViewport(host, host.getByText(name, { exact: true }));
    }
  }
  // (c) Score page on a phone: no sideways scroll.
  await atTop(host, PHONE);
  await noSideScroll(host);
  await noSideScroll(players[0]);
});
