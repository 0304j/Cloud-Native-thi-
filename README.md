# Cloud-Native-thi-

Ein cloud-natives Microservices-Projekt für die Hochschule, entwickelt mit Go, MongoDB und Docker.

## Features

- Authentifizierungsservice (JWT, User-Registrierung/Login)
- Shopping-Service
- Checkout-Service
- MongoDB als zentrale Datenbank
- Docker-Compose für einfaches Setup

## Voraussetzungen

- [Docker](https://www.docker.com/) und [Docker Compose](https://docs.docker.com/compose/)

## Projekt starten

Im Projektordner ausführen:

```sh
docker-compose up --build -d
```

Alle Services (inkl. MongoDB) laufen dann im Hintergrund.

## Wichtige Endpunkte

### Auth-Service (http://localhost:8081)

- `POST /register` – User registrieren (Body: email, password, role)
- `POST /login` – Login, gibt JWT-Token zurück

### Shopping-Service (http://localhost:8080)

- `GET /products` – Alle Produkte anzeigen
- `POST /products` – Produkt anlegen (nur mit JWT-Token, Rolle admin)

## Beispiel-Requests

### User registrieren

```json
POST http://localhost:8081/register
{
	"email": "admin@example.com",
	"password": "deinPasswort",
	"role": "admin"
}
```

### Login (Token erhalten)

```json
POST http://localhost:8081/login
{
	"email": "admin@example.com",
	"password": "deinPasswort"
}
```

### Produkt anlegen (mit Token)

```json
POST http://localhost:8080/products
Header: Authorization: Bearer <JWT_TOKEN>
{
	"name": "Pizza Salami",
	"price": 7.99,
	"description": "Leckere Pizza"
}
```

## Datenbankzugriff

- MongoDB läuft auf `localhost:27017` (Standard-User: root/rootpass)
- Datenbanken: `shopping` (Produkte), `auth_db` (User)
- Zugriff z.B. mit [MongoDB Compass](https://www.mongodb.com/try/download/compass) oder `mongosh`

## Nützliche Tools

- [Bruno](https://www.usebruno.com/) oder Postman für API-Tests
- [MongoDB Compass](https://www.mongodb.com/try/download/compass) für Datenbank-Visualisierung

## Lizenz

MIT License – siehe LICENSE-Datei.
