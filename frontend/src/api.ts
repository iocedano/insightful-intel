import type {
  DomainSearchResult,
  Entity,
  ScjCase,
  Register,
  PgrNews,
  GoogleDockingResult,
  PipelineResponse,
  DynamicPipelineResult,
} from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  let response: Response;
  try {
    // Build headers - only set Content-Type for requests with a body
    const headers: Record<string, string> = {};
    
    // Copy existing headers if they exist
    if (options?.headers) {
      if (options.headers instanceof Headers) {
        options.headers.forEach((value, key) => {
          headers[key] = value;
        });
      } else if (Array.isArray(options.headers)) {
        options.headers.forEach(([key, value]) => {
          headers[key] = value;
        });
      } else {
        Object.assign(headers, options.headers);
      }
    }
    
    // Only set Content-Type for requests that have a body (POST, PUT, PATCH, etc.)
    headers['Content-Type'] = 'application/json';

    
    response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    console.log('response', response.ok);
    console.log('response', response.status);
    console.log('response', response.statusText);
    console.log('response', response.headers);
    console.log('response', response.body);
    console.log('response', response.redirected);
    console.log('response', response.url);
    console.log('response', response.type);
    console.log('response', response.headers.get('Content-Type'));
    console.log('response', response.headers.get('Content-Length'));
    console.log('response', response.headers.get('Content-Encoding'));
    console.log('response', response.headers.get('Content-Language'));

    if (!response.ok) { 
      // const error = await response.text();
      throw new Error(`HTTP error! status: ${response.status} ${response.statusText}`);
    }


    return response.json();
  } catch (error) {
    console.error('Error fetching API:', error);
    throw error;
  }
}

const API_FUNCTIONS = Object.freeze({
  // Pipeline endpoints

  getPipelineByID: (id: string) =>
    fetchAPI<{ success: boolean; data: DynamicPipelineResult }>(`/api/pipeline?id=${id}`),

  getPipelines: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: any[]; count: number }>(
      `/api/pipeline?offset=${offset}&limit=${limit}`
    ),

  savePipeline: (data: any) =>
    fetchAPI<{ success: boolean; message: string; type: string }>('/api/pipeline/save', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  createPipeline: (data: any) =>
    fetchAPI<{ success: boolean; message: string; type: string }>('/api/pipeline', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  getPipelineSteps: (pipelineID: string) =>
    fetchAPI<{ success: boolean; steps: any[]; count: number }>(`/api/pipeline/steps?pipeline_id=${pipelineID}`),


  // Dynamic pipeline
  executeDynamicPipeline: (query: string, depth = 3, skipDuplicates = true, stream = false) => {
    const params = new URLSearchParams({
      q: query,
      depth: depth.toString(),
      skip_duplicates: skipDuplicates.toString(),
      stream: stream.toString(),
    });
    return fetchAPI<PipelineResponse>(`/dynamic?${params}`);
  },

  // Search endpoint
  search: (query: string, domain?: string) => {
    const params = new URLSearchParams({ q: query });
    if (domain) params.set('domain', domain);
    return fetchAPI<DomainSearchResult>(`/search?${params}`);
  },

  // Domain endpoints
  getOnapi: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: Entity[]; count: number }>(
      `/api/onapi?offset=${offset}&limit=${limit}`
    ),

  createOnapi: (entity: Entity) =>
    fetchAPI<{ success: boolean; message: string }>('/api/onapi', {
      method: 'POST',
      body: JSON.stringify(entity),
    }),

  getScj: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: ScjCase[]; count: number }>(
      `/api/scj?offset=${offset}&limit=${limit}`
    ),

  getDgii: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: Register[]; count: number }>(
      `/api/dgii?offset=${offset}&limit=${limit}`
    ),

  getPgr: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: PgrNews[]; count: number }>(
      `/api/pgr?offset=${offset}&limit=${limit}`
    ),

  getDocking: (offset = 0, limit = 10) =>
    fetchAPI<{ success: boolean; data: GoogleDockingResult[]; count: number }>(
      `/api/docking?offset=${offset}&limit=${limit}`
    ),
});

export const api = API_FUNCTIONS;

// SSE (Server-Sent Events) for streaming pipeline
export function createPipelineStream(
  query: string,
  depth: number,
  skipDuplicates: boolean,
  onStep: (step: any) => void,
  onComplete: () => void,
  onError: (error: Error) => void
) {
  const params = new URLSearchParams({
    q: query,
    depth: depth.toString(),
    skip_duplicates: skipDuplicates.toString(),
    stream: 'true',
  });

  const eventSource = new EventSource(`${API_BASE_URL}/dynamic?${params}`);

  eventSource.addEventListener('step', (event) => {
    try {
      const data = JSON.parse(event.data);
      onStep(data);
    } catch (error) {
      console.error('Error parsing step data:', error);
    }
  });

  eventSource.addEventListener('sumary', (event) => {
    try {
      const data = JSON.parse(event.data);
      onStep(data);
    } catch (error) {
      console.error('Error parsing summary data:', error);
    }
  });

  eventSource.addEventListener('complete', () => {
    onComplete();
    eventSource.close();
  });

  eventSource.addEventListener('error', () => {
    onError(new Error('Stream error'));
    eventSource.close();
  });

  return () => eventSource.close();
}

