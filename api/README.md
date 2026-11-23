# API «Тёплый дом»

Документация контрактов взаимодействия микросервисов (задание 4).

| Спецификация | Тип | Назначение | Просмотр |
| --- | --- | --- | --- |
| [openapi.yaml](openapi.yaml) | OpenAPI 3.0 (REST) | 5 синхронных эндпоинтов через API Gateway | [Swagger Editor](https://editor.swagger.io/) — File → Import file |
| [asyncapi.yaml](asyncapi.yaml) | AsyncAPI 2.6 | События Message Broker | [AsyncAPI Studio](https://studio.asyncapi.com/) — открыть файл |

## REST-эндпоинты

| # | Метод | Путь | Сервис |
| --- | --- | --- | --- |
| 1 | `POST` | `/api/v1/devices` | device-service |
| 2 | `GET` | `/api/v1/devices/{deviceId}` | device-service |
| 3 | `PATCH` | `/api/v1/devices/{deviceId}` | device-service |
| 4 | `POST` | `/api/v1/heating/{deviceId}/on` | heating-service |
| 5 | `GET` | `/api/v1/telemetry/devices/{deviceId}/latest` | telemetry-service |

## Async-каналы

| Канал | Событие | Publisher → Consumer |
| --- | --- | --- |
| `heating.commands` | `HeatingCommandIssued` | heating-service → device-adapter |
| `telemetry.received` | `TelemetryReceived` | device-adapter → telemetry-service, scenario-service |
| `devices.status` | `DeviceStatusChanged` | device-adapter → device-service, heating-service |
| `devices.lifecycle` | `DeviceRegistered` / `DeviceUpdated` | device-service → подписчики |
