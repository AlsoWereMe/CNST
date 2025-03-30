# Auto-generated for order service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/order && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/order

CMD ["bash", "output/bootstrap.sh"]
