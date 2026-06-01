#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT_DIR/docs/benchmarks"
RESULT_FILE="$OUT_DIR/latest.txt"

cd "$ROOT_DIR"
mkdir -p "$OUT_DIR"

go test -tags purego -run '^$' -bench . -benchmem ./... | tee "$RESULT_FILE"
rm -f "$OUT_DIR/tinylib-generated.svg"
go run ./cmd/benchsvg -in "$RESULT_FILE" -out-dir "$OUT_DIR" -title "msgpack benchmark (purego)"
