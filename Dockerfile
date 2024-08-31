FROM alpine

WORKDIR /root

COPY build/${PROJECT_NAME} /root/
COPY config/config.toml /root/config/

CMD ["/root/promisys", "-c", "/root/config/config.toml"]
