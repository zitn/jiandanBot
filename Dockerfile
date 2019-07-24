# 引入golang编译环境
FROM golang as gobuilder
# 切换路径
WORKDIR /bot
# 复制go mod 文件至容器内
COPY go.mod go.sum ./
# 单独下载所需mod为一层,避免每次编译重复拉取依赖
RUN go mod download
# 复制源文件进入容器
COPY . .
#编译
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .


# 第二层,引入alpine环境
FROM scratch
# 添加证书
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# 从第一层中复制已经编译好的二进制文件
COPY --from=gobuilder /bot/main .
# 将config.json复制进容器内
COPY config.json .
# 运行容器
CMD ["/main"]
