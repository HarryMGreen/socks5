ARG GOLANG_VERSION="1.25"

FROM golang:$GOLANG_VERSION-alpine AS builder
RUN apk --no-cache add tzdata
WORKDIR /tmp/socks5
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o socks5 docker/main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /tmp/socks5/socks5 /
ENTRYPOINT ["/socks5"]
