# CLI Usage Guide

This guide explains how to use the Insightful Intel Command-Line Interface (CLI) tool to execute dynamic pipeline searches across multiple data sources.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Basic Usage](#basic-usage)
4. [Command Reference](#command-reference)
5. [Options and Flags](#options-and-flags)
6. [Examples](#examples)
7. [Docker Usage](#docker-usage)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- Go 1.24.2 or later
- MySQL database (or Docker)
- Environment variables configured (see [Environment Variables](#environment-variables))

---

## Installation

### Build from Source

```bash
# Build the CLI binary
make build-cli

# Or use go build directly
go build -o cli cmd/cli/main.go
```

The binary will be created in the current directory as `cli` (or `cli.exe` on Windows).

### Using Docker

See the [Docker Usage](#docker-usage) section for containerized execution.

---

## Basic Usage

### Run Command

The main command is `run`, which executes a dynamic pipeline search:

```bash
./cli run "Novasco"
```

This will:
1. Initialize the database connection
2. Run migrations if needed
3. Execute a dynamic pipeline search across all available domains
4. Display the results

---

## Command Reference

### `run [query]`

Executes a dynamic pipeline search with the specified query.

**Syntax:**
```bash
./cli run <query> [flags]
```

**Arguments:**
- `query` (required): The search query string (e.g., company name, RNC, person name)

**Example:**
```bash
./cli run "ABC Company"
```

---

## Options and Flags

### `--max-depth, -d`

Sets the maximum depth for pipeline execution.

- **Type**: Integer
- **Default**: `5`
- **Range**: 1-10 (recommended)
- **Description**: Controls how many levels deep the pipeline will explore. Higher values mean more comprehensive searches but longer execution times.

**Examples:**
```bash
# Shallow search (depth 2)
./cli run "Novasco" -d 2

# Deep search (depth 7)
./cli run "Novasco" --max-depth 7
```

### `--skip-duplicates, -s`

Controls whether to skip duplicate keyword searches across domains.

- **Type**: Boolean
- **Default**: `true`
- **Description**: When enabled, prevents searching the same keyword multiple times across different domains, improving performance and reducing redundant API calls.

**Examples:**
```bash
# Skip duplicates (default)
./cli run "Novasco" -s true

# Allow duplicates
./cli run "Novasco" --skip-duplicates false
```

### Combined Flags

You can combine multiple flags:

```bash
./cli run "Novasco" -d 7 -s false
```

---

## Examples

### Example 1: Basic Search

Simple search with default settings:

```bash
./cli run "ABC Company"
```

**What happens:**
- Maximum depth: 5
- Skip duplicates: true
- Searches across: ONAPI, SCJ, DGII, PGR, Google Docking

### Example 2: Deep Search

Comprehensive investigation with maximum depth:

```bash
./cli run "ABC Company" --max-depth 10
```

**Use case:** When you need a thorough investigation and have time for longer execution.

### Example 3: Shallow Quick Search

Quick search with minimal depth:

```bash
./cli run "ABC Company" -d 2
```

**Use case:** When you need quick results without deep exploration.

### Example 4: Allow Duplicates

Search that allows duplicate keyword searches:

```bash
./cli run "ABC Company" -d 5 -s false
```

**Use case:** When you want to ensure all possible connections are explored, even if it means redundant searches.

### Example 5: Search with RNC

Search using a tax identification number:

```bash
./cli run "123456789"
```

---

## Docker Usage

### Prerequisites

- Docker and Docker Compose installed
- Environment variables configured in `.env` file

### Building the CLI Docker Image

```bash
docker-compose -f docker-compose.cli.yml build
```

### Running with Docker

#### Basic Usage

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco"
```

#### With Custom Parameters

```bash
# Set maximum depth
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --max-depth 3

# Disable skipping duplicates
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --skip-duplicates false

# Combined options
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" -d 7 -s false
```

### Getting Help

```bash
# General help
docker-compose -f docker-compose.cli.yml run --rm cli --help

# Command-specific help
docker-compose -f docker-compose.cli.yml run --rm cli run --help
```

### Development Mode

For development with live reload using Air:

```bash
docker-compose -f docker-compose.cli.yml up
```

---

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=insightful_intel

# Application Environment
APP_ENV=development
```

### Docker Environment Variables

For Docker usage, ensure your `.env` file contains:

```env
APP_ENV=development
BLUEPRINT_DB_HOST=mysql_bp
BLUEPRINT_DB_PORT=3306
BLUEPRINT_DB_DATABASE=your_database
BLUEPRINT_DB_USERNAME=your_username
BLUEPRINT_DB_PASSWORD=your_password
BLUEPRINT_DB_ROOT_PASSWORD=root_password
```

---

## How It Works

1. **Initialization**: The CLI connects to the MySQL database and runs migrations
2. **Pipeline Creation**: Creates a dynamic pipeline structure based on available domains
3. **Domain Search**: Executes searches across multiple domains:
   - **ONAPI**: Commercial names and trademarks
   - **SCJ**: Court cases and legal records
   - **DGII**: Tax registry and RNC information
   - **PGR**: News and public announcements
   - **Google Docking**: Web search results with relevance scoring
4. **Keyword Extraction**: Extracts keywords from search results
5. **Dynamic Expansion**: Creates new search steps based on discovered keywords
6. **Result Aggregation**: Combines and displays all results

### Pipeline Flow

```
Initial Query
    ↓
ONAPI Search → Extract Keywords
    ↓
SCJ Search (using person names) → Extract Keywords
    ↓
DGII Search (using company names) → Extract Keywords
    ↓
PGR Search → Extract Keywords
    ↓
Google Docking Search → Extract Keywords
    ↓
Repeat up to MaxDepth
```

---

## Getting Help

### Command Help

```bash
# General help
./cli --help

# Help for run command
./cli run --help
```

### Available Commands

```bash
./cli --help
```

Output:
```
A CLI tool for running dynamic pipeline searches across multiple domains

Usage:
  cli [command]

Available Commands:
  run         Run dynamic pipeline search

Flags:
  -h, --help   help for cli

Use "cli [command] --help" for more information about a command.
```

---

## Troubleshooting

### Database Connection Issues

**Problem**: Cannot connect to database

**Solutions**:
1. Verify MySQL is running:
   ```bash
   # Check MySQL service
   mysql -u root -p -e "SELECT 1"
   ```

2. Check environment variables:
   ```bash
   # Verify .env file exists and has correct values
   cat .env
   ```

3. Test database connection:
   ```bash
   mysql -h localhost -u root -p insightful_intel
   ```

### Migration Errors

**Problem**: Migration failures

**Solutions**:
1. Check database schema:
   ```bash
   mysql -u root -p insightful_intel -e "SHOW TABLES;"
   ```

2. Verify database user permissions:
   ```bash
   mysql -u root -p -e "SHOW GRANTS FOR 'your_user'@'localhost';"
   ```

### Docker Issues

**Problem**: Docker container fails to start

**Solutions**:
1. Check container status:
   ```bash
   docker-compose -f docker-compose.cli.yml ps
   ```

2. View logs:
   ```bash
   docker-compose -f docker-compose.cli.yml logs mysql_bp
   docker-compose -f docker-compose.cli.yml logs cli
   ```

3. Rebuild after code changes:
   ```bash
   docker-compose -f docker-compose.cli.yml build --no-cache
   ```

### Performance Issues

**Problem**: CLI execution is slow

**Solutions**:
1. Reduce max depth:
   ```bash
   ./cli run "query" -d 2
   ```

2. Enable skip duplicates:
   ```bash
   ./cli run "query" -s true
   ```

3. Check database performance:
   ```bash
   mysql -u root -p -e "SHOW PROCESSLIST;"
   ```

### No Results Returned

**Problem**: Pipeline returns no results

**Solutions**:
1. Verify query is correct
2. Check if domains are accessible
3. Review logs for errors
4. Try a different query format
5. Verify database has data from previous searches

---

## Advanced Usage

### Scripting

You can use the CLI in scripts:

```bash
#!/bin/bash
QUERIES=("Company A" "Company B" "Company C")

for query in "${QUERIES[@]}"; do
    echo "Searching for: $query"
    ./cli run "$query" -d 3
    echo "---"
done
```

### Output Formatting

The CLI outputs results in a structured format. You can pipe output to files:

```bash
./cli run "Novasco" > results.json
```

### Integration with Other Tools

The CLI can be integrated with monitoring tools, schedulers, or other automation systems:

```bash
# Example: Run with cron
0 2 * * * /path/to/cli run "Daily Search" -d 3
```

---

## Best Practices

1. **Start with Shallow Searches**: Begin with `-d 2` or `-d 3` to get quick results
2. **Use Appropriate Depth**: Higher depth (5+) should be used for comprehensive investigations
3. **Enable Skip Duplicates**: Keep `-s true` unless you specifically need duplicate searches
4. **Monitor Execution Time**: Deep searches can take significant time
5. **Check Database Health**: Ensure MySQL is properly configured and has adequate resources
6. **Review Logs**: Check logs for errors or warnings during execution

---

## Related Documentation

- [Dynamic Pipeline Guide](DYNAMIC_PIPELINE_GUIDE.md) - Detailed explanation of the dynamic pipeline system
- [Implementing New Domain](IMPLEMENTING_NEW_DOMAIN.md) - How to add new domain types
- [Project Documentation](PROJECT_DOCUMENTATION.md) - Complete project overview

---

## Support

For issues, questions, or contributions, please refer to the main project documentation or open an issue in the repository.
