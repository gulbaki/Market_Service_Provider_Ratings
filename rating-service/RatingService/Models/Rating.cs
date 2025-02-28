using System;

namespace RatingService.Models
{
    public class Rating
    {
        public int Id { get; set; }
        public int ProviderId { get; set; } 
        public int Score { get; set; }       
        public string? Comment { get; set; } 
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    }
}
