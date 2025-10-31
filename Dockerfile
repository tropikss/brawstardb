# Étape de build Go
FROM golang:1.25.3-alpine AS builder
WORKDIR /app

# Copier les fichiers go.mod/go.sum et télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code et compiler le binaire
COPY . .
RUN go build -o app .

# Étape finale
FROM alpine:latest
WORKDIR /app

# Installer Python et bash
RUN apk add --no-cache python3 py3-pip bash

# Copier le binaire Go
COPY --from=builder /app/app .

# Copier ton script Python
COPY brawlstar.py .

# Exposer le port du serveur Go
EXPOSE 8000

# Commande par défaut : lancer le serveur Go
CMD ["./app"]
