"""Репозиторий телеметрии (Database per Service)."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

import psycopg2
import psycopg2.extras

from app.domain import Measurement


class TelemetryRepository:
    def __init__(self, database_url: str) -> None:
        self._database_url = database_url
        self._ensure_schema()

    def _connect(self):
        return psycopg2.connect(self._database_url)

    def _ensure_schema(self) -> None:
        last_error: Exception | None = None
        for attempt in range(1, 16):
            try:
                self._create_schema()
                return
            except Exception as exc:  # noqa: BLE001 — retry until postgres/db ready
                last_error = exc
                import time

                time.sleep(2)
        raise RuntimeError(f"telemetry db not ready: {last_error}")

    def _create_schema(self) -> None:
        with self._connect() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    CREATE TABLE IF NOT EXISTS telemetry_data (
                        id SERIAL PRIMARY KEY,
                        device_id INTEGER NOT NULL,
                        metric_type VARCHAR(50) NOT NULL,
                        value DOUBLE PRECISION NOT NULL,
                        unit VARCHAR(20) NOT NULL DEFAULT 'C',
                        location VARCHAR(100) NOT NULL DEFAULT '',
                        status VARCHAR(20) NOT NULL DEFAULT 'active',
                        recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
                    );
                    CREATE INDEX IF NOT EXISTS idx_telemetry_device_id
                        ON telemetry_data(device_id);
                    CREATE INDEX IF NOT EXISTS idx_telemetry_recorded_at
                        ON telemetry_data(recorded_at DESC);
                    """
                )
            conn.commit()

    def save(self, measurement: Measurement) -> Measurement:
        with self._connect() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    INSERT INTO telemetry_data
                        (device_id, metric_type, value, unit, location, status, recorded_at)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                    RETURNING id
                    """,
                    (
                        measurement.device_id,
                        measurement.metric_type,
                        measurement.value,
                        measurement.unit,
                        measurement.location,
                        measurement.status,
                        measurement.recorded_at,
                    ),
                )
                measurement.id = cur.fetchone()[0]
            conn.commit()
        return measurement

    def list_by_device(self, device_id: int, limit: int = 50) -> List[Measurement]:
        with self._connect() as conn:
            with conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as cur:
                cur.execute(
                    """
                    SELECT id, device_id, metric_type, value, unit, location, status, recorded_at
                    FROM telemetry_data
                    WHERE device_id = %s
                    ORDER BY recorded_at DESC
                    LIMIT %s
                    """,
                    (device_id, limit),
                )
                rows = cur.fetchall()
        return [self._to_measurement(row) for row in rows]

    def latest_by_device(self, device_id: int) -> Optional[Measurement]:
        items = self.list_by_device(device_id, limit=1)
        return items[0] if items else None

    def list_all(self, limit: int = 100) -> List[Measurement]:
        with self._connect() as conn:
            with conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as cur:
                cur.execute(
                    """
                    SELECT id, device_id, metric_type, value, unit, location, status, recorded_at
                    FROM telemetry_data
                    ORDER BY recorded_at DESC
                    LIMIT %s
                    """,
                    (limit,),
                )
                rows = cur.fetchall()
        return [self._to_measurement(row) for row in rows]

    @staticmethod
    def _to_measurement(row) -> Measurement:
        recorded_at = row["recorded_at"]
        if isinstance(recorded_at, datetime) and recorded_at.tzinfo is None:
            recorded_at = recorded_at.replace(tzinfo=None)
        return Measurement(
            id=row["id"],
            device_id=row["device_id"],
            metric_type=row["metric_type"],
            value=row["value"],
            unit=row["unit"],
            location=row["location"],
            status=row["status"],
            recorded_at=recorded_at,
        )
