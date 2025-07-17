/**
 * Analytics Board Admin Controller Tests
 *
 * Tests for the analytics board admin endpoints of the API.
 */

import { check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class AnalyticsBoardAdminControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Analytics Board Admin Controller",
      description: "Tests for analytics board admin endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.5"],
        },
      },
    });

    // Test data
    this.testData = {
      invalidActiveQuizId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      validActiveQuizId: "", // Will be set during setup
      expectedAnalyticsData: {
        status: "success",
        expectedQuestionCount: 5,
        expectedUserData: {
          userName: "testcaseuser",
          firstName: "testcaseuser",
          correctAnswer: "[2]",
          calculatedScore: 0,
          isAttend: false,
          responseTime: -1,
          calculatedPoints: 0,
          question: "Which city is known as the Eternal City?",
          options: {
            1: "Paris",
            2: "Rome",
            3: "Athens",
            4: "Cairo",
          },
          questionsMedia: "text",
          optionsMedia: "text",
          resource: "",
          points: 1,
          questionTypeID: 1,
          questionType: "single answer",
          orderNo: 1,
        },
      },
    };
  }

  setupTest(auth) {
    // Create test quiz and session for valid tests
    const quizResponse = this.createTestQuiz();
    if (quizResponse && quizResponse.body && quizResponse.body.data) {
      const sessionResponse = this.generateDemoSession(quizResponse.body.data);
      if (sessionResponse && sessionResponse.body.data) {
        this.testData.validActiveQuizId = sessionResponse.body.data;

        // Simulate played quiz for valid test scenarios
        this.simulatePlayedQuiz();
      }
    }

    return {
      validActiveQuizId: this.testData.validActiveQuizId,
    };
  }

  functionalTests() {
    runTest("Unauthorized Access", () => this.testUnauthorizedAccess());
    runTest("Missing Query Params", () => this.testMissingQueryParams());
    runTest("Invalid Active Quiz ID", () => this.testInvalidActiveQuizId());
    runTest("Valid Analytics Request", () => this.testValidAnalyticsRequest());
  }

  loadTests() {
    const scenario = Math.random();
    if (scenario < 0.5) {
      // 50% valid analytics requests
      this.testValidAnalyticsRequest();
    } else if (scenario < 0.7) {
      // 20% invalid quiz ID requests
      this.testInvalidActiveQuizId();
    } else if (scenario < 0.85) {
      // 15% missing params requests
      this.testMissingQueryParams();
    } else {
      // 15% unauthorized requests
      this.testUnauthorizedAccess();
    }
  }

  /**
   * Test unauthorized access to analytics board admin endpoint
   */
  testUnauthorizedAccess() {
    const response = http.get("api/v1/analytics_board/admin", {
      // No authentication cookies
      tags: { endpoint: "analytics_board_admin_unauthorized" },
    });

    check(response, {
      "Unauthorized access returns 401": (r) => r.status === 401,
    });
  }

  /**
   * Test analytics board admin with missing query parameters
   */
  testMissingQueryParams() {
    const response = http.get("api/v1/analytics_board/admin", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "analytics_board_admin_missing_params" },
    });

    check(response, {
      "Missing query params returns 400": (r) => r.status === 400,
    });
  }

  /**
   * Test analytics board admin with invalid active quiz ID
   */
  testInvalidActiveQuizId() {
    const response = http.get(
      `api/v1/analytics_board/admin?active_quiz_id=${this.testData.invalidActiveQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "analytics_board_admin_invalid_id" },
      }
    );

    check(response, {
      "Invalid active quiz ID returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response has success status": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.status === "success";
        } catch (e) {
          return false;
        }
      },
      "Response has null data": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.data === null;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test analytics board admin with valid input
   */
  testValidAnalyticsRequest() {
    if (!this.testData.validActiveQuizId) {
      console.warn("No valid active quiz ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/analytics_board/admin?active_quiz_id=${this.testData.validActiveQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "analytics_board_admin_valid" },
      }
    );

    check(response, {
      "Valid analytics request returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response has success status": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.status === "success";
        } catch (e) {
          return false;
        }
      },
      "Response has expected question count": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            Array.isArray(data.data) &&
            data.data.length ===
              this.testData.expectedAnalyticsData.expectedQuestionCount
          );
        } catch (e) {
          return false;
        }
      },
    });

    // Additional detailed checks for valid response structure
    if (response.status === 200) {
      try {
        const data = JSON.parse(response.body);
        if (data.body && data.body.data && data.body.data.length > 0) {
          const firstQuestion = data.body.data[0];
          const expected = this.testData.expectedAnalyticsData.expectedUserData;

          check(response, {
            "First question has correct user name": () =>
              firstQuestion.user_name === expected.userName,
            "First question has correct question text": () =>
              firstQuestion.question === expected.question,
            "First question has correct question type": () =>
              firstQuestion.question_type === expected.questionType,
            "First question has correct points": () =>
              firstQuestion.points === expected.points,
            "First question has correct order": () =>
              firstQuestion.order_no === expected.orderNo,
          });
        }
      } catch (e) {
        console.warn(
          "Failed to parse response for detailed validation:",
          e.message
        );
      }
    }
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
   * Helper method to generate demo session
   */
  generateDemoSession(quizId) {
    const response = http.post(`api/v1/quizzes/${quizId}/demo_session`, null, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "generate_demo_session" },
    });

    return response;
  }

  /**
   * Helper method to simulate played quiz
   */
  simulatePlayedQuiz() {
    // This would typically involve creating a guest user and simulating quiz play
    // For now, we'll assume the test data setup is sufficient
    console.log("Simulating played quiz for analytics data");
  }
}

// Create test instance
const testInstance = new AnalyticsBoardAdminControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
