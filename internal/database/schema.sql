-- Database schema for Insightful Intel repositories
-- This file contains all the necessary table definitions for the repository layer

-- Domain search results table (must be created first as it's referenced by other tables)
CREATE TABLE IF NOT EXISTS  domain_search_results (
    id CHAR(36) PRIMARY KEY,
    success BOOLEAN DEFAULT FALSE,
    error_message TEXT,
    domain_type VARCHAR(50) NOT NULL,
    search_parameter VARCHAR(255),
    keywords_per_category JSON,
    output JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_domain_type (domain_type),
    INDEX idx_success (success),
    INDEX idx_search_parameter (search_parameter),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ONAPI entities table
CREATE TABLE IF NOT EXISTS onapi_entities (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    serie_expediente INT NOT NULL,
    numero_expediente INT NOT NULL,
    certificado VARCHAR(255),
    tipo VARCHAR(100),
    subtipo VARCHAR(100),
    texto TEXT,
    clases TEXT,
    aplicado_a_proteger TEXT,
    expedicion VARCHAR(100),
    vencimiento VARCHAR(100),
    en_tramite BOOLEAN DEFAULT FALSE,
    titular VARCHAR(255),
    gestor VARCHAR(255),
    domicilio TEXT,
    status VARCHAR(100),
    tipo_signo VARCHAR(100),
    imagenes JSON,
    lista_clases JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_serie_numero (serie_expediente, numero_expediente),
    INDEX idx_texto (texto(100)),
    INDEX idx_titular (titular),
    INDEX idx_gestor (gestor),
    INDEX idx_domicilio (domicilio(100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- SCJ cases table
CREATE TABLE IF NOT EXISTS scj_cases (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    linea INT,
    agno_cabecera INT,
    mes_cabecera INT,
    url_cabecera TEXT,
    url_cuerpo TEXT,
    id_expediente INT NOT NULL,
    no_expediente VARCHAR(255),
    no_sentencia VARCHAR(255),
    no_unico VARCHAR(255),
    no_interno VARCHAR(255),
    id_tribunal VARCHAR(100),
    desc_tribunal TEXT,
    id_materia VARCHAR(100),
    desc_materia TEXT,
    fecha_fallo VARCHAR(100),
    involucrados TEXT,
    guid_blob VARCHAR(255),
    tipo_documento_adjunto VARCHAR(100),
    total_filas INT,
    url_blob TEXT,
    extension VARCHAR(10),
    origen INT,
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    UNIQUE KEY unique_expediente (id_expediente),
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_no_expediente (no_expediente),
    INDEX idx_no_sentencia (no_sentencia),
    INDEX idx_involucrados (involucrados(100)),
    INDEX idx_desc_tribunal (desc_tribunal(100)),
    INDEX idx_desc_materia (desc_materia(100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- DGII registers table
CREATE TABLE IF NOT EXISTS  dgii_registers (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    rnc VARCHAR(20) NOT NULL UNIQUE,
    razon_social VARCHAR(255),
    nombre_comercial VARCHAR(255),
    categoria VARCHAR(100),
    regimen_pagos VARCHAR(100),
    facturador_electronico VARCHAR(100),
    licencia_comercial VARCHAR(100),
    estado VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_rnc (rnc),
    INDEX idx_razon_social (razon_social),
    INDEX idx_nombre_comercial (nombre_comercial),
    INDEX idx_categoria (categoria),
    INDEX idx_estado (estado)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- PGR news table
CREATE TABLE IF NOT EXISTS  pgr_news (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    url TEXT NOT NULL,
    title TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    UNIQUE KEY unique_url (url(255)),
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_title (title(100))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Google Docking results table
CREATE TABLE IF NOT EXISTS  google_docking_results (
    id CHAR(36) PRIMARY KEY,
    domain_search_result_id CHAR(36),
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    relevance DECIMAL(3,2) DEFAULT 0.00,
    search_rank INT DEFAULT 0,
    keywords JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (domain_search_result_id) REFERENCES domain_search_results(id) ON DELETE CASCADE,
    INDEX idx_domain_search_result_id (domain_search_result_id),
    INDEX idx_url (url(255)),
    INDEX idx_title (title(100)),
    INDEX idx_relevance (relevance),
    INDEX idx_rank (search_rank)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dynamic pipeline results table
CREATE TABLE IF NOT EXISTS  dynamic_pipeline_results (
    id CHAR(36) PRIMARY KEY,
    total_steps INT DEFAULT 0,
    successful_steps INT DEFAULT 0,
    failed_steps INT DEFAULT 0,
    max_depth_reached INT DEFAULT 0,
    config JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_total_steps (total_steps),
    INDEX idx_successful_steps (successful_steps),
    INDEX idx_failed_steps (failed_steps),
    INDEX idx_max_depth (max_depth_reached),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dynamic pipeline steps table
CREATE TABLE IF NOT EXISTS  dynamic_pipeline_steps  (
    id CHAR(36) PRIMARY KEY,
    pipeline_id CHAR(36) NOT NULL,
    domain_type VARCHAR(50) NOT NULL,
    search_parameter VARCHAR(255),
    category VARCHAR(50),
    keywords JSON,
    success BOOLEAN DEFAULT FALSE,
    error_message TEXT,
    output JSON,
    keywords_per_category JSON,
    depth INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (pipeline_id) REFERENCES dynamic_pipeline_results(id) ON DELETE CASCADE,
    INDEX idx_pipeline_id (pipeline_id),
    INDEX idx_domain_type (domain_type),
    INDEX idx_category (category),
    INDEX idx_success (success),
    INDEX idx_depth (depth),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
