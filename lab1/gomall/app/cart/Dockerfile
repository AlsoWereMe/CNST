# Auto-generated for cart service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/cart && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/cart

CMD ["bash", "output/bootstrap.sh"]
