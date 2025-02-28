using System.Text.Json;
using Confluent.Kafka;
using Microsoft.Extensions.Configuration;
using RatingService.Events;
using RatingService.Messaging.Interfaces;

namespace RatingService.Messaging.Kafka
{
    public class KafkaEventPublisher : IEventPublisher<RatingCreatedEvent>
    {
        private readonly string _bootstrapServers;
        private readonly string _topicName;

        public KafkaEventPublisher(IConfiguration configuration)
        {
             System.Diagnostics.Debug.WriteLine(configuration["Kafka:BootstrapServers"]);

            _bootstrapServers = configuration["Kafka:BootstrapServers"] ?? "localhost:9093";
            _topicName = configuration["Kafka:TopicName"] ?? "rating-created";
        }

        public async Task PublishAsync(RatingCreatedEvent @event)
        {
            var config = new ProducerConfig
            {
                BootstrapServers = _bootstrapServers
            };

            using var producer = new ProducerBuilder<Null, string>(config).Build();
            
            // Convert event to JSON
            string messageValue = JsonSerializer.Serialize(@event);

            var message = new Message<Null, string>
            {
                Value = messageValue
            };

            // Publish message to Kafka
            var deliveryResult = await producer.ProduceAsync(_topicName, message);

            // Log delivery status to console
            Console.WriteLine($"Kafka delivery to topic {_topicName}: {deliveryResult.Status}");
        }
    }
}
