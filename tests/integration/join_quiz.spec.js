import { test, expect } from '../fixtures';

test.describe('Join Quiz (Integration)', () => {

  test('Player waits in lobby until admin starts quiz', async ({ activeSession, browser }) => {
    const code = activeSession.code;

    const { request } = await import('@playwright/test');
    const playerRequest = await request.newContext();

    const username = `P_${Math.floor(1000 + Math.random() * 9000)}`;
    const guestRes = await playerRequest.post(`/api/v1/user/${username}?avatar_name=1.webp`);
    expect(guestRes.status()).toBe(200);

    const guestCookies = guestRes.headers()['set-cookie'];
    const cookieHeader = Array.isArray(guestCookies) ? guestCookies.join('; ') : (guestCookies || '');

    const joinRes = await playerRequest.post(`/api/v1/user_played_quizes/${code}`, {
      headers: {
        'Cookie': cookieHeader
      }
    });
    expect(joinRes.status()).toBe(200);
    const joinData = await joinRes.json();

    expect(joinData.data).toHaveProperty('user_played_quiz');

    await playerRequest.dispose();
  });

});
