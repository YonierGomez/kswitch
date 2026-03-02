#!/usr/bin/env bash
# Generates all demo GIFs for ksw README
# Requirements: vhs (brew install vhs)

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

if ! command -v vhs &>/dev/null; then
  echo "❌ vhs not found. Install with: brew install vhs"
  exit 1
fi

if ! command -v ksw &>/dev/null; then
  echo "❌ ksw not found. Install with: brew install ksw"
  exit 1
fi

export KUBECONFIG="$SCRIPT_DIR/kubeconfig-demo.yaml"

tapes=(demo-tui demo-ai demo-pins demo-groups demo-aliases demo-history demo-previous)

for tape in "${tapes[@]}"; do
  echo "🎬 Generating $tape..."
  vhs "$SCRIPT_DIR/$tape.tape"
done

echo ""
echo "✔ Done! GIFs saved to demo/"
