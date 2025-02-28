using System;

namespace RatingService.Events
{
    public class RatingCreatedEvent
    {
        public int ProviderId { get; set; }
        public int Score { get; set; }
        public string? Comment { get; set; }
        public DateTime CreatedAt { get; set; }
    }
}
