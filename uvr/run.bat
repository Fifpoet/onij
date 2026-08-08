@echo off
setlocal
cd /d "%~dp0"

if not exist ".venv\Scripts\python.exe" (
  echo [uvr] 创建虚拟环境...
  python -m venv .venv
  call .venv\Scripts\python.exe -m pip install -U pip
)

set "TORCH_REQ=requirements-cpu.txt"
where nvidia-smi >nul 2>&1
if not errorlevel 1 set "TORCH_REQ=requirements-gpu.txt"

echo [uvr] 检查依赖...
call .venv\Scripts\python.exe scripts\check_deps.py >nul 2>&1
if errorlevel 1 goto install_deps

call .venv\Scripts\python.exe -c "import torch,sys; sys.exit(0 if (('%TORCH_REQ%'=='requirements-cpu.txt') or torch.cuda.is_available()) else 1)" >nul 2>&1
if errorlevel 1 goto install_deps
goto start_server

:install_deps
if /i "%TORCH_REQ%"=="requirements-gpu.txt" (
  call scripts\repair_torch.bat
) else (
  echo [uvr] 安装依赖: %TORCH_REQ%
  call .venv\Scripts\pip.exe install -r %TORCH_REQ%
  call .venv\Scripts\python.exe scripts\check_deps.py
)
if errorlevel 1 (
  echo [uvr] 依赖仍不可用，请检查网络后重试
  exit /b 1
)
if /i "%TORCH_REQ%"=="requirements-gpu.txt" (
  call .venv\Scripts\pip.exe uninstall -y onnxruntime 2>nul
)

:start_server
set "UVR_ROOT=C:\Users\onij\Downloads\Application\Ultimate Vocal Remover"
set "UVR_EXE=%UVR_ROOT%\UVR.exe"
set "WORK_DIR=C:\Users\onij\Downloads\Workbench\data\uvr"
if not exist "%WORK_DIR%" mkdir "%WORK_DIR%"

if not defined HF_ENDPOINT set "HF_ENDPOINT=https://hf-mirror.com"
echo [uvr] 启动 API http://0.0.0.0:5555  data=%WORK_DIR%  HF_ENDPOINT=%HF_ENDPOINT%
call .venv\Scripts\python.exe -m app.main
