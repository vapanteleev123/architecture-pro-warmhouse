"""Application-сервис телеметрии."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from app.domain import Measurement
from app.repository import TelemetryRepository


class TelemetryService:
    def __init__(self, repository: TelemetryRepository) -> None:
        self._repository = repository

    def ingest(
        self,
        *,
        device_id: int,
        value: float,
        metric_type: str = "temperature",
        unit: str = "C",
        location: str = "",
        status: str = "active",
        recorded_at: Optional[datetime] = None,
    ) -> Measurement:
        measurement = Measurement.create(
            device_id=device_id,
            metric_type=metric_type,
            value=value,
            unit=unit,
            location=location,
            status=status,
            recorded_at=recorded_at,
        )
        return self._repository.save(measurement)

    def list_all(self, limit: int = 100) -> List[Measurement]:
        return self._repository.list_all(limit=limit)

    def list_by_device(self, device_id: int, limit: int = 50) -> List[Measurement]:
        return self._repository.list_by_device(device_id, limit=limit)

    def latest(self, device_id: int) -> Optional[Measurement]:
        return self._repository.latest_by_device(device_id)
