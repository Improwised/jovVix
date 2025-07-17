/**
 * Test Helper Functions
 *
 * Common utility functions for test scripts.
 */

import { sleep, group, check } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";
import * as http from "../http/client.js";

/**
 * Generate a unique identifier
 * @param {string} prefix - Optional prefix for the ID
 * @param {number} length - Length of the random part (default: 8)
 * @returns {string} Unique identifier
 */
export function generateId(prefix = "", length = 8) {
  return `${prefix}${uuidv4().substring(0, length)}`;
}

/**
 * Run a test with proper error handling
 * @param {string} name - Test name
 * @param {Function} testFn - Test function
 * @param {Object} options - Test options
 */
export function runTest(name, testFn, options = {}) {
  group(name, () => {
    try {
      testFn();
    } catch (error) {
      console.error(`Error in test "${name}": ${error.message}`);
      if (options.failOnError) {
        throw error;
      }
    }
  });
}

/**
 * Generate test CSV content for quiz upload
 * @param {Object} options - Options for CSV generation
 * @returns {string} CSV content as string
 */
export function generateTestCSV(options = {}) {
  const numQuestions = options.numQuestions || 5;

  let csvContent = "Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource\n";

  for (let i = 1; i <= numQuestions; i++) {
    const questionText = `Test Question ${i}?`;
    const optionA = `Option A for question ${i}`;
    const optionB = `Option B for question ${i}`;
    const optionC = `Option C for question ${i}`;
    const optionD = `Option D for question ${i}`;
    
    csvContent += `${questionText},single answer,1,${optionA},${optionB},${optionC},${optionD},,2,text,text,\n`;
  }

  return csvContent;
}

/**
 * Create a quiz for testing
 * @param {string} sessionCookie - Admin session cookie
 * @param {Object} quizData - Quiz data
 * @returns {Object} Created quiz data
 */
export function createTestQuiz(sessionCookie, quizData = {}) {
  // Generate quiz title if not provided
  const quizTitle = quizData.title || `quiz_${generateId("", 8)}`;

  // Generate CSV content
  const csvContent = generateTestCSV({
    numQuestions: quizData.numQuestions || 5
  });

  // Create the multipart form data manually
  const boundary = '----formdata-k6-' + Math.random().toString(36);
  let body = '';
  
  // Add description field
  body += `--${boundary}\r\n`;
  body += `Content-Disposition: form-data; name="description"\r\n\r\n`;
  body += `${quizData.description || "Test quiz description"}\r\n`;
  
  // Add file field
  body += `--${boundary}\r\n`;
  body += `Content-Disposition: form-data; name="attachment"; filename="test_quiz.csv"\r\n`;
  body += `Content-Type: text/csv\r\n\r\n`;
  body += csvContent + '\r\n';
  body += `--${boundary}--\r\n`;

  // Create quiz
  const response = http.post(
    `api/v1/quizzes/${quizTitle}/upload`,
    body,
    {
      headers: {
        'Content-Type': `multipart/form-data; boundary=${boundary}`
      },
      cookies: { ory_kratos_session: sessionCookie }
    }
  );


  // Check response
  const success = check(response, {
    "Quiz creation successful": (r) => r.status === 202,
  });

  let quizId = "";
  if (success && data.data) {
    quizId = data.data;
  }

  return {
    success,
    quizId,
    quizTitle,
    response,
  };
}

/**
 * Create a quiz session for testing
 * @param {string} sessionCookie - Admin session cookie
 * @param {string} quizId - Quiz ID
 * @param {boolean} isDemoSession - Whether to create a demo session
 * @returns {Object} Created session data
 */
export function createQuizSession(sessionCookie, quizId, isDemoSession = true) {
  // Determine endpoint based on session type
  const endpoint = isDemoSession
    ? `api/v1/quizzes/${quizId}/demo_session`
    : `api/v1/quizzes/${quizId}/session`;

  // Create session
  const response = http.post(endpoint, null, {
    cookies: { ory_kratos_session: sessionCookie },
  });


  // Check response
  const success = check(response, {
    "Session creation successful": (r) => r.status === 202,
  });

  let sessionId = "";
  if (success && response.data) {
    sessionId = response.data;
  }

  return {
    success,
    sessionId,
    quizId,
    response,
  };
}

/**
 * Wait for a condition to be true
 * @param {Function} conditionFn - Function that returns true when condition is met
 * @param {Object} options - Options for waiting
 * @returns {boolean} Whether condition was met within timeout
 */
export function waitFor(conditionFn, options = {}) {
  const timeout = options.timeout || 10; // seconds
  const interval = options.interval || 0.5; // seconds
  const startTime = Date.now();

  while (Date.now() - startTime < timeout * 1000) {
    if (conditionFn()) {
      return true;
    }
    sleep(interval);
  }

  return false;
}

/**
 * Parse JSON safely
 * @param {string} text - JSON string to parse
 * @param {*} defaultValue - Default value if parsing fails
 * @returns {*} Parsed JSON or default value
 */
export function safeJsonParse(text, defaultValue = null) {
  try {
    return JSON.parse(text);
  } catch (e) {
    return defaultValue;
  }
}

/**
 * Format duration in milliseconds to human-readable string
 * @param {number} ms - Duration in milliseconds
 * @returns {string} Formatted duration
 */
export function formatDuration(ms) {
  if (ms < 1000) {
    return `${ms}ms`;
  } else if (ms < 60000) {
    return `${(ms / 1000).toFixed(2)}s`;
  } else {
    const minutes = Math.floor(ms / 60000);
    const seconds = ((ms % 60000) / 1000).toFixed(2);
    return `${minutes}m ${seconds}s`;
  }
}
