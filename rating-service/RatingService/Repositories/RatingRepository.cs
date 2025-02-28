using Microsoft.EntityFrameworkCore;
using RatingService.Data;
using RatingService.Models;

namespace RatingService.Repositories
{
    public class RatingRepository : IRatingRepository
    {
        private readonly RatingDbContext _context;

        public RatingRepository(RatingDbContext context)
        {
            _context = context;
        }

        public async Task<Rating> CreateAsync(Rating rating)
        {
            _context.Ratings.Add(rating);
            await _context.SaveChangesAsync();
            return rating;
        }

        public async Task<Rating?> GetByIdAsync(int id)
        {
            return await _context.Ratings.FindAsync(id);
        }

        public async Task<IEnumerable<Rating>> GetByProviderIdAsync(int providerId)
        {
            return await _context.Ratings
                .Where(r => r.ProviderId == providerId)
                .ToListAsync();
        }

        public async Task<double> GetAverageScoreByProviderIdAsync(int providerId)
        {
            var ratings = await _context.Ratings
                .Where(r => r.ProviderId == providerId)
                .ToListAsync();

            if (ratings.Count == 0)
                return 0.0;

            return ratings.Average(r => r.Score);
        }
    }
}
