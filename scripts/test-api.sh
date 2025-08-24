#!/bin/bash

# Beer API Test Script
# This script tests all the API endpoints automatically

set -e  # Exit on any error

API_BASE="http://localhost:8080"
API_V1="$API_BASE/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_status="$5"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    printf "${BLUE}Test $TESTS_RUN: $test_name${NC}\n"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" "$url" \
            -H "Accept: application/json")
    fi
    
    http_code=$(echo "$response" | grep -o "HTTPSTATUS:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed -E 's/HTTPSTATUS:[0-9]*$//')
    
    if [ "$http_code" = "$expected_status" ]; then
        printf "${GREEN}✅ PASS${NC} - Status: $http_code\n"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        if [ -n "$body" ] && command -v jq >/dev/null 2>&1; then
            echo "$body" | jq . 2>/dev/null || echo "$body"
        else
            echo "$body"
        fi
    else
        printf "${RED}❌ FAIL${NC} - Expected: $expected_status, Got: $http_code\n"
        echo "Response: $body"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    printf "\n"
}

# Function to check if server is running
check_server() {
    printf "${YELLOW}Checking if server is running...${NC}\n"
    if curl -s "$API_BASE/ping" > /dev/null; then
        printf "${GREEN}✅ Server is running${NC}\n\n"
    else
        printf "${RED}❌ Server is not running. Please start the server first:${NC}\n"
        printf "${YELLOW}make dev${NC}\n"
        printf "${YELLOW}# or${NC}\n"
        printf "${YELLOW}export DB_TYPE=inmemory && go run cmd/main.go${NC}\n\n"
        exit 1
    fi
}

# Main test execution
main() {
    printf "${BLUE}=== Beer API Test Suite ===${NC}\n\n"
    
    check_server
    
    # Test 1: Health Check
    run_test "Health Check" "GET" "$API_BASE/ping" "" "200"
    
    # Test 2: Get All Beers (Empty)
    run_test "Get All Beers (Empty)" "GET" "$API_V1/beers" "" "200"
    
    # Test 3: Create Beer 1
    run_test "Create Corona Extra" "POST" "$API_V1/beers" '{
        "id": 1,
        "name": "Corona Extra",
        "brewery": "Cervecería Modelo",
        "country": "Mexico",
        "price": 1200,
        "currency": "CLP"
    }' "201"
    
    # Test 4: Create Beer 2
    run_test "Create Heineken" "POST" "$API_V1/beers" '{
        "id": 2,
        "name": "Heineken",
        "brewery": "Heineken N.V.",
        "country": "Netherlands",
        "price": 1500,
        "currency": "CLP"
    }' "201"
    
    # Test 5: Create Beer 3 (USD)
    run_test "Create Budweiser (USD)" "POST" "$API_V1/beers" '{
        "id": 3,
        "name": "Budweiser",
        "brewery": "Anheuser-Busch",
        "country": "USA",
        "price": 25.99,
        "currency": "USD"
    }' "201"
    
    # Test 6: Get All Beers (With Data)
    run_test "Get All Beers (With Data)" "GET" "$API_V1/beers" "" "200"
    
    # Test 7: Get Beer by ID
    run_test "Get Beer by ID (Corona)" "GET" "$API_V1/beers/1" "" "200"
    
    # Test 8: Get Non-existent Beer
    run_test "Get Non-existent Beer" "GET" "$API_V1/beers/999" "" "404"
    
    # Test 9: Calculate Box Price (Same Currency)
    run_test "Calculate Box Price (CLP)" "GET" "$API_V1/beers/1/boxprice?quantity=6&currency=CLP" "" "200"
    
    # Test 10: Calculate Box Price (Currency Conversion)
    run_test "Calculate Box Price (USD)" "GET" "$API_V1/beers/1/boxprice?quantity=12&currency=USD" "" "200"
    
    # Test 11: Create Duplicate Beer (Error)
    run_test "Create Duplicate Beer" "POST" "$API_V1/beers" '{
        "id": 1,
        "name": "Duplicate Corona",
        "brewery": "Another Brewery",
        "country": "Chile",
        "price": 1000,
        "currency": "CLP"
    }' "409"
    
    # Test 12: Create Beer with Invalid Data
    run_test "Create Beer (Missing Name)" "POST" "$API_V1/beers" '{
        "id": 10,
        "brewery": "Test Brewery",
        "country": "Chile",
        "price": 1000,
        "currency": "CLP"
    }' "400"
    
    # Test 13: Create Beer with Negative Price
    run_test "Create Beer (Negative Price)" "POST" "$API_V1/beers" '{
        "id": 11,
        "name": "Invalid Beer",
        "brewery": "Test Brewery",
        "country": "Chile",
        "price": -100,
        "currency": "CLP"
    }' "400"
    
    # Test 14: Create Beer with Invalid Currency
    run_test "Create Beer (Invalid Currency)" "POST" "$API_V1/beers" '{
        "id": 12,
        "name": "Invalid Currency Beer",
        "brewery": "Test Brewery",
        "country": "Chile",
        "price": 1000,
        "currency": "INVALID"
    }' "400"
    
    # Test 15: Box Price with Invalid Quantity
    run_test "Box Price (Invalid Quantity)" "GET" "$API_V1/beers/1/boxprice?quantity=0&currency=CLP" "" "400"
    
    # Test 16: Box Price with Invalid Currency
    run_test "Box Price (Invalid Currency)" "GET" "$API_V1/beers/1/boxprice?quantity=6&currency=INVALID" "" "400"
    
    # Test 17: Box Price for Non-existent Beer
    run_test "Box Price (Non-existent Beer)" "GET" "$API_V1/beers/999/boxprice?quantity=6&currency=CLP" "" "404"
    
    # Test 18: Legacy Route - Create Beer
    run_test "Legacy Route - Create Beer" "POST" "$API_BASE/beers" '{
        "id": 100,
        "name": "Legacy Beer",
        "brewery": "Legacy Brewery",
        "country": "Chile",
        "price": 800,
        "currency": "CLP"
    }' "201"
    
    # Test 19: Legacy Route - Get All Beers
    run_test "Legacy Route - Get All Beers" "GET" "$API_BASE/beers" "" "200"
    
    # Test 20: Legacy Route - Get Beer by ID
    run_test "Legacy Route - Get Beer by ID" "GET" "$API_BASE/beers/100" "" "200"
    
    # Test Results Summary
    printf "${BLUE}=== Test Results Summary ===${NC}\n"
    printf "Total Tests: $TESTS_RUN\n"
    printf "${GREEN}Passed: $TESTS_PASSED${NC}\n"
    printf "${RED}Failed: $TESTS_FAILED${NC}\n"
    
    if [ "$TESTS_FAILED" -eq 0 ]; then
        printf "\n${GREEN}🎉 All tests passed!${NC}\n"
        exit 0
    else
        printf "\n${RED}❌ Some tests failed${NC}\n"
        exit 1
    fi
}

# Check for dependencies
check_dependencies() {
    if ! command -v curl >/dev/null 2>&1; then
        printf "${RED}Error: curl is required but not installed${NC}\n"
        exit 1
    fi
    
    if ! command -v jq >/dev/null 2>&1; then
        printf "${YELLOW}Warning: jq is not installed. JSON responses will not be formatted.${NC}\n"
        printf "${YELLOW}Install jq for better output: brew install jq${NC}\n\n"
    fi
}

# Show usage
show_usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -v, --verbose  Show verbose output"
    echo ""
    echo "Prerequisites:"
    echo "  - Server must be running on http://localhost:8080"
    echo "  - curl must be installed"
    echo "  - jq is recommended for JSON formatting"
    echo ""
    echo "Start the server with:"
    echo "  make dev"
    echo "  # or"
    echo "  export DB_TYPE=inmemory && go run cmd/main.go"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -v|--verbose)
            set -x
            shift
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Run the tests
check_dependencies
main
