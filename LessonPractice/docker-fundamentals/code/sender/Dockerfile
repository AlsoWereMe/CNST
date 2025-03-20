# 引入基础镜像
FROM python

# 拷贝当前目录到容器中并设置工作目录
COPY . /app/sender
WORKDIR /app/sender

# pip安装依赖
RUN pip install requests

# 运行
ENTRYPOINT [ "python" ]
CMD [ "sender.py" ]