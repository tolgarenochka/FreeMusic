FROM golang:alpine AS builder

WORKDIR /FreeMusic

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

COPY --from=builder /FreeMusic/main /build/main
COPY --from=builder /FreeMusic/configs/local_config.json /configs/local_config.json

CMD ["./build/main", "-config-path", "/configs/local_config.json"]