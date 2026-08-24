import { test, expect } from '../fixtures';

test.describe('Quizzes (Integration)', () => {

  test('Create, Verify, and Delete Quiz', async ({ authRequest }) => {
    const csvContent = `Question,Type,Points,Opt1,Opt2,Correct\nQ1,single,10,A,B,A`;
    const buffer = Buffer.from(csvContent);
    const title = `Integration_Quiz_${Date.now()}`;

    const createRes = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: {
        description: 'Integration Flow',
        attachment: { name: 'test.csv', mimeType: 'text/csv', buffer }
      }
    });
    expect(createRes.status()).toBe(202);
    const createData = await createRes.json();
    const quizId = createData.data;

    const listRes = await authRequest.get('/api/v1/quizzes');
    const listData = await listRes.json();
    const found = listData.data.find(q => q.id === quizId);
    expect(found).toBeTruthy();
    expect(found.title).toBe(title);

    const deleteRes = await authRequest.delete(`/api/v1/quizzes/${quizId}`);
    expect(deleteRes.status()).toBe(200);

    const verifyRes = await authRequest.get('/api/v1/quizzes');
    const verifyData = await verifyRes.json();
    const notFound = verifyData.data.find(q => q.id === quizId);
    expect(notFound).toBeFalsy();
  });

});
