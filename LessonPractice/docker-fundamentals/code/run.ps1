docker network create mynet

# 构建镜像
docker build -t receiver:v1 -f ./receiver/Dockerfile ./receiver
docker build -t sender:v1 -f ./sender/Dockerfile ./sender

docker run -d -p 8081:8081 --network mynet --name receiver -v $pwd/file:/file receiver:v1
# 仿照receiver的例子运行sender
docker run -d -p 8080:8080 --network mynet --name sender -v $pwd/file:/file sender:v1