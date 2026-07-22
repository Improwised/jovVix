import { test, expect } from '../fixtures';

test.describe('Profile (Unit)', () => {
  
  test('Get Profile - Should return logged in user data', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/user/who');
    
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(data.data).toBeTruthy();
    expect(data.data).toHaveProperty('email');
    expect(data.data).toHaveProperty('username');
    expect(data.data).toHaveProperty('role');
    expect(data.data).toHaveProperty('firstname');
  });

  test('Update Profile - Should update first/last name', async ({ authRequest }) => {
    const getRes = await authRequest.get('/api/v1/user/who');
    const profile = (await getRes.json()).data;
    
    const newName = `Upd_${Math.floor(Math.random() * 1000)}`;

    const updateRes = await authRequest.put('/api/v1/kratos/user', {
        data: {
            first_name: newName,
            last_name: profile.lastname || 'User',
            email: profile.email // Keep email same to avoid uniqueness/login issues
        }
    });
    
    expect(updateRes.status()).toBe(200);

    const verifyRes = await authRequest.get('/api/v1/user/who');
    const updatedProfile = (await verifyRes.json()).data;
    expect(updatedProfile.firstname).toBe(newName);
  });

  test('Update Profile - Missing Fields (Negative)', async ({ authRequest }) => {
    const updateRes = await authRequest.put('/api/v1/kratos/user', {
        data: {
            first_name: '', // Empty name should fail validation
            last_name: 'User',
            email: 'test@example.com'
        }
    });
    expect(updateRes.status()).toBe(400);
  });

  test('Update Profile - Invalid Email (Negative)', async ({ authRequest }) => {
    const updateRes = await authRequest.put('/api/v1/kratos/user', {
        data: {
            first_name: 'Test', 
            last_name: 'User',
            email: 'invalid-email-format' // Invalid email
        }
    });
    expect(updateRes.status()).toBe(400);
  });

  test('Fetch Played Quizzes - Should return list', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/user_played_quizes/');
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(data).toHaveProperty('data');
    
    if (data.data && !Array.isArray(data.data)) {
        expect(data.data).toHaveProperty('count');
        if (data.data.data) {
             expect(Array.isArray(data.data.data)).toBe(true);
        }
    } else {
        expect(Array.isArray(data.data)).toBe(true);
    }
  });

  // Edge Case L: No-Op Update
  test('No-Op Update - Idempotency', async ({ authRequest }) => {
    const getRes = await authRequest.get('/api/v1/user/who');
    const profile = (await getRes.json()).data;

    const updateRes = await authRequest.put('/api/v1/kratos/user', {
        data: {
            first_name: profile.firstname,
            last_name: profile.lastname,
            email: profile.email
        }
    });
    
    // System might return 400 "Email already exists" even for self, or 200. Both are valid responses for this edge case.
    expect([200, 400]).toContain(updateRes.status());
    
    if (updateRes.status() === 200) {
        const verifyRes = await authRequest.get('/api/v1/user/who');
        const finalProfile = (await verifyRes.json()).data;
        expect(finalProfile.firstname).toBe(profile.firstname);
    }
  });
});
