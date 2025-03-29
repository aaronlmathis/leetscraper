#!/bin/bash

set -e
set -x  # add this
set -euo pipefail

echo "🧪 Running LeetScraper Integration Tests..."

BIN="./dist/leetscraper"
TEST_DIR="./.test_output"
CONFIG_FILE="$HOME/.leetscraper.json"

# Backup existing config
if [[ -f "$CONFIG_FILE" ]]; then
  cp "$CONFIG_FILE" "$CONFIG_FILE.bak"
  RESTORE_CONFIG=true
else
  RESTORE_CONFIG=false
fi

cleanup() {
  echo "🧹 Cleaning up..."
  rm -rf "$TEST_DIR"
  if [[ "$RESTORE_CONFIG" == true ]]; then
    mv "$CONFIG_FILE.bak" "$CONFIG_FILE"
  fi
}
trap cleanup EXIT

mkdir -p "$TEST_DIR"

echo "🧪 Test 1: No flags (should use default config and fetch daily)"
$BIN > "$TEST_DIR/stdout1.txt" 2>&1

grep -q "📘 Problem:" "$TEST_DIR/stdout1.txt"

echo "✅ Passed Test 1"

echo "🧪 Test 2: Custom output dir and filename format"
$BIN \
  --out "$TEST_DIR/custom_output" \
  --langs "golang" \
  --format "{id}-{slug}.{ext}" \
  > "$TEST_DIR/stdout2.txt"

EXPECTED_FILES=$(find "$TEST_DIR/custom_output" -type f)
if [[ -z "$EXPECTED_FILES" ]]; then
  echo "❌ No output files written in custom_output"
  exit 1
fi

echo "✅ Passed Test 2"

echo "🧪 Test 3: Using a slug explicitly"
SLUG="two-sum"
$BIN --slug "$SLUG" --out "$TEST_DIR/slug_test" --langs "python3,golang" > "$TEST_DIR/stdout3.txt"
grep -qi "two sum" "$TEST_DIR/slug_test/"* || {
  echo "❌ Expected slug output file to contain 'two sum'"
  exit 1
}

echo "✅ Passed Test 3"

echo "🧪 Test 4: Config file fallback"
cat > "$CONFIG_FILE" <<EOF
{
  "outputDir": "$TEST_DIR/config_output",
  "filenameFormat": "{id}-{slug}.{ext}",
  "languages": ["python3", "golang"]
}
EOF

$BIN > "$TEST_DIR/stdout4.txt"

CONFIG_FILES=$(find "$TEST_DIR/config_output" -type f)
if [[ -z "$CONFIG_FILES" ]]; then
  echo "❌ Config-based output not created"
  exit 1
fi

echo "✅ Passed Test 4"

echo "🎉 All tests passed!"
exit 0