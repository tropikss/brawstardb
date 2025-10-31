# Étape de build Go
FROM golang:1.25.3-alpine AS builder
WORKDIR /app

# Copier les fichiers go.mod/go.sum et télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tous les fichiers Go
COPY . .

# Compiler le serveur principal
RUN go build -o server main.go

# Compiler le script brawlstar
RUN go build -o brawlstar brawlstar.go

# Étape finale : image Alpine légère
FROM alpine:latest
WORKDIR /app

# Installer bash si nécessaire
RUN apk add --no-cache bash

# Copier les binaires Go compilés
COPY --from=builder /app/server .
COPY --from=builder /app/brawlstar .

# Exposer le port du serveur
EXPOSE 8000

# Lancer le serveur par défaut
CMD ["./server"]
