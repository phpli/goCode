#基础镜像
FROM ubuntu:latest
LABEL authors="changdong"
#把编译后的打包进来的这个镜像放到工作目录/app
COPY webook /app/webook
WORKDIR /app
RUN chmod +x /app/webook
ENTRYPOINT ["./webook"]