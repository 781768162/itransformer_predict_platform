FROM golang:1.24-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway ./applications

FROM alpine:3.19
WORKDIR /app

RUN mkdir -p /app/tmp/uploads
VOLUME ["/app/tmp/uploads"]

COPY --from=builder /app/gateway /app/gateway

EXPOSE 8080

CMD ["/app/gateway"]
