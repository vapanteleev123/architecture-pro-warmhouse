package ru.warmhouse.device.api.dto;

public record DeviceRequest(
        Long id,
        String name,
        String type,
        String location,
        String status,
        Long sourceSensorId
) {
}
