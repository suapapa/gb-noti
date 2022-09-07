# build stage
FROM golang:1.19 as builder

ENV CGO_ENABLED=0

RUN apt-get -qq update && \
    apt-get install -yqq upx

COPY . /build
WORKDIR /build

ARG BUILD_TIME=unknown
ARG GITHASH=unknown
ARG BUILD_TAG=dev

RUN go build \
    -ldflags "-X main.buildStamp=${BUILD_TIME} -X main.gitHash=${GITHASH} -X main.buildTag=${BUILD_TAG}"
RUN strip /build/gb-noti
RUN upx -q -9 /build/gb-noti

# ---
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/gb-noti .

ENV MQTT_USERNAME=secret
ENV MQTT_PASSWORD=secret
ENV TELEGRAM_APITOKEN=secret
ENV TELEGRAM_ROOM_ID=secret

EXPOSE 8080

ENTRYPOINT ["./gb-noti"]
