@echo off
setlocal
cd /d "%~dp0"

if not exist "uvr-api.exe" (
  echo 请在 dist\uvr-api 目录运行，或先执行 scripts\build_exe.bat
  exit /b 1
)

if not exist ".env" if exist ".env.example" copy /Y ".env.example" ".env" >nul

set "WORK_DIR=C:\Users\onij\Downloads\Workbench\data\uvr"
if not exist "%WORK_DIR%" mkdir "%WORK_DIR%"

echo [uvr-api] http://0.0.0.0:5555  data=%WORK_DIR%
echo [uvr-api] 关闭本窗口即停止服务
uvr-api.exe
if errorlevel 1 (
  echo.
  echo [uvr-api] 启动失败，请确认在 dist\uvr-api 目录运行（勿用 build\uvr-api）
  pause
)
