FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn

WORKDIR /opt/server

COPY go.mod .
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/server/main.go

FROM alpine:latest
WORKDIR  /
COPY --from=builder /opt/server/server .

# ENV SERVER_MODE="dev" 
# ENV CONSUL_ADDR="192.168.0.223:8500"

ENTRYPOINT ["/server"]