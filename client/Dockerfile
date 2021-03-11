FROM golang:1.16 as build
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY client.go .
RUN go build -o client

FROM debian:10-slim
WORKDIR /app
COPY --from=build /build/client .
