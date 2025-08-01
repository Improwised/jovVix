/**
 * HTTP Client Wrapper
 * 
 * Provides a consistent interface for making HTTP requests with proper error handling,
 * authentication, and logging.
 */

import http from 'k6/http';
import { check } from 'k6';
import { getBaseUrls } from '../utils/config.js';

// Get base URLs from configuration
const { baseUrl, kratosUrl } = getBaseUrls();

/**
 * Make a GET request
 * @param {string} path - The endpoint path (without base URL)
 * @param {Object} params - Request parameters
 * @returns {Object} Response object
 */
export function get(path, params = {}) {
  const url = buildUrl(path);
  const response = http.get(url, prepareParams(params));
  return handleResponse(response, 'GET', url);
}

/**
 * Make a POST request
 * @param {string} path - The endpoint path (without base URL)
 * @param {Object|string} payload - Request payload
 * @param {Object} params - Request parameters
 * @returns {Object} Response object
 */
export function post(path, payload = null, params = {}) {
  const url = buildUrl(path);
  const response = http.post(url, payload, prepareParams(params));
  return handleResponse(response, 'POST', url);
}

/**
 * Make a PUT request
 * @param {string} path - The endpoint path (without base URL)
 * @param {Object|string} payload - Request payload
 * @param {Object} params - Request parameters
 * @returns {Object} Response object
 */
export function put(path, payload = null, params = {}) {
  const url = buildUrl(path);
  const response = http.put(url, payload, prepareParams(params));
  return handleResponse(response, 'PUT', url);
}

/**
 * Make a DELETE request
 * @param {string} path - The endpoint path (without base URL)
 * @param {Object|string} payload - Request payload
 * @param {Object} params - Request parameters
 * @returns {Object} Response object
 */
export function del(path, payload = null, params = {}) {
  const url = buildUrl(path);
  const response = http.del(url, payload, prepareParams(params));
  return handleResponse(response, 'DELETE', url);
}

/**
 * Make a request to Kratos authentication service
 * @param {string} path - The endpoint path (without Kratos base URL)
 * @param {string} method - HTTP method
 * @param {Object|string} payload - Request payload
 * @param {Object} params - Request parameters
 * @returns {Object} Response object
 */
export function kratosRequest(path, method = 'GET', payload = null, params = {}) {
  const url = `${kratosUrl}${path}`;
  let response;
  
  switch (method.toUpperCase()) {
    case 'POST':
      response = http.post(url, payload, prepareParams(params));
      break;
    case 'PUT':
      response = http.put(url, payload, prepareParams(params));
      break;
    case 'DELETE':
      response = http.del(url, payload, prepareParams(params));
      break;
    default:
      response = http.get(url, prepareParams(params));
  }
  
  return handleResponse(response, method, url);
}

/**
 * Build a complete URL from a path
 * @param {string} path - The endpoint path
 * @returns {string} Complete URL
 */
function buildUrl(path) {
  // If path is already a complete URL, return it as is
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path;
  }
  
  // Ensure path starts with a slash
  if (!path.startsWith('/')) {
    path = '/' + path;
  }
  
  return `${baseUrl}${path}`;
}

/**
 * Prepare request parameters
 * @param {Object} params - Request parameters
 * @returns {Object} Prepared parameters
 */
function prepareParams(params = {}) {
  // Set default headers if not provided
  if (!params.headers) {
    params.headers = {};
  }
  
  // Set default Content-Type for JSON payloads if not specified
  if (!params.headers['Content-Type'] && !params.files) {
    params.headers['Content-Type'] = 'application/json';
  }
  
  // Add default tags for metrics
  if (!params.tags) {
    params.tags = {};
  }
  
  return params;
}

/**
 * Handle HTTP response
 * @param {Object} response - HTTP response object
 * @param {string} method - HTTP method used
 * @param {string} url - Request URL
 * @returns {Object} Processed response
 */
function handleResponse(response, method, url) {
  // Add parsed JSON body if response is JSON
  try {
    if (response.headers['Content-Type'] && 
        response.headers['Content-Type'].includes('application/json')) {
      response.json = JSON.parse(response.body);
    }
  } catch (e) {
    // Not JSON or invalid JSON, ignore
  }
  
  // Log failed requests (5xx errors) for debugging
  if (response.status >= 500) {
    console.error(`${method} ${url} failed with status ${response.status}`);
    console.error(`Response: ${response.body.substring(0, 200)}${response.body.length > 200 ? '...' : ''}`);
  }
  
  return response;
}

/**
 * Check response for common success conditions
 * @param {Object} response - HTTP response object
 * @param {string} description - Check description
 * @param {Object} additionalChecks - Additional checks to perform
 * @returns {boolean} Whether all checks passed
 */
export function checkResponse(response, description, additionalChecks = {}) {
  const defaultChecks = {
    [`${description} returns successful status`]: (r) => r.status >= 200 && r.status < 300,
    [`${description} has valid response time`]: (r) => r.timings.duration < 10000, // 10s timeout
  };
  
  // Combine default checks with additional checks
  const checks = { ...defaultChecks, ...additionalChecks };
  
  return check(response, checks);
}