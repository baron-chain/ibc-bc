#!/usr/bin/env bash

# Fail fast on any error or undefined variable
set -euo pipefail

# Define constants
readonly COVERAGE_FILE="coverage.txt"
readonly PROFILE_FILE="profile.out"
readonly TIMEOUT="30m"
readonly TEST_TAGS="ledger test_ledger_mock"
readonly EXCLUDE_PATTERN="/simapp"

main() {
    initialize_coverage_file
    run_tests
    cleanup
}

initialize_coverage_file() {
    echo "mode: atomic" > "$COVERAGE_FILE"
}

get_packages() {
    go list ./... | grep -v "$EXCLUDE_PATTERN"
}

run_tests() {
    local packages
    packages=$(get_packages)
    
    echo "Running tests for packages..."
    
    local count=0
    local total
    total=$(echo "$packages" | wc -l)
    
    while IFS= read -r pkg; do
        ((count++))
        echo "[${count}/${total}] Testing $pkg"
        
        # Run test with coverage for each package
        go test -v \
            -timeout "$TIMEOUT" \
            -race \
            -coverprofile="$PROFILE_FILE" \
            -covermode=atomic \
            -tags="$TEST_TAGS" \
            "$pkg"
        
        # Append coverage data if profile exists
        if [[ -f "$PROFILE_FILE" ]]; then
            tail -n +2 "$PROFILE_FILE" >> "$COVERAGE_FILE"
        fi
    done <<< "$packages"
}

cleanup() {
    rm -f "$PROFILE_FILE"
    echo "Coverage report generated in $COVERAGE_FILE"
}

# Run main function
main
