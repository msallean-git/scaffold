#!/usr/bin/env bash
set -e
echo "Building scaffold..."
go mod tidy
go build -o scaffold .
echo ""
echo "Build successful: ./scaffold"
echo ""
echo "Quick start:"
echo "  ./scaffold init"
echo "  ./scaffold list"
echo "  ./scaffold use general-dev"
