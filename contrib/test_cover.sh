#!/usr/bin/env bash

set -euo pipefail

readonly COVERAGE_FILE="coverage.txt"
readonly PROFILE_FILE="profile.out"
readonly TIMEOUT="30m"
readonly TEST_TAGS="ledger test_ledger_mock"
readonly EXCLUDE_PATTERN="/simapp"
readonly LOG_FILE="test.log"
readonly PARALLEL_JOBS=${PARALLEL_JOBS:-4}

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

initialize_files() {
    echo "mode: atomic" > "$COVERAGE_FILE"
    : > "$LOG_FILE"
}

get_packages() {
    go list ./... | grep -v "$EXCLUDE_PATTERN" || true
}

run_package_test() {
    local pkg=$1
    local count=$2
    local total=$3
    
    log "[$count/$total] Testing $pkg"
    
    if go test -v \
        -timeout "$TIMEOUT" \
        -race \
        -coverprofile="$PROFILE_FILE.$count" \
        -covermode=atomic \
        -tags="$TEST_TAGS" \
        "$pkg" >> "$LOG_FILE" 2>&1; then
        
        if [[ -f "$PROFILE_FILE.$count" ]]; then
            tail -n +2 "$PROFILE_FILE.$count" >> "$COVERAGE_FILE"
            rm -f "$PROFILE_FILE.$count"
        fi
        return 0
    else
        log "Failed testing $pkg"
        return 1
    fi
}

export -f run_package_test
export -f log
export COVERAGE_FILE PROFILE_FILE TIMEOUT TEST_TAGS LOG_FILE

run_tests() {
    local packages count=0
    mapfile -t packages < <(get_packages)
    local total=${#packages[@]}
    
    log "Running tests for $total packages using $PARALLEL_JOBS parallel jobs..."
    
    local failures=0
    for pkg in "${packages[@]}"; do
        ((count++))
        {
            run_package_test "$pkg" "$count" "$total" || {
                failures=$((failures + 1))
            }
        } &
        
        # Limit parallel jobs
        if ((count % PARALLEL_JOBS == 0)); then
            wait
        fi
    done
    
    # Wait for remaining jobs
    wait
    
    if ((failures > 0)); then
        log "Test suite failed with $failures package failures"
        return 1
    fi
}

cleanup() {
    local status=$?
    rm -f "$PROFILE_FILE"*
    
    if ((status == 0)); then
        log "Coverage report generated successfully in $COVERAGE_FILE"
        log "Test logs available in $LOG_FILE"
    else
        log "Tests failed - check $LOG_FILE for details"
    fi
    
    return $status
}

main() {
    trap cleanup EXIT
    initialize_files
    run_tests
}

main
