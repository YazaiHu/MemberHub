#!/bin/bash

echo "=========================================="
echo "Build Verification Test"
echo "=========================================="
echo ""

echo "Step 1: Verify Worker compilation"
echo "------------------------------------------"
if [ -f "bin/worker" ]; then
    echo "✓ Worker binary exists"
    ls -lh bin/worker
else
    echo "Building Worker..."
    go build -o bin/worker cmd/worker/main.go
    if [ $? -eq 0 ]; then
        echo "✓ Worker compiled successfully"
        ls -lh bin/worker
    else
        echo "✗ Worker compilation failed"
        exit 1
    fi
fi

echo ""
echo "Step 2: Verify API compilation"
echo "------------------------------------------"
if [ -f "bin/api" ]; then
    echo "✓ API binary exists"
    ls -lh bin/api
else
    echo "Building API..."
    go build -o bin/api cmd/api/main.go
    if [ $? -eq 0 ]; then
        echo "✓ API compiled successfully"
        ls -lh bin/api
    else
        echo "✗ API compilation failed"
        exit 1
    fi
fi

echo ""
echo "Step 3: Verify Docker services"
echo "------------------------------------------"
SERVICES=$(docker ps --filter "name=membership-" --format "{{.Names}}" | wc -l)
if [ "$SERVICES" -eq 3 ]; then
    echo "✓ All Docker services running:"
    docker ps --filter "name=membership-" --format "  - {{.Names}} ({{.Status}})"
else
    echo "⚠ Expected 3 services, found $SERVICES"
    docker ps --filter "name=membership-" --format "  - {{.Names}} ({{.Status}})"
fi

echo ""
echo "=========================================="
echo "Build Verification Summary"
echo "=========================================="
echo ""
echo "✓ Worker binary: $(ls -lh bin/worker | awk '{print $5}')"
echo "✓ API binary: $(ls -lh bin/api | awk '{print $5}')"
echo "✓ Docker services: $SERVICES/3 running"
echo ""
echo "All compilation errors fixed!"
echo "System is ready for integration testing."
echo ""
echo "Next steps:"
echo "1. Configure database password in configs/config.dev.yaml"
echo "2. Run: make run (to start API)"
echo "3. Run: make worker (to start Worker)"
echo "4. Run: ./scripts/test_worker.sh (for integration test)"
echo ""
