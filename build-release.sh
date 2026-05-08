#!/usr/bin/env sh
# 在本机（需 Node/pnpm + Go）执行；产物在 release/，可 git add 提交后服务器 pull 即可运行，无需在云机装前端工具链。
# ARM 云主机：TARGET_GOARCH=arm64 sh build-release.sh
set -e
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

GOARCH="${TARGET_GOARCH:-amd64}"

mkdir -p release/web release/bin

if command -v pnpm >/dev/null 2>&1; then
  pnpm run build:vite
else
  npm run build:vite
fi

# 清空旧静态资源后写入本次构建
rm -rf release/web/*
cp -R dist/. release/web/

cd server
CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build -trimpath -ldflags="-s -w" -o "$ROOT/release/bin/onij-server" .

echo "Done: $ROOT/release/web (frontend), $ROOT/release/bin/onij-server (backend, linux/$GOARCH)"
