# 在仓库根目录执行（PowerShell）：产出 release/web、release/bin/onij-server（默认 Linux amd64）
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $Root

$GoArch = if ($env:TARGET_GOARCH) { $env:TARGET_GOARCH } else { "amd64" }

New-Item -ItemType Directory -Force -Path "$Root\release\web", "$Root\release\bin" | Out-Null

if (Get-Command pnpm -ErrorAction SilentlyContinue) {
  pnpm run build:vite
} else {
  npm run build:vite
}

Get-ChildItem "$Root\release\web" -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force
Copy-Item -Path "$Root\dist\*" -Destination "$Root\release\web" -Recurse -Force

Push-Location "$Root\server"
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = $GoArch
go build -trimpath -ldflags="-s -w" -o "$Root/release/bin/onij-server" .
Pop-Location

Write-Host "Done: $Root\release\web , $Root\release\bin\onij-server (linux/$GoArch)"
