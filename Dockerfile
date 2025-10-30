FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod init myapi && go mod tidy
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8000
CMD ["./app"]
