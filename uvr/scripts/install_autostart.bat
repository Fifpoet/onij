@echo off
setlocal
cd /d "%~dp0"

set "EXE=%~dp0uvr-api.exe"
set "SILENT=%~dp0start_uvr_silent.vbs"
if not exist "%EXE%" (
  echo 未找到 uvr-api.exe，请先 build_exe.bat
  pause
  exit /b 1
)
if not exist "%SILENT%" (
  echo 未找到 start_uvr_silent.vbs
  pause
  exit /b 1
)

set "TASK=ONIJ UVR API"
schtasks /Delete /TN "%TASK%" /F >nul 2>&1
schtasks /Create /TN "%TASK%" /TR "wscript.exe \"%SILENT%\"" /SC ONLOGON /RL HIGHEST /F
if errorlevel 1 (
  echo 创建计划任务失败，请右键「以管理员身份运行」
  pause
  exit /b 1
)

echo 已注册开机自启（用户登录时）: %TASK%
echo 后台启动: %SILENT%
echo 手动调试（有终端）: start_uvr.bat
echo 取消自启: uninstall_autostart.bat
pause
exit /b 0
