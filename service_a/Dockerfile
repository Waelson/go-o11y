FROM golang:1.22.1-alpine as builder

WORKDIR /app
COPY .. .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]