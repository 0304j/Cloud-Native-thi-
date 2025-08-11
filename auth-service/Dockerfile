# Basis-Image
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Go-Module kopieren und Abhängigkeiten laden
COPY go.mod go.sum ./
RUN go mod download

# Restlichen Code kopieren
COPY . .

# Binary bauen
RUN go build -o auth-service ./cmd

# Runtime-Image
FROM alpine:latest
WORKDIR /root/

# Binary kopieren
COPY --from=builder /app/auth-service .

# Port freigeben
EXPOSE 8081

# Startbefehl
CMD ["./auth-service"]
