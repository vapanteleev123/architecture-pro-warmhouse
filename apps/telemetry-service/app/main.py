"""Точка входа telemetry-service."""

from __future__ import annotations

import logging
import os

from fastapi import FastAPI

from app.api import build_router
from app.consumer import TelemetryEventConsumer
from app.repository import TelemetryRepository
from app.service import TelemetryService

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s: %(message)s")


def create_app() -> FastAPI:
    database_url = os.getenv(
        "DATABASE_URL",
        "postgresql://postgres:postgres@localhost:5432/telemetry_db",
    )
    rabbitmq_url = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

    repository = TelemetryRepository(database_url)
    service = TelemetryService(repository)

    app = FastAPI(title="telemetry-service", version="0.1.0")
    app.include_router(build_router(service))

    @app.get("/health")
    def health():
        return {"status": "ok", "service": "telemetry-service"}

    @app.on_event("startup")
    def startup() -> None:
        consumer = TelemetryEventConsumer(rabbitmq_url, service)
        consumer.start()

    return app


app = create_app()
