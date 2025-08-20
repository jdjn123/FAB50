@echo off
echo 测试FAB50客户端Excel功能...
echo.

echo 1. 检查客户端文件...
if exist "client.exe" (
    echo [✓] client.exe 存在
) else (
    echo [✗] client.exe 不存在
    goto :error
)

echo.
echo 2. 启动客户端（将自动打开服务器信息表格）...
echo 注意：客户端将在后台运行，收到服务端确认后自动删除
start "" "client.exe"

echo.
echo 3. 等待5秒让客户端启动...
timeout /t 5 /nobreak >nul

echo.
echo 4. 检查进程状态...
tasklist /FI "IMAGENAME eq client.exe" 2>NUL | find /I /N "client.exe">NUL
if "%ERRORLEVEL%"=="0" (
    echo [✓] 客户端进程正在运行
    echo.
    echo 进程信息:
    tasklist /FI "IMAGENAME eq client.exe" /FO TABLE
) else (
    echo [✗] 客户端进程未运行（可能已自删除）
)

echo.
echo 5. 检查临时Excel文件...
set "temp_dir=%TEMP%"
if exist "%temp_dir%\服务器信息表.csv" (
    echo [✓] 服务器信息表格已生成
    echo.
    echo 表格内容预览:
    type "%temp_dir%\服务器信息表.csv"
) else (
    echo [✗] 服务器信息表格未生成
)

echo.
echo ========================================
echo 测试完成！
echo ========================================
echo.
echo 如果服务器信息表格已打开，说明功能正常
echo 客户端正在后台运行，等待服务端确认
echo 收到确认后将自动删除自身
echo.
echo 按任意键退出...
pause >nul
goto :end

:error
echo.
echo 测试过程中出现错误，请检查上述问题
echo.

:end
