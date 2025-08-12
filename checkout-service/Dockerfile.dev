# Vorher (veraltet)
# FROM golang:1.21
# oder
# FROM golang:1.22

# ✅ Neuer, funktionierender Base Image
FROM golang:1.24.5

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

CMD ["./main"]

