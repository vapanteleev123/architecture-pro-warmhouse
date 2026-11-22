package ru.warmhouse.device.messaging;

import org.springframework.amqp.core.Binding;
import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.TopicExchange;
import org.springframework.amqp.support.converter.MessageConverter;
import org.springframework.amqp.support.converter.SimpleMessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitConfig {

    public static final String EXCHANGE = "warmhouse.events";
    public static final String QUEUE = "device-service.devices";
    public static final String ROUTING_PATTERN = "device.*";

    @Bean
    public TopicExchange warmhouseExchange() {
        return new TopicExchange(EXCHANGE, true, false);
    }

    @Bean
    public Queue deviceQueue() {
        return new Queue(QUEUE, true);
    }

    @Bean
    public Binding deviceBinding(Queue deviceQueue, TopicExchange warmhouseExchange) {
        return BindingBuilder.bind(deviceQueue).to(warmhouseExchange).with(ROUTING_PATTERN);
    }

    @Bean
    public MessageConverter messageConverter() {
        // Plain JSON from Go monolith — без __TypeId__ заголовков Spring.
        return new SimpleMessageConverter();
    }
}
