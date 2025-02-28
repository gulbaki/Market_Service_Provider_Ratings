using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using RatingService.Models;
using RatingService.Services;

namespace RatingService.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class RatingsController : ControllerBase
    {
        private readonly IRatingService _ratingService;

        public RatingsController(IRatingService ratingService)
        {
            _ratingService = ratingService;
        }

        // POST /api/ratings
        [HttpPost]
        public async Task<IActionResult> Create([FromBody] Rating rating)
        {
            if (rating == null)
            {
                return BadRequest("Rating data is required.");
            }

            try
            {
                var created = await _ratingService.CreateRatingAsync(rating);
                return Ok(created);
            }
            catch (ArgumentException ex)
            {
                return BadRequest(ex.Message);
            }
        }

        // GET /api/ratings/provider/{providerId}/average
        [HttpGet("provider/{providerId}/average")]
        public async Task<IActionResult> GetAverage(int providerId)
        {
            var averageScore = await _ratingService.GetAverageScoreAsync(providerId);
            return Ok(new { ProviderId = providerId, AverageScore = averageScore });
        }
    }
}
