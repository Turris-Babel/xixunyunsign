# 使用 Apache 官方镜像（基于 Debian）
FROM httpd:latest

# 更新系统并安装 glibc
RUN apt-get update && apt-get install -y libc6 && apt-get clean

# 复制可执行文件到容器内
COPY output/xixunyunsign /usr/local/bin/

# 确保文件有可执行权限
RUN chmod +x /usr/local/bin/xixunyunsign

# 使用默认的 Apache CMD，可以通过额外命令运行二进制文件

CMD ["sh","-c","xixunyunsign","serve"]
