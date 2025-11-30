# 阶段1：构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用 - 确保输出文件名与CMD一致
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shorturl-app .

# 阶段2：运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/shorturl-app .

EXPOSE 8080

CMD ["./shorturl-app"]
