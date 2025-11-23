# CLI Usage Guide

This guide explains how to use the Insightful Intel CLI with Docker.

## Prerequisites

- Docker and Docker Compose installed
- Environment variables configured in `.env` file

## Building the CLI Docker Image

```bash
docker-compose -f docker-compose.cli.yml build
```

## Running the CLI

### Basic Usage

Run a dynamic pipeline search with a query:

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco"
```

### With Custom Parameters

#### Set Maximum Depth

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --max-depth 3
```

Or use the short flag:

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" -d 3
```

#### Disable Skipping Duplicates

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --skip-duplicates false
```

Or use the short flag:

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" -s false
```

#### Combined Options

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" -d 7 -s false
```

## Available Flags

- `-d, --max-depth`: Maximum depth for pipeline execution (default: 5)
- `-s, --skip-duplicates`: Skip duplicate searches (default: true)

## Getting Help

```bash
docker-compose -f docker-compose.cli.yml run --rm cli --help
```

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run --help
```

## Development Mode (with live reload)

For development with live reload using Air, you can use the original configuration:

```bash
docker-compose -f docker-compose.cli.yml up
```

Note: The current Docker setup is optimized for production use. For development, you may want to mount volumes and use the Air configuration.

## Environment Variables

Make sure your `.env` file contains:

```env
APP_ENV=development
BLUEPRINT_DB_HOST=mysql_bp
BLUEPRINT_DB_PORT=3306
BLUEPRINT_DB_DATABASE=your_database
BLUEPRINT_DB_USERNAME=your_username
BLUEPRINT_DB_PASSWORD=your_password
BLUEPRINT_DB_ROOT_PASSWORD=root_password
```

## Examples

### Example 1: Basic Search

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "ABC Company"
```

### Example 2: Deep Search with No Duplicate Skipping

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "ABC Company" -d 10 -s false
```

### Example 3: Shallow Search

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "ABC Company" -d 2
```

## How It Works

1. The CLI initializes database connections
2. Runs migrations
3. Executes a dynamic pipeline search across multiple domains:
   - ONAPI (Commercial names)
   - SCJ (Court cases)
   - DGII (Tax registry)
   - PGR (News)
   - Google Docking (Web scraping)

4. The pipeline dynamically expands based on discovered keywords and relationships

## Troubleshooting

### Database Connection Issues

Make sure the MySQL service is healthy:

```bash
docker-compose -f docker-compose.cli.yml ps
```

Check logs:

```bash
docker-compose -f docker-compose.cli.yml logs mysql_bp
```

### Rebuilding After Code Changes

If you've modified the code, rebuild the image:

```bash
docker-compose -f docker-compose.cli.yml build --no-cache
```

