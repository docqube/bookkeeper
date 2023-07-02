# bookkeeper

<img align="right" alt="gopher" width="180px" src="docs/assets/gopher.png">

Bookkeeper is an application set of an Angular WebApp and simple API written in Go to manage your
income and expenses of a month. You define some categories with RegEx rules to automatically
match your transactions and upload your bank statement as CSV file. The backend will then
automatically assign the categories to the transactions.

## Features

### API

- Import your bank statement as CSV file (see supported banks below)
- Define categories in your PostgreSQL database
- Create matching rules with RegEx for your categories
- Completely hosted by **yourself**, nothing leaves your system

### Supported banks

- ING

### WebApp

- View your transactions in a nice dashboard
- Beautiful in light and dark mode
- Developed for a mobile first experience

#### Mobile

<div style="display:flex;flex-direction:row">
  <img alt="Mobile Dashboard (light)" width="350px" src="docs/assets/dashboard_iphone_light.png">
  <img alt="Mobile Dashboard (dark)" width="350px" src="docs/assets/dashboard_iphone_dark.png">
</div>

#### Desktop

![](docs/assets/dashboard.png)

## Installation

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)

The recommended way to run bookkeeper is a containerized setup using Docker. As we need a running
PostgreSQL instance as a datastore for the API, it is recommended to use a docker-compose setup.

### API Environment Variables

| Name | Description | Default | Required |
|---|---|---|---|
| `BOOKKEEPER_PORT` | Port of the HTTP server | `8080` | ❌ |
| `BOOKKEEPER_DATABASE_HOST` | Hostname of the PostgreSQL database | - | ✅ |
| `BOOKKEEPER_DATABASE_PORT` | Port of the PostgreSQL database | - | ✅ |
| `BOOKKEEPER_DATABASE_USERNAME` | Username used to connect to PostgreSQL database | - | ✅ |
| `BOOKKEEPER_DATABASE_PASSWORD` | Password used to connect to PostgreSQL database | - | ✅ |
| `BOOKKEEPER_DATABASE_NAME` | Used PostgreSQL database name | - | ✅ |
| `BOOKKEEPER_DATABASE_SSLMODE` | Used PostgreSQL SSL mode (e.g. `"disabled"` if running on localhost) | - | ❌ |

### Docker Compose Setup

The following `docker-compose.yaml` file will setup a minimal bookkeeper instance, with a PostgreSQL instance,
Traefik as a reverse proxy and the API and WebApp containers. You'll need to create a `web` network using
`docker network create web`.

```yaml
services:

  traefik:
    image: traefik:v2.10
    restart: always
    ports:
      - 80:80
      - 443:443
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    expose:
      - "8080"
    networks:
      web:
        aliases:
          - traefik

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=bookkeeper
      - POSTGRES_USER=bookkeeper
      - POSTGRES_PASSWORD=s3cr3t
    restart: always
    networks:
      web:
        aliases: 
          - postgres

  backend:
    image: ghcr.io/docqube/bookkeeper/backend:latest
    labels:
      - traefik.enable=true
      - "traefik.http.routers.backend.rule=Host(`bookkeeper.example.com`) && PathPrefix(`/api`)"
      - traefik.http.services.backend.loadbalancer.server.port=8080
    environment:
      - BOOKKEEPER_DATABASE_HOST=postgres
      - BOOKKEEPER_DATABASE_PORT=5432
      - BOOKKEEPER_DATABASE_NAME=bookkeeper
      - BOOKKEEPER_DATABASE_USER=bookeeper
      - BOOKKEEPER_DATABASE_PASSWORD=s3cr3t
    expose:
      - "8080"
    restart: always
    depends_on:
      - postgres
    networks:
      web:
        aliases: 
          - backend
  
  frontend:
    image: ghcr.io/docqube/bookkeeper/frontend:latest
    labels:
      - traefik.enable=true
      - traefik.http.routers.frontend.rule=Host(`bookkeeper.example.com`)
    expose:
      - "80"
    restart: always
    depends_on:
      - backend
    networks:
      web:
        aliases: 
          - frontend

networks:
  web:
    external: true
```
