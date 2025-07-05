# PostgreSQL with pgvector Docker Setup

Docker Compose configuration for PostgreSQL with pgvector extension for vector similarity search.

## Quick Start

```bash
# Start database
docker-compose up -d

# Stop database
docker-compose down

# View logs
docker-compose logs -f postgres

# Connect to database
docker exec -it rag_pg psql -U postgres -d postgres
```

## Configuration

- **Image**: `ankane/pgvector:latest`
- **Port**: `5432`
- **Database**: `postgres`
- **Username**: `postgres`
- **Password**: `test123`
- **Connection**: `postgresql://postgres:test123@localhost:5432/postgres`

### Environment Variables

To customize the configuration, you can modify the `docker-compose.yaml` file or use environment variables:

```bash
# Set environment variables before running docker-compose
export POSTGRES_USER=myuser
export POSTGRES_PASSWORD=mypassword
export POSTGRES_DB=mydatabase
export POSTGRES_PORT=5433

# Or create a .env file in this directory
echo "POSTGRES_USER=myuser" > .env
echo "POSTGRES_PASSWORD=mypassword" >> .env
echo "POSTGRES_DB=mydatabase" >> .env
echo "POSTGRES_PORT=5433" >> .env
```

## Features

- PostgreSQL with pgvector extension for vector operations
- Persistent data storage via Docker volume
- Health checks for container monitoring
- UTF-8 encoding support

## Security Note

⚠️ Default credentials for development only. Change password for production use.

## Troubleshooting

- **Port conflict**: Change port mapping in docker-compose.yaml
- **Container issues**: Check logs with `docker-compose logs postgres` 