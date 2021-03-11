FROM golang:1.16 as build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY server.go .
RUN go build -o server

FROM debian:10-slim
WORKDIR /app
COPY --from=build /build/server .
CMD ["./server"]
