# Auto-generated for frontend service
FROM golang:1.21.13-bullseye

WORKDIR /gomall

COPY . .

RUN cd app/frontend && \
    chmod +x build.sh && \
    ./build.sh

WORKDIR /gomall/app/frontend

CMD ["bash", "output/bootstrap.sh"]
