FROM golang:1.18.3-alpine AS builder

COPY . /github.com/iamtonydev/url-shortener/
WORKDIR /github.com/iamtonydev/url-shortener/

RUN go mod download
RUN go build -o ./bin/url_shortener cmd/url_shortener/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/iamtonydev/url-shortener/bin/url_shortener .
COPY --from=builder /github.com/iamtonydev/url-shortener/config.yml .

EXPOSE 50051
EXPOSE 8000

CMD ["./url_shortener"]
