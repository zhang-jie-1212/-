#指定基础镜像
FROM golang:alpine AS builder

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#生成二进制文件
WORKDIR /build
COPY . .
RUN go build -o example.

# 创建小镜像运行二进制代码
FROM scratch

# 二进制文件copy到当前目录
COPY --from=builder /build/example /

# 声明服务端口
EXPOSE 8888

# 需要运行的命令
ENTRYPOINT ["/example"]