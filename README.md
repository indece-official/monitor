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
    image: indece/monitor:latest
    restart: always
    depends_on:
      - postgres
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DATABASE=monitor
      - POSTGRES_USER=myuser
      - POSTGRES_PASSWORD=mypassword
      - SMTP_HOST=mymailserver.com
      - SMTP_PORT=465
      - SMTP_USER=bot@mymailserver.com
      - SMTP_PASSWORD=mymailpassword
      - SMTP_FROM=bot@mymailserver.com
    ports:
      - 0.0.0.0:9440:9440
    exposes:
      - 8080
```
