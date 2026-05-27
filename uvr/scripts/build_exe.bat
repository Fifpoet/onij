@echo off
setlocal
cd /d "%~dp0\.."

echo [build] 检查依赖...
call .venv\Scripts\python.exe scripts\check_deps.py
if errorlevel 1 (
  echo [build] 请先 run.bat 安装完整依赖
  exit /b 1
)

call .venv\Scripts\pip.exe install pyinstaller
if errorlevel 1 exit /b 1

echo [build] PyInstaller 打包中，约需数分钟...
call .venv\Scripts\pyinstaller.exe --noconfirm --clean uvr-api.spec
if errorlevel 1 (
  echo.
  echo [build] PyInstaller 失败（常见于 onnx/torch 分析阶段崩溃）
  echo [build] 可改用便携方案: scripts\build_portable.bat
  exit /b 1
)

set "DIST=dist\uvr-api"
if not exist "%DIST%" (
  echo [build] 未找到输出目录 %DIST%
  exit /b 1
)

if exist .env.example copy /Y .env.example "%DIST%\.env.example" >nul
copy /Y scripts\start_uvr.bat "%DIST%\start_uvr.bat" >nul
copy /Y scripts\start_uvr_silent.vbs "%DIST%\start_uvr_silent.vbs" >nul
copy /Y scripts\install_autostart.bat "%DIST%\install_autostart.bat" >nul
copy /Y scripts\uninstall_autostart.bat "%DIST%\uninstall_autostart.bat" >nul

echo.
echo [build] 完成: %CD%\%DIST%
echo [build] 请运行 dist\uvr-api\start_uvr.bat 或 dist\uvr-api\uvr-api.exe
echo [build] 注意: build\uvr-api\ 是中间目录，不能直接双击其中的 exe
exit /b 0
