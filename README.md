# Insightful Intel

Una plataforma inteligente de agregación e investigación de datos que realiza búsquedas exhaustivas en múltiples fuentes de datos gubernamentales y públicas de la República Dominicana. El sistema permite la recopilación automatizada de inteligencia entre dominios mediante la creación dinámica de pipelines de búsqueda que extraen palabras clave de un dominio y las utilizan para buscar en otros dominios relacionados.

## 🎯 Resumen

Insightful Intel automatiza el proceso de recopilación de inteligencia desde múltiples fuentes de datos públicas, incluyendo:

- **ONAPI** (Oficina Nacional de la Propiedad Industrial) - Registros de marcas y patentes
- **SCJ** (Suprema Corte de Justicia) - Registros de casos de la Corte Suprema
- **DGII** (Dirección General de Impuestos Internos) - Registros de la autoridad fiscal
- **PGR** (Procuraduría General de la República) - Noticias de la Procuraduría General
- **Google Docking** - Resultados de búsqueda web con puntuación de relevancia
- **Redes Sociales** - Búsquedas en plataformas de redes sociales
- **Búsquedas por Tipo de Archivo** - Búsquedas de documentos y archivos

## ✨ Características Principales

- **Sistema de Pipeline Dinámico**: Genera automáticamente pasos de búsqueda basados en palabras clave extraídas
- **Búsqueda Multi-Dominio**: Interfaz unificada para buscar en 7+ fuentes de datos
- **Extracción y Categorización de Palabras Clave**: Extracción y categorización automática de palabras clave relevantes
- **Transmisión en Tiempo Real**: Server-Sent Events (SSE) para actualizaciones en vivo del pipeline
- **Almacenamiento Persistente**: Todos los resultados de búsqueda y ejecuciones de pipeline almacenados en MySQL
- **API RESTful**: Endpoints REST limpios para todas las operaciones
- **Herramienta CLI**: Interfaz de línea de comandos para uso automatizado/scripted
- **Diseño Dirigido por Dominio**: Arquitectura bien estructurada siguiendo principios DDD

## 🏗️ Arquitectura

El proyecto sigue una arquitectura de **Diseño Dirigido por Dominio (DDD)** con clara separación de responsabilidades:

```
┌─────────────────────────────────────────┐
│      Capa de Presentación               │
│  (Manejadores HTTP, Frontend React)    │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Capa de Aplicación                  │
│  (Interactores, Casos de Uso)          │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Capa de Dominio                  │
│  (Entidades, Objetos de Valor, Servicios)│
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Capa de Infraestructura            │
│  (Repositorios, Base de Datos, HTTP)   │
└─────────────────────────────────────────┘
```

## 🛠️ Tecnologías

### Backend
- **Go 1.24.2** - Lenguaje principal del backend
- **MySQL** - Base de datos relacional
- **Colly** - Framework de web scraping
- **Cobra** - Framework CLI

### Frontend
- **React 19** - Framework de UI
- **TypeScript** - Desarrollo con tipos seguros
- **Vite** - Herramienta de construcción y servidor de desarrollo
- **Tailwind CSS** - Framework CSS utility-first

### Infraestructura
- **Docker & Docker Compose** - Containerización
- **Make** - Automatización de construcción

## 🚀 Inicio Rápido

### Prerrequisitos

- Go 1.24.2 o posterior
- Node.js 18+ y npm
- Docker y Docker Compose
- MySQL (o usar Docker)

### Instalación

1. **Clonar el repositorio**
   ```bash
   git clone <repository-url>
   cd insightful-intel
   ```

2. **Iniciar la base de datos**
   ```bash
   make docker-run
   ```

3. **Instalar dependencias del frontend**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Configurar variables de entorno**
   Crear un archivo `.env` en el directorio raíz:
   ```env
   BLUEPRINT_DB_HOST=localhost
   BLUEPRINT_DB_PORT=3306
   BLUEPRINT_DB_USER=root
   BLUEPRINT_DB_PASSWORD=password
   BLUEPRINT_DB_NAME=insightful_intel
   ```

5. **Ejecutar la aplicación**
   ```bash
   make run
   ```
   Esto iniciará tanto el servidor API backend como el servidor de desarrollo del frontend.

### Desarrollo

**Iniciar backend con recarga automática:**
```bash
make watch
```

**Ejecutar pruebas:**
```bash
make test          # Pruebas unitarias
make itest         # Pruebas de integración
```

**Construir la aplicación:**
```bash
make build         # Construir servidor API
make build-cli     # Construir herramienta CLI
```

## 📚 Uso

### Endpoints de API

#### Operaciones de Búsqueda
- `GET /search?q={query}&domain={domain}` - Buscar un dominio específico
- `GET /search?q={query}` - Buscar en todos los dominios por defecto
- `GET /dynamic?q={query}&depth={depth}&skip_duplicates={bool}&stream={bool}` - Ejecutar pipeline dinámico

#### Datos Específicos por Dominio
- `GET /api/onapi` - Entidades ONAPI
- `GET /api/scj` - Casos SCJ
- `GET /api/dgii` - Registros DGII
- `GET /api/pgr` - Noticias PGR
- `GET /api/docking` - Resultados de Google Docking

#### Operaciones de Pipeline
- `GET /api/pipeline` - Listar todos los pipelines
- `GET /api/pipeline/steps?pipeline_id={id}` - Obtener pasos del pipeline
- `POST /api/pipeline/save` - Guardar ejecución del pipeline

### Uso de CLI

```bash
# Construir CLI
make build-cli

# Ejecutar pipeline dinámico
./cli run "Novasco" --max-depth 5 --skip-duplicates

# O usar go run
go run cmd/cli/main.go run "Novasco" --max-depth 5
```

### Ejemplo: Pipeline Dinámico

Ejecutar un pipeline dinámico que explora automáticamente entidades relacionadas:

```bash
curl "http://localhost:8080/dynamic?q=Novasco&depth=5&skip_duplicates=true&stream=true"
```

El pipeline:
1. Comienza con la consulta inicial "Novasco"
2. Busca en todos los dominios disponibles
3. Extrae palabras clave de los resultados
4. Crea nuevas búsquedas usando las palabras clave extraídas
5. Continúa hasta la profundidad especificada
6. Transmite resultados en tiempo real

## 📁 Estructura del Proyecto

```
insightful-intel/
├── cmd/                    # Puntos de entrada de la aplicación
│   ├── api/               # Servidor HTTP API
│   └── cli/               # Interfaz de línea de comandos
├── internal/              # Código privado de la aplicación
│   ├── domain/           # Modelos de dominio y lógica de negocio
│   ├── repositories/     # Capa de acceso a datos
│   ├── interactor/       # Casos de uso de la aplicación
│   ├── module/           # Servicios de dominio
│   ├── server/           # Manejadores HTTP y rutas
│   ├── database/         # Conexión a base de datos y migraciones
│   └── infra/            # Preocupaciones de infraestructura
├── frontend/             # Aplicación frontend React
│   ├── src/
│   │   ├── components/  # Componentes React
│   │   ├── pages/       # Componentes de página
│   │   └── api.ts       # Cliente API
├── config/              # Gestión de configuración
├── docs/                # Documentación
└── vendor/              # Dependencias Go
```

## 📖 Documentación

- **[Documentación Completa del Proyecto](docs/PROJECT_DOCUMENTATION.md)** - Guía completa que cubre arquitectura, implementación DDD, casos de uso y más
- **[Implementando un Nuevo Tipo de Dominio](docs/IMPLEMENTING_NEW_DOMAIN_ES.md)** - Guía paso a paso para agregar un nuevo tipo de dominio al sistema
- **[Guía de Mejoras del Proyecto](docs/PROJECT_IMPROVEMENTS.md)** - Recomendaciones para migraciones, estructura de base de datos, pruebas y más
- **[Guía de Pipeline Dinámico](docs/DYNAMIC_PIPELINE_GUIDE_ES.md)** - Explicación detallada del sistema de pipeline dinámico
- **[Uso de Búsqueda por Dominio](docs/DOMAIN_SEARCH_USAGE.md)** - Cómo usar las funciones de búsqueda por dominio
- **[Uso de CLI](docs/CLI_USAGE_ES.md)** - Documentación de la interfaz de línea de comandos

## 🧪 Pruebas

```bash
# Ejecutar todas las pruebas
make test

# Ejecutar pruebas de integración (requiere Docker)
make itest

# Ejecutar pruebas para un paquete específico
go test ./internal/domain/... -v
```

**Build the application:**
```bash
# Iniciar contenedor de base de datos
make docker-run

# Detener contenedor de base de datos
make docker-down

# Construir imagen Docker de API
docker build -f API.Dockerfile -t insightful-intel-api .

# Construir imagen Docker de CLI
docker build -f CLI.Dockerfile -t insightful-intel-cli .
```

## 🤝 Contribuir

1. Fork el repositorio
2. Crear una rama de funcionalidad (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add some amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abrir un Pull Request

## 📝 Licencia

Ver archivo [LICENSE](LICENSE) para más detalles.

## 🔗 Documentación Relacionada

- [README de Capa de Repositorio](docs/REPO_ES.md) - Documentación de la capa de repositorio

---
