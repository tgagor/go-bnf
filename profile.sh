#!/bin/bash

# Profile the bnf tool parsing purge.sh
echo "Building with profiling support..."
go build -o bin/bnf-profile .

echo "Running with CPU profiling..."
timeout 30s bin/bnf-profile -g examples/bash-quoted.bnf -i examples/purge.sh -s SCRIPT 2>&1 | head -20

echo ""
echo "If it's still running after 30s, it's definitely stuck in exponential backtracking"
