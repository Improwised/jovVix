import { test as setup, request } from '@playwright/test';
import fs from 'fs';
import path from 'path';

// Helper function for service readiness
async function waitForService(request, url, timeout = 60000) {
  const start = Date.now();
  while (Date.now() - start < timeout) {
    try {
      const res = await request.get(url, { headers: { 'Accept': 'application/json' } });
      // console.log(`Checking ${url}, status: ${res.status()}`);
      if (res.ok() || res.status() === 404 || res.status() === 400) return; // Accept 200, 404, or 400 as service responding
    } catch (e) {
      // console.log(`Checking ${url}, error: ${e.message}`);
      // Ignore connection errors
    }
    await new Promise(r => setTimeout(r, 500));
  }
  throw new Error(`Service at ${url} did not become ready in ${timeout}ms`);
}

setup('authenticate', async () => {
  setup.setTimeout(120000);

  const kratosUrl = process.env.KRATOS_PUBLIC_URL || 'http://127.0.0.1:4433';
  const apiUrl = process.env.BASE_URL || 'http://127.0.0.1:3000';
  // frontendUrl not needed for API flow

  // console.log(`Global Setup: Using Kratos at ${kratosUrl}, API at ${apiUrl}`);

  // console.log('Global Setup: Starting Kratos Login (Reference Implementation)...');

  // Use a new request context (no browser overhead)
  const apiContext = await request.newContext({
    baseURL: kratosUrl
  });

  // Wait for Kratos readiness
  // console.log('Global Setup: Waiting for Kratos readiness...');
  await waitForService(apiContext, '/self-service/login/browser'); // Wait for root endpoint

  try {
    // 1. Initialize Login Flow
    const initResponse = await apiContext.get('/self-service/login/browser?refresh=true', {
      headers: { 'Accept': 'application/json' }
    });
    
    if (!initResponse.ok()) {
      const text = await initResponse.text();
      // console.log('Init response text:', text);
      throw new Error(`Failed to initialize login flow: ${initResponse.status()}`);
    }

    const flowData = await initResponse.json();
    const actionUrl = flowData.ui.action; // The URL to POST credentials to
    
    // Extract CSRF Token from UI nodes
    const csrfNode = flowData.ui.nodes.find(n => n.attributes.name === 'csrf_token');
    if (!csrfNode) {
      throw new Error('CSRF token not found in Kratos response');
    }
    const csrfToken = csrfNode.attributes.value;
    
    // Extract CSRF Cookie explicitly (like auth.js)
    let state = await apiContext.storageState();
    const csrfCookie = state.cookies.find(c => c.name.startsWith('csrf_token_'));
    // console.log('CSRF Cookie found:', csrfCookie ? csrfCookie.name : 'NONE');

    // console.log(`Global Setup: Flow initialized. Submitting credentials via JSON...`);

    // 2. Submit Login
    const loginResponse = await apiContext.post(actionUrl, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken,
        'Cookie': `${csrfCookie.name}=${csrfCookie.value}`
      },
      data: {
        password_identifier: process.env.TEST_USERNAME || 'xasoj88044@daikoa.com',
        password: process.env.TEST_PASSWORD || 'p;iuli3#%c(',
        method: 'password',
        csrf_token: csrfToken
      }
    });

    if (!loginResponse.ok()) {
      const errorText = await loginResponse.text();
      console.error('Login Failed Response:', errorText);
      throw new Error(`Login submission failed: ${loginResponse.status()}`);
    }

    const loginResult = await loginResponse.json();
    // console.log('Login successful.');

    // 3. Capture Cookies
    state = await apiContext.storageState();
    
    // Verify we captured the session cookie
    const sessionCookie = state.cookies.find(c => c.name.includes('ory_kratos_session'));
    if (!sessionCookie) {
         console.warn('Warning: ory_kratos_session cookie not found in storage state.');
         if (loginResult.session_token) {
             // console.log('Using session_token from response body as cookie.');
             state.cookies.push({
                 name: 'ory_kratos_session',
                 value: loginResult.session_token,
                 domain: '127.0.0.1',
                 path: '/',
                 expires: -1,
                 httpOnly: true,
                 secure: false,
                 sameSite: 'Lax'
             });
         }
    } else {
        // console.log('Session cookie captured successfully.');
    }

    // 4. Patch Cookies for HTTP Testing
    state.cookies = state.cookies.map(c => {
        c.secure = false;
        c.sameSite = 'Lax';
        return c;
    });

    // 5. Save to file
    const authDir = 'tests/.auth';
    if (!fs.existsSync(authDir)){
        fs.mkdirSync(authDir, { recursive: true });
    }
    
    fs.writeFileSync(path.join(authDir, 'admin.json'), JSON.stringify(state, null, 2));
    // console.log('Global Setup: Auth state saved to tests/.auth/admin.json');

    // 6. Now create test data: quiz, session, code
    // console.log('Global Setup: Creating test data (quiz, session, code)...');

    // Wait for API readiness
    // console.log('Global Setup: Waiting for API readiness...');
    const tempApiContext = await request.newContext();
    await waitForService(tempApiContext, '/');
    await tempApiContext.dispose();

    // Create a new request context for API calls
    const adminApiContext = await request.newContext({
      storageState: path.join(authDir, 'admin.json'),
      baseURL: apiUrl
    });

    // Create Quiz
    const quizTitle = `TestQuiz_${Date.now()}`;
    const csvContent = `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"What is 2+2?",single answer,10,1,2,3,4,,4,text,text,
"Is the sky blue?",single answer,5,Yes,No,,,,1,text,text,`;

    const buffer = Buffer.from(csvContent);

    const response = await adminApiContext.post(`/api/v1/quizzes/${quizTitle}/upload`, {
      multipart: {
        description: 'Test Quiz for Automated Tests',
        attachment: {
          name: `test_${Date.now()}.csv`,
          mimeType: 'text/csv',
          buffer: buffer
        }
      }
    });

    if (!response.ok()) {
      throw new Error(`Failed to create test quiz: ${response.status()}`);
    }

    const quizIdRes = await response.json();
    const quizId = quizIdRes.data || quizIdRes;

    // Create Demo Session
    const sessionRes = await adminApiContext.post(`/api/v1/quizzes/${quizId}/demo_session`);
    if (!sessionRes.ok()) throw new Error('Failed to create session');
    const sessionData = await sessionRes.json();
    const sessionId = sessionData.data;

    // Activate Session (navigate to page)
    const { chromium } = await import('playwright');
    const browser = await chromium.launch();
    const page = await browser.newPage({
      storageState: path.join(authDir, 'admin.json')
    });

    const frontendUrl = process.env.FRONTEND_URL || 'http://127.0.0.1:5000';
    await page.goto(`${frontendUrl}/admin/arrange/${sessionId}`);

    // Wait for activation (minimal)
    await page.waitForTimeout(500);

    // Query Code from DB
    const { execSync } = await import('child_process');
    const query = `SELECT invitation_code FROM active_quizzes WHERE id = '${sessionId}' AND is_active = true`;
    const output = execSync(`docker exec jovvix-db psql -U jovvix -d jovvix -t -c "${query}"`).toString().trim();

    if (!output || output.includes('row')) {
      throw new Error('Failed to retrieve invitation code from DB');
    }

    const code = output;

    // Save Test Data
    const testDataDir = 'tests/.test-data';
    if (!fs.existsSync(testDataDir)){
        fs.mkdirSync(testDataDir, { recursive: true });
    }
    const testData = { quizId, quizTitle, sessionId, code };
    fs.writeFileSync(path.join(testDataDir, 'shared.json'), JSON.stringify(testData, null, 2));
    // console.log('Global Setup: Test data saved to tests/.test-data/shared.json');

    await browser.close();
    await adminApiContext.dispose();

  } catch (error) {
    console.error('Global Setup: Fatal Error during login or test data creation', error);
    throw error;
  }
});
