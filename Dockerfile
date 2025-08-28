# --------------------------- 构建阶段 ---------------------------
FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 根据你的项目结构，cd 到正确目录
RUN cd cmd/exam_api && go build -ldflags="-s -w" -o /app/exam_api .

# --------------------------- 运行阶段 ---------------------------
FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates
RUN adduser -D -s /bin/sh appuser
USER appuser

WORKDIR /home/appuser

# 复制二进制和配置文件
COPY --from=builder /app/exam_api .
#COPY --from=builder /app/configs ./configs

EXPOSE 8000

# 推荐使用 exec 格式
CMD ["./exam_api", "-conf", "./configs"]