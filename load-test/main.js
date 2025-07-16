/**
 * Main Entry Point for K6 Load Tests
 * 
 * This file serves as the main entry point for running tests.
 * It dynamically imports the specified test module and executes it.
 */

import { sleep } from 'k6';

// Get test module from environment variable
const testModule = __ENV.TEST_MODULE || 'health-controller';

// Import the test module dynamically
export { default as setup } from `./tests/${testModule}.js`;
export { default as teardown } from `./tests/${testModule}.js`;

// Export options from the test module
export const options = setup.getOptions();

// Main test function
export default function(data) {
  // Call the default function of the test module
  setup.default(data);
}