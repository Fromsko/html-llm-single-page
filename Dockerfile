FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 安装GCC和SQLite开发包，支持CGO
RUN apk add --no-cache gcc musl-dev sqlite-dev
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o html-manager main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata sqlite
WORKDIR /app
COPY --from=builder /app/html-manager .
# 创建数据目录
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./html-manager"]
