#!/usr/bin/env sh
# 在仓库根目录执行：产出 release/web（静态前端）、release/bin/onij-server（Linux amd64 后端二进制）
set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

mkdir -p release/web release/bin

if command -v pnpm >/dev/null 2>&1; then
  pnpm run build
else
  npm run build
fi

# 清空旧静态资源后写入本次构建
rm -rf release/web/*
cp -R dist/. release/web/

cd server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o "$ROOT/release/bin/onij-server" .

echo "Done: $ROOT/release/web (frontend), $ROOT/release/bin/onij-server (backend, linux/amd64)"
