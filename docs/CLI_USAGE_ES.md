# Guía de Uso de CLI

Esta guía explica cómo usar la herramienta de Interfaz de Línea de Comandos (CLI) de Insightful Intel para ejecutar búsquedas de pipeline dinámico a través de múltiples fuentes de datos.

## Tabla de Contenidos

1. [Prerrequisitos](#prerrequisitos)
2. [Instalación](#instalación)
3. [Uso Básico](#uso-básico)
4. [Referencia de Comandos](#referencia-de-comandos)
5. [Opciones y Banderas](#opciones-y-banderas)
6. [Ejemplos](#ejemplos)
7. [Uso con Docker](#uso-con-docker)
8. [Solución de Problemas](#solución-de-problemas)

---

## Prerrequisitos

- Go 1.24.2 o posterior
- Base de datos MySQL (o Docker)
- Variables de entorno configuradas (ver [Variables de Entorno](#variables-de-entorno))

---

## Instalación

### Construir desde el Código Fuente

```bash
# Construir el binario CLI
make build-cli

# O usar go build directamente
go build -o cli cmd/cli/main.go
```

El binario se creará en el directorio actual como `cli` (o `cli.exe` en Windows).

### Usando Docker

Ver la sección [Uso con Docker](#uso-con-docker) para ejecución containerizada.

---

## Uso Básico

### Comando Run

El comando principal es `run`, que ejecuta una búsqueda de pipeline dinámico:

```bash
./cli run "Novasco"
```

Esto:
1. Inicializa la conexión a la base de datos
2. Ejecuta migraciones si es necesario
3. Ejecuta una búsqueda de pipeline dinámico a través de todos los dominios disponibles
4. Muestra los resultados

---

## Referencia de Comandos

### `run [query]`

Ejecuta una búsqueda de pipeline dinámico con la consulta especificada.

**Sintaxis:**
```bash
./cli run <query> [flags]
```

**Argumentos:**
- `query` (requerido): La cadena de consulta de búsqueda (por ejemplo, nombre de empresa, RNC, nombre de persona)

**Ejemplo:**
```bash
./cli run "ABC Company"
```

---

## Opciones y Banderas

### `--max-depth, -d`

Establece la profundidad máxima para la ejecución del pipeline.

- **Tipo**: Entero
- **Por Defecto**: `5`
- **Rango**: 1-10 (recomendado)
- **Descripción**: Controla cuántos niveles profundos explorará el pipeline. Valores más altos significan búsquedas más completas pero tiempos de ejecución más largos.

**Ejemplos:**
```bash
# Búsqueda superficial (profundidad 2)
./cli run "Novasco" -d 2

# Búsqueda profunda (profundidad 7)
./cli run "Novasco" --max-depth 7
```

### `--skip-duplicates, -s`

Controla si se omiten búsquedas de palabras clave duplicadas entre dominios.

- **Tipo**: Booleano
- **Por Defecto**: `true`
- **Descripción**: Cuando está habilitado, previene buscar la misma palabra clave múltiples veces a través de diferentes dominios, mejorando el rendimiento y reduciendo llamadas redundantes a la API.

**Ejemplos:**
```bash
# Omitir duplicados (por defecto)
./cli run "Novasco" -s true

# Permitir duplicados
./cli run "Novasco" --skip-duplicates false
```

### Banderas Combinadas

Puedes combinar múltiples banderas:

```bash
./cli run "Novasco" -d 7 -s false
```

---

## Ejemplos

### Ejemplo 1: Búsqueda Básica

Búsqueda simple con configuración por defecto:

```bash
./cli run "ABC Company"
```

**Qué sucede:**
- Profundidad máxima: 5
- Omitir duplicados: true
- Busca en: ONAPI, SCJ, DGII, PGR, Google Docking

### Ejemplo 2: Búsqueda Profunda

Investigación completa con profundidad máxima:

```bash
./cli run "ABC Company" --max-depth 10
```

**Caso de uso:** Cuando necesitas una investigación exhaustiva y tienes tiempo para una ejecución más larga.

### Ejemplo 3: Búsqueda Rápida Superficial

Búsqueda rápida con profundidad mínima:

```bash
./cli run "ABC Company" -d 2
```

**Caso de uso:** Cuando necesitas resultados rápidos sin exploración profunda.

### Ejemplo 4: Permitir Duplicados

Búsqueda que permite búsquedas de palabras clave duplicadas:

```bash
./cli run "ABC Company" -d 5 -s false
```

**Caso de uso:** Cuando quieres asegurar que todas las conexiones posibles sean exploradas, incluso si significa búsquedas redundantes.

### Ejemplo 5: Búsqueda con RNC

Búsqueda usando un número de identificación fiscal:

```bash
./cli run "123456789"
```

---

## Uso con Docker

### Prerrequisitos

- Docker y Docker Compose instalados
- Variables de entorno configuradas en archivo `.env`

### Construir la Imagen Docker de CLI

```bash
docker-compose -f docker-compose.cli.yml build
```

### Ejecutar con Docker

#### Uso Básico

```bash
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco"
```

#### Con Parámetros Personalizados

```bash
# Establecer profundidad máxima
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --max-depth 3

# Deshabilitar omisión de duplicados
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" --skip-duplicates false

# Opciones combinadas
docker-compose -f docker-compose.cli.yml run --rm cli run "Novasco" -d 7 -s false
```

### Obtener Ayuda

```bash
# Ayuda general
docker-compose -f docker-compose.cli.yml run --rm cli --help

# Ayuda específica de comando
docker-compose -f docker-compose.cli.yml run --rm cli run --help
```

### Modo de Desarrollo

Para desarrollo con recarga automática usando Air:

```bash
docker-compose -f docker-compose.cli.yml up
```

---

## Variables de Entorno

Crear un archivo `.env` en el directorio raíz con las siguientes variables:

```env
# Configuración de Base de Datos
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=tu_contraseña
DB_NAME=insightful_intel

# Entorno de Aplicación
APP_ENV=development
```

### Variables de Entorno para Docker

Para uso con Docker, asegurar que tu archivo `.env` contenga:

```env
APP_ENV=development
BLUEPRINT_DB_HOST=mysql_bp
BLUEPRINT_DB_PORT=3306
BLUEPRINT_DB_DATABASE=tu_base_de_datos
BLUEPRINT_DB_USERNAME=tu_usuario
BLUEPRINT_DB_PASSWORD=tu_contraseña
BLUEPRINT_DB_ROOT_PASSWORD=contraseña_root
```

---

## Cómo Funciona

1. **Inicialización**: El CLI se conecta a la base de datos MySQL y ejecuta migraciones
2. **Creación de Pipeline**: Crea una estructura de pipeline dinámico basada en dominios disponibles
3. **Búsqueda de Dominio**: Ejecuta búsquedas a través de múltiples dominios:
   - **ONAPI**: Nombres comerciales y marcas
   - **SCJ**: Casos judiciales y registros legales
   - **DGII**: Registro fiscal e información de RNC
   - **PGR**: Noticias y anuncios públicos
   - **Google Docking**: Resultados de búsqueda web con puntuación de relevancia
4. **Extracción de Palabras Clave**: Extrae palabras clave de los resultados de búsqueda
5. **Expansión Dinámica**: Crea nuevos pasos de búsqueda basados en palabras clave descubiertas
6. **Agregación de Resultados**: Combina y muestra todos los resultados

### Flujo del Pipeline

```
Consulta Inicial
    ↓
Búsqueda ONAPI → Extraer Palabras Clave
    ↓
Búsqueda SCJ (usando nombres de personas) → Extraer Palabras Clave
    ↓
Búsqueda DGII (usando nombres de empresas) → Extraer Palabras Clave
    ↓
Búsqueda PGR → Extraer Palabras Clave
    ↓
Búsqueda Google Docking → Extraer Palabras Clave
    ↓
Repetir hasta MaxDepth
```

---

## Obtener Ayuda

### Ayuda de Comandos

```bash
# Ayuda general
./cli --help

# Ayuda para comando run
./cli run --help
```

### Comandos Disponibles

```bash
./cli --help
```

Salida:
```
Una herramienta CLI para ejecutar búsquedas de pipeline dinámico a través de múltiples dominios

Uso:
  cli [comando]

Comandos Disponibles:
  run         Ejecutar búsqueda de pipeline dinámico

Banderas:
  -h, --help   ayuda para cli

Usa "cli [comando] --help" para más información sobre un comando.
```

---

## Solución de Problemas

### Problemas de Conexión a Base de Datos

**Problema**: No se puede conectar a la base de datos

**Soluciones**:
1. Verificar que MySQL esté ejecutándose:
   ```bash
   # Verificar servicio MySQL
   mysql -u root -p -e "SELECT 1"
   ```

2. Verificar variables de entorno:
   ```bash
   # Verificar que el archivo .env existe y tiene valores correctos
   cat .env
   ```

3. Probar conexión a base de datos:
   ```bash
   mysql -h localhost -u root -p insightful_intel
   ```

### Errores de Migración

**Problema**: Fallos en migraciones

**Soluciones**:
1. Verificar esquema de base de datos:
   ```bash
   mysql -u root -p insightful_intel -e "SHOW TABLES;"
   ```

2. Verificar permisos de usuario de base de datos:
   ```bash
   mysql -u root -p -e "SHOW GRANTS FOR 'tu_usuario'@'localhost';"
   ```

### Problemas con Docker

**Problema**: El contenedor Docker falla al iniciar

**Soluciones**:
1. Verificar estado del contenedor:
   ```bash
   docker-compose -f docker-compose.cli.yml ps
   ```

2. Ver logs:
   ```bash
   docker-compose -f docker-compose.cli.yml logs mysql_bp
   docker-compose -f docker-compose.cli.yml logs cli
   ```

3. Reconstruir después de cambios en el código:
   ```bash
   docker-compose -f docker-compose.cli.yml build --no-cache
   ```

### Problemas de Rendimiento

**Problema**: La ejecución del CLI es lenta

**Soluciones**:
1. Reducir profundidad máxima:
   ```bash
   ./cli run "query" -d 2
   ```

2. Habilitar omitir duplicados:
   ```bash
   ./cli run "query" -s true
   ```

3. Verificar rendimiento de base de datos:
   ```bash
   mysql -u root -p -e "SHOW PROCESSLIST;"
   ```

### No se Retornan Resultados

**Problema**: El pipeline no retorna resultados

**Soluciones**:
1. Verificar que la consulta sea correcta
2. Verificar que los dominios sean accesibles
3. Revisar logs para errores
4. Intentar un formato de consulta diferente
5. Verificar que la base de datos tenga datos de búsquedas anteriores

---

## Uso Avanzado

### Scripting

Puedes usar el CLI en scripts:

```bash
#!/bin/bash
QUERIES=("Company A" "Company B" "Company C")

for query in "${QUERIES[@]}"; do
    echo "Buscando: $query"
    ./cli run "$query" -d 3
    echo "---"
done
```

### Formato de Salida

El CLI muestra resultados en un formato estructurado. Puedes redirigir la salida a archivos:

```bash
./cli run "Novasco" > resultados.json
```

### Integración con Otras Herramientas

El CLI puede integrarse con herramientas de monitoreo, programadores u otros sistemas de automatización:

```bash
# Ejemplo: Ejecutar con cron
0 2 * * * /ruta/a/cli run "Búsqueda Diaria" -d 3
```

---

## Mejores Prácticas

1. **Comenzar con Búsquedas Superficiales**: Comienza con `-d 2` o `-d 3` para obtener resultados rápidos
2. **Usar Profundidad Apropiada**: Profundidad más alta (5+) debe usarse para investigaciones completas
3. **Habilitar Omitir Duplicados**: Mantener `-s true` a menos que específicamente necesites búsquedas duplicadas
4. **Monitorear Tiempo de Ejecución**: Las búsquedas profundas pueden tomar tiempo significativo
5. **Verificar Salud de Base de Datos**: Asegurar que MySQL esté correctamente configurado y tenga recursos adecuados
6. **Revisar Logs**: Verificar logs para errores o advertencias durante la ejecución

---

## Documentación Relacionada

- [Guía de Pipeline Dinámico](DYNAMIC_PIPELINE_GUIDE_ES.md) - Explicación detallada del sistema de pipeline dinámico
- [Implementando Nuevo Dominio](IMPLEMENTING_NEW_DOMAIN_ES.md) - Cómo agregar nuevos tipos de dominio
- [Documentación del Proyecto](../PROJECT_DOCUMENTATION.md) - Resumen completo del proyecto

---

## Soporte

Para problemas, preguntas o contribuciones, por favor consulta la documentación principal del proyecto o abre un issue en el repositorio.
