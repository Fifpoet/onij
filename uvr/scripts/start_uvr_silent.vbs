' 后台启动 uvr-api.exe，不显示终端窗口（供计划任务 / install_autostart 使用）
Set fso = CreateObject("Scripting.FileSystemObject")
Set shell = CreateObject("WScript.Shell")
dir = fso.GetParentFolderName(WScript.ScriptFullName)
shell.CurrentDirectory = dir

envExample = dir & "\.env.example"
envFile = dir & "\.env"
If Not fso.FileExists(envFile) And fso.FileExists(envExample) Then
  fso.CopyFile envExample, envFile, True
End If

exe = dir & "\uvr-api.exe"
If Not fso.FileExists(exe) Then
  MsgBox "未找到 uvr-api.exe", vbCritical, "ONIJ UVR API"
  WScript.Quit 1
End If

shell.Run Chr(34) & exe & Chr(34), 0, False
