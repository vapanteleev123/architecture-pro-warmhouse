package ru.warmhouse.device.messaging;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record DeviceEvent(
        Long deviceId,
        String name,
        String type,
        String location,
        String status,
        String eventType
) {
}
