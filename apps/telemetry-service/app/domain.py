"""Доменная модель телеметрии (ООП)."""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timezone
from typing import Optional


@dataclass
class Measurement:
    """Одно измерение телеметрии от устройства."""

    device_id: int
    metric_type: str
    value: float
    unit: str
    location: str
    status: str
    recorded_at: datetime
    id: Optional[int] = None

    @classmethod
    def create(
        cls,
        *,
        device_id: int,
        metric_type: str,
        value: float,
        unit: str = "C",
        location: str = "",
        status: str = "active",
        recorded_at: Optional[datetime] = None,
    ) -> "Measurement":
        return cls(
            device_id=device_id,
            metric_type=metric_type or "temperature",
            value=value,
            unit=unit or "C",
            location=location or "",
            status=status or "active",
            recorded_at=recorded_at or datetime.now(timezone.utc),
        )

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "device_id": self.device_id,
            "metric_type": self.metric_type,
            "value": self.value,
            "unit": self.unit,
            "location": self.location,
            "status": self.status,
            "recorded_at": self.recorded_at.isoformat(),
        }
