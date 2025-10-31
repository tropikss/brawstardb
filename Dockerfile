FROM golang:1.25.3-alpine AS builder
WORKDIR /app

# Copier et compiler le binaire Go
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

# Étape finale avec Python
FROM python:3.11-slim
WORKDIR /app

# Copier le binaire Go
COPY --from=builder /app/app .

# Copier le script Python
COPY brawlstar.py .


# Étape finale
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8000
CMD ["./app"]
