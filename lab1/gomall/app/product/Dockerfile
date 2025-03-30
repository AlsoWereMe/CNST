# Auto-generated for product service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/product && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/product

CMD ["bash", "output/bootstrap.sh"]
