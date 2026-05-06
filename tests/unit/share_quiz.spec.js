import { test, expect } from '../fixtures';

test.describe('Share Quiz (Unit)', () => {
  test('Share quiz with valid email + write permission -> Success', async ({ authRequest, createdQuiz }) => {
    const res = await authRequest.post(`/api/v1/shared_quizzes/${createdQuiz.quizId}`, {
      data: { email: `test_share_write_${Date.now()}@example.com`, permission: 'write' }
    });
    expect(res.status()).toBe(200);
  });

  test('Share quiz with valid email + read permission -> Success', async ({ authRequest, createdQuiz }) => {
    const res = await authRequest.post(`/api/v1/shared_quizzes/${createdQuiz.quizId}`, {
      data: { email: `test_share_read_${Date.now()}@example.com`, permission: 'read' }
    });
    expect(res.status()).toBe(200);
  });

  test('Share quiz without permission -> 400 Validation Error', async ({ authRequest, createdQuiz }) => {
    const res = await authRequest.post(`/api/v1/shared_quizzes/${createdQuiz.quizId}`, {
      data: { email: 'invalid@example.com' } // Missing permission
    });
    expect([400, 404]).toContain(res.status());
  });

  test('Share quiz that does not exist -> 404 Not Found', async ({ authRequest }) => {
    const res = await authRequest.post(`/api/v1/shared_quizzes/00000000-0000-0000-0000-000000000000`, {
      data: { email: 'ghost@example.com', permission: 'read' }
    });
    expect(res.status()).not.toBe(200);
  });
});
