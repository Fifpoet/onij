#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$ROOT_DIR/server"

if [ ! -d "$BACKEND_DIR" ]; then
  echo "未找到后端目录: $BACKEND_DIR" >&2
  exit 1
fi

echo "[dev] 启动后端: go run . (目录: $BACKEND_DIR)"
(
  cd "$BACKEND_DIR"
  go run .
) &
BACKEND_PID=$!

cleanup() {
  echo "[dev] 正在停止后端进程..."
  if kill -0 "$BACKEND_PID" 2>/dev/null; then
    kill "$BACKEND_PID" 2>/dev/null || true
    wait "$BACKEND_PID" 2>/dev/null || true
  fi
  echo "[dev] 已清理后端进程"
}

trap cleanup EXIT INT TERM

echo "[dev] 启动前端: pnpm run dev"
cd "$ROOT_DIR"
pnpm run dev
