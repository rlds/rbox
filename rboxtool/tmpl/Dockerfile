FROM alpine:latest

# 日志存储文件夹
RUN mkdir {/var/{{.BoxConf.Name}},/var/{{.BoxConf.Name}}/{{.BoxConf.LogDir}} }
ADD ./{{.BoxConf.Name}} /var/{{.BoxConf.Name}}/
# 对外可访问端口
EXPOSE {{.EXPOSE}}

# 启动进程
CMD ["/var/{{.BoxConf.Name}}","-mode","http"]
