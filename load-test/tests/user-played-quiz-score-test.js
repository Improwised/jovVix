/**
 * User Played Score Quiz Controller Tests
 *
 * Tests for the user played quiz score endpoints of the API.
 */

import { group, check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class UserPlayedScoreQuizControllerTest extends BaseTest {
  constructor() {
    super({
      name: "User Played Score Quiz Controller",
      description: "Tests for user played quiz score endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.5"],
        },
      },
    });

    // Test data
    this.testData = {
      invalidInvitationCode: 62674,
      invalidUserPlayedQuizId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      validInvitationCode: null, // Will be set during setup
      validUserPlayedQuizId: "", // Will be set during setup
      validSessionId: "", // Will be set during setup
      expectedPlayedQuizData: {
        status: "success",
        expectedQuestionCount: 5,
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
        sessionResponse.body &&
        sessionResponse.body.data
      ) {
        this.testData.validSessionId = sessionResponse.body.data;

        // Get invitation code for the session
        const invitationCode = this.getInvitationCode(
          this.testData.validSessionId
        );
        if (invitationCode) {
          this.testData.validInvitationCode = invitationCode;

          // Simulate played quiz validation to get user played quiz ID
          const playedQuizResponse =
            this.simulatePlayedQuizValidation(invitationCode);
          if (
            playedQuizResponse &&
            playedQuizResponse.body &&
            playedQuizResponse.body.data &&
            playedQuizResponse.body.data.user_played_quiz
          ) {
            this.testData.validUserPlayedQuizId =
              playedQuizResponse.body.data.user_played_quiz;
          }
        }
      }
    }

    return {
      validInvitationCode: this.testData.validInvitationCode,
      validUserPlayedQuizId: this.testData.validUserPlayedQuizId,
      validSessionId: this.testData.validSessionId,
    };
  }

  functionalTests() {
    group("Played Quiz Validation Tests", () => {
      runTest("Played Quiz Validation - Invalid Code", () =>
        this.testPlayedQuizValidationInvalidCode()
      );
      runTest("Played Quiz Validation - Host Cannot Play", () =>
        this.testPlayedQuizValidationHostCannotPlay()
      );
      runTest("Played Quiz Validation - Valid Input", () =>
        this.testPlayedQuizValidationValid()
      );
    });

    group("List User Played Quizzes Tests", () => {
      runTest("List User Played Quizzes", () =>
        this.testListUserPlayedQuizzes()
      );
    });

    group("List User Played Quizzes With Questions Tests", () => {
      runTest("List With Questions - Invalid ID", () =>
        this.testListUserPlayedQuizzesWithQuestionsInvalidId()
      );
      runTest("List With Questions - Valid ID", () =>
        this.testListUserPlayedQuizzesWithQuestionsValid()
      );
    });
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.4) {
      // 40% played quiz validation requests
      this.testPlayedQuizValidationValid();
    } else if (scenario < 0.6) {
      // 20% list user played quizzes requests
      this.testListUserPlayedQuizzes();
    } else if (scenario < 0.8) {
      // 20% list with questions requests
      this.testListUserPlayedQuizzesWithQuestionsValid();
    } else {
      // 20% invalid requests (mixed)
      const invalidScenario = Math.random();
      if (invalidScenario < 0.5) {
        this.testPlayedQuizValidationInvalidCode();
      } else {
        this.testListUserPlayedQuizzesWithQuestionsInvalidId();
      }
    }
  }

  /**
   * Test played quiz validation with invalid invitation code
   */
  testPlayedQuizValidationInvalidCode() {
    const response = http.post(
      `api/v1/user_played_quizes/${this.testData.invalidInvitationCode}`,
      null,
      {    
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "played_quiz_validation_invalid_code" },
      }
    );

    console.log(response.status)

    check(response, {
      "Played quiz validation with invalid code returns 400": (r) =>
        r.status === 400,
    });
  }

  /**
   * Test played quiz validation where host tries to play their own quiz
   */
  testPlayedQuizValidationHostCannotPlay() {
    if (!this.testData.validInvitationCode) {
      console.warn("No valid invitation code available for testing");
      return;
    }

    const response = http.post(
      `api/v1/user_played_quizes/${this.testData.validInvitationCode}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie }, // Admin tries to play their own quiz
        tags: { endpoint: "played_quiz_validation_host_cannot_play" },
      }
    );

    check(response, {
      "Host cannot play their own quiz returns 500": (r) => r.status === 500,
    });
  }

  /**
   * Test played quiz validation with valid input
   */
  testPlayedQuizValidationValid() {
    if (!this.testData.validInvitationCode) {
      console.warn("No valid invitation code available for testing");
      return;
    }

    const response = http.post(
      `api/v1/user_played_quizes/${this.testData.validInvitationCode}`,
      null,
      {
        // Use user client (no admin session) to simulate guest user
        tags: { endpoint: "played_quiz_validation_valid" },
      }
    );

    check(response, {
      "Played quiz validation with valid input returns 200": (r) =>
        r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains session ID": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            data.body &&
            data.body.data &&
            data.body.data.session_id === this.testData.validSessionId
          );
        } catch (e) {
          return false;
        }
      },
      "Response contains user played quiz ID": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.body && data.body.data && data.body.data.user_played_quiz;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test listing user played quizzes
   */
  testListUserPlayedQuizzes() {
    const response = http.get("api/v1/user_played_quizes", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "list_user_played_quizzes" },
    });

    check(response, {
      "List user played quizzes returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains count data": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            data &&
            data.data &&
            typeof data.data.count === "number"
          );
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test listing user played quizzes with questions by invalid ID
   */
  testListUserPlayedQuizzesWithQuestionsInvalidId() {
    const response = http.get(
      `api/v1/user_played_quizes/${this.testData.invalidUserPlayedQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: {
          endpoint: "list_user_played_quizzes_with_questions_invalid_id",
        },
      }
    );

    check(response, {
      "List with questions (invalid ID) returns 200": (r) => r.status === 200,
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
   * Test listing user played quizzes with questions by valid ID
   */
  testListUserPlayedQuizzesWithQuestionsValid() {
    if (!this.testData.validUserPlayedQuizId) {
      console.warn("No valid user played quiz ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/user_played_quizes/${this.testData.validUserPlayedQuizId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "list_user_played_quizzes_with_questions_valid" },
      }
    );

    check(response, {
      "List with questions (valid ID) returns 200": (r) => r.status === 200,
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
              this.testData.expectedPlayedQuizData.expectedQuestionCount
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
   * Helper method to get invitation code for a session
   */
  getInvitationCode(sessionId) {
    // In a real implementation, this would query the database or API
    // For now, we'll return a mock invitation code
    return Math.floor(Math.random() * 900000) + 100000; // 6-digit code
  }

  /**
   * Helper method to simulate played quiz validation
   */
  simulatePlayedQuizValidation(invitationCode) {
    const response = http.post(
      `api/v1/user_played_quizes/${invitationCode}`,
      null,
      {
        // Use user client (no admin session) to simulate guest user
        tags: { endpoint: "simulate_played_quiz_validation" },
      }
    );

    return response;
  }
}

// Export test instance
const testInstance = new UserPlayedScoreQuizControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
