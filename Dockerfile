FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copier les fichiers du module avant tout
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code
COPY . .

# Compiler le binaire
RUN go build -o app .

# Étape finale
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8000
CMD ["./app"]
