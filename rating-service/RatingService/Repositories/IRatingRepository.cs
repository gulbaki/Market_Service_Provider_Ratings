using RatingService.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace RatingService.Repositories
{
    public interface IRatingRepository
    {
        Task<Rating> CreateAsync(Rating rating);
        Task<Rating?> GetByIdAsync(int id);
        Task<IEnumerable<Rating>> GetByProviderIdAsync(int providerId);
        Task<double> GetAverageScoreByProviderIdAsync(int providerId);
    }
}
