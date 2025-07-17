/**
 * Analytics Board User Controller Tests
 *
 * Tests for the analytics board user endpoints of the API.
 */

import { check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class AnalyticsBoardUserControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Analytics Board User Controller",
      description: "Tests for analytics board user endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.3"],
        },
      },
    });

    // Test data
    this.testData = {
      invalidUserPlayedQuizId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      validUserPlayedQuizId: "", // Will be set during setup
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
      if (
        sessionResponse &&
        sessionResponse &&
        sessionResponse.body &&
        sessionResponse.body.data
      ) {
        // Simulate played quiz and get user played quiz ID
        const userPlayedQuizId = this.simulatePlayedQuiz(
          sessionResponse.body.data
        );
        this.testData.validUserPlayedQuizId = userPlayedQuizId;
      }
    }

    return {
      validUserPlayedQuizId: this.testData.validUserPlayedQuizId,
    };
  }

  functionalTests() {
    runTest("Missing Query Params", () => this.testMissingQueryParams());
    runTest("Invalid User Played Quiz ID", () =>
      this.testInvalidUserPlayedQuizId()
    );
    runTest("Valid User Analytics Request", () =>
      this.testValidUserAnalyticsRequest()
    );
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.6) {
      // 60% valid user analytics requests
      this.testValidUserAnalyticsRequest();
    } else if (scenario < 0.8) {
      // 20% invalid user played quiz ID requests
      this.testInvalidUserPlayedQuizId();
    } else {
      // 20% missing params requests
      this.testMissingQueryParams();
    }
  }

  /**
   * Test analytics board user with missing query parameters
   */
  testMissingQueryParams() {
    const response = http.get("api/v1/analytics_board/user", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "analytics_board_user_missing_params" },
    });

    check(response, {
      "Missing query params returns 400": (r) => r.status === 400,
    });
  }

  /**
   * Test analytics board user with invalid user played quiz ID
   */
  testInvalidUserPlayedQuizId() {
    const response = http.get(
      `api/v1/analytics_board/user?user_played_quiz=${this.testData.invalidUserPlayedQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "analytics_board_user_invalid_id" },
      }
    );

    check(response, {
      "Invalid user played quiz ID returns 200": (r) => r.status === 200,
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
      "Response has empty data array": (r) => {
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
   * Test analytics board user with valid input
   */
  testValidUserAnalyticsRequest() {
    if (!this.testData.validUserPlayedQuizId) {
      console.warn("No valid user played quiz ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/analytics_board/user?user_played_quiz=${this.testData.validUserPlayedQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "analytics_board_user_valid" },
      }
    );

    check(response, {
      "Valid user analytics request returns 200": (r) => r.status === 200,
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
  simulatePlayedQuiz(sessionId) {
    // This would typically involve creating a guest user and simulating quiz play
    // For now, we'll return a mock user played quiz ID
    console.log("Simulating played quiz for user analytics data");
    return `user-played-${sessionId}`;
  }
}

// Export test instance
const testInstance = new AnalyticsBoardUserControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
