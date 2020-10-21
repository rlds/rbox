FROM alpine:latest

# 建立需要的文件夹
RUN mkdir {/var/www,/var/log}

# 添加基础文件
ADD ./csbox /usr/local/bin/
ADD static /var/www

# 对外访问的端口
EXPOSE 8080

# 启动进程
CMD ["/usr/local/bin/csbox" -host ":8080" -log "/var/log" -static "/var/www/static"]
