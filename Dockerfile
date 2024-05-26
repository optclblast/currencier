FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/currencier cmd/app/main.go

EXPOSE 8080

CMD ["/app/currencier", "-log-level=info","-http-port=8080","-cache-url=cache:6379"]