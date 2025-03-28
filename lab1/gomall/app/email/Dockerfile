# Auto-generated for email service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/email && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/email

CMD ["bash", "output/bootstrap.sh"]
