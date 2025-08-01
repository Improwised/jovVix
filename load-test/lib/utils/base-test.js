/**
 * Base Test Class
 * 
 * Provides a foundation for all test scripts with common functionality.
 */

import { group, sleep } from 'k6';
import { buildOptions } from './config.js';
import { setupAuth } from '../auth/auth.js';

/**
 * Base test class to be extended by specific test implementations
 */
export class BaseTest {
  /**
   * Constructor
   * @param {Object} params - Test parameters
   */
  constructor(params = {}) {
    this.name = params.name || 'Unnamed Test';
    this.description = params.description || '';
    this.options = params.options || {};
    this.authOptions = params.authOptions || {};
    this.state = {};
  }
  
  /**
   * Get K6 options for this test
   * @returns {Object} K6 options
   */
  getOptions() {
    return buildOptions(this.options);
  }
  
  /**
   * Setup function - runs once before all test iterations
   * @returns {Object} Setup data
   */
  setup() {
    console.log(`Setting up test: ${this.name}`);
    
    // Setup authentication
    const auth = setupAuth(this.authOptions);
    
    // Call test-specific setup
    const setupData = this.setupTest(auth) || {};
    
    return {
      ...auth,
      ...setupData,
      testName: this.name
    };
  }
  
  /**
   * Test-specific setup - to be implemented by subclasses
   * @param {Object} auth - Authentication data
   * @returns {Object} Additional setup data
   */
  setupTest(auth) {
    // Default implementation does nothing
    return {};
  }
  
  /**
   * Main test function - called for each VU iteration
   * @param {Object} data - Data from setup
   */
  default(data) {
    // Store setup data in state
    this.state = { ...data };
    
    // Determine test type based on scenario
    const testType = __ENV.K6_SCENARIO_NAME || 'functional_test';
    
    if (testType === 'functional_test') {
      this.runFunctionalTests();
    } else {
      this.runLoadTests();
    }
  }
  
  /**
   * Run comprehensive functional tests
   */
  runFunctionalTests() {
    group(`${this.name} - Functional Tests`, () => {
      // Call test-specific implementation
      this.functionalTests();
    });
  }
  
  /**
   * Run load tests
   */
  runLoadTests() {
    // Call test-specific implementation
    this.loadTests();
    
    // Add a small sleep to prevent hammering the server
    sleep(0.3);
  }
  
  /**
   * Functional tests - to be implemented by subclasses
   */
  functionalTests() {
    console.warn(`${this.name}: functionalTests() not implemented`);
  }
  
  /**
   * Load tests - to be implemented by subclasses
   */
  loadTests() {
    console.warn(`${this.name}: loadTests() not implemented`);
  }
  
  /**
   * Teardown function - runs once after all test iterations
   * @param {Object} data - Data from setup
   */
  teardown(data) {
    console.log(`Tearing down test: ${this.name}`);
    
    // Call test-specific teardown
    this.teardownTest(data);
  }
  
  /**
   * Test-specific teardown - to be implemented by subclasses
   * @param {Object} data - Data from setup
   */
  teardownTest(data) {
    // Default implementation does nothing
  }
}