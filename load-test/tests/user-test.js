/**
 * User Controller Tests
 *
 * Tests for the user management endpoints of the API.
 */

import { group, check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class UserControllerTest extends BaseTest {
  constructor() {
    super({
      name: "User Controller",
      description: "Tests for user management endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.3"],
        },
      },
    });

    // Test data
    this.testData = {
      testUsername: "testuser",
      testAvatarName: "Chase",
      createdUsers: [], // Track created users for cleanup
    };
  }

  functionalTests() {
    group("Create Guest User Tests", () => {
      runTest("Create User - Missing Avatar", () =>
        this.testCreateUserMissingAvatar()
      );
      runTest("Create User - Valid Input", () => this.testCreateUserValid());
    });

    group("Get User Meta Tests", () => {
      runTest("Get User Meta", () => this.testGetUserMeta());
    });
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.4) {
      // 40% create user requests
      this.testCreateUserValid();
    } else if (scenario < 0.7) {
      // 30% get user meta requests
      this.testGetUserMeta();
    } else {
      // 30% invalid create user requests
      this.testCreateUserMissingAvatar();
    }
  }

  /**
   * Test creating user without avatar parameter
   */
  testCreateUserMissingAvatar() {
    const username = `user_${Date.now().toString().slice(-6)}`;

    const response = http.post(`api/v1/user/${username}`, null, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "create_user_missing_avatar" },
    });

    check(response, {
      "Create user without avatar returns 400": (r) => r.status === 400,
    });
  }

  /**
   * Test creating user with valid input
   */
  testCreateUserValid() {
    const username = `user_${Date.now().toString().slice(-6)}`;

    const response = http.post(
      `api/v1/user/${username}?avatar_name=${this.testData.testAvatarName}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "create_user_valid" },
      }
    );

    console.log("response: ", response.body);

    check(response, {
      "Create user with valid input returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
    });

    // Track created user for potential cleanup
    if (response.status === 200) {
      this.testData.createdUsers.push(username);
    }
  }

  /**
   * Test getting user meta information
   */
  testGetUserMeta() {
    const response = http.get("api/v1/user/who", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_user_meta" },
    });

    check(response, {
      "Get user meta returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains user data": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data && data.data;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Cleanup method to remove created test users
   */
  teardownTest(data) {
    // Note: In a real implementation, you might want to clean up created users
    // However, this would require database access or a delete user endpoint
    // For now, we'll just log the cleanup need
    if (this.testData.createdUsers.length > 0) {
      console.log(
        `Test created ${this.testData.createdUsers.length} users that may need cleanup:`,
        this.testData.createdUsers
      );
    }
  }
}

// Export test instance
const testInstance = new UserControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
