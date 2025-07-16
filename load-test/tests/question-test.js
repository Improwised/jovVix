/**
 * Question Controller Tests
 *
 * Tests for the question management endpoints of the API.
 */

import { group, check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class QuestionControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Question Controller",
      description: "Tests for question management endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<2000"],
          http_req_failed: ["rate<0.7"],
        },
      },
    });

    // Test data
    this.testData = {
      invalidQuizId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      invalidQuestionId: "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a",
      validQuizId: "", // Will be set during setup
      validQuestionId: "", // Will be set during setup
      expectedQuestionData: {
        expectedQuestionCount: 5,
        expectedFirstQuestion: {
          correctAnswer: "[2]",
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
          durationInSeconds: 60,
        },
      },
      updateQuestionPayload: {
        question: "What is the capital of France?",
        type: 1,
        options: {
          1: "Paris",
          2: "London",
          3: "Berlin",
          4: "Madrid",
        },
        answers: [1],
        points: 10,
        duration_in_seconds: 30,
        question_media: "text",
        options_media: "text",
        resource: "",
      },
    };
  }

  setupTest(auth) {
    // Create test quiz for question tests
    const quizResponse = this.createTestQuiz();
    if (quizResponse && quizResponse.body && quizResponse.body.data) {
      this.testData.validQuizId = quizResponse.body.data;

      // Get questions for the quiz to get a valid question ID
      const questionsResponse = this.listQuestions(this.testData.validQuizId);
      if (
        questionsResponse &&
        questionsResponse.body &&
        questionsResponse.body.data &&
        questionsResponse.body.data.data &&
        questionsResponse.body.data.data.length > 0
      ) {
        this.testData.validQuestionId =
          questionsResponse.body.data.data[0].question_id;
      }
    }

    return {
      validQuizId: this.testData.validQuizId,
      validQuestionId: this.testData.validQuestionId,
    };
  }

  functionalTests() {
    group("List Questions Tests", () => {
      runTest("List Questions - Invalid Quiz ID", () =>
        this.testListQuestionsInvalidQuizId()
      );
      runTest("List Questions - Valid Quiz ID", () =>
        this.testListQuestionsValidQuizId()
      );
    });

    group("Get Question By ID Tests", () => {
      runTest("Get Question - Invalid Quiz ID", () =>
        this.testGetQuestionInvalidQuizId()
      );
      runTest("Get Question - Invalid Question ID", () =>
        this.testGetQuestionInvalidQuestionId()
      );
      runTest("Get Question - Valid IDs", () => this.testGetQuestionValidIds());
    });

    group("Update Question Tests", () => {
      runTest("Update Question - Invalid Quiz ID", () =>
        this.testUpdateQuestionInvalidQuizId()
      );
      runTest("Update Question - Missing Body", () =>
        this.testUpdateQuestionMissingBody()
      );
      runTest("Update Question - Valid Input", () =>
        this.testUpdateQuestionValid()
      );
    });

    group("Delete Question Tests", () => {
      runTest("Delete Question - Invalid Quiz ID", () =>
        this.testDeleteQuestionInvalidQuizId()
      );
      runTest("Delete Question - Active Quiz Present", () =>
        this.testDeleteQuestionActiveQuiz()
      );
      runTest("Delete Question - Valid Input", () =>
        this.testDeleteQuestionValid()
      );
    });
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.4) {
      // 40% list questions requests
      this.testListQuestionsValidQuizId();
    } else if (scenario < 0.6) {
      // 20% get question by ID requests
      this.testGetQuestionValidIds();
    } else if (scenario < 0.8) {
      // 20% update question requests
      this.testUpdateQuestionValid();
    } else {
      // 20% invalid requests (mixed)
      const invalidScenario = Math.random();
      if (invalidScenario < 0.5) {
        this.testListQuestionsInvalidQuizId();
      } else {
        this.testGetQuestionInvalidQuizId();
      }
    }
  }

  /**
   * Test listing questions with invalid quiz ID
   */
  testListQuestionsInvalidQuizId() {
    const response = http.get(
      `api/v1/quizzes/${this.testData.invalidQuizId}/questions`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "list_questions_invalid_quiz_id" },
      }
    );

    check(response, {
      "List questions with invalid quiz ID returns 401": (r) =>
        r.status === 401,
    });
  }

  /**
   * Test listing questions with valid quiz ID
   */
  testListQuestionsValidQuizId() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/quizzes/${this.testData.validQuizId}/questions`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "list_questions_valid" },
      }
    );

    check(response, {
      "List questions with valid quiz ID returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains expected question count": (r) => {
        try {
          const data = JSON.parse(r.body);
          return (
            data.body &&
            data.body.data &&
            data.body.data.data &&
            data.body.data.data.length ===
              this.testData.expectedQuestionData.expectedQuestionCount
          );
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test getting question by ID with invalid quiz ID
   */
  testGetQuestionInvalidQuizId() {
    const response = http.get(
      `api/v1/quizzes/${this.testData.invalidQuizId}/questions/${
        this.testData.validQuestionId || "test-question-id"
      }`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "get_question_invalid_quiz_id" },
      }
    );

    check(response, {
      "Get question with invalid quiz ID returns 401": (r) => r.status === 401,
    });
  }

  /**
   * Test getting question by ID with invalid question ID
   */
  testGetQuestionInvalidQuestionId() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/quizzes/${this.testData.validQuizId}/questions/${this.testData.invalidQuestionId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "get_question_invalid_question_id" },
      }
    );

    check(response, {
      "Get question with invalid question ID returns 500": (r) =>
        r.status === 500,
    });
  }

  /**
   * Test getting question by ID with valid IDs
   */
  testGetQuestionValidIds() {
    if (!this.testData.validQuizId || !this.testData.validQuestionId) {
      console.warn("No valid quiz ID or question ID available for testing");
      return;
    }

    const response = http.get(
      `api/v1/quizzes/${this.testData.validQuizId}/questions/${this.testData.validQuestionId}`,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "get_question_valid" },
      }
    );

    check(response, {
      "Get question with valid IDs returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains question data": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.body && data.body.data && data.body.data.question_id;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test updating question with invalid quiz ID
   */
  testUpdateQuestionInvalidQuizId() {
    const response = http.put(
      `api/v1/quizzes/${this.testData.invalidQuizId}/questions/${
        this.testData.validQuestionId || "test-question-id"
      }`,
      JSON.stringify(this.testData.updateQuestionPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "update_question_invalid_quiz_id" },
      }
    );

    check(response, {
      "Update question with invalid quiz ID returns 401": (r) =>
        r.status === 401,
    });
  }

  /**
   * Test updating question with missing body
   */
  testUpdateQuestionMissingBody() {
    if (!this.testData.validQuizId || !this.testData.validQuestionId) {
      console.warn("No valid quiz ID or question ID available for testing");
      return;
    }

    const response = http.put(
      `api/v1/quizzes/${this.testData.validQuizId}/questions/${this.testData.validQuestionId}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "update_question_missing_body" },
      }
    );

    check(response, {
      "Update question with missing body returns 400": (r) => r.status === 400,
    });
  }

  /**
   * Test updating question with valid input
   */
  testUpdateQuestionValid() {
    if (!this.testData.validQuizId || !this.testData.validQuestionId) {
      console.warn("No valid quiz ID or question ID available for testing");
      return;
    }

    const response = http.put(
      `api/v1/quizzes/${this.testData.validQuizId}/questions/${this.testData.validQuestionId}`,
      JSON.stringify(this.testData.updateQuestionPayload),
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "update_question_valid" },
      }
    );

    check(response, {
      "Update question with valid input returns 200": (r) => r.status === 200,
    });
  }

  /**
   * Test deleting question with invalid quiz ID
   */
  testDeleteQuestionInvalidQuizId() {
    const response = http.del(
      `api/v1/quizzes/${this.testData.invalidQuizId}/questions/${
        this.testData.validQuestionId || "test-question-id"
      }`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_question_invalid_quiz_id" },
      }
    );

    check(response, {
      "Delete question with invalid quiz ID returns 401": (r) =>
        r.status === 401,
    });
  }

  /**
   * Test deleting question when active quiz is present
   */
  testDeleteQuestionActiveQuiz() {
    if (!this.testData.validQuizId || !this.testData.validQuestionId) {
      console.warn("No valid quiz ID or question ID available for testing");
      return;
    }

    // First generate a demo session to make the quiz active
    this.generateDemoSession(this.testData.validQuizId);

    const response = http.del(
      `api/v1/quizzes/${this.testData.validQuizId}/questions/${this.testData.validQuestionId}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_question_active_quiz" },
      }
    );

    check(response, {
      "Delete question with active quiz returns 400": (r) => r.status === 400,
    });
  }

  /**
   * Test deleting question with valid input
   */
  testDeleteQuestionValid() {
    // Create a new quiz for deletion test to avoid conflicts
    const quizResponse = this.createTestQuiz();
    if (!quizResponse || !quizResponse.body || !quizResponse.body.data) {
      console.warn("Failed to create quiz for deletion test");
      return;
    }

    const quizId = quizResponse.body.data;
    const questionsResponse = this.listQuestions(quizId);

    if (
      !questionsResponse ||
      !questionsResponse.body ||
      !questionsResponse.body.data ||
      !questionsResponse.body.data.data ||
      questionsResponse.body.data.data.length === 0
    ) {
      console.warn("No questions available for deletion test");
      return;
    }

    const questionId = questionsResponse.body.data.data[0].question_id;

    const response = http.del(
      `api/v1/quizzes/${quizId}/questions/${questionId}`,
      null,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "delete_question_valid" },
      }
    );

    check(response, {
      "Delete question with valid input returns 200": (r) => r.status === 200,
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
   * Helper method to list questions
   */
  listQuestions(quizId) {
    const response = http.get(`api/v1/quizzes/${quizId}/questions`, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "list_questions_helper" },
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
}

// Export test instance
const testInstance = new QuestionControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
