package ru.warmhouse.device.application;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.warmhouse.device.domain.Device;
import ru.warmhouse.device.domain.DeviceRepository;
import ru.warmhouse.device.api.dto.DeviceRequest;
import ru.warmhouse.device.api.dto.DeviceResponse;
import ru.warmhouse.device.messaging.DeviceEvent;

import java.util.List;
import java.util.NoSuchElementException;

@Service
public class DeviceService {

    private final DeviceRepository deviceRepository;

    public DeviceService(DeviceRepository deviceRepository) {
        this.deviceRepository = deviceRepository;
    }

    @Transactional(readOnly = true)
    public List<DeviceResponse> listDevices() {
        return deviceRepository.findAll().stream().map(DeviceResponse::from).toList();
    }

    @Transactional(readOnly = true)
    public DeviceResponse getDevice(Long id) {
        return DeviceResponse.from(findOrThrow(id));
    }

    @Transactional
    public DeviceResponse registerDevice(DeviceRequest request) {
        Long id = request.id() != null ? request.id() : nextSyntheticId();
        Device device = new Device(
                id,
                request.name(),
                request.type(),
                request.location(),
                request.status() == null || request.status().isBlank() ? "inactive" : request.status(),
                request.sourceSensorId() != null ? request.sourceSensorId() : id
        );
        return DeviceResponse.from(deviceRepository.save(device));
    }

    @Transactional
    public DeviceResponse updateDevice(Long id, DeviceRequest request) {
        Device device = findOrThrow(id);
        device.applyUpdate(request.name(), request.type(), request.location(), request.status());
        return DeviceResponse.from(deviceRepository.save(device));
    }

    @Transactional
    public void deleteDevice(Long id) {
        if (!deviceRepository.existsById(id)) {
            throw new NoSuchElementException("Device not found: " + id);
        }
        deviceRepository.deleteById(id);
    }

    @Transactional
    public void upsertFromEvent(DeviceEvent event) {
        String status = event.status() == null || event.status().isBlank() ? "inactive" : event.status();
        deviceRepository.findById(event.deviceId()).ifPresentOrElse(
                existing -> {
                    existing.applyUpdate(event.name(), event.type(), event.location(), status);
                    deviceRepository.save(existing);
                },
                () -> deviceRepository.save(new Device(
                        event.deviceId(),
                        event.name(),
                        event.type(),
                        event.location(),
                        status,
                        event.deviceId()
                ))
        );
    }

    @Transactional
    public void removeFromEvent(Long deviceId) {
        deviceRepository.deleteById(deviceId);
    }

    private Device findOrThrow(Long id) {
        return deviceRepository.findById(id)
                .orElseThrow(() -> new NoSuchElementException("Device not found: " + id));
    }

    private Long nextSyntheticId() {
        return deviceRepository.findAll().stream()
                .map(Device::getId)
                .max(Long::compareTo)
                .orElse(0L) + 1;
    }
}
