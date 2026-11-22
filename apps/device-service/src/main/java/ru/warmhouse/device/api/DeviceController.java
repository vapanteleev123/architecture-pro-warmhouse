package ru.warmhouse.device.api;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ru.warmhouse.device.api.dto.DeviceRequest;
import ru.warmhouse.device.api.dto.DeviceResponse;
import ru.warmhouse.device.application.DeviceService;

import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;

@RestController
@RequestMapping("/api/v1/devices")
public class DeviceController {

    private final DeviceService deviceService;

    public DeviceController(DeviceService deviceService) {
        this.deviceService = deviceService;
    }

    @GetMapping
    public List<DeviceResponse> listDevices() {
        return deviceService.listDevices();
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> getDevice(@PathVariable Long id) {
        try {
            return ResponseEntity.ok(deviceService.getDevice(id));
        } catch (NoSuchElementException ex) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", ex.getMessage()));
        }
    }

    @PostMapping
    public ResponseEntity<DeviceResponse> createDevice(@RequestBody DeviceRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(deviceService.registerDevice(request));
    }

    @PutMapping("/{id}")
    public ResponseEntity<?> updateDevice(@PathVariable Long id, @RequestBody DeviceRequest request) {
        try {
            return ResponseEntity.ok(deviceService.updateDevice(id, request));
        } catch (NoSuchElementException ex) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", ex.getMessage()));
        }
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<?> deleteDevice(@PathVariable Long id) {
        try {
            deviceService.deleteDevice(id);
            return ResponseEntity.ok(Map.of("message", "Device deleted successfully"));
        } catch (NoSuchElementException ex) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", ex.getMessage()));
        }
    }
}
