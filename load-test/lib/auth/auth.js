/**
 * Authentication Helper
 *
 * Provides functions for authentication and session management.
 */

import { check } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import * as http from "../http/client.js";
import { getTestData } from "../utils/config.js";

// Test data
const testData = getTestData();

/**
 * Create a guest user
 * @param {string} username - Optional username (generated if not provided)
 * @param {string} avatarName - Optional avatar name
 * @returns {Object} Object containing user info and cookies
 */
export function createGuestUser(username = null, avatarName = null) {
  // Generate username if not provided
  if (!username) {
    username = `guest_${uuidv4().substring(0, 8)}`;
  }

  // Use default avatar if not provided
  if (!avatarName) {
    avatarName = testData.guestAvatarName || "Chase";
  }

  const response = http.post(
    `api/v1/user/${username}?avatar_name=${avatarName}`
  );

  check(response, {
    "Guest user creation successful": (r) => r.status === 200,
    "Guest user cookie received": (r) => r.cookies && r.cookies.user,
  });

  let userCookie = "";
  if (response.cookies && response.cookies.user) {
    userCookie = response.cookies.user[0].value;
  }

  return {
    username,
    avatarName,
    userCookie,
    response,
  };
}

/**
 * Extract session cookie from response
 * @param {Object} response - HTTP response object
 * @param {string} cookieName - Name of the cookie to extract
 * @returns {string} Cookie value or empty string if not found
 */
export function extractCookie(response, cookieName = "ory_kratos_session") {
  if (response.cookies && response.cookies[cookieName]) {
    return response.cookies[cookieName][0].value;
  }
  return "";
}

/**
 * Perform Kratos registration flow
 * @param {Object} userData - User data for registration
 * @returns {Object} Object containing session info and cookies
 */
export function registerKratosUser(userData = {}) {
  // Generate unique email if not provided
  if (!userData.email) {
    userData.email = `testuser_${uuidv4().substring(0, 8)}@example.com`;
  }

  // Set default password if not provided
  if (!userData.password) {
    userData.password = "dockerisnotpodman!";
  }

    if (!userData.firstName) {
    userData.firstName = `user${uuidv4().substring(0, 4)}`;
  }

  if (!userData.lastName) {
    userData.lastName = "Test";
  }

  // 1. Initiate registration flow with proper headers for API access
  let response = http.kratosRequest(
    "/self-service/registration/browser",
    "GET",
    null,
    {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    }
  );

  check(response, {
    "Initiate registration successful": (r) => r.status === 200,
  });

  const responseBody = JSON.parse(response.body);
  if (!response || !responseBody.id) {
    console.error("Failed to get registration flow ID");
    return { success: false };
  }

  // 2. Extract flow data
  const flowId = responseBody.id;

  // 3. Extract CSRF token
  let csrfToken = "";
  if (responseBody.ui && responseBody.ui.nodes) {
    responseBody.ui.nodes.forEach((node) => {
      if (node.attributes && node.attributes.name === "csrf_token") {
        csrfToken = node.attributes.value;
      }
    });
  }

  // 4. Extract CSRF token cookie
  let csrfCookie = "";
  if (response.cookies) {
    Object.keys(response.cookies).forEach((name) => {
      if (name.startsWith("csrf_token_")) {
        csrfCookie = `${name}=${response.cookies[name][0].value}`;
      }
    });
  }

  // 5. Submit registration
  const regPayload = JSON.stringify({
    "traits.email": userData.email,
    password: userData.password,
    "traits.name.first": userData.firstName || "Test",
    "traits.name.last": userData.lastName || "User",
    method: "password",
    csrf_token: csrfToken,
  });

  response = http.kratosRequest(
    `/self-service/registration?flow=${flowId}`,
    "POST",
    regPayload,
    {
      headers: {
        "Content-Type": "application/json",
        Cookie: csrfCookie,
      },
    }
  );

  check(response, {
    "Registration submission successful": (r) =>
      r.status === 200 || r.status === 303,
    "Session cookie received": (r) => r.cookies && r.cookies.ory_kratos_session,
  });

  const sessionCookie = extractCookie(response);

  // Make auth verification request before returning
  const authResponse = http.get("/api/v1/kratos/auth", {
    headers: {
      Cookie: `ory_kratos_session=${sessionCookie}`,
    },
  });

  check(authResponse, {
    "Auth verification successful": (r) => r.status === 200,
  });

  return {
    success: sessionCookie !== "",
    email: userData.email,
    password: userData.password,
    sessionCookie,
    response,
    authResponse,
  };
}

/**
 * Perform Kratos login flow
 * @param {string} email - User email
 * @param {string} password - User password
 * @returns {Object} Object containing session info and cookies
 */
export function loginKratosUser(email, password) {
  // 1. Initiate login flow with proper headers for API access
  let response = http.kratosRequest(
    "/self-service/login/browser",
    "GET",
    null,
    {
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    }
  );

  check(response, {
    "Initiate login successful": (r) => r.status === 200,
  });

  if (!response.id) {
    console.error("Failed to get login flow ID");
    return { success: false };
  }

  // 2. Extract flow data
  const flowId = response.id;

  // 3. Extract CSRF token
  let csrfToken = "";
  if (response.ui && response.ui.nodes) {
    response.ui.nodes.forEach((node) => {
      if (node.attributes && node.attributes.name === "csrf_token") {
        csrfToken = node.attributes.value;
      }
    });
  }

  // 4. Extract CSRF token cookie
  let csrfCookie = "";
  if (response.cookies) {
    Object.keys(response.cookies).forEach((name) => {
      if (name.startsWith("csrf_token_")) {
        csrfCookie = `${name}=${response.cookies[name][0].value}`;
      }
    });
  }

  // 5. Submit login
  const loginPayload = JSON.stringify({
    password_identifier: email,
    password: password,
    method: "password",
    csrf_token: csrfToken,
  });

  response = http.kratosRequest(
    `/self-service/login?flow=${flowId}`,
    "POST",
    loginPayload,
    {
      headers: {
        "Content-Type": "application/json",
        Cookie: csrfCookie,
      },
    }
  );

  check(response, {
    "Login submission successful": (r) => r.status === 200 || r.status === 303,
    "Session cookie received": (r) => r.cookies && r.cookies.ory_kratos_session,
  });

  const sessionCookie = extractCookie(response);

  return {
    success: sessionCookie !== "",
    email,
    sessionCookie,
    response,
  };
}

/**
 * Use existing session token from environment variable
 * @returns {Object} Object containing session info
 */
export function useExistingSession() {
  const sessionToken = __ENV.ADMIN_SESSION_TOKEN || "";

  if (!sessionToken) {
    console.warn("No ADMIN_SESSION_TOKEN environment variable found");
    return { success: false };
  }

  return {
    success: true,
    sessionCookie: sessionToken,
  };
}

/**
 * Setup authentication for tests
 * @param {Object} options - Authentication options
 * @returns {Object} Authentication data for tests
 */
export function setupAuth(options = {}) {
  // Try to use existing session first
  const existingSession = useExistingSession();
  if (existingSession.success) {
    return {
      adminSessionCookie: existingSession.sessionCookie,
      authMethod: "existing_session",
    };
  }

  // If existing session not available, try to login
  if (options.email && options.password) {
    const loginResult = loginKratosUser(options.email, options.password);
    if (loginResult.success) {
      return {
        adminSessionCookie: loginResult.sessionCookie,
        authMethod: "login",
      };
    }
  }

  // If login not available, try to register
  if (options.allowRegistration !== false) {
    const registerResult = registerKratosUser(options);
    if (registerResult.success) {
      return {
        adminSessionCookie: registerResult.sessionCookie,
        email: registerResult.email,
        authMethod: "registration",
      };
    }
  }

  // If all auth methods fail, return mock session for testing
  console.warn("All authentication methods failed, using mock session");
  return {
    adminSessionCookie: "mock_session_for_testing",
    authMethod: "mock",
  };
}
