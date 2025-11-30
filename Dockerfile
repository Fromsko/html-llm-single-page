FROM golang:1.21-alpine AS builder
WORKDIR /app

# 安装GCC和SQLite开发包，支持CGO
# 使用更快的镜像源和并行下载
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache --virtual .build-deps gcc musl-dev sqlite-dev git

# 设置 Go 构建缓存
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-mod

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

# 优化构建：使用并行编译和更好的缓存
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o html-manager main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata sqlite
WORKDIR /app
COPY --from=builder /app/html-manager .
# 创建数据目录
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./html-manager"]
