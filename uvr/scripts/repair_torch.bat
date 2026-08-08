@echo off
setlocal
cd /d "%~dp0\.."

echo [repair] 请先确保已 Ctrl+C 停止 UVR 服务
echo [repair] 清理损坏的 torch...
if exist ".venv\Lib\site-packages\~orch" rmdir /s /q ".venv\Lib\site-packages\~orch"
if exist ".venv\Lib\site-packages\~orch-2.6.0+cpu.dist-info" rmdir /s /q ".venv\Lib\site-packages\~orch-2.6.0+cpu.dist-info"
if exist ".venv\Lib\site-packages\~unctorch" rmdir /s /q ".venv\Lib\site-packages\~unctorch"
if exist ".venv\Lib\site-packages\torch" rmdir /s /q ".venv\Lib\site-packages\torch"
if exist ".venv\Lib\site-packages\torchvision" rmdir /s /q ".venv\Lib\site-packages\torchvision"

call .venv\Scripts\pip.exe uninstall -y torch torchvision torchaudio 2>nul

echo [repair] 安装 CUDA 版 PyTorch，约 2.5GB...
call .venv\Scripts\pip.exe uninstall -y onnxruntime onnxruntime-gpu 2>nul
call .venv\Scripts\pip.exe install --default-timeout=300 -r requirements-gpu.txt
if errorlevel 1 (
  echo [repair] 安装失败
  exit /b 1
)
REM faster-whisper 可能再次拉入 CPU 版 onnxruntime，卸掉以免冲掉 GPU
call .venv\Scripts\pip.exe uninstall -y onnxruntime 2>nul
call .venv\Scripts\pip.exe install --default-timeout=300 "onnxruntime-gpu==1.20.2"
if errorlevel 1 (
  echo [repair] 恢复 onnxruntime-gpu 失败
  exit /b 1
)
call .venv\Scripts\pip.exe install "setuptools>=69,<81" "numpy>=1.26,<2.5"

call .venv\Scripts\python.exe scripts\check_deps.py
exit /b %errorlevel%
