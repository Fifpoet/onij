@echo off
cd /d "%~dp0"

if not defined UVR_ROOT set "UVR_ROOT=C:\Users\onij\Downloads\Application\Ultimate Vocal Remover"
if not defined UVR_EXE set "UVR_EXE=%UVR_ROOT%\UVR.exe"
if not defined WORK_DIR set "WORK_DIR=C:\Users\onij\Downloads\Workbench\data\uvr"
if not defined OUTPUT_DIR set "OUTPUT_DIR=%WORK_DIR%\_output"
if not defined HF_ENDPOINT set "HF_ENDPOINT=https://hf-mirror.com"
if not defined HF_HUB_DISABLE_XET set "HF_HUB_DISABLE_XET=1"
if not defined WHISPER_MODEL set "WHISPER_MODEL=small"
if not defined WHISPER_DEVICE set "WHISPER_DEVICE=cuda"
if not defined WHISPER_COMPUTE_TYPE set "WHISPER_COMPUTE_TYPE=float16"

if not exist "%WORK_DIR%" mkdir "%WORK_DIR%"
uvr-api.exe
