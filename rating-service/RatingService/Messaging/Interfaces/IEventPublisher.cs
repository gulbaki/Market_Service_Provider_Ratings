using System.Threading.Tasks;

namespace RatingService.Messaging.Interfaces
{
    public interface IEventPublisher<TEvent>
    {
        Task PublishAsync(TEvent @event);
    }
}
