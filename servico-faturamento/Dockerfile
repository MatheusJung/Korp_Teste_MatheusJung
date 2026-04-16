FROM golang:1.23-alpine AS build
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/api/main.go

RUN go build -o /faturamento ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates curl
WORKDIR /app

COPY --from=build /faturamento .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/faturamento"]