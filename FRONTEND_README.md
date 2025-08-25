# Cloud-Native E-Commerce Frontend

## 📋 Übersicht

Dieses Projekt enthält responsive HTML-Frontends für jeden Ihrer Backend-Services, die nahtlos miteinander kommunizieren und eine vollständige E-Commerce-Erfahrung bieten.

## 🎯 Funktionalitäten

### 🔐 Auth Service Frontend (`auth-service/index.html`)
- **Benutzerregistrierung** mit E-Mail, Passwort und Rolle (user/admin)
- **Sichere Anmeldung** mit JWT-Token-Authentifizierung
- **Automatische Weiterleitung** zu anderen Services nach erfolgreicher Anmeldung
- **Responsive Design** für alle Geräte

### 🛒 Shopping Service Frontend (`shopping-service/index.html`)
- **Produktkatalog** anzeigen und durchsuchen
- **Admin-Funktionen** für Produkterstellung (nur für Administratoren)
- **Warenkorb-Management** mit Produkten hinzufügen/entfernen
- **Checkout-Prozess** mit Kafka-Integration
- **JWT-geschützte Bereiche** mit automatischer Token-Validierung

### 💳 Checkout Service Frontend (`checkout-service/index.html`)
- **Echtzeit-Bestellungsübersicht** (Simulation)
- **Kafka-Message-Monitoring** Dashboard
- **Bestellungsstatistiken** und -verarbeitung
- **Live-Updates** für neue Bestellungen

### 🏠 Haupt-Dashboard (`index.html`)
- **Service-Übersicht** mit Status-Monitoring
- **Architektur-Diagramm** des Systems
- **Quick-Start-Guide** für Entwickler
- **Technologie-Stack** Übersicht

## 🚀 Verwendung

### 1. Services starten
Stelle sicher, dass deine Backend-Services laufen:
```bash
# Docker-Container starten
docker-compose up -d

# Services sind verfügbar unter:
# Auth Service: http://localhost:8081
# Shopping Service: http://localhost:8080
```

### 2. Frontend öffnen
1. **Haupt-Dashboard öffnen**: Öffne `index.html` in deinem Browser
2. **Oder direkt zu einem Service**: Öffne die jeweilige `index.html` Datei

### 3. Kompletter Workflow testen

#### Als Administrator:
1. **Auth Service** → Registrierung mit Rolle "admin"
2. **Anmelden** und zum Shopping Service wechseln  
3. **Neue Produkte erstellen** (nur als Admin möglich)
4. **Checkout Service** für Monitoring öffnen

#### Als normaler Kunde:
1. **Auth Service** → Registrierung mit Rolle "user"
2. **Anmelden** und zum Shopping Service wechseln
3. **Produkte durchsuchen** und zum Warenkorb hinzufügen
4. **Bestellung aufgeben** → automatische Weiterleitung an Checkout Service
5. **Bestellungsstatus** im Checkout Dashboard verfolgen

## 🔧 Features

### ✨ Cross-Service-Kommunikation
- **Automatische Token-Übertragung** zwischen Services
- **Nahtlose Navigation** durch URL-Parameter
- **Persistente Anmeldung** über localStorage
- **Rollenbasierte UI-Anpassungen**

### 📱 Responsive Design
- **Mobile-First** Approach
- **Flexible Layouts** für alle Bildschirmgrößen
- **Touch-optimierte** Bedienelemente
- **Progressive Enhancement**

### 🔒 Sicherheit
- **JWT-Token-Validierung** in Echtzeit
- **Automatische Abmeldung** bei ungültigen Tokens
- **Rollenbasierte Zugriffskontrolle**
- **Sichere API-Kommunikation**

### 🎨 Benutzerfreundlichkeit
- **Intuitive Navigation** mit Tabs und Buttons
- **Live-Feedback** für alle Aktionen
- **Fehlerbehandlung** mit aussagekräftigen Meldungen
- **Loading-Spinner** und Status-Indikatoren

## 📁 Dateistruktur

```
Cloud-Native-thi/
├── index.html                    # Haupt-Dashboard
├── auth-service/
│   └── index.html                # Auth Service Frontend
├── shopping-service/
│   └── index.html                # Shopping Service Frontend
└── checkout-service/
    └── index.html                # Checkout Service Frontend
```

## 🌐 API-Endpunkte

Die Frontends kommunizieren mit folgenden Backend-Endpunkten:

### Auth Service (Port 8081)
- `POST /register` - Benutzerregistrierung
- `POST /login` - Benutzeranmeldung

### Shopping Service (Port 8080)
- `GET /products` - Alle Produkte abrufen
- `POST /products` - Neues Produkt erstellen (Admin)
- `POST /cart` - Artikel zum Warenkorb hinzufügen
- `GET /cart` - Warenkorb abrufen
- `POST /checkout` - Bestellung aufgeben

### Checkout Service
- Konsumiert Kafka-Messages vom Shopping Service
- Frontend zeigt Simulation der empfangenen Bestellungen

## ⚡ Quick Start

1. **Browser öffnen** und zu `index.html` navigieren
2. **Services-Status prüfen** im Dashboard
3. **Auth Service öffnen** und Account erstellen
4. **Shopping Service erkunden** und Produkte hinzufügen
5. **Bestellung testen** und Checkout Service überwachen

## 🎯 Besonderheiten

- **Keine Backend-Änderungen erforderlich** - alle Frontends arbeiten mit deinen bestehenden APIs
- **Vollständig eigenständig** - jeder Service hat sein eigenes Frontend
- **Service-übergreifende Navigation** - nahtlose Benutzererfahrung
- **Echtzeit-Updates** - automatische Aktualisierung von Daten
- **Entwicklerfreundlich** - alle APIs werden getestet und Fehler werden angezeigt

## 🔍 Debugging

- **Browser-Konsole** öffnen für detaillierte Logs
- **Netzwerk-Tab** für API-Aufrufe überwachen
- **Service-Status** im Haupt-Dashboard prüfen
- **CORS-Probleme** können bei localhost auftreten (Normal bei Entwicklung)

## 📝 Anpassungen

Die HTML-Dateien können einfach angepasst werden:
- **API-URLs** in den JavaScript-Konstanten ändern
- **Styling** über CSS-Variablen anpassen
- **Neue Funktionen** durch zusätzliche JavaScript-Funktionen hinzufügen

Viel Spaß mit deinem Cloud-Native E-Commerce System! 🚀
