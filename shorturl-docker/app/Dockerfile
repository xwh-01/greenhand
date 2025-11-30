# 阶段1：构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 阶段2：运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户（安全最佳实践）
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser/

# 从构建阶段复制二进制文件
COPY --from=builder --chown=appuser:appuser /app/main .

# 切换到非root用户
USER appuser

EXPOSE 8080

ENV MYSQL_DSN="root:123456@tcp(mysql:3306)/short_url?charset=utf8mb4&parseTime=True&loc=Local"

CMD ["./main"]
