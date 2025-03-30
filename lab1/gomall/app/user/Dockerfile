# Auto-generated for user service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/user && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/user

CMD ["bash", "output/bootstrap.sh"]
