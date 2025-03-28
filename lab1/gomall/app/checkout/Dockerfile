# Auto-generated for checkout service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/checkout && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/checkout

CMD ["bash", "output/bootstrap.sh"]
