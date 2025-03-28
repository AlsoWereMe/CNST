# Auto-generated for payment service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/payment && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/payment

CMD ["bash", "output/bootstrap.sh"]
