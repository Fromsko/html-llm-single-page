#!/bin/bash

# HTML Manager 构建脚本
# 支持多平台构建

set -e

# 版本信息
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建信息
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# 输出目录
OUTPUT_DIR="dist"
mkdir -p ${OUTPUT_DIR}

echo "开始构建 HTML Manager..."
echo "版本: ${VERSION}"
echo "构建时间: ${BUILD_TIME}"
echo "Git 提交: ${GIT_COMMIT}"
echo ""

# 构建函数
build() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "构建 ${os}/${arch}..."
    
    # SQLite需要CGO支持，但在交叉编译时需要特殊处理
    if [ "${os}" = "linux" ]; then
        # Linux平台启用CGO支持SQLite
        GOOS=${os} GOARCH=${arch} go build \
            -ldflags="${LDFLAGS} -s -w" \
            -o ${OUTPUT_DIR}/html-manager-${os}-${arch}${ext} \
            .
    else
        # 其他平台禁用CGO，使用纯Go的SQLite实现
        GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build \
            -tags=sqlite_omit_load_extension \
            -ldflags="${LDFLAGS} -s -w" \
            -o ${OUTPUT_DIR}/html-manager-${os}-${arch}${ext} \
            .
    fi
    
    echo "✓ ${os}/${arch} 构建完成"
}

# 构建各个平台
build linux amd64 ""
build linux arm64 ""
build darwin amd64 ""
build darwin arm64 ""
build windows amd64 ".exe"

# 创建压缩包
echo ""
echo "创建压缩包..."
cd ${OUTPUT_DIR}

for file in html-manager-*; do
    if [[ $file == *.exe ]]; then
        7z a -tzip "${file%.exe}.zip" "$file" >/dev/null
        echo "✓ 创建 ${file%.exe}.zip"
    else
        tar -czf "${file}.tar.gz" "$file"
        echo "✓ 创建 ${file}.tar.gz"
    fi
done

cd ..

echo ""
echo "构建完成！输出目录: ${OUTPUT_DIR}"
echo ""
echo "文件列表:"
ls -la ${OUTPUT_DIR}/