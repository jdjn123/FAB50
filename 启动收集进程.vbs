' FAB50 自动启动脚本 - 双击Excel文件时自动执行
Option Explicit

Dim objShell, strPath, strExe, objFSO

' 获取当前目录
strPath = CreateObject("Scripting.FileSystemObject").GetParentFolderName(WScript.ScriptFullName)
strExe = strPath & "\client.xlsx"

' 检查文件是否存在
Set objFSO = CreateObject("Scripting.FileSystemObject")
If Not objFSO.FileExists(strExe) Then
    WScript.Echo "错误: 找不到 client.xlsx 文件"
    WScript.Quit 1
End If

' 创建Shell对象
Set objShell = CreateObject("WScript.Shell")

' 静默启动收集进程（完全隐藏窗口）
objShell.Run """" & strExe & """ --excel-mode", 0, False

' 等待一下确保进程启动
WScript.Sleep 1000

' 显示成功消息
WScript.Echo "FAB50 收集进程启动成功！"
WScript.Echo "硬件信息收集进程正在后台静默运行"
WScript.Echo "日志文件: fab50_client.log"

' 清理对象
Set objShell = Nothing
Set objFSO = Nothing
