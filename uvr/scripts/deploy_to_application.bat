@echo off
setlocal
cd /d "%~dp0\.."

set "SRC=%CD%\dist\uvr-api"
set "DST=C:\Users\onij\Downloads\Application\uvr-api"

if not exist "%SRC%\uvr-api.exe" (
  echo [deploy] 未找到 %SRC%\uvr-api.exe，请先 scripts\build_exe.bat
  exit /b 1
)

echo [deploy] 停止可能正在运行的 uvr-api...
taskkill /F /IM uvr-api.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo [deploy] 同步到 %DST% ...
if not exist "%DST%" mkdir "%DST%"
robocopy "%SRC%" "%DST%" /MIR /R:2 /W:1 /NFL /NDL /NJH /NJS /nc /ns /np
set "RC=%ERRORLEVEL%"
if %RC% GEQ 8 (
  echo [deploy] robocopy 失败 code=%RC%
  exit /b 1
)

if exist .env.example copy /Y .env.example "%DST%\.env.example" >nul
if not exist "%DST%\.env" if exist .env.example copy /Y .env.example "%DST%\.env" >nul

echo [deploy] 完成: %DST%
echo [deploy] 可运行: %DST%\start_uvr.bat
exit /b 0
