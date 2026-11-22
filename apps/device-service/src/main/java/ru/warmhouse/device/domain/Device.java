package ru.warmhouse.device.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

import java.time.Instant;

@Entity
@Table(name = "devices")
public class Device {

    @Id
    private Long id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false)
    private String type;

    @Column(nullable = false)
    private String location;

    @Column(nullable = false)
    private String status;

    @Column(name = "source_sensor_id")
    private Long sourceSensorId;

    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    @Column(name = "created_at", nullable = false)
    private Instant createdAt;

    protected Device() {
    }

    public Device(Long id, String name, String type, String location, String status, Long sourceSensorId) {
        Instant now = Instant.now();
        this.id = id;
        this.name = name;
        this.type = type;
        this.location = location;
        this.status = status;
        this.sourceSensorId = sourceSensorId;
        this.createdAt = now;
        this.updatedAt = now;
    }

    public void applyUpdate(String name, String type, String location, String status) {
        if (name != null && !name.isBlank()) {
            this.name = name;
        }
        if (type != null && !type.isBlank()) {
            this.type = type;
        }
        if (location != null && !location.isBlank()) {
            this.location = location;
        }
        if (status != null && !status.isBlank()) {
            this.status = status;
        }
        this.updatedAt = Instant.now();
    }

    public Long getId() {
        return id;
    }

    public String getName() {
        return name;
    }

    public String getType() {
        return type;
    }

    public String getLocation() {
        return location;
    }

    public String getStatus() {
        return status;
    }

    public Long getSourceSensorId() {
        return sourceSensorId;
    }

    public Instant getUpdatedAt() {
        return updatedAt;
    }

    public Instant getCreatedAt() {
        return createdAt;
    }
}
