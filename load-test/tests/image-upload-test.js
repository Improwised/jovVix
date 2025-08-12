/**
 * Image Controller Tests
 *
 * Tests for the image upload endpoints of the API.
 */

import { check } from "k6";
import { BaseTest } from "../lib/utils/base-test.js";
import * as http from "../lib/http/client.js";
import { runTest } from "../lib/utils/test-helpers.js";

export class ImageControllerTest extends BaseTest {
  constructor() {
    super({
      name: "Image Controller",
      description: "Tests for image upload endpoints",
      options: {
        thresholds: {
          http_req_duration: ["p(95)<5000"], // Higher threshold for file uploads
          http_req_failed: ["rate<0.3"],
        },
      },
  });

    // Test data
    this.testData = {
      validQuizId: "", // Will be set during setup
      validQuestionId: "", // Will be set during setup
      testImageData: null, // Mock image data for testing
      expectedResponses: {
        missingQuizId: 400,
        missingFile: 400,
        validUpload: 200,
        internalError: 500,
      },
    };
  }

  setupTest(auth) {
    // Create test quiz for image upload tests
    const quizResponse = this.createTestQuiz();
    if (quizResponse && quizResponse.body && quizResponse.body.data) {
      this.testData.validQuizId = quizResponse.body.data;

      // Get questions for the quiz to get a valid question ID
      const questionsResponse = this.getQuizQuestions(
        this.testData.validQuizId
      );
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

    // Create mock image data for testing
    this.testData.testImageData = this.createMockImageData();

    return {
      validQuizId: this.testData.validQuizId,
      validQuestionId: this.testData.validQuestionId,
    };
  }

  functionalTests() {
    runTest("Missing Quiz ID", () => this.testMissingQuizId());
    runTest("Missing File", () => this.testMissingFile());
    runTest("Invalid Filename Format", () => this.testInvalidFilenameFormat());
    runTest("Valid Image Upload", () => this.testValidImageUpload());
  }

  loadTests() {
    const scenario = Math.random();

    if (scenario < 0.4) {
      // 40% valid image uploads
      this.testValidImageUpload();
    } else if (scenario < 0.6) {
      // 20% missing quiz ID requests
      this.testMissingQuizId();
    } else if (scenario < 0.8) {
      // 20% missing file requests
      this.testMissingFile();
    } else {
      // 20% invalid filename format requests
      this.testInvalidFilenameFormat();
    }
  }

  /**
   * Test image upload with missing quiz ID
   */
  testMissingQuizId() {
    if (!this.testData.testImageData) {
      console.warn("No test image data available");
      return;
    }

    const formData = {
      "image-attachment": http.file(
        this.testData.testImageData,
        "test-image.png",
        "image/png"
      ),
    };

    const response = http.post("api/v1/images", formData, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "image_upload_missing_quiz_id" },
    });

    check(response, {
      "Missing quiz ID returns 400": (r) => r.status === 400,
      "Response contains error message": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.message && data.message.includes("quiz_id");
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test image upload with missing file
   */
  testMissingFile() {
    if (!this.testData.validQuizId) {
      console.warn("No valid quiz ID available for testing");
      return;
    }

    const response = http.post(
      `api/v1/images?quiz_id=${this.testData.validQuizId}`,
      {},
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "image_upload_missing_file" },
      }
    );

    check(response, {
      "Missing file returns 400": (r) => r.status === 400,
      "Response contains error message": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.message && data.message.includes("files");
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test image upload with invalid filename format (should cause internal error)
   */
  testInvalidFilenameFormat() {
    if (!this.testData.validQuizId || !this.testData.testImageData) {
      console.warn("No valid quiz ID or test image data available for testing");
      return;
    }

    const formData = {
      "image-attachment": http.file(
        this.testData.testImageData,
        "invalid-filename.png",
        "image/png"
      ),
    };

    const response = http.post(
      `api/v1/images?quiz_id=${this.testData.validQuizId}`,
      formData,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "image_upload_invalid_filename" },
      }
    );

    check(response, {
      "Invalid filename format returns 500": (r) => r.status === 500,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
    });
  }

  /**
   * Test valid image upload
   */
  testValidImageUpload() {
    if (
      !this.testData.validQuizId ||
      !this.testData.validQuestionId ||
      !this.testData.testImageData
    ) {
      console.warn("Missing required data for valid image upload test");
      return;
    }

    // Use question ID as filename to match the expected format
    const filename = this.testData.validQuestionId;

    const formData = {
      "image-attachment": http.file(
        this.testData.testImageData,
        filename,
        "image/png"
      ),
    };

    const response = http.post(
      `api/v1/images?quiz_id=${this.testData.validQuizId}`,
      formData,
      {
        cookies: { ory_kratos_session: this.state.adminSessionCookie },
        tags: { endpoint: "image_upload_valid" },
      }
    );

    check(response, {
      "Valid image upload returns 200": (r) => r.status === 200,
      "Response is valid JSON": (r) => {
        try {
          JSON.parse(r.body);
          return true;
        } catch (e) {
          return false;
        }
      },
      "Response contains success message": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.body && data.body.includes("uploaded successfully");
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
   * Helper method to get quiz questions
   */
  getQuizQuestions(quizId) {
    const response = http.get(`api/v1/quizzes/${quizId}/questions`, {
      cookies: { ory_kratos_session: this.state.adminSessionCookie },
      tags: { endpoint: "get_quiz_questions" },
    });

    return response;
  }

  /**
   * Helper method to create mock image data
   */
  createMockImageData() {
    // Create a simple 1x1 PNG image data (base64 encoded)
    // This is a minimal valid PNG file
    const pngData =
      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==";

    // Convert base64 to binary data
    try {
      return Uint8Array.from(atob(pngData), (c) => c.charCodeAt(0));
    } catch (e) {
      console.warn("Failed to create mock image data:", e.message);
      return null;
    }
  }
}

// Export test instance
const testInstance = new ImageControllerTest();

// Export K6 functions
export const options = testInstance.getOptions();
export const setup = testInstance.setup.bind(testInstance);
export const teardown = testInstance.teardown.bind(testInstance);
export default testInstance.default.bind(testInstance);
