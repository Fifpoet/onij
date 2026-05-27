@echo off
setlocal
cd /d "%~dp0\.."

echo [portable] 检查依赖...
call .venv\Scripts\python.exe scripts\check_deps.py
if errorlevel 1 (
  echo [portable] 请先 run.bat 安装完整依赖
  exit /b 1
)

set "OUT=dist\uvr-portable"
if exist "%OUT%" rmdir /S /Q "%OUT%"
mkdir "%OUT%"
mkdir "%OUT%\app"

echo [portable] 复制应用文件...
xcopy /E /I /Y app "%OUT%\app" >nul
copy /Y launcher.py "%OUT%\launcher.py" >nul
if exist .env.example copy /Y .env.example "%OUT%\.env.example" >nul

echo [portable] 生成启动脚本...
(
  echo @echo off
  echo setlocal
  echo cd /d "%%~dp0"
  echo set "VENV=%%~dp0..\..\.venv"
  echo if not exist "%%VENV%%\Scripts\python.exe" (
  echo   echo 未找到 uvr\.venv，请先在 uvr 目录 run.bat
  echo   exit /b 1
  echo ^)
  echo "%%VENV%%\Scripts\python.exe" launcher.py %%*
) > "%OUT%\start_uvr.bat"

copy /Y scripts\install_autostart_portable.bat "%OUT%\install_autostart.bat" >nul
copy /Y scripts\uninstall_autostart.bat "%OUT%\uninstall_autostart.bat" >nul

echo.
echo [portable] 完成: %CD%\%OUT%
echo [portable] 运行: dist\uvr-portable\start_uvr.bat
echo [portable] 说明: 依赖 uvr\.venv，无需 PyInstaller，适合 ML 栈
exit /b 0
