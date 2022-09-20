FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn

WORKDIR /opt/server

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR  /
COPY --from=builder /opt/server/server .
COPY --from=builder /opt/server/configs ./configs

ENV ServerMode=":8889" ConsulAddr=":9110"

ENTRYPOINT ["/server"]