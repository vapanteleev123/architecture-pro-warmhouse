"""Consumer событий TelemetryReceived из RabbitMQ."""

from __future__ import annotations

import json
import logging
import threading
from datetime import datetime
from typing import Optional

import pika

from app.service import TelemetryService

logger = logging.getLogger(__name__)

EXCHANGE = "warmhouse.events"
QUEUE = "telemetry-service.telemetry"
ROUTING_KEY = "telemetry.received"


class TelemetryEventConsumer:
    def __init__(self, rabbitmq_url: str, service: TelemetryService) -> None:
        self._rabbitmq_url = rabbitmq_url
        self._service = service
        self._thread: Optional[threading.Thread] = None

    def start(self) -> None:
        self._thread = threading.Thread(target=self._run, name="telemetry-consumer", daemon=True)
        self._thread.start()

    def _run(self) -> None:
        while True:
            try:
                self._consume_forever()
            except Exception as exc:  # noqa: BLE001 — reconnect loop for MVP
                logger.exception("Telemetry consumer error: %s; reconnecting in 5s", exc)
                threading.Event().wait(5)

    def _consume_forever(self) -> None:
        params = pika.URLParameters(self._rabbitmq_url)
        connection = pika.BlockingConnection(params)
        channel = connection.channel()
        channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
        channel.queue_declare(queue=QUEUE, durable=True)
        channel.queue_bind(queue=QUEUE, exchange=EXCHANGE, routing_key=ROUTING_KEY)

        def on_message(ch, method, properties, body):  # noqa: ANN001
            try:
                payload = json.loads(body.decode("utf-8"))
                device_id = int(payload.get("deviceId") or payload.get("device_id"))
                value = float(payload["value"])
                recorded_at = self._parse_time(payload.get("recordedAt") or payload.get("recorded_at"))
                self._service.ingest(
                    device_id=device_id,
                    value=value,
                    metric_type=payload.get("metricType") or payload.get("metric_type") or "temperature",
                    unit=payload.get("unit") or "C",
                    location=payload.get("location") or "",
                    status=payload.get("status") or "active",
                    recorded_at=recorded_at,
                )
                logger.info("Ingested telemetry for device_id=%s value=%s", device_id, value)
                ch.basic_ack(delivery_tag=method.delivery_tag)
            except Exception:
                logger.exception("Failed to process telemetry event: %s", body)
                ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

        channel.basic_qos(prefetch_count=10)
        channel.basic_consume(queue=QUEUE, on_message_callback=on_message)
        logger.info("Telemetry consumer started on queue=%s", QUEUE)
        channel.start_consuming()

    @staticmethod
    def _parse_time(value: Optional[str]) -> Optional[datetime]:
        if not value:
            return None
        try:
            return datetime.fromisoformat(value.replace("Z", "+00:00"))
        except ValueError:
            return None
