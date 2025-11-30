# HTML Page Manager
AI 时代单页面 HTML 管理系统

## 功能特性

- 🚀 **拖拽上传**: 支持拖拽上传 HTML 文件
- ✏️ **在线编辑**: 集成 Monaco 编辑器，支持在线编辑 HTML 页面
- 📱 **响应式设计**: 适配各种设备屏幕
- 🗄️ **SQLite 存储**: 使用 SQLite 数据库存储页面数据
- 📥 **反向下载**: 支持将页面下载为原始 HTML 文件
- 🔧 **环境配置**: 支持 .env 文件配置环境变量
- 🌐 **双模式支持**: 同时支持 API 调用和 Web 端管理
- 🛣️ **路由展示**: 每个页面都有独立的访问路由

## 快速开始

### 1. 环境要求

- Go 1.21+
- Git

### 2. 安装依赖

```bash
cd html-manager
go mod tidy
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
# 服务器配置
PORT=8080
GIN_MODE=release

# 网站信息
SITE_NAME=HTML Page Manager
AUTHOR_NAME=Your Name

# 支持的站点（逗号分隔）
SUPPORTED_SITES=github.com,gitea.com,gitlab.com

# 数据库配置
DB_PATH=./data/pages.db
```

### 4. 运行项目

```bash
go run main.go
```

访问 http://localhost:8080 开始使用。

## API 文档

### 页面管理

#### 获取所有页面
```
GET /api/pages
```

#### 获取单个页面
```
GET /api/pages/:slug
```

#### 创建页面
```
POST /api/pages
Content-Type: application/json

{
  "title": "页面标题",
  "content": "HTML内容",
  "description": "页面描述（可选）"
}
```

#### 更新页面
```
PUT /api/pages/:id
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "description": "更新后的描述"
}
```

#### 删除页面
```
DELETE /api/pages/:id
```

### 文件操作

#### 上传文件
```
POST /api/upload
Content-Type: multipart/form-data

file: HTML文件
```

#### 下载文件
```
GET /api/download/:id
```

### 页面访问

#### 访问页面
```
GET /page/:slug
```

## 项目结构

```
html-manager/
├── config/          # 配置文件
│   └── config.go
├── handlers/        # 路由处理器
│   └── api.go
├── models/          # 数据模型
│   └── page.go
├── static/          # 静态文件
│   └── custom.css
├── templates/       # HTML 模板
│   └── admin.html
├── data/            # 数据目录（自动创建）
│   └── pages.db     # SQLite 数据库
├── .env.example     # 环境变量示例
├── go.mod           # Go 模块文件
├── main.go          # 主程序入口
└── README.md        # 项目说明
```

## 部署

### 使用 Docker 部署

1. 创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./main"]
```

2. 构建和运行:

```bash
docker build -t html-manager .
docker run -p 8080:8080 -e GIN_MODE=release html-manager
```

### 直接部署

```bash
# 编译
go build -o html-manager main.go

# 运行
./html-manager
```

## 构建和发布

### 本地构建

项目提供了构建脚本，支持多平台构建：

#### Linux/macOS
```bash
chmod +x build.sh
./build.sh
```

#### Windows
```cmd
build.bat
```

构建完成后，所有二进制文件和压缩包将在 `dist/` 目录中。

### 自动构建和发布

项目使用 GitHub Actions 进行自动化构建和发布：

1. **触发条件**：
   - 推送标签（如 `v1.0.0`）时自动构建并发布 Release
   - 推送到 main/master 分支时运行测试

2. **支持的平台**：
   - Linux (amd64, arm64)
   - Windows (amd64)
   - macOS (amd64, arm64)

3. **发布内容**：
   - 各平台的二进制文件压缩包
   - Docker 镜像（多架构支持）

4. **Docker 镜像**：
   - 镜像名：`html-manager/html-manager`
   - 标签：版本号和 `latest`

### 版本管理

项目使用语义化版本控制（Semantic Versioning）：

- 主版本号：不兼容的 API 修改
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

发布新版本：

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 使用说明

### 1. 上传 HTML 文件

- 在管理界面拖拽 HTML 文件到上传区域
- 或点击选择文件按钮选择文件

### 2. 创建新页面

- 输入页面标题和描述
- 在 Monaco 编辑器中编写 HTML 代码
- 点击"创建页面"保存

### 3. 编辑现有页面

- 在页面列表中点击"编辑"按钮
- 在弹出的编辑器中修改内容
- 点击"保存"更新页面

### 4. 访问页面

- 通过 `/page/:slug` 路径访问页面
- 例如: `http://localhost:8080/page/my-awesome-page`

### 5. 下载页面

- 在页面列表中点击"下载"链接
- 或通过 API `/api/download/:id` 下载

## 配置说明

| 环境变量        | 说明                         | 默认值                          |
| --------------- | ---------------------------- | ------------------------------- |
| PORT            | 服务器端口                   | 8080                            |
| GIN_MODE        | Gin 运行模式 (debug/release) | debug                           |
| SITE_NAME       | 网站名称                     | HTML Page Manager               |
| AUTHOR_NAME     | 作者名称                     | Your Name                       |
| SUPPORTED_SITES | 支持的站点列表               | github.com,gitea.com,gitlab.com |
| DB_PATH         | 数据库文件路径               | ./data/pages.db                 |

## 许可证

MIT License
