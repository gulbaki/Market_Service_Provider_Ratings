using Microsoft.EntityFrameworkCore;
using RatingService.Data;
using RatingService.Repositories;
using RatingService.Messaging.Interfaces;
using RatingService.Messaging.Kafka;
using RatingService.Events;
using RatingService.Services;

var builder = WebApplication.CreateBuilder(args);

// PostgreSQL Connection
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection");
builder.Services.AddDbContext<RatingDbContext>(options =>
    options.UseNpgsql(connectionString));

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddScoped<IRatingRepository, RatingRepository>();
builder.Services.AddScoped<IRatingService, RatingService.Services.RatingService>();
builder.Services.AddSingleton<IEventPublisher<RatingCreatedEvent>, KafkaEventPublisher>();

var app = builder.Build();

// Swagger middleware
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}


app.MapControllers();

app.Run();
