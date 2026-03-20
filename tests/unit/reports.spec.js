import { test, expect } from '../fixtures';

test.describe('Reports (Unit)', () => {

  test('List Reports - Should return paginated reports', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/admin/reports/list?page=1&limit=10');
    expect(response.status()).toBe(200);
    
    const data = await response.json();

    expect(data).toHaveProperty('data');
    if (data.data) {
        expect(data.data).toHaveProperty('count');
        if (data.data.count > 0) {
            expect(Array.isArray(data.data.data)).toBe(true);
        } else {
             expect(data.data.data === null || Array.isArray(data.data.data)).toBe(true);
        }
    }
  });

  test('Filter Reports - By Date (Valid)', async ({ authRequest }) => {
    const date = new Date().toISOString().split('T')[0];
    const response = await authRequest.get(`/api/v1/admin/reports/list?date=${date}`);
    expect(response.status()).toBe(200);
  });

  test('Filter Reports - By Date (Invalid Format)', async ({ authRequest }) => {
    // Invalid date format should trigger DB error (500) based on backend logic
    const response = await authRequest.get('/api/v1/admin/reports/list?date=invalid-date');
    expect(response.status()).toBe(500); 
  });

  test('Search Reports - By Quiz Name', async ({ authRequest }) => {
    const searchParams = new URLSearchParams({
        name: 'NonExistentQuizName_Likely'
    });
    
    const response = await authRequest.get(`/api/v1/admin/reports/list?${searchParams.toString()}`);
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    if (data.data) {
        expect(data.data.count).toBe(0);
    }
  });

  // Edge Case J: Pagination Overflow
  test('Pagination Overflow - Request Page 9999', async ({ authRequest }) => {
    const response = await authRequest.get('/api/v1/admin/reports/list?page=9999&limit=10');
    expect(response.status()).toBe(200);
    const data = await response.json();
    
    if (data.data && data.data.data) {
        expect(data.data.data.length).toBe(0);
    }
  });

  // Edge Case K: Combined Filters
  test('Combined Filters - Date + Search', async ({ authRequest }) => {
    const params = new URLSearchParams({
        date: new Date().toISOString().split('T')[0],
        name: 'Quiz',
        page: '1',
        limit: '10'
    });
    const response = await authRequest.get(`/api/v1/admin/reports/list?${params.toString()}`);
    expect(response.status()).toBe(200);
  });
});
