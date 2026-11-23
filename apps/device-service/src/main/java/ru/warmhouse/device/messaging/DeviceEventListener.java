package ru.warmhouse.device.messaging;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;
import ru.warmhouse.device.application.DeviceService;

@Component
public class DeviceEventListener {

    private static final Logger log = LoggerFactory.getLogger(DeviceEventListener.class);

    private final DeviceService deviceService;
    private final ObjectMapper objectMapper;

    public DeviceEventListener(DeviceService deviceService, ObjectMapper objectMapper) {
        this.deviceService = deviceService;
        this.objectMapper = objectMapper;
    }

    @RabbitListener(queues = RabbitConfig.QUEUE)
    public void onDeviceEvent(byte[] payload) throws Exception {
        DeviceEvent event = objectMapper.readValue(payload, DeviceEvent.class);
        log.info("Received device event: type={}, id={}", event.eventType(), event.deviceId());
        if (event.deviceId() == null) {
            return;
        }

        String type = event.eventType() == null ? "" : event.eventType();
        if ("DeviceRemoved".equalsIgnoreCase(type)) {
            deviceService.removeFromEvent(event.deviceId());
            return;
        }

        deviceService.upsertFromEvent(event);
    }
}
