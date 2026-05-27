@echo off
schtasks /Delete /TN "ONIJ UVR API" /F
echo 已移除计划任务 ONIJ UVR API
exit /b 0
