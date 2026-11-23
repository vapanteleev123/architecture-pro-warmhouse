"""REST API телеметрии."""

from __future__ import annotations

from datetime import datetime
from typing import Optional

from fastapi import APIRouter, HTTPException, Query
from pydantic import BaseModel, Field

from app.service import TelemetryService


class TelemetryIngestRequest(BaseModel):
    device_id: int = Field(..., ge=1)
    value: float
    metric_type: str = "temperature"
    unit: str = "C"
    location: str = ""
    status: str = "active"
    recorded_at: Optional[datetime] = None


def build_router(service: TelemetryService) -> APIRouter:
    router = APIRouter(prefix="/api/v1/telemetry", tags=["telemetry"])

    @router.get("")
    def list_telemetry(
        device_id: Optional[int] = Query(default=None),
        limit: int = Query(default=100, ge=1, le=500),
    ):
        if device_id is not None:
            items = service.list_by_device(device_id, limit=limit)
        else:
            items = service.list_all(limit=limit)
        return [item.to_dict() for item in items]

    @router.get("/devices/{device_id}/latest")
    def latest_for_device(device_id: int):
        item = service.latest(device_id)
        if item is None:
            raise HTTPException(status_code=404, detail="No telemetry for device")
        return item.to_dict()

    @router.post("", status_code=201)
    def ingest(request: TelemetryIngestRequest):
        item = service.ingest(
            device_id=request.device_id,
            value=request.value,
            metric_type=request.metric_type,
            unit=request.unit,
            location=request.location,
            status=request.status,
            recorded_at=request.recorded_at,
        )
        return item.to_dict()

    return router
