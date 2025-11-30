#!/bin/bash
echo "测试跨平台构建..."

# 创建测试输出目录
TEST_OUTPUT_DIR="test-dist"
mkdir -p ${TEST_OUTPUT_DIR}

# 测试 Linux 构建（启用CGO）
echo "测试 Linux 构建（启用CGO）..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${TEST_OUTPUT_DIR}/html-manager-linux-amd64 .
if [ $? -eq 0 ]; then
    echo "✓ Linux 构建成功"
else
    echo "✗ Linux 构建失败"
fi

# 测试 macOS 构建（禁用CGO）
echo "测试 macOS 构建（禁用CGO）..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags=sqlite_omit_load_extension -ldflags="-s -w" -o ${TEST_OUTPUT_DIR}/html-manager-darwin-amd64 .
if [ $? -eq 0 ]; then
    echo "✓ macOS 构建成功"
else
    echo "✗ macOS 构建失败"
fi

# 测试 Windows 构建（禁用CGO）
echo "测试 Windows 构建（禁用CGO）..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags=sqlite_omit_load_extension -ldflags="-s -w" -o ${TEST_OUTPUT_DIR}/html-manager-windows-amd64.exe .
if [ $? -eq 0 ]; then
    echo "✓ Windows 构建成功"
else
    echo "✗ Windows 构建失败"
fi

echo ""
echo "测试构建完成！输出目录: ${TEST_OUTPUT_DIR}"
ls -la ${TEST_OUTPUT_DIR}/

# 清理测试文件
# rm -rf ${TEST_OUTPUT_DIR}
