# Smart Home Sensor Management API

## Prerequisites

- Docker and Docker Compose

## Getting Started

### Option 1: Using Docker Compose (Recommended)

The easiest way to start the application is to use Docker Compose:

```bash
./init.sh
```

This script will:

1. Build and start PostgreSQL, RabbitMQ, temperature-api, device-service, telemetry-service and smart_home
2. Wait for the services to be ready
3. Display information about how to access the API

Alternatively, you can run Docker Compose directly:

```bash
docker compose up --build -d
```

Services:

| Service | URL |
| --- | --- |
| smart_home API | http://localhost:8080 |
| temperature-api | http://localhost:8081 |
| device-service (Java) | http://localhost:8082 |
| telemetry-service (Python) | http://localhost:8083 |
| RabbitMQ management | http://localhost:15672 (guest/guest) |
| PostgreSQL | localhost:5432 |

Example: `GET http://localhost:8081/temperature?location=Living%20Room` returns a random temperature each call.

MVP integration: CRUD/чтение датчиков в монолите публикует события в RabbitMQ → device-service / telemetry-service.

### Option 2: Manual setup

If you prefer to run the application without Docker:

1. Start the PostgreSQL database:

```bash
docker-compose up -d postgres
```

2. Build and run the application:

```bash
go build -o smarthome
./smarthome
```

## API Testing

A Postman collection is provided for testing the API. Import the `smarthome-api.postman_collection.json` file into Postman to get started.

## API Endpoints

### Monolith (smart_home)

- `GET /health` - Health check
- `GET /api/v1/sensors` - Get all sensors
- `GET /api/v1/sensors/:id` - Get a specific sensor
- `POST /api/v1/sensors` - Create a new sensor
- `PUT /api/v1/sensors/:id` - Update a sensor
- `DELETE /api/v1/sensors/:id` - Delete a sensor
- `PATCH /api/v1/sensors/:id/value` - Update a sensor's value and status

### device-service

- `GET /health`
- `GET/POST /api/v1/devices`
- `GET/PUT/DELETE /api/v1/devices/{id}`

### telemetry-service

- `GET /health`
- `GET /api/v1/telemetry`
- `GET /api/v1/telemetry?device_id=1`
- `GET /api/v1/telemetry/devices/{id}/latest`
- `POST /api/v1/telemetry`
