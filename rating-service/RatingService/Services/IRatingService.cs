using RatingService.Models;

namespace RatingService.Services;

public interface IRatingService
{
    Task<Rating> CreateRatingAsync(Rating rating);
    
    Task<double> GetAverageScoreAsync(int providerId);
}