/**
 * Shared Quizzes Controller Tests
 *
 * Tests for the quiz sharing endpoints of the API.
 */

import { group, check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class SharedQuizzesControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Shared Quizzes Controller",
      description: "Tests for quiz sharing endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.6"],
        },
      },
    });

    // Test data
    this.testData = {
      invalidQuizId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      validQuizId: "", // Will be set during setup
      shareQuizId: "", // Will be set when sharing a quiz
      shareQuizPayload: {
        email: "demosharequizuser@example.com",
        permission: "read",
      },
      shareQuizPayload2: {
        email: "demosharequizuser2@example.com",
        permission: "share",
      },
      updatePermissionPayload: {
        email: "demosharequizuser@example.com",
        permission: "share",
      },
    };
  }

  setupTest(auth) {
    // Create test quiz for sharing tests
    const quizResponse = this.createTestQuiz();
    if (quizResponse && quizResponse.body && quizResponse.body.data) {
      this.testData.validQuizId = quizResponse.body.data;
    }

    return {
      validQuizId: this.testData.validQuizId,
    };
  }

  functionalTests() {
    group("Share Quiz Tests", () => {
      runTest("Share Quiz - Invalid Quiz ID", () =>
        this.testShareQuizInvalidId()
      );
      runTest("Share Quiz - Valid Input", () => this.testShareQuizValid());
    });

    group("List Quiz Authorized Users Tests", () => {
      runTest("List Authorized Users - Invalid Quiz ID", () =>
        this.testListAuthorizedUsersInvalidId()
      );
      runTest("List Authorized Users - Valid Input", () =>
        this.testListAuthorizedUsersValid()
      );
    });

    group("Update User Permission Tests", () => {
      runTest("Update Permission - Invalid Quiz ID", () =>
        this.testUpdatePermissionInvalidId()
      );
      runTest("Update Permission - Valid Input", () =>
        this.testUpdatePermissionValid()
      );
    });

    group("Delete User Permission Tests", () => {
      runTest("Delete Permission - Invalid Quiz ID", () =>
        this.testDeletePermissionInvalidId()
      );
      runTest("Delete Permission - Valid Input", () =>
        this.testDeletePermissionValid()
      );
    });

    group("List Shared Quizzes Tests", () => {
      runTest("List Shared Quizzes - Missing Type", () =>
        this.testListSharedQuizzesMissingType()
      );
      runTest("List Shared Quizzes - Shared With Me", () =>
        this.testListSharedQuizzesSharedWithMe()
      );
      runTest("List Shared Quizzes - Shared By Me", () =>
        this.testListSharedQuizzesSharedByMe()
      );
    });
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.3) {
      // 30% share quiz requests
      this.testShareQuizValid();
    } else if (scenario < 0.5) {
      // 20% list authorized users requests
      this.testListAuthorizedUsersValid();
    } else if (scenario < 0.7) {
      // 20% list shared quizzes requests
      this.testListSharedQuizzesSharedByMe();
    } else if (scenario < 0.85) {
      // 15% update permission requests
      this.testUpdatePermissionValid();
    } else {
      // 15% invalid requests (mixed)
      const invalidScenario = Math.random();
      if (invalidScenario < 0.5) {
        this.testShareQuizInvalidId();
      } else {
        this.testListAuthorizedUsersInvalidId();
      }
    }
  }

  /**
   * Test sharing quiz with invalid quiz ID
   */
  testShareQuizInvalidId() {
    const response = http.post(
      `api/v1/shared_quizzes/${this.testData.invalidQuizId}`,
      JSON.stringify(this.testData.shareQuizPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "share_quiz_invalid_id" },
      }
    );

    check(response, {
      "Share quiz with invalid ID returns 401": (r) => r.status === 401,
    });
  }

  /**
   * Test sharing quiz with valid input
   */
  testShareQuizValid() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    const response = http.post(
      `api/v1/shared_quizzes/${this.testData.validQuizId}`,
      JSON.stringify(this.testData.shareQuizPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "share_quiz_valid" },
      }
    );

    check(response, {
      "Share quiz with valid input returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains share quiz ID": (r) => {
        try {
          const data = JSON.parse(r.body);
          if (data.body && data.body.data) {
            this.testData.shareQuizId = data.body.data;
            return true;
          }
          return false;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test listing authorized users with invalid quiz ID
   */
  testListAuthorizedUsersInvalidId() {
    const response = http.get(
      `api/v1/shared_quizzes/${this.testData.invalidQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "list_authorized_users_invalid_id" },
      }
    );

    check(response, {
      "List authorized users with invalid ID returns 401": (r) =>
        r.status === 401,
    });
  }

  /**
   * Test listing authorized users with valid input
   */
  testListAuthorizedUsersValid() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    // First share the quiz to create authorized users
    this.shareQuizHelper(this.testData.shareQuizPayload);
    this.shareQuizHelper(this.testData.shareQuizPayload2);

    const response = http.get(
      `api/v1/shared_quizzes/${this.testData.validQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "list_authorized_users_valid" },
      }
    );

    check(response, {
      "List authorized users with valid input returns 200": (r) =>
        r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains authorized users": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            data.body &&
            data.body.data &&
            Array.isArray(data.body.data) &&
            data.body.data.length >= 2
          );
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test updating permission with invalid quiz ID
   */
  testUpdatePermissionInvalidId() {
    const response = http.put(
      `api/v1/shared_quizzes/${this.testData.invalidQuizId}?shared_quiz_id=test-id`,
      JSON.stringify(this.testData.updatePermissionPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "update_permission_invalid_id" },
      }
    );

    check(response, {
      "Update permission with invalid ID returns 401": (r) => r.status === 401,
    });
  }

  /**
   * Test updating permission with valid input
   */
  testUpdatePermissionValid() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    // First share the quiz to get a share quiz ID
    const shareResponse = this.shareQuizHelper(this.testData.shareQuizPayload);
    if (!shareResponse || !this.testData.shareQuizId) {
      console.warn("Failed to share quiz for update permission test");
      return;
    }

    const response = http.put(
      `api/v1/shared_quizzes/${this.testData.validQuizId}?shared_quiz_id=${this.testData.shareQuizId}`,
      JSON.stringify(this.testData.updatePermissionPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "update_permission_valid" },
      }
    );

    check(response, {
      "Update permission with valid input returns 200": (r) => r.status === 200,
    });
  }

  /**
   * Test deleting permission with invalid quiz ID
   */
  testDeletePermissionInvalidId() {
    const response = http.del(
      `api/v1/shared_quizzes/${this.testData.invalidQuizId}?shared_quiz_id=test-id`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_permission_invalid_id" },
      }
    );

    check(response, {
      "Delete permission with invalid ID returns 401": (r) => r.status === 401,
    });
  }

  /**
   * Test deleting permission with valid input
   */
  testDeletePermissionValid() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    // First share the quiz to get a share quiz ID
    const shareResponse = this.shareQuizHelper(this.testData.shareQuizPayload);
    if (!shareResponse || !this.testData.shareQuizId) {
      console.warn("Failed to share quiz for delete permission test");
      return;
    }

    const response = http.del(
      `api/v1/shared_quizzes/${this.testData.validQuizId}?shared_quiz_id=${this.testData.shareQuizId}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_permission_valid" },
      }
    );

    check(response, {
      "Delete permission with valid input returns 200": (r) => r.status === 200,
    });
  }

  /**
   * Test listing shared quizzes with missing type parameter
   */
  testListSharedQuizzesMissingType() {
    const response = http.get("api/v1/shared_quizzes", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "list_shared_quizzes_missing_type" },
    });

    check(response, {
      "List shared quizzes with missing type returns 400": (r) =>
        r.status === 400,
    });
  }

  /**
   * Test listing shared quizzes with type=shared_with_me
   */
  testListSharedQuizzesSharedWithMe() {
    const response = http.get("api/v1/shared_quizzes?type=shared_with_me", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "list_shared_quizzes_shared_with_me" },
    });

    check(response, {
      "List shared quizzes (shared with me) returns 200": (r) =>
        r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains data array": (r) => {
        try {
          const data = JSON.parse(r.body);
          return Array.isArray(data.data) && data.data.length === 0;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test listing shared quizzes with type=shared_by_me
   */
  testListSharedQuizzesSharedByMe() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    // First share some quizzes to create data
    this.shareQuizHelper(this.testData.shareQuizPayload);
    this.shareQuizHelper(this.testData.shareQuizPayload2);

    const response = http.get("api/v1/shared_quizzes?type=shared_by_me", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "list_shared_quizzes_shared_by_me" },
    });

    check(response, {
      "List shared quizzes (shared by me) returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains shared quizzes": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            data.body &&
            data.body.data &&
            Array.isArray(data.body.data) &&
            data.body.data.length >= 2
          );
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Helper method to create a test quiz
   */
  createTestQuiz() {
    const formData = {
      title: `Test Quiz ${Date.now()}`,
      description: "This Quiz is created for test cases",
    };

    const response = http.post("api/v1/quizzes", formData, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "create_test_quiz" },
    });

    return response;
  }

  /**
   * Helper method to share a quiz
   */
  shareQuizHelper(sharePayload) {
    if (!this.testData.validQuizId) {
      return null;
    }

    const response = http.post(
      `api/v1/shared_quizzes/${this.testData.validQuizId}`,
      JSON.stringify(sharePayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "share_quiz_helper" },
      }
    );

    // Extract share quiz ID from response
    if (response.status === 200) {
      try {
        const data = JSON.parse(response.body);
        if (data.body && data.body.data) {
          this.testData.shareQuizId = data.body.data;
        }
      } catch (e) {
        console.warn("Failed to parse share quiz response:", e.message);
      }
    }

    return response;
  }
}

// Export test instance
const testInstance = new SharedQuizzesControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
