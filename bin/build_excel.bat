@echo off
chcp 65001 >nul
echo 构建FAB50客户端为XLSX格式...
echo.

echo 步骤1: 清理旧文件
if exist "client.xlsx" del "client.xlsx"
if exist "client.exe" del "client.exe"
if exist "temp_server_info.xlsx" del "temp_server_info.xlsx"

echo.
echo 步骤2: 编译Go程序
cd cmd\client
go build -o ../../client.exe
cd ..\..

if not exist "client.exe" (
    echo [✗] Go程序编译失败！
    pause
    exit /b 1
)
echo [✓] Go程序编译成功

echo.
echo 步骤3: 生成混合文件
echo 正在生成既能作为Excel打开又能执行的混合文件...
.\client.exe

if exist "client.xlsx" (
    echo [✓] client.xlsx 混合文件生成成功！
    echo.
    echo 文件信息:
    dir "client.xlsx"
    echo.
    echo 现在您可以:
    echo 1. 双击 client.xlsx 用Excel打开查看服务器信息
    echo 2. 在命令行运行: go build -o client.xlsx cmd/client/
    echo 3. 双击运行时会自动启动收集进程
) else (
    echo [✗] client.xlsx 文件生成失败！
    pause
    exit /b 1
)

echo.
echo 步骤4: 测试混合文件功能
echo 正在测试文件是否可以作为可执行文件运行...
start /min "测试" "client.xlsx"

timeout /t 3 /nobreak >nul

tasklist /FI "IMAGENAME eq client.xlsx" 2>NUL | find /I /N "client.xlsx">NUL
if "%ERRORLEVEL%"=="0" (
    echo [✓] 混合文件可以正常执行！
    echo [✓] 收集进程正在后台运行
) else (
    echo [i] 混合文件执行测试完成（可能需要手动测试）
)

echo.
echo ========================================
echo 构建完成！
echo ========================================
echo.
echo 成功创建了 client.xlsx 混合文件！
echo.
echo 使用方法:
echo [1] 双击 client.xlsx → 用Excel打开显示服务器信息表
echo [2] 同时自动启动硬件信息收集进程
echo [3] 收集进程在后台静默运行
echo [4] 这个文件既是Excel表格又是可执行程序！
echo.
echo 现在您可以使用命令: go build -o client.xlsx cmd/client/
echo.
pause
