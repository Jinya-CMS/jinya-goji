---
title: Docker Compose
---

# Installation with Docker Compose

This guide will walk you through the process of setting up Jinya CMS using Docker Compose. This is the easiest way to get Jinya CMS up and running.

## Prerequisites

- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher

## Installation Steps

### Create a `docker-compose.yml`

Create a new directory for your Jinya CMS installation and create a `docker-compose.yml` file with the following content:

```yaml
services:
  database:
    image: postgres:18-alpine
    restart: always
    environment:
      POSTGRES_DB: jinya
      POSTGRES_PASSWORD: jinya
      POSTGRES_ROOT_PASSWORD: jinya
      POSTGRES_USER: jinya
    volumes:
      - db_data:/var/lib/postgres

  jinya-goji:
    image: jinyacms/jinya-goji:latest
    restart: always
    ports:
      - "8090:8090"
    environment:
      DATABASE_URL: postgres://jinya:jinya@database:5432/jinya
      STORAGE_PATH: /var/data/storage
      THEME_PATH: /var/data/themes
      CACHE_PATH: /var/data/cache
    volumes:
      - jinya_data:/var
    depends_on:
      - database

volumes:
  db_data:
  jinya_data:
```

### Start the Containers

Run the following command in the directory where you created the `docker-compose.yml` file:

```bash
docker compose up -d
```

### Run the Installer

Once the containers are started, navigate to `http://localhost:8090` in your web browser. Jinya CMS includes a web-based installer that will guide you through the remaining steps:

- Administrative user creation
- Initial site setup

Follow the on-screen instructions to complete your installation.
