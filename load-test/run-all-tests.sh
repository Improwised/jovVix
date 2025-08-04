#!/bin/bash

# K6 Load Test Runner - Sequential Execution
# Usage: ./run-all-tests.sh [ENVIRONMENT] [INTENSITY] [STRICTNESS]
# Example: ./run-all-tests.sh local light default

set -e

# Configuration
ENVIRONMENT=${1:-local}
INTENSITY=${2:-functional}
STRICTNESS=${3:-default}
RESULTS_DIR="results/$(date +%Y%m%d_%H%M%S)"
LOG_FILE="$RESULTS_DIR/execution.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test files (excluding non-test files)
TESTS=(
    "analytics-board-admin-test"
    "analytics-board-user-test"
    "final-scoreboard-admin-test"
    "final-scoreboard-test"
    "image-test"
    "question-test"
    "quiz-test"
    "quiz-websocket-test"
    "shared-quizzes-test"
    "user-test"
    "user-played-score-quiz-test"
)

# Create results directory
mkdir -p "$RESULTS_DIR"

echo -e "${BLUE}=== K6 Load Test Suite - Sequential Execution ===${NC}"
echo -e "${YELLOW}Environment: $ENVIRONMENT${NC}"
echo -e "${YELLOW}Intensity: $INTENSITY${NC}"
echo -e "${YELLOW}Strictness: $STRICTNESS${NC}"
echo -e "${YELLOW}Results Directory: $RESULTS_DIR${NC}"
echo ""

# Initialize summary
START_TIME=$(date +%s)

# Function to run a single test
run_test() {
    local test_name=$1
    local test_file="tests/${test_name}.js"
    local result_file="$RESULTS_DIR/${test_name}_result.json"
    local summary_file="$RESULTS_DIR/${test_name}_summary.txt"
    
    echo -e "${BLUE}Running: $test_name${NC}" | tee -a "$LOG_FILE"
    echo "Started at: $(date)" | tee -a "$LOG_FILE"
    
    # Run K6 test
    k6 run \
        --env ENVIRONMENT="$ENVIRONMENT" \
        --env INTENSITY="$INTENSITY" \
        --env STRICTNESS="$STRICTNESS" \
        --env TEST_MODULE="$test_name" \
        --out json="$result_file" \
        --summary-export="$summary_file" \
        "$test_file" 2>&1 | tee -a "$LOG_FILE";
    
    echo "Completed at: $(date)" | tee -a "$LOG_FILE"
    echo "----------------------------------------" | tee -a "$LOG_FILE"
}

# Run all tests sequentially
for test in "${TESTS[@]}"; do
    run_test "$test"
    # Small delay between tests to allow system recovery
    sleep 2
done

# Calculate execution time
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Generate final summary
echo -e "\n${BLUE}=== EXECUTION SUMMARY ===${NC}" | tee -a "$LOG_FILE"
echo "Total Duration: ${DURATION}s" | tee -a "$LOG_FILE"
echo "Total Tests: ${#TESTS[@]}" | tee -a "$LOG_FILE"

echo -e "\n${YELLOW}Results saved to: $RESULTS_DIR${NC}"
echo -e "${YELLOW}Log file: $LOG_FILE${NC}"
