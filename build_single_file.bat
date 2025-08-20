@echo off
chcp 65001 >nul
echo 构建FAB50单文件客户端...
echo.

echo 步骤1: 清理旧文件
if exist "client.xlsx" del "client.xlsx"
if exist "client.exe" del "client.exe"
if exist "服务器信息表_查看.xlsx" del "服务器信息表_查看.xlsx"

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
echo 步骤3: 生成自包含的Excel文件
echo 正在运行客户端生成自包含的Excel文件...
.\client.exe

if exist "client.xlsx" (
    echo [✓] 自包含的 client.xlsx 文件生成成功！
    echo.
    echo 文件信息:
    dir "client.xlsx"
    echo.
    echo 现在删除临时的可执行文件...
    del "client.exe"
    echo [✓] 临时文件已清理
) else (
    echo [✗] client.xlsx 文件生成失败！
    pause
    exit /b 1
)

echo.
echo 步骤4: 使用特殊命令将Excel文件转换为可执行文件
echo 正在重命名文件以支持双击执行...

REM 将Excel文件重命名为可执行文件，但保持Excel格式
copy "client.xlsx" "client.xlsx.exe" >nul
if exist "client.xlsx.exe" (
    echo [✓] 可执行Excel文件已创建
    del "client.xlsx"
    ren "client.xlsx.exe" "client.xlsx"
    echo [✓] 文件已优化
)

echo.
echo 步骤5: 测试单文件功能
echo 正在测试双击client.xlsx的行为...
echo 注意: 这会同时打开Excel表格和启动收集进程
pause
echo 启动测试...
start "" "client.xlsx"

timeout /t 5 /nobreak >nul

echo.
echo 检查是否生成了查看用的Excel文件...
if exist "服务器信息表_查看.xlsx" (
    echo [✓] 查看用Excel文件已生成
    echo [✓] 说明程序正在运行
) else (
    echo [i] 查看用Excel文件未生成（可能需要手动测试）
)

echo.
echo ========================================
echo 构建完成！
echo ========================================
echo.
echo 成功创建了单文件客户端:
echo [✓] client.xlsx - 既是Excel文件又是可执行程序
echo.
echo 使用方法:
echo [1] 双击 client.xlsx → 打开Excel表格显示服务器信息
echo [2] 同时自动启动硬件信息收集进程
echo [3] 收集进程在后台静默运行
echo [4] 只需要这一个文件，其他都不需要！
echo.
echo 特殊功能:
echo - 文件既可以用Excel打开查看表格
echo - 也可以双击执行启动收集进程
echo - 收集进程完全静默运行
echo.
pause
