package ru.warmhouse.device.api.dto;

import ru.warmhouse.device.domain.Device;

import java.time.Instant;

public record DeviceResponse(
        Long id,
        String name,
        String type,
        String location,
        String status,
        Long sourceSensorId,
        Instant updatedAt,
        Instant createdAt
) {
    public static DeviceResponse from(Device device) {
        return new DeviceResponse(
                device.getId(),
                device.getName(),
                device.getType(),
                device.getLocation(),
                device.getStatus(),
                device.getSourceSensorId(),
                device.getUpdatedAt(),
                device.getCreatedAt()
        );
    }
}
