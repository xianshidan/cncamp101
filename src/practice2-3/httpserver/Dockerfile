#源镜像
FROM golang:latest
#作者
MAINTAINER danny "xianshidan@gmail.com"
#设置工作目录
WORKDIR $GOPATH/src/github.com/xianshidan/cncamp101/src/practice2-2/httpserver
#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/github.com/httpserver
#go构建可执行文件
RUN go build .
#暴露端口
EXPOSE 6064
#最终运行docker的命令
ENTRYPOINT  ["./httpserver"]