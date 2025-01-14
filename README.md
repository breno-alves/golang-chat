# Chat and Bot Apps

This repository contains a Go (Golang) application designed to run seamlessly within a Docker container. Follow this README to set up, build, and run the application.

## Prerequisites

Ensure you have the following installed on your system:

1. [Docker](https://docs.docker.com/get-docker/)
2. [Git](https://git-scm.com/)

## Getting Started

### Running the Application

To start the application with Docker Compose, run:

```bash
docker-compose up -d
```

This will build the Docker image (if not already built), start the container, and map the appropriate ports as defined in the `docker-compose.yml` file.

To stop the application:

```bash
docker-compose down
```

## Environment Variables

You can configure the application using environment variables defined in a `.env.docker` file or directly in the `docker-compose.yml` file.

### Example `.env` File

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASS=postgres
DB_NAME=mhe
REDIS_HOST=localhost:6379
REDIS_PASS=
EXCHANGER_HOST=amqp://rabbitmq:rabbitmq@localhost:5672/
```

### Passing Environment Variables

Docker Compose will automatically load variables from a `.env.docker` file located in the root directory. Alternatively, you can edit the `docker-compose.yml` file to pass variables directly.

```yaml
services:
  chat:
    environment:
      - ...
```

## TODOs

- [ ] Cache request to external APIs
- [ ] Change broadcasting approach to one worker per room
- [ ] Add load balancer with persistent WS connection


