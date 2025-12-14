export const DOMAIN_TYPE_MAP = {
  ONAPI: "onapi",
  SCJ: "scj",
  DGII: "dgii",
  PGR: "pgr",
  GOOGLE_DOCKING: "google_docking",
  DOCKING: "docking",
  SOCIAL_MEDIA: "social_media",
  X_SOCIAL_MEDIA: "x_social_media",
  FILE_TYPE: "file_type",
  ERROR: "error",
} as const;

export type DomainType = typeof DOMAIN_TYPE_MAP[keyof typeof DOMAIN_TYPE_MAP];

export interface DomainSearchResult {
  id?: string;
  success: boolean;
  error?: string;
  name: string;
  search_parameter: string;
  keywords_per_category?: Record<string, string[]>;
  output?: Entity[] | ScjCase[] | PgrNews[] | GoogleDockingResult[];
}


export interface DynamicPipelineStep {
  id?: string;
  domain_type: DomainType;
  search_parameter: string;
  category?: string;
  keywords?: string[];
  success: boolean;
  error?: string;
  output?: any;
  keywords_per_category?: Record<string, string[]>;
  depth: number;
}

export interface DynamicPipelineResult {
  id: string;
  steps: DynamicPipelineStep[];
  created_at: Date;
  updated_at: Date;
  total_steps: number;
  successful_steps: number;
  failed_steps: number;
  max_depth_reached: number;
  config: {
    max_depth: number;
    max_concurrent_steps: number;
    delay_between_steps: number;
    skip_duplicates: boolean;
    available_domains: DomainType[];
    query: string;
  };
}

export interface PipelineResponse {
  execution_id: string;
  message: string;
  status: string;
}

// Legacy response structure (for streaming or completed pipelines)
export interface PipelineResponseComplete {
  execution_id: string;
  dynamic_result?: DynamicPipelineResult;
  pipeline?: DomainSearchResult[];
  summary?: {
    total_steps: number;
    successful_steps: number;
    failed_steps: number;
    max_depth_reached: number;
  };
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  count?: number;
  message?: string;
  type?: string;
}

export interface Entity {
  id: string;
  domain_search_result_id: string;
  serie_expediente: number;
  numero_expediente: number;
  certificado: string;
  tipo: string;
  subtipo: string;
  texto: string;
  [key: string]: any;
}

export interface ScjCase {
  id: string;
  domain_search_result_id: string;
  linea: number;
  id_expediente: number;
  no_expediente: string;
  no_sentencia: string;
  [key: string]: any;
}

export interface Register {
  id: string;
  domain_search_result_id: string;
  rnc: string;
  razon_social: string;
  nombre_comercial: string;
  [key: string]: any;
}

export interface PgrNews {
  id: string;
  domain_search_result_id: string;
  url: string;
  title: string;
  [key: string]: any;
}

export interface GoogleDockingResult {
  id: string;
  domain_search_result_id: string;
  search_parameter: string;
  link: string;
  title: string;
  description: string;
  relevance: number;
  search_rank: number;
  [key: string]: any;
}
