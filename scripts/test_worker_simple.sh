#!/bin/bash

echo "=========================================="
echo "Worker System Compilation Test"
echo "=========================================="
echo ""

echo "Step 1: Test Worker compilation"
echo "------------------------------------------"
go build -o bin/worker cmd/worker/main.go 2>&1

if [ $? -eq 0 ]; then
    echo "✓ Worker compiled successfully!"
    ls -lh bin/worker
    exit 0
else
    echo "✗ Worker compilation failed"
    exit 1
fi
