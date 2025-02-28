using RatingService.Models;
using RatingService.Repositories;
using RatingService.Events;
using RatingService.Messaging.Interfaces;

namespace RatingService.Services;

public class RatingService : IRatingService
{
    private readonly IRatingRepository _ratingRepository;
    private readonly IEventPublisher<RatingCreatedEvent> _eventPublisher;
    
    public RatingService(IRatingRepository ratingRepository, IEventPublisher<RatingCreatedEvent> eventPublisher)
    {
        _ratingRepository = ratingRepository;
        _eventPublisher = eventPublisher;
    }

    public async Task<Rating> CreateRatingAsync(Rating rating)
    {

        if (rating.Score < 1 || rating.Score > 5)
        {
            throw new ArgumentException("Score must be between 1 and 5.");
        }

        var created = await _ratingRepository.CreateAsync(rating);

        var ratingEvent = new RatingCreatedEvent
        {
            ProviderId = created.ProviderId,
            Score = created.Score,
            Comment = created.Comment,
            CreatedAt = created.CreatedAt
        };
        await _eventPublisher.PublishAsync(ratingEvent);

        return created;
    }

    public async Task<double> GetAverageScoreAsync(int providerId)
    {
        return await _ratingRepository.GetAverageScoreByProviderIdAsync(providerId);
    }
}

