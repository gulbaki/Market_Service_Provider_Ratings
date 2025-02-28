

# Marketplace - Microservices (Rating & Notification)

This repository contains two microservices:

1. **Rating Service** (written in .NET 9)  
2. **Notification Service** (written in Go)

The system allows customers to submit ratings for service providers. When a new rating is created, it is saved in a PostgreSQL database and an event is published to Kafka. The Notification Service consumes these events from Kafka and stores them in-memory. Clients can then fetch and clear new notifications for a specific provider.

## Table of Contents

- [Overview](#overview)  
- [Why .NET 9 for Rating and Go for Notification](#why-net-9-for-rating-and-go-for-notification)  
- [Architecture](#architecture)  
- [Sequence Diagram](#sequence-diagram)  
- [Services & Responsibilities](#services--responsibilities)  
  - [Rating Service (.NET 9)](#rating-service-net-9)  
  - [Notification Service (Go)](#notification-service-go)  
- [Getting Started](#getting-started)  
  - [Prerequisites](#prerequisites)  
  - [Build & Run](#build--run)  
  - [Endpoints](#endpoints)  
- [Configuration](#configuration)  
- [Project Structure](#project-structure)  
- [Future Improvements](#future-improvements)  
---

## Overview

The primary goal of this project is to demonstrate a simple **service marketplace** approach, where:

- **Customers** can rate a **Service Provider**.  
- A new **rating** triggers an event publication to **Kafka**.  
- The **Notification Service** (in Go) listens to these events, stores them in memory, and offers a REST endpoint to retrieve notifications for a given provider.

Data is persisted in a PostgreSQL database for ratings, and notifications are stored in-memory. For real-world scalability or durability requirements, we could adapt the in-memory storage to a more robust solution (e.g., Redis etc.).

---

## Why .NET 9 for Rating and Go for Notification

- **Rating Service (.NET 9)**  
  - **.NET** offers a mature ecosystem for building CRUD-heavy applications, featuring robust data access libraries like Entity Framework Core, comprehensive validation/annotation mechanisms, and integrated dependency injection.  
  - Building a RESTful API is straightforward with ASP.NET Core, and the tooling (Visual Studio, CLI) accelerates development.  
  - .NET is well-suited for enterprise scenarios that involve complex business logic, validations, and layered architectures.

- **Notification Service (Go)**  
  - **Go (Golang)** is optimized for concurrency and efficient memory usage, making it an excellent choice for high-throughput or real-time components.  
  - Its lightweight runtime and straightforward syntax make containerization simpler, resulting in smaller images.  
  - It handles I/O-bound tasks such as message consumption from Kafka gracefully, thanks to goroutines and channels.  
  - Perfect for a microservice that may need to scale horizontally in response to high notification traffic.

Together, these choices showcase the strengths of each language platform: .NET for data-driven and robust enterprise solutions, and Go for performance-focused services requiring high concurrency and minimal overhead.

---

## Architecture

- **Microservices**:
  - Each microservice has its own codebase and data layer.
  - They communicate asynchronously via **Kafka**.
- **Databases**:
  - **Rating Service** uses **PostgreSQL** for persistent rating storage.
  - **Notification Service** uses an **in-memory** approach by default.
- **Event-Driven**:
  - The `Rating Service` publishes a `RatingCreated` event to Kafka whenever a new rating is created.
  - The `Notification Service` consumes this event and creates a new notification record.

Below is a high-level diagram of the components:

```
               ┌───────────────┐
               │   PostgreSQL  │
               └───────┬───────┘
                       │
┌─────────────┐  POST  │  ┌───────────────────┐
│  Client App  │───► Rating Service (.NET 9)  │
└─────────────┘       │  └───────────────────┘
                       │         Publish
                       │         to Kafka
                       ▼
                   ┌─────────┐
                   │  Kafka  │
                   └─────────┘
                       ▲
   ┌───────────────────┘
   │   Consume
   │
┌──────────────────────┐     GET
│ Notification Service │◄──────── Client (fetch notifications)
│        (Go)          │
└──────────────────────┘
```

---

## Sequence Diagram

Here is a more detailed sequence diagram using a standard flow:

```
 Customer                      RatingService             KafkaBroker              NotificationService
     | 1) POST /ratings            |                          |                             |
     |---------------------------->|                          |                             |
     |                             | 2) DB record (Ratings)   |                             |
     |                             |------------------------->|                             |
     |                             |<-------------------------|                             |
     |                             | 3) Produce "rating-created" event                       |
     |                             |------------------------->|                             |
     |                             |                          | 4) "rating-created" stored  |
     |                             |                          |----------------------------->|
     |                             |                          |<-----------------------------|
     |                             |                          | 5) Consumer reads message    |
     |                             |                          |----------------------------->|
     |                             |                          |                             | 6) Store in memory
     |                             |                          |                             |
     | 7) GET /notifications/{id}  |                          |                             |
     |---------------------------->|                          |                             |
     |                             |                          |                             | 8) Return JSON & clear
     |<----------------------------|                          |                             |
     | 9) GET /notifications/{id} again (empty)               |                             |
     |---------------------------->|                          |                             |
     |<----------------------------|                          |                             |
```

---

## Services & Responsibilities

### Rating Service (.NET 9)

- **Endpoints**:
  - `POST /api/ratings`: Creates a new rating in PostgreSQL and publishes a `RatingCreatedEvent` to Kafka.
  - `GET /api/ratings/provider/{providerId}/average`: Retrieves the average rating for a specific provider.
- **Technologies**:
  - ASP.NET Core 9 (WebAPI)
  - Entity Framework Core (PostgreSQL)
  - Confluent.Kafka (for producing Kafka messages)

### Notification Service (Go)

- **Function**:
  - Consumes `RatingCreated` messages from Kafka.
  - Stores notifications **in-memory**.
  - Returns and clears notifications via `GET /notifications/{providerId}`.
- **Technologies**:
  - Go 1.20
  - [segmentio/kafka-go](https://github.com/segmentio/kafka-go) for consumer logic
  - Gorilla Mux (REST API routing)

---

## Getting Started

### Prerequisites

- **Docker** and **Docker Compose** (version 3.8 or above recommended).
- Optionally, if you want to run locally without containers:
  - .NET 9 SDK
  - Go 1.20+
  - PostgreSQL

### Build & Run

1. **Clone this repository** (if not already).

2. **In the project root**, you can use the provided scripts (located in the `scripts` folder) to build and run:

   ```bash
   # Make scripts executable (Linux/Mac)
   chmod +x scripts/*.sh

   # Build Docker images for both services
   ./scripts/build_and_run.sh
   ```

   Alternatively, you can manually run:
   ```bash
   docker-compose up --build
   ```

3. **Docker Compose** will spin up the following containers:
   - **Postgres** (on port 5432)
   - **Zookeeper** (on port 2181)
   - **Kafka** (on port 9092)
   - **rating-service** (exposed on port 8181 → internally 80, per your setup)
   - **notification-service** (exposed on port 9191 → internally 8080, per your setup)

4. **Test Endpoints**:
   - **Rating Service**:  
     - `POST http://localhost:8181/api/ratings`
     - `GET http://localhost:8181/api/ratings/provider/{providerId}/average`
   - **Notification Service**:  
     - `GET http://localhost:9191/notifications/{providerId}`

### Endpoints

1. **Create a Rating**  
   ```
   POST /api/ratings
   Content-Type: application/json

   {
     "providerId": 101,
     "score": 5,
     "comment": "Excellent service!"
   }
   ```
   - Saves the rating to PostgreSQL.
   - Publishes a `RatingCreated` event to Kafka.

2. **Get Average Rating**  
   ```
   GET /api/ratings/provider/{providerId}/average
   ```
   - Returns an object containing the `providerId` and `averageScore`.

3. **Get Notifications**  
   ```
   GET /notifications/{providerId}
   ```
   - Returns a list of new notifications for the given `providerId`.
   - Clears them afterward so subsequent calls return an empty list until new notifications arrive.

---

## Configuration

- **docker-compose.yml** defines environment variables, ports, and container dependencies.
- **Rating Service** uses the following environment variables (or `appsettings.json` keys):
  - `ConnectionStrings__DefaultConnection` (PostgreSQL)
  - `Kafka__BootstrapServers`
  - `Kafka__TopicName`
- **Notification Service** uses environment variables:
  - `KAFKA_BOOTSTRAP_SERVERS`
  - `KAFKA_TOPIC`
  - `KAFKA_GROUP_ID`

You can modify the `.NET` or Go code to handle these differently if desired.

---

## Project Structure

```
.
├── docker-compose.yml
├── README.md
├── scripts/
│   ├── build-all.sh      # Builds Docker images for both services
│   └── run.sh            # Runs docker-compose up --build
├── rating-service/
│   ├── Dockerfile
│   ├── RatingService.sln
│   ├── ...
│   └── src/
│       ├── Models/
│       ├── Controllers/
│       ├── Repositories/
│       └── Services/
└── notification-service/
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── cmd/
    │   └── notification-service/
    │       └── main.go
    └── internal/
        ├── api/
        ├── consumer/
        ├── domain/
        └── service/
```

- **rating-service**: Contains .NET 9 WebAPI, EF Core code, Dockerfile.  
- **notification-service**: Contains Go consumer logic and REST API, Dockerfile.  
- **scripts**: Contains build and run scripts for convenience.

---

## Future Improvements

Below are some recommended enhancements and design considerations. We designed the solution to be **simple**, so some production-level features are either omitted or simplified. If quality was partially sacrificed for speed (e.g., no persistent notifications), these choices are documented here along with possible improvements.

- **Use Redis/NoSQL/RDBMS**: Instead of in-memory storage, we can leverage Redis, a NoSQL database, or an RDBMS to enable horizontal scaling and persistent data.
- **Load Balancing**: Run multiple instances of the Notification Service behind a load balancer or ingress to handle high traffic.
- **Rate Limiting & Caching**: Manage heavy polling requests more efficiently and prevent resource exhaustion.
- **Push Mechanism (Optional)**: Replace or complement polling with real-time push notifications for immediate delivery.
- **Enhanced Observability**: Incorporate logging, metrics, and tracing to pinpoint performance bottlenecks.
- **Advanced Error Handling**: Implement retry policies, dead-letter queues, and other mechanisms for fault tolerance.

These improvements make the Notification Service more robust and scalable in **large-scale** environments. Although our prototype (in-memory + single instance) is sufficient for basic needs, such enhancements become crucial if we assume **multiple clients or services** are polling frequently at high volume.

### Additional Considerations

1. **Integration Tests with Testcontainers**  
   Leverage Testcontainers to spin up ephemeral instances of PostgreSQL, Kafka, etc., for real-world integration tests without manual local setups.

2. **CI/CD Pipeline**  
   Set up a pipeline (e.g., GitHub Actions) to automate builds, tests, deployments, and versioning (with semantic release or similar).

3. **Document Main Design Decisions**  
   Note key trade-offs (e.g., in-memory notifications are lost on restart) and justify them in terms of sustainability, reliability, and scalability.

4. **Handling an Untrusted Network**  
   Enable TLS/SSL for Kafka and REST endpoints if the network is untrusted. Consider using a service mesh or API gateway to manage certificates and security policies.

5. **High Traffic Considerations**  
   For very large request volumes, adopt distributed caches or high-performance data stores. Add logging/monitoring solutions (e.g., ELK stack) for better diagnostics.

6. **Monitoring & Observability**  
   Integrate tools such as Prometheus and Grafana, and use correlation IDs for full traceability across microservices.

7. **Other Reliability Enhancements**  
   Consider **dead-letter queues** (DLQ) for unprocessable messages, retry mechanisms, and persistent storage if losing notifications is unacceptable in case of failures.