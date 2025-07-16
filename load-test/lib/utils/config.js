/**
 * Configuration Manager
 * 
 * Centralizes all configuration loading and management for K6 tests.
 * Handles environment-specific settings, test options, and thresholds.
 */

import { SharedArray } from 'k6/data';

// Load config once and share across VUs
const configData = new SharedArray('config', function() {
  return [JSON.parse(open('../../config.json'))];
});

// Get the active environment from ENV or default to local
const activeEnv = __ENV.ENVIRONMENT || 'local';

/**
 * Get base configuration for the active environment
 * @returns {Object} Environment configuration
 */
export function getEnvironmentConfig() {
  const config = configData[0];
  
  if (!config.environments[activeEnv]) {
    console.error(`Environment '${activeEnv}' not found in config. Using 'local' instead.`);
    return config.environments.local;
  }
  
  return config.environments[activeEnv];
}

/**
 * Get test data configuration
 * @returns {Object} Test data configuration
 */
export function getTestData() {
  return configData[0].testData;
}

/**
 * Get load test options based on intensity level
 * @param {string} intensity - The intensity level (functional, light, moderate, heavy)
 * @returns {Object} K6 options for the specified intensity
 */
export function getLoadTestOptions(intensity = 'functional') {
  const config = configData[0];
  const validIntensities = Object.keys(config.loadTestOptions);
  
  if (!validIntensities.includes(intensity)) {
    console.warn(`Invalid intensity '${intensity}'. Using 'functional' instead.`);
    intensity = 'functional';
  }
  
  return config.loadTestOptions[intensity];
}

/**
 * Get threshold configuration based on strictness level
 * @param {string} strictness - The strictness level (default, strict, relaxed)
 * @returns {Object} Thresholds configuration
 */
export function getThresholds(strictness = 'default') {
  const config = configData[0];
  const validStrictness = Object.keys(config.thresholds);
  
  if (!validStrictness.includes(strictness)) {
    console.warn(`Invalid strictness '${strictness}'. Using 'default' instead.`);
    strictness = 'default';
  }
  
  return config.thresholds[strictness];
}

/**
 * Build complete K6 options object
 * @param {Object} customOptions - Custom options to override defaults
 * @returns {Object} Complete K6 options object
 */
export function buildOptions(customOptions = {}) {
  const intensity = __ENV.INTENSITY || 'functional';
  const strictness = __ENV.STRICTNESS || 'default';
  
  const loadOptions = getLoadTestOptions(intensity);
  const thresholds = getThresholds(strictness);
  
  // Build scenarios based on intensity
  let scenarios = {};
  
  if (intensity === 'functional') {
    scenarios.functional_test = {
      executor: 'shared-iterations',
      vus: loadOptions.vus,
      iterations: loadOptions.iterations,
      maxDuration: loadOptions.maxDuration,
      tags: { test_type: 'functional' },
    };
  } else {
    scenarios.load_test = {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: loadOptions.stages,
      tags: { test_type: 'load' },
    };
  }
  
  // Combine default options with custom options
  const options = {
    scenarios,
    thresholds,
    ...customOptions
  };
  
  return options;
}

/**
 * Get base URLs for API and authentication
 * @returns {Object} Object containing baseUrl and kratosUrl
 */
export function getBaseUrls() {
  const envConfig = getEnvironmentConfig();
  
  return {
    baseUrl: __ENV.BASE_URL || envConfig.baseUrl,
    kratosUrl: __ENV.KRATOS_URL || envConfig.kratosUrl
  };
}

/**
 * Get test metadata for a specific test
 * @param {string} testName - The name of the test
 * @returns {Object|null} Test metadata or null if not found
 */
export function getTestMetadata(testName) {
  const config = configData[0];
  
  if (config.testScripts && config.testScripts[testName]) {
    return config.testScripts[testName];
  }
  
  return null;
}