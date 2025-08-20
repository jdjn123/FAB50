@echo off
chcp 65001 >nul
echo FAB50 硬件信息收集进程启动器
echo ================================
echo.

REM 获取当前目录
set "CURRENT_DIR=%~dp0"
set "EXE_FILE=%CURRENT_DIR%client.xlsx"

REM 检查文件是否存在
if not exist "%EXE_FILE%" (
    echo 错误: 找不到 client.xlsx 文件
    echo 请确保此批处理文件与 client.xlsx 在同一目录
    echo.
    pause
    exit /b 1
)

echo 正在启动硬件信息收集进程...
echo 收集进程将在后台静默运行
echo 日志文件: fab50_client.log
echo.

REM 静默启动收集进程
start /min "" "%EXE_FILE%" --excel-mode

echo 收集进程已启动！
echo 请查看 fab50_client.log 文件确认运行状态
echo.
echo 按任意键退出...
pause >nul
