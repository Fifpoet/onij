@echo off
setlocal
cd /d "%~dp0"

set "START=%~dp0start_uvr.bat"
if not exist "%START%" (
  echo 未找到 start_uvr.bat
  exit /b 1
)

set "TASK=ONIJ UVR API"
schtasks /Delete /TN "%TASK%" /F >nul 2>&1
schtasks /Create /TN "%TASK%" /TR "\"%START%\"" /SC ONLOGON /RL HIGHEST /F
if errorlevel 1 (
  echo 创建计划任务失败，请右键「以管理员身份运行」
  exit /b 1
)

echo 已注册开机自启（用户登录时）: %TASK%
echo 启动脚本: %START%
echo 取消自启: uninstall_autostart.bat
exit /b 0
