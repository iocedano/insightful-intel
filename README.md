# Insightful Intel

An intelligent data aggregation and investigation platform that performs comprehensive searches across multiple Dominican Republic government and public data sources. The system enables automated cross-domain intelligence gathering by dynamically creating search pipelines that extract keywords from one domain and use them to search other related domains.

## 🎯 Overview

Insightful Intel automates the process of gathering intelligence from multiple public data sources, including:

- **ONAPI** (Oficina Nacional de la Propiedad Industrial) - Trademark and patent registrations
- **SCJ** (Suprema Corte de Justicia) - Supreme Court case records
- **DGII** (Dirección General de Impuestos Internos) - Tax authority registrations
- **PGR** (Procuraduría General de la República) - Attorney General's Office news
- **Google Docking** - Web search results with relevance scoring
- **Social Media** - Social media platform searches
- **File Type Searches** - Document and file searches

## ✨ Key Features

- **Dynamic Pipeline System**: Automatically generates search steps based on extracted keywords
- **Multi-Domain Search**: Unified interface for searching across 7+ data sources
- **Keyword Extraction & Categorization**: Automatic extraction and categorization of relevant keywords
- **Real-Time Streaming**: Server-Sent Events (SSE) for live pipeline updates
- **Persistent Storage**: All search results and pipeline executions stored in MySQL
- **RESTful API**: Clean REST endpoints for all operations
- **CLI Tool**: Command-line interface for automated/scripted usage
- **Domain-Driven Design**: Well-structured architecture following DDD principles

## 🏗️ Architecture

The project follows a **Domain-Driven Design (DDD)** architecture with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│         Presentation Layer               │
│  (HTTP Handlers, React Frontend)        │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Application Layer                  │
│  (Interactors, Use Cases)              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Domain Layer                     │
│  (Entities, Value Objects, Services)   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Infrastructure Layer               │
│  (Repositories, Database, HTTP)        │
└─────────────────────────────────────────┘
```

## 🛠️ Technologies

### Backend
- **Go 1.24.2** - Core backend language
- **MySQL** - Relational database
- **Colly** - Web scraping framework
- **Cobra** - CLI framework

### Frontend
- **React 19** - UI framework
- **TypeScript** - Type-safe development
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework

### Infrastructure
- **Docker & Docker Compose** - Containerization
- **Make** - Build automation

## 🚀 Getting Started

### Prerequisites

- Go 1.24.2 or later
- Node.js 18+ and npm
- Docker and Docker Compose
- MySQL (or use Docker)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd insightful-intel
   ```

2. **Start the database**
   ```bash
   make docker-run
   ```

3. **Install frontend dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=password
   DB_NAME=insightful_intel
   ```

5. **Run the application**
   ```bash
   make run
   ```
   This will start both the backend API server and frontend development server.

### Development

**Start backend with live reload:**
```bash
make watch
```

**Run tests:**
```bash
make test          # Unit tests
make itest         # Integration tests
```

**Build the application:**
```bash
make build         # Build API server
make build-cli     # Build CLI tool
```

## 📚 Usage

### API Endpoints

#### Search Operations
- `GET /search?q={query}&domain={domain}` - Search a specific domain
- `GET /search?q={query}` - Search all default domains
- `GET /dynamic?q={query}&depth={depth}&skip_duplicates={bool}&stream={bool}` - Execute dynamic pipeline

#### Domain-Specific Data
- `GET /api/onapi` - ONAPI entities
- `GET /api/scj` - SCJ cases
- `GET /api/dgii` - DGII registers
- `GET /api/pgr` - PGR news
- `GET /api/docking` - Google Docking results

#### Pipeline Operations
- `GET /api/pipeline` - List all pipelines
- `GET /api/pipeline/steps?pipeline_id={id}` - Get pipeline steps
- `POST /api/pipeline/save` - Save pipeline execution

### CLI Usage

```bash
# Build CLI
make build-cli

# Run dynamic pipeline
./cli run "Novasco" --max-depth 5 --skip-duplicates

# Or use go run
go run cmd/cli/main.go run "Novasco" --max-depth 5
```

### Example: Dynamic Pipeline

Execute a dynamic pipeline that automatically explores related entities:

```bash
curl "http://localhost:8080/dynamic?q=Novasco&depth=5&skip_duplicates=true&stream=true"
```

The pipeline will:
1. Start with initial query "Novasco"
2. Search across all available domains
3. Extract keywords from results
4. Create new searches using extracted keywords
5. Continue up to the specified depth
6. Stream results in real-time

## 📁 Project Structure

```
insightful-intel/
├── cmd/                    # Application entry points
│   ├── api/               # HTTP API server
│   └── cli/               # Command-line interface
├── internal/              # Private application code
│   ├── domain/           # Domain models and business logic
│   ├── repositories/     # Data access layer
│   ├── interactor/       # Application use cases
│   ├── module/           # Domain services
│   ├── server/           # HTTP handlers and routes
│   ├── database/         # Database connection and migrations
│   └── infra/            # Infrastructure concerns
├── frontend/             # React frontend application
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   └── api.ts       # API client
├── config/              # Configuration management
├── doc/                 # Documentation
└── vendor/              # Go dependencies
```

## 📖 Documentation

- **[Complete Project Documentation](PROJECT_DOCUMENTATION.md)** - Comprehensive guide covering architecture, DDD implementation, use cases, and more
- **[Dynamic Pipeline Guide](DYNAMIC_PIPELINE_GUIDE.md)** - Detailed explanation of the dynamic pipeline system
- **[Google Docking Builder](GOOGLE_DOCKING_BUILDER.md)** - Google Docking search system documentation
- **[Domain Search Usage](DOMAIN_SEARCH_USAGE.md)** - How to use domain search functions
- **[CLI Usage](CLI_USAGE.md)** - Command-line interface documentation

## 🧪 Testing

```bash
# Run all tests
make test

# Run integration tests (requires Docker)
make itest

# Run tests for specific package
go test ./internal/domain/... -v
```

## 🐳 Docker

```bash
# Start database container
make docker-run

# Stop database container
make docker-down

# Build API Docker image
docker build -f API.Dockerfile -t insightful-intel-api .

# Build CLI Docker image
docker build -f CLI.Dockerfile -t insightful-intel-cli .
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

See [LICENSE](LICENSE) file for details.

## 🔗 Related Documentation

- [Sources](doc/Sources.md) - List of data sources and their purposes
- [Repository Layer README](internal/repositories/README.md) - Repository layer documentation
- [Pipeline Save Endpoint](PIPELINE_SAVE_ENDPOINT.md) - Pipeline save endpoint documentation
- [Streaming Pipeline Guide](STREAMING_PIPELINE_GUIDE.md) - Real-time streaming implementation

---

**Built with ❤️ using Domain-Driven Design principles**
