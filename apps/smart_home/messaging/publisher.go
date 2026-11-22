package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "warmhouse.events"

	RoutingDeviceRegistered = "device.registered"
	RoutingDeviceUpdated    = "device.updated"
	RoutingDeviceRemoved    = "device.removed"
	RoutingTelemetry        = "telemetry.received"
)

// EventPublisher публикует доменные события в RabbitMQ (Strangler Fig / outbox-lite).
type EventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewEventPublisher(rabbitURL string) (*EventPublisher, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq exchange: %w", err)
	}

	log.Printf("Event publisher connected to RabbitMQ (%s)", ExchangeName)
	return &EventPublisher{conn: conn, channel: ch}, nil
}

func (p *EventPublisher) Close() {
	if p == nil {
		return
	}
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

type DeviceEvent struct {
	DeviceID  int64     `json:"deviceId"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Location  string    `json:"location"`
	Status    string    `json:"status"`
	EventType string    `json:"eventType"`
	EmittedAt time.Time `json:"emittedAt"`
}

type TelemetryEvent struct {
	DeviceID   int64     `json:"deviceId"`
	MetricType string    `json:"metricType"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Location   string    `json:"location"`
	Status     string    `json:"status"`
	RecordedAt time.Time `json:"recordedAt"`
	EventType  string    `json:"eventType"`
}

func (p *EventPublisher) PublishDeviceRegistered(deviceID int64, name, typ, location, status string) {
	p.publishDevice(RoutingDeviceRegistered, "DeviceRegistered", deviceID, name, typ, location, status)
}

func (p *EventPublisher) PublishDeviceUpdated(deviceID int64, name, typ, location, status string) {
	p.publishDevice(RoutingDeviceUpdated, "DeviceUpdated", deviceID, name, typ, location, status)
}

func (p *EventPublisher) PublishDeviceRemoved(deviceID int64) {
	p.publishDevice(RoutingDeviceRemoved, "DeviceRemoved", deviceID, "", "", "", "")
}

func (p *EventPublisher) PublishTelemetry(deviceID int64, value float64, unit, location, status string, recordedAt time.Time) {
	if p == nil || p.channel == nil {
		return
	}
	event := TelemetryEvent{
		DeviceID:   deviceID,
		MetricType: "temperature",
		Value:      value,
		Unit:       unit,
		Location:   location,
		Status:     status,
		RecordedAt: recordedAt,
		EventType:  "TelemetryReceived",
	}
	p.publish(RoutingTelemetry, event)
}

func (p *EventPublisher) publishDevice(routingKey, eventType string, deviceID int64, name, typ, location, status string) {
	if p == nil || p.channel == nil {
		return
	}
	event := DeviceEvent{
		DeviceID:  deviceID,
		Name:      name,
		Type:      typ,
		Location:  location,
		Status:    status,
		EventType: eventType,
		EmittedAt: time.Now().UTC(),
	}
	p.publish(routingKey, event)
}

func (p *EventPublisher) publish(routingKey string, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal event %s: %v", routingKey, err)
		return
	}

	err = p.channel.Publish(
		ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
			Body:         body,
		},
	)
	if err != nil {
		log.Printf("failed to publish event %s: %v", routingKey, err)
		return
	}
	log.Printf("published event routing_key=%s body=%s", routingKey, string(body))
}

// NoopPublisher используется, если брокер недоступен — монолит продолжает работать.
type NoopPublisher struct{}

func (n *NoopPublisher) Close() {}
func (n *NoopPublisher) PublishDeviceRegistered(deviceID int64, name, typ, location, status string) {
}
func (n *NoopPublisher) PublishDeviceUpdated(deviceID int64, name, typ, location, status string) {
}
func (n *NoopPublisher) PublishDeviceRemoved(deviceID int64) {}
func (n *NoopPublisher) PublishTelemetry(deviceID int64, value float64, unit, location, status string, recordedAt time.Time) {
}

// Publisher — общий контракт для handlers.
type Publisher interface {
	Close()
	PublishDeviceRegistered(deviceID int64, name, typ, location, status string)
	PublishDeviceUpdated(deviceID int64, name, typ, location, status string)
	PublishDeviceRemoved(deviceID int64)
	PublishTelemetry(deviceID int64, value float64, unit, location, status string, recordedAt time.Time)
}
