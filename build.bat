@echo off
REM HTML Manager 构建脚本 (Windows)
REM 支持多平台构建

setlocal enabledelayedexpansion

REM 获取版本信息
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul || echo dev') do set VERSION=%%i
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul || echo unknown') do set GIT_COMMIT=%%i

REM 获取当前时间
for /f "tokens=1-3 delims=/ " %%a in ('date /t') do set BUILD_DATE=%%c-%%a-%%b
for /f "tokens=1-2 delims=: " %%a in ('time /t') do set BUILD_TIME=%%a:%%b
set BUILD_TIME=%BUILD_DATE%_%BUILD_TIME%

REM 构建信息
set LDFLAGS=-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%

REM 输出目录
set OUTPUT_DIR=dist
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

echo 开始构建 HTML Manager...
echo 版本: %VERSION%
echo 构建时间: %BUILD_TIME%
echo Git 提交: %GIT_COMMIT%
echo.

REM 构建函数
:build
set os=%1
set arch=%2
set ext=%3

echo 构建 %os%/%arch%...

set GOOS=%os%
set GOARCH=%arch%
REM SQLite需要CGO支持，所以不能设置CGO_ENABLED=0
REM set CGO_ENABLED=0
go build -ldflags="%LDFLAGS% -s -w" -o %OUTPUT_DIR%/html-manager-%os%-%arch%%ext% .

echo ✓ %os%/%arch% 构建完成
goto :eof

REM 构建各个平台
call :build linux amd64 ""
call :build linux arm64 ""
call :build darwin amd64 ""
call :build darwin arm64 ""
call :build windows amd64 ".exe"

echo.
echo 创建压缩包...
cd %OUTPUT_DIR%

REM 创建ZIP文件
for %%f in (html-manager-*.exe) do (
    powershell -command "Compress-Archive -Path '%%f' -DestinationPath '%%~nf.zip' -Force"
    echo ✓ 创建 %%~nf.zip
)

for %%f in (html-manager-*) do (
    if not %%f == *.exe (
        if not exist %%f.tar.gz (
            tar -czf %%f.tar.gz %%f
            echo ✓ 创建 %%f.tar.gz
        )
    )
)

cd ..

echo.
echo 构建完成！输出目录: %OUTPUT_DIR%
echo.
echo 文件列表:
dir %OUTPUT_DIR%

pause