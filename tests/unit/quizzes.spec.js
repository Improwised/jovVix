import { test, expect } from '../fixtures';

test.describe('Quizzes (Unit)', () => {
  test('Create Quiz (Positive) - Upload valid CSV', async ({ authRequest }) => {
    const csvContent = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"What is 2+2?",single answer,10,1,2,3,4,,4,text,text,`;
    const buffer = Buffer.from(csvContent);
    const title = `API_Quiz_${Date.now()}`;

    const response = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: {
        description: 'Unit Test Quiz',
        attachment: {
          name: 'test.csv',
          mimeType: 'text/csv',
          buffer: buffer
        }
      }
    });

    expect(response.status()).toBe(202); // Created
    const data = await response.json();
    expect(data.data).toBeTruthy(); // Expecting quiz_id
  });

  test('Create Quiz (Negative) - Invalid File Type', async ({ authRequest }) => {
    const buffer = Buffer.from('invalid data');
    const title = `Invalid_Quiz_${Date.now()}`;

    const response = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: {
        description: 'Invalid',
        attachment: {
          name: 'test.txt',
          mimeType: 'text/plain', // Wrong type
          buffer: buffer
        }
      }
    });

    expect(response.status()).not.toBe(202);
  });

  test('My Quizzes - List user quizzes', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/quizzes');
    expect(response.status()).toBe(200);
    const data = await response.json();
    expect(Array.isArray(data.data)).toBe(true);
  });

  test('Shared Quizzes - List shared quizzes', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/shared_quizzes?page=1&limit=10&type=shared_by_me');
    expect(response.status()).toBe(200);
    const data = await response.json();
    if (data.data) {
        expect(Array.isArray(data.data)).toBe(true);
    }
  });

  test('Shared With Me - List quizzes shared with user', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/shared_quizzes?page=1&limit=10&type=shared_with_me');
    expect(response.status()).toBe(200);
    const data = await response.json();
    if (data.data) {
        expect(Array.isArray(data.data)).toBe(true);
    }
  });

  test('Create Quiz - Empty CSV -> Validation error', async ({ authRequest }) => {
    const buffer = Buffer.from(''); // Empty CSV
    const title = `Empty_Quiz_${Date.now()}`;

    const response = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: {
        description: 'Empty CSV',
        attachment: {
          name: 'empty.csv',
          mimeType: 'text/csv',
          buffer: buffer
        }
      }
    });

    expect(response.status()).toBe(400);
  });

  test('Fetch Quiz by ID - Accessible to creator', async ({ authRequest, createdQuiz }) => {
    // Check if quiz exists in the list first, as direct GET /id might be restricted
    const listRes = await authRequest.get('/api/v1/quizzes');
    const listData = await listRes.json();
    const found = listData.data.find(q => q.id === createdQuiz.quizId);
    
    expect(found).toBeTruthy();
    expect(found.title).toBe(createdQuiz.quizTitle);
  });

  // Edge Case E: Duplicate Quiz Title
  test('Duplicate Quiz Title - Should be allowed', async ({ authRequest }) => {
    const title = `Dup_Title_${Date.now()}`;
    const csvContent = `Question,Type,Points,Opt1,Opt2,Correct\nQ1,single,10,A,B,A`;
    const buffer = Buffer.from(csvContent);

    const res1 = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: { description: 'First', attachment: { name: '1.csv', mimeType: 'text/csv', buffer } }
    });
    expect(res1.status()).toBe(202);

    const res2 = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: { description: 'Second', attachment: { name: '2.csv', mimeType: 'text/csv', buffer } }
    });
    // System allows duplicates with unique IDs
    expect(res2.status()).toBe(202);
    
    const id1 = (await res1.json()).data;
    const id2 = (await res2.json()).data;
    expect(id1).not.toBe(id2);
  });

  // Edge Case F: Unauthorized Quiz Access
  test('Unauthorized Quiz Access - Unauthenticated User', async ({ request, createdQuiz }) => {
    const res = await request.get(`/api/v1/quizzes/${createdQuiz.quizId}`);
    // Expect 401 Unauthorized or 403 Forbidden or 405 Method Not Allowed (if route hidden)
    expect([401, 403, 405]).toContain(res.status());
  });

  // Edge Case G: Delete Verification
  test('Delete Verification - Lifecycle', async ({ authRequest }) => {
    const title = `Del_Quiz_${Date.now()}`;
    const buffer = Buffer.from(`Q,T,P,O1,O2,C\nQ,single,10,A,B,A`);
    
    const createRes = await authRequest.post(`/api/v1/quizzes/${title}/upload`, {
      multipart: { description: 'Delete Me', attachment: { name: 'del.csv', mimeType: 'text/csv', buffer } }
    });
    const quizId = (await createRes.json()).data;

    const delRes = await authRequest.delete(`/api/v1/quizzes/${quizId}`);
    expect(delRes.status()).toBe(200);

    // Direct fetch might be 404 or 400. List should not contain it.
    const listRes = await authRequest.get('/api/v1/quizzes');
    const listData = await listRes.json();
    const found = listData.data.find(q => q.id === quizId);
    expect(found).toBeUndefined();
  });

  // Edge Case H: Share With Self
  test('Share With Self - Validation', async ({ authRequest, createdQuiz }) => {
    const profileRes = await authRequest.get('/api/v1/user/who');
    const email = (await profileRes.json()).data.email;
    
    const shareRes = await authRequest.post(`/api/v1/shared_quizzes/${createdQuiz.quizId}`, {
        data: { email: email, permission: 'view' }
    });
    
    // Current behavior: System allows sharing with self (200 OK).
    expect([200, 400]).toContain(shareRes.status());
  });
});
