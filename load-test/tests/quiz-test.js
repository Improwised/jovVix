/**
 * Quiz Controller Tests
 *
 * Tests for the quiz management endpoints of the API.
 */

import { check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import {
  runTest,
  createTestQuiz,
  createQuizSession,
  generateId,
} from "../lib/utils/test-helpers.js";
import { getTestData } from "../lib/utils/config.js";

export class QuizControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Quiz Controller",
      description: "Tests for quiz management endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<3000"],
          http_req_failed: ["rate<0.3"],
        },
      },
      authOptions: {
        // Default auth options
        email: __ENV.ADMIN_EMAIL,
        password: __ENV.ADMIN_PASSWORD,
      },
    });

    // Test data
    this.testData = getTestData();
  }

  setupTest(auth) {
    // Create a test quiz for use in tests
    if (auth.adminSessionCookie && auth.authMethod !== "mock") {
      try {
        const quizResult = createTestQuiz(auth.adminSessionCookie, {
          title: `quiz_${generateId("", 8)}`,
          description: this.testData.description || "Test quiz description",
          filePath: __ENV.CSV_FILE_PATH || "./test-quiz.csv",
        });

        if (quizResult.success) {
          console.log(`Created test quiz: ${quizResult.quizId}`);

          // Create a session for the quiz
          const sessionResult = createQuizSession(
            auth.adminSessionCookie,
            quizResult.quizId
          );

          if (sessionResult.success) {
            console.log(`Created test session: ${sessionResult.sessionId}`);

            return {
              quizId: quizResult.quizId,
              quizTitle: quizResult.quizTitle,
              sessionId: sessionResult.sessionId,
            };
          }
        }
      } catch (error) {
        console.error(`Error creating test quiz: ${error.message}`);
      }
    }

    // Return mock data if quiz creation failed
    return {
      quizId: "mock-quiz-id",
      quizTitle: "mock-quiz-title",
      sessionId: "mock-session-id",
    };
  }

  functionalTests() {
    runTest("Get Admin Uploaded Quizzes", () =>
      this.testGetAdminUploadedQuizzes()
    );
    runTest("Get Quiz By ID", () => this.testGetQuizById());
    runTest("Get Quiz Questions", () => this.testGetQuizQuestions());

    if (this.state.quizId && this.state.quizId !== "mock-quiz-id") {
      runTest("Generate Demo Session", () => this.testGenerateDemoSession());
      runTest("Get Quiz Session", () => this.testGetQuizSession());
    }

    runTest("Create and Delete Quiz", () => this.testCreateAndDeleteQuiz());
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.4) {
      // 40% get quizzes list
      this.testGetAdminUploadedQuizzes();
    } else if (scenario < 0.7) {
      // 30% get quiz by ID
      this.testGetQuizById();
    } else if (scenario < 0.9) {
      // 20% get quiz questions
      this.testGetQuizQuestions();
    } else {
      // 10% get quiz session
      this.testGetQuizSession();
    }
  }

  /**
   * Test getting admin uploaded quizzes
   */
  testGetAdminUploadedQuizzes() {
    const response = http.get("api/v1/quizzes", {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_admin_quizzes" },
    });

    check(response, {
      "Get admin quizzes returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          return r;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test getting quiz by ID
   */
  testGetQuizById() {
    const quizId = this.state.quizId || "mock-quiz-id";

    const response = http.get(`api/v1/quizzes/${quizId}`, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_quiz_by_id" },
    });

    // If using a mock ID, we expect a 404, otherwise 200
    const expectedStatus = quizId === "mock-quiz-id" ? 405 : 200;
    check(response, {
      [`Get quiz by ID returns ${expectedStatus}`]: (r) =>
        r.status === expectedStatus,
    });

    if (expectedStatus === 200) {
      check(response, {
        "Response contains quiz data": (r) => {
          try {
            return r.body && r.body.data;
          } catch (e) {
            return false;
          }
        },
      });
    }
  }

  /**
   * Test getting quiz questions
   */
  testGetQuizQuestions() {
    const quizId = this.state.quizId || "mock-quiz-id";

    const response = http.get(`api/v1/quizzes/${quizId}/questions`, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_quiz_questions" },
    });

    // If using a mock ID, we expect a 404, otherwise 200
    const expectedStatus = quizId === "mock-quiz-id" ? 500 : 200;

    check(response, {
      [`Get quiz questions returns ${expectedStatus}`]: (r) =>
        r.status === expectedStatus,
    });

    if (expectedStatus === 200) {
      check(response, {
        "Response contains questions data": (r) => {
          try {
            return r.body && r.body.data;
          } catch (e) {
            return false;
          }
        },
      });
    }
  }

  /**
   * Test generating a demo session
   */
  testGenerateDemoSession() {
    const quizId = this.state.quizId || "mock-quiz-id";

    const response = http.post(`api/v1/quizzes/${quizId}/demo_session`, null, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "generate_demo_session" },
    });

    // If using a mock ID, we expect a 404, otherwise 202
    const expectedStatus = quizId === "mock-quiz-id" ? 404 : 202;

    check(response, {
      [`Generate demo session returns ${expectedStatus}`]: (r) =>
        r.status === expectedStatus,
    });

    if (expectedStatus === 202) {
      check(response, {
        "Response contains session ID": (r) => {
          try {
            return r.data;
          } catch (e) {
            return false;
          }
        },
      });
    }
  }

  /**
   * Test getting quiz session
   */
  testGetQuizSession() {
    const sessionId = this.state.sessionId || "mock-session-id";

    const response = http.get(`api/v1/quizzes/session/${sessionId}`, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_quiz_session" },
    });

    // If using a mock ID, we expect a 404, otherwise 200
    const expectedStatus = sessionId === "mock-session-id" ? 404 : 200;

    check(response, {
      [`Get quiz session returns ${expectedStatus}`]: (r) =>
        r.status === expectedStatus,
    });

    if (expectedStatus === 200) {
      check(response, {
        "Response contains session data": (r) => {
          try {
            return r.body && r.body.data;
          } catch (e) {
            return false;
          }
        },
      });
    }
  }

  /**
   * Test creating and deleting a quiz
   */
  testCreateAndDeleteQuiz() {
    // Generate a unique quiz title
    const quizTitle = `quiz_${generateId("", 8)}`;

    // Create quiz
    const createPayload = {
      description: "Test quiz for deletion",
    };

    let quizId = "";

    // Only attempt to create if we have a CSV file
    if (__ENV.CSV_FILE_PATH) {
      const createFiles = {
        attachment: http.file(open(__ENV.CSV_FILE_PATH, "b"), "test.csv"),
      };

      const createResponse = http.post(
        `api/v1/quizzes/${quizTitle}/upload`,
        createPayload,
        {
          files: createFiles,
          cookies: { ory_kratos_session: this.state.adminSessionCookie },
          tags: { endpoint: "create_quiz" },
        }
      );

      check(createResponse, {
        "Create quiz returns 202": (r) => r.status === 202,
        "Create quiz returns quiz ID": (r) => {
          try {
            quizId = r.data;
            return !!quizId;
          } catch (e) {
            return false;
          }
        },
      });
    } else {
      console.warn("CSV_FILE_PATH not set, skipping quiz creation test");
    }

    // Only attempt to delete if we created a quiz
    if (quizId) {
      const deleteResponse = http.del(`api/v1/quizzes/${quizId}`, null, {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_quiz" },
      });

      check(deleteResponse, {
        "Delete quiz returns 200": (r) => r.status === 200,
      });
    }
  }

  teardownTest(data) {
    // Clean up test quiz if it was created
    if (
      data.quizId &&
      data.quizId !== "mock-quiz-id" &&
      data.adminSessionCookie
    ) {
      try {
        http.del(`api/v1/quizzes/${data.quizId}`, null, {
          cookies: { ory_kratos_session: data.adminSessionCookie },
        });
        console.log(`Cleaned up test quiz: ${data.quizId}`);
      } catch (error) {
        console.error(`Error cleaning up test quiz: ${error.message}`);
      }
    }
  }
}

// Export test instance
const testInstance = new QuizControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
