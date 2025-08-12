FROM golang:1.24.5-alpine

# benötigte Tools installieren
RUN apk add --no-cache git curl

# Arbeitsverzeichnis setzen
WORKDIR /app

# Quellcode ins Image kopieren
COPY . .


# Air installieren
RUN go install github.com/air-verse/air@v1.62.0

# Starte Air zum Hot Reload
CMD ["air"]




