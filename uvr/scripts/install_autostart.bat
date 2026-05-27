@echo off
setlocal
cd /d "%~dp0"

set "EXE=%~dp0uvr-api.exe"
if not exist "%EXE%" (
  echo 未找到 uvr-api.exe
  pause
  exit /b 1
)

set "TASK=ONIJ UVR API"
schtasks /Delete /TN "%TASK%" /F >nul 2>&1
schtasks /Create /TN "%TASK%" /TR "\"%EXE%\"" /SC ONLOGON /RL HIGHEST /F
if errorlevel 1 (
  echo 创建计划任务失败
  pause
  exit /b 1
)

echo 已注册: %TASK%
echo 程序: %EXE%
echo 取消: uninstall_autostart.bat
pause
