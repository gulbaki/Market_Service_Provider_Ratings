using System;
using System.Threading.Tasks;
using Xunit;
using RatingService.Services;
using RatingService.Repositories;
using RatingService.Models;
using Moq;
using RatingService.Events;
using RatingService.Messaging.Interfaces;

namespace RatingService.Tests.Services
{
    public class RatingServiceTests
    {
        [Fact]
        public async Task CreateRatingAsync_ValidScore_ShouldSaveAndPublishEvent()
        {
            // Arrange
            var mockRepo = new Mock<IRatingRepository>();
            mockRepo.Setup(r => r.CreateAsync(It.IsAny<Rating>()))
                .ReturnsAsync((Rating r) => r); // DB'ye eklenmiş gibi aynı objeyi dönsün

            var mockPublisher = new Mock<IEventPublisher<RatingCreatedEvent>>();
            var sut = new RatingService.Services.RatingService(mockRepo.Object, mockPublisher.Object);

            var newRating = new Rating { ProviderId = 101, Score = 4, Comment = "Test" };

            // Act
            var result = await sut.CreateRatingAsync(newRating);

            // Assert
            mockRepo.Verify(r => r.CreateAsync(It.IsAny<Rating>()), Times.Once);
            mockPublisher.Verify(p => p.PublishAsync(It.IsAny<RatingCreatedEvent>()), Times.Once);

            Assert.NotNull(result);
            Assert.Equal(101, result.ProviderId);
            Assert.Equal("Test", result.Comment);
        }
        [Theory]
        [InlineData(0)]
        [InlineData(6)]
        public async Task CreateRatingAsync_InvalidScore_ShouldThrowArgumentException(int invalidScore)
        {
            // Arrange
            var mockRepo = new Mock<IRatingRepository>();
            var mockPublisher = new Mock<IEventPublisher<RatingCreatedEvent>>();
            var sut = new RatingService.Services.RatingService(mockRepo.Object, mockPublisher.Object);

            var newRating = new Rating { ProviderId = 100, Score = invalidScore };

            // Act & Assert
            await Assert.ThrowsAsync<ArgumentException>(() => sut.CreateRatingAsync(newRating));

            // Publisher'ın hiç çağrılmadığını doğrulayalım
            mockPublisher.Verify(p => p.PublishAsync(It.IsAny<RatingCreatedEvent>()), Times.Never);
        }
    }
}
