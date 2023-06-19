# indece Monitor
Simple & secure server monitoring

Docker: https://hub.docker.com/r/indece/monitor

# Usage

docker-compose.yml
```
version: '3'
services:
  postgres:
    image: postgres:14
    restart: always
    environment:
      - POSTGRES_USER=myuser
      - POSTGRES_PASSWORD=mypassword
      - POSTGRES_DB=monitor
    volumes:
      - /opt/postgres-data:/var/lib/postgresql/data

  monitor:
    image: indece/monitor:v1.0.0-alpha.3
    restart: always
    depends_on:
      - postgres
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DATABASE=monitor
      - POSTGRES_USER=myuser
      - POSTGRES_PASSWORD=mypassword
    ports:
      - 0.0.0.0:9440:9440
    exposes:
      - 8080
```
