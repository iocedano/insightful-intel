import { useState, useMemo } from 'react';
import type { DynamicPipelineStep, DynamicPipelineResult, GoogleDockingResult, Register, ScjCase, Entity, PgrNews } from '../types';
import { DOMAIN_TYPE_MAP } from '../types';

interface PipelineDetailsProps {
  pipeline: DynamicPipelineResult;
  steps: DynamicPipelineStep[];
  onBack?: () => void;
}

type TabType = 'overview' | 'timeline' | 'results' | 'keywords';

// Group steps by depth
const groupStepsByDepth = (steps: DynamicPipelineStep[]) => {
  const grouped: Record<number, DynamicPipelineStep[]> = {};
  steps.forEach((step) => {
    const depth = step.depth || 0;
    if (!grouped[depth]) {
      grouped[depth] = [];
    }
    grouped[depth].push(step);
  });
  return grouped;
};

// Aggregate keywords_per_category from all steps
const aggregateKeywordsByCategory = (steps: DynamicPipelineStep[]) => {
  const aggregated: Record<string, Set<string>> = {};
  steps.forEach((step) => {
    if (step.keywords_per_category) {
      Object.entries(step.keywords_per_category).forEach(([category, keywords]) => {
        if (!aggregated[category]) {
          aggregated[category] = new Set();
        }
        (keywords as string[]).forEach((keyword) => {
          aggregated[category].add(keyword);
        });
      });
    }
  });
  // Convert Sets to Arrays
  const result: Record<string, string[]> = {};
  Object.entries(aggregated).forEach(([category, keywordsSet]) => {
    result[category] = Array.from(keywordsSet);
  });
  return result;
};

// Helper to extract URLs from step output
const extractUrls = (output: any): GoogleDockingResult[] => {
  if (!output) return [];
  if (Array.isArray(output)) {
    return output.filter((item) => item && (item.link || item.url));
  }
  return [];
};

// Helper to get domain type display name
const getDomainTypeDisplay = (domainType: string): string => {
  return domainType
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
};

export default function PipelineDetails({
  pipeline,
  steps,
  onBack,
}: PipelineDetailsProps) {
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedDepths, setExpandedDepths] = useState<Set<number>>(new Set([0]));

  const stepsByDepth = useMemo(() => groupStepsByDepth(steps), [steps]);
  const keywordsByCategory = useMemo(() => aggregateKeywordsByCategory(steps), [steps]);
  const depthLevels = useMemo(
    () => Object.keys(stepsByDepth).map(Number).sort((a, b) => a - b),
    [stepsByDepth]
  );

  // Filter steps by domain type
  const socialMediaSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.SOCIAL_MEDIA.toUpperCase() && step.output !== null
    ),
    [steps]
  );

  const googleDockingSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.GOOGLE_DOCKING.toUpperCase() && step.output !== null
    ),
    [steps]
  );

  const fileTypeSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.FILE_TYPE.toUpperCase() && step.output !== null
    ),
    [steps]
  );

  const xSocialMediaSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA.toUpperCase() && step.output !== null
    ),
    [steps]
  );

  const dockingSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.GOOGLE_DOCKING.toUpperCase() && step.output !== null
    ),
    [steps]
  );

  const dgiiSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.DGII.toUpperCase() && step.output !== null && step.output.length > 0
    ),
    [steps]
  );

  const pgrSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.PGR.toUpperCase() && step.output !== null && step.output.length > 0
    ),
    [steps]
  );

  const scjSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.SCJ.toUpperCase() && step.output !== null && step.output.length > 0
    ),
    [steps]
  );

  const onapiSteps = useMemo(
    () => steps.filter(
      (step) => step.domain_type === DOMAIN_TYPE_MAP.ONAPI.toUpperCase() && step.output !== null && step.output.length > 0
    ),
    [steps]
  );  



  // Filter steps based on search query
  const filteredSteps = useMemo(() => {
    if (!searchQuery.trim()) return steps;
    const query = searchQuery.toLowerCase();
    return steps.filter((step) => {
      return (
        step.search_parameter?.toLowerCase().includes(query) ||
        step.domain_type?.toLowerCase().includes(query) ||
        step.keywords?.some((kw) => kw.toLowerCase().includes(query)) ||
        step.error?.toLowerCase().includes(query)
      );
    });
  }, [steps, searchQuery]);

  const toggleDepth = (depth: number) => {
    const newExpanded = new Set(expandedDepths);
    if (newExpanded.has(depth)) {
      newExpanded.delete(depth);
    } else {
      newExpanded.add(depth);
    }
    setExpandedDepths(newExpanded);
  };

  // Render URL results as cards
  const renderUrlResults = (results: GoogleDockingResult[]) => {
    if (results.length === 0) return null;

    return (
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 mb-2">
          Results ({results.length})
        </p>
        <div className="grid grid-cols-1 gap-2 max-h-96 overflow-y-auto">
          {results.map((result, index) => (
            <a
              key={index}
              href={result.link || result.url}
              target="_blank"
              rel="noopener noreferrer"
              className="block p-3 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                  <h4 className="text-sm font-semibold text-gray-900 truncate mb-1">
                    {result.title || 'Untitled'}
                  </h4>
                  {result.description && (
                    <p className="text-xs text-gray-600 line-clamp-2 mb-2">
                      {result.description}
                    </p>
                  )}
                  <p className="text-xs text-blue-600 truncate">
                    {result.link || result.url}
                  </p>
                  {result.relevance !== undefined && (
                    <span className="inline-block mt-1 px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded">
                      Relevance: {result.relevance.toFixed(2)}
                    </span>
                  )}
                </div>
                <svg
                  className="w-4 h-4 text-gray-400 flex-shrink-0 mt-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                  />
                </svg>
              </div>
            </a>
          ))}
        </div>
      </div>
    );
  };

  const renderDgiiResults = (results: Register[]) => {
    if (results.length === 0) return null;
    return (
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 mb-3">
          Results ({results.length})
        </p>
        <div className="grid grid-cols-1 gap-3 max-h-96 overflow-y-auto">
          {results.map((result, index) => (
            <div
              key={result.id || index}
              className="p-4 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all"
            >
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <h4 className="text-sm font-semibold text-gray-900 mb-2">
                      {result.razon_social || 'N/A'}
                    </h4>
                    {result.nombre_comercial && (
                      <p className="text-xs text-gray-600 mb-2">
                        Commercial Name: <span className="font-medium text-gray-900">{result.nombre_comercial}</span>
                      </p>
                    )}
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-2 text-xs">
                  <div>
                    <span className="text-gray-500">RNC:</span>{' '}
                    <span className="font-medium text-gray-900">{result.rnc || 'N/A'}</span>
                  </div>
                  {result.categoria && (
                    <div>
                      <span className="text-gray-500">Category:</span>{' '}
                      <span className="font-medium text-gray-900">{result.categoria}</span>
                    </div>
                  )}
                  {result.estado && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Status:</span>{' '}
                      <span
                        className={`font-medium px-2 py-0.5 rounded ${
                          result.estado.toLowerCase().includes('activo')
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {result.estado}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  const renderPgrResults = (results: PgrNews[]) => {
    if (results.length === 0) return null;
    return (
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 mb-3">
          Results ({results.length})
        </p>
        <div className="grid grid-cols-1 gap-3 max-h-96 overflow-y-auto">
          {results.map((result, index) => (
            <a
              key={result.id || index}
              href={result.url}
              target="_blank"
              rel="noopener noreferrer"
              className="block p-4 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all"
            >
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                  <h4 className="text-sm font-semibold text-gray-900 mb-2 line-clamp-2">
                    {result.title || 'Untitled Article'}
                  </h4>
                  <p className="text-xs text-blue-600 truncate">
                    {result.url}
                  </p>
                </div>
                <svg
                  className="w-4 h-4 text-gray-400 flex-shrink-0 mt-1"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                  />
                </svg>
              </div>
            </a>
          ))}
        </div>
      </div>
    );
  };

  const renderScjResults = (results: ScjCase[]) => {
    if (results.length === 0) return null;
    return (
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 mb-3">
        Results ({results.length})
        </p>
        <div className="grid grid-cols-1 gap-3 max-h-96 overflow-y-auto">
          {results.map((result, index) => (
            <div
              key={result.id || index}
              className="p-4 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all"
            >
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <h4 className="text-sm font-semibold text-gray-900">
                        Case #{result.no_expediente || result.id_expediente || 'N/A'}
                      </h4>
                      {result.no_sentencia && (
                        <span className="px-2 py-0.5 bg-blue-100 text-blue-800 text-xs rounded">
                          Sentence: {result.no_sentencia}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs">
                  {result.id_expediente && (
                    <div>
                      <span className="text-gray-500">Expediente ID:</span>{' '}
                      <span className="font-medium text-gray-900">{result.id_expediente}</span>
                    </div>
                  )}
                  {result.fecha_fallo && (
                    <div>
                      <span className="text-gray-500">Date:</span>{' '}
                      <span className="font-medium text-gray-900">{result.fecha_fallo}</span>
                    </div>
                  )}
                  {result.desc_tribunal && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Court:</span>{' '}
                      <span className="font-medium text-gray-900">{result.desc_tribunal}</span>
                    </div>
                  )}
                  {result.desc_materia && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Subject:</span>{' '}
                      <span className="font-medium text-gray-900">{result.desc_materia}</span>
                    </div>
                  )}
                  {result.involucrados && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Involved Parties:</span>{' '}
                      <span className="font-medium text-gray-900">{result.involucrados}</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  const renderOnapiResults = (results: Entity[]) => {
    if (results.length === 0) return null;
    return (
      <div className="space-y-2">
        <p className="text-sm font-medium text-gray-700 mb-3">
        Results ({results.length})
        </p>
        <div className="grid grid-cols-1 gap-3 max-h-96 overflow-y-auto">
          {results.map((result, index) => (
            <div
              key={result.id || index}
              className="p-4 bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-md transition-all"
            >
              <div className="space-y-2">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <h4 className="text-sm font-semibold text-gray-900">
                        {result.serie_expediente}-{result.numero_expediente}
                      </h4>
                      {result.certificado && (
                        <span className="px-2 py-0.5 bg-green-100 text-green-800 text-xs rounded">
                          Cert: {result.certificado}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs">
                  {result.tipo && (
                    <div>
                      <span className="text-gray-500">Type:</span>{' '}
                      <span className="font-medium text-gray-900">{result.tipo}</span>
                    </div>
                  )}
                  {result.subtipo && (
                    <div>
                      <span className="text-gray-500">Subtype:</span>{' '}
                      <span className="font-medium text-gray-900">{result.subtipo}</span>
                    </div>
                  )}
                  {result.titular && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Holder:</span>{' '}
                      <span className="font-medium text-gray-900">{result.titular}</span>
                    </div>
                  )}
                  {result.texto && (
                    <div className="col-span-2">
                      <span className="text-gray-500">Description:</span>
                      <p className="mt-1 text-gray-900 line-clamp-3">{result.texto}</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };
  // Render step card
  const renderStepCard = (step: DynamicPipelineStep, showDetails = true) => {
    const results = extractUrls(step.output);
    const hasResults = results.length > 0;

    return (
      <div
        key={step.id}
        className={`p-4 border rounded-lg transition-all ${
          step.success
            ? 'border-green-300 bg-green-50/50 hover:bg-green-50'
            : 'border-red-300 bg-red-50/50 hover:bg-red-50'
        }`}
      >
        <div className="flex items-start justify-between mb-2">
          <div className="flex items-center gap-2 flex-1">
            <span className="font-semibold text-gray-900 text-sm">
              {getDomainTypeDisplay(step.domain_type)}
            </span>
            <span
              className={`px-2 py-1 text-xs rounded-full ${
                step.success
                  ? 'bg-green-200 text-green-800'
                  : 'bg-red-200 text-red-800'
              }`}
            >
              {step.success ? '✓ Success' : '✗ Failed'}
            </span>
            {step.depth !== undefined && (
              <span className="px-2 py-1 bg-gray-200 text-gray-700 text-xs rounded">
                Depth {step.depth}
              </span>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <p className="text-sm text-gray-700">
            <span className="font-medium">Query:</span>{' '}
            <span className="text-gray-900">{step.search_parameter}</span>
          </p>

          {step.error && (
            <div className="p-2 bg-red-100 border border-red-300 rounded text-sm text-red-800">
              <span className="font-medium">Error:</span> {step.error}
            </div>
          )}

          {showDetails && step.keywords && step.keywords.length > 0 && (
            <div>
              <p className="text-xs font-medium text-gray-700 mb-1">Keywords:</p>
              <div className="flex flex-wrap gap-1">
                {step.keywords.map((keyword, idx) => (
                  <span
                    key={idx}
                    className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded"
                  >
                    {keyword}
                  </span>
                ))}
              </div>
            </div>
          )}

          {showDetails &&
            step.keywords_per_category &&
            Object.keys(step.keywords_per_category).length > 0 && (
              <div>
                <p className="text-xs font-medium text-gray-700 mb-1">Keywords by Category:</p>
                <div className="space-y-1">
                  {Object.entries(step.keywords_per_category).map(([category, keywords]) => (
                    <div key={category}>
                      <span className="text-xs font-medium text-gray-600 capitalize">
                        {category.replace(/_/g, ' ')}:
                      </span>{' '}
                      <div className="inline-flex flex-wrap gap-1">
                        {keywords.map((kw, idx) => (
                          <span
                            key={idx}
                            className="px-1.5 py-0.5 bg-gray-100 text-gray-700 text-xs rounded"
                          >
                            {kw}
                          </span>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

          {hasResults && showDetails && renderUrlResults(results)}
        </div>
      </div>
    );
  };

  const tabs: { id: TabType; label: string; count?: number }[] = [
    { id: 'overview', label: 'Overview' },
    { id: 'timeline', label: 'Timeline', count: steps.length },
    { id: 'results', label: 'Results' },
    { id: 'keywords', label: 'Keywords', count: Object.keys(keywordsByCategory).length },
  ];

  return (
    <div className="max-w-7xl mx-auto p-4 md:p-6 space-y-6">
      <div className="bg-white rounded-lg shadow-lg">
        {/* Header */}
        <div className="border-b border-gray-200 p-4 md:p-6">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl md:text-3xl font-bold text-gray-900">Pipeline Details</h1>
              <p className="text-sm text-gray-500 mt-1">Execution ID: {pipeline.id}</p>
            </div>
            <button
              onClick={onBack}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
            >
              ← Back to Dashboard
            </button>
          </div>

          {/* Summary Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4 mt-6">
            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg border border-blue-200">
              <p className="text-xs md:text-sm text-blue-700 font-medium mb-1">Total Steps</p>
              <p className="text-xl md:text-2xl font-bold text-blue-900">{pipeline.total_steps}</p>
            </div>
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-lg border border-green-200">
              <p className="text-xs md:text-sm text-green-700 font-medium mb-1">Successful</p>
              <p className="text-xl md:text-2xl font-bold text-green-900">{pipeline.successful_steps}</p>
              <p className="text-xs text-green-600 mt-1">
                {pipeline.total_steps > 0
                  ? Math.round((pipeline.successful_steps / pipeline.total_steps) * 100)
                  : 0}%
              </p>
            </div>
            <div className="bg-gradient-to-br from-red-50 to-red-100 p-4 rounded-lg border border-red-200">
              <p className="text-xs md:text-sm text-red-700 font-medium mb-1">Failed</p>
              <p className="text-xl md:text-2xl font-bold text-red-900">{pipeline.failed_steps}</p>
              <p className="text-xs text-red-600 mt-1">
                {pipeline.total_steps > 0
                  ? Math.round((pipeline.failed_steps / pipeline.total_steps) * 100)
                  : 0}%
              </p>
            </div>
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-lg border border-purple-200">
              <p className="text-xs md:text-sm text-purple-700 font-medium mb-1">Max Depth</p>
              <p className="text-xl md:text-2xl font-bold text-purple-900">{pipeline.max_depth_reached}</p>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="border-b border-gray-200">
          <div className="flex overflow-x-auto scrollbar-hide">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`px-4 md:px-6 py-3 font-medium text-sm whitespace-nowrap border-b-2 transition-colors ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600 bg-blue-50/50'
                    : 'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300'
                }`}
              >
                {tab.label}
                {tab.count !== undefined && (
                  <span className="ml-2 px-2 py-0.5 bg-gray-200 text-gray-700 text-xs rounded-full">
                    {tab.count}
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>

        {/* Tab Content */}
        <div className="p-4 md:p-6">

          {/* Overview Tab */}
          {activeTab === 'overview' && (
            <div className="space-y-6">
              {/* Keywords by Category */}
              {Object.keys(keywordsByCategory).length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Keywords by Category</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {Object.entries(keywordsByCategory).map(([category, keywords]) => (
                      <div
                        key={category}
                        className="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-4 border border-gray-200 hover:shadow-md transition-shadow"
                      >
                        <h3 className="font-semibold text-gray-900 mb-3 capitalize text-sm">
                          {category.replace(/_/g, ' ')}
                        </h3>
                        <div className="flex flex-wrap gap-2 mb-2">
                          {keywords.slice(0, 10).map((keyword, index) => (
                            <span
                              key={index}
                              className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-md font-medium"
                            >
                              {keyword}
                            </span>
                          ))}
                          {keywords.length > 10 && (
                            <span className="px-2 py-1 bg-gray-200 text-gray-600 text-xs rounded-md">
                              +{keywords.length - 10} more
                            </span>
                          )}
                        </div>
                        <p className="text-xs text-gray-500 mt-2">
                          {keywords.length} keyword{keywords.length !== 1 ? 's' : ''} total
                        </p>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Results Summary */}
              <div>
                <h2 className="text-xl font-bold mb-4 text-gray-900">Results Summary</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {socialMediaSteps.length > 0 && (
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                      <h3 className="font-semibold text-blue-900 mb-2">Social Media</h3>
                      <p className="text-2xl font-bold text-blue-700">{socialMediaSteps.length}</p>
                      <p className="text-sm text-blue-600">steps with results</p>
                    </div>
                  )}
                  {xSocialMediaSteps.length > 0 && (
                    <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
                      <h3 className="font-semibold text-purple-900 mb-2">X Social Media</h3>
                      <p className="text-2xl font-bold text-purple-700">{xSocialMediaSteps.length}</p>
                      <p className="text-sm text-purple-600">steps with results</p>
                    </div>
                  )}
                  {fileTypeSteps.length > 0 && (
                    <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                      <h3 className="font-semibold text-green-900 mb-2">File Types</h3>
                      <p className="text-2xl font-bold text-green-700">{fileTypeSteps.length}</p>
                      <p className="text-sm text-green-600">steps with results</p>
                    </div>
                  )}
                  {googleDockingSteps.length > 0 && (
                    <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
                      <h3 className="font-semibold text-orange-900 mb-2">Google Docking</h3>
                      <p className="text-2xl font-bold text-orange-700">{googleDockingSteps.length}</p>
                      <p className="text-sm text-orange-600">steps with results</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Timeline Tab */}
          {activeTab === 'timeline' && (
            <div className="space-y-4">
              <div className="flex items-center gap-4 mb-4">
                <input
                  type="text"
                  placeholder="Search steps by query, domain type, or keyword..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="px-4 py-2 text-gray-600 hover:text-gray-900"
                  >
                    Clear
                  </button>
                )}
              </div>

              <div className="space-y-4">
                {depthLevels.map((depth) => {
                  const depthSteps = filteredSteps.filter((s) => s.depth === depth);
                  if (depthSteps.length === 0) return null;

                  const isExpanded = expandedDepths.has(depth);

                  return (
                    <div key={depth} className="border border-gray-200 rounded-lg overflow-hidden">
                      <button
                        onClick={() => toggleDepth(depth)}
                        className="w-full flex items-center justify-between p-4 bg-gray-50 hover:bg-gray-100 transition-colors"
                      >
                        <div className="flex items-center gap-3">
                          <svg
                            className={`w-5 h-5 text-gray-600 transition-transform ${
                              isExpanded ? 'transform rotate-90' : ''
                            }`}
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M9 5l7 7-7 7"
                            />
                          </svg>
                          <h3 className="text-lg font-semibold text-gray-900">Depth {depth}</h3>
                          <span className="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full">
                            {depthSteps.length} step{depthSteps.length !== 1 ? 's' : ''}
                          </span>
                        </div>
                        <div className="flex items-center gap-2 text-sm text-gray-600">
                          <span className="text-green-600">
                            {depthSteps.filter((s) => s.success).length} successful
                          </span>
                          <span className="text-red-600">
                            {depthSteps.filter((s) => !s.success).length} failed
                          </span>
                        </div>
                      </button>

                      {isExpanded && (
                        <div className="p-4 space-y-3 bg-white">
                          {depthSteps.map((step) => renderStepCard(step))}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>

              {filteredSteps.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  <p>No steps match your search query.</p>
                </div>
              )}
            </div>
          )}

          {/* Results Tab */}
          {activeTab === 'results' && (
            <div className="space-y-6">
              {dockingSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Google Docking Results</h2>
                  <div className="space-y-4">
                    {dockingSteps.map((step) => {
                      const results = extractUrls(step.output);
                      if (results.length === 0) return null;
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderUrlResults(results)}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {dgiiSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Dirección General de Impuestos Internos (DGII) Results</h2>
                  <div className="space-y-4">
                    {dgiiSteps.map((step) => {
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderDgiiResults(step.output as Register[])}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {pgrSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Ministerio de Procuraduría General de la República (PGR) Results</h2>
                  <div className="space-y-4">
                    {pgrSteps.map((step) => {
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderPgrResults(step.output as PgrNews[])}
                        </div>
                      );
                    })}   
                  </div>
                </div>
              )}

              {scjSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Sala de Casación de la Suprema Corte de Justicia (SCJ) Results</h2>
                  <div className="space-y-4">
                    {scjSteps.map((step) => {
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderScjResults(step.output as ScjCase[])}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {onapiSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Oficina Nacional de la Propiedad Industrial (ONAPI)  Results</h2>
                  <div className="space-y-4">
                    {onapiSteps.map((step) => {
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderOnapiResults(step.output as Entity[])}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {socialMediaSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Social Media Results</h2>
                  <div className="space-y-4">
                    {socialMediaSteps.map((step) => {
                      const results = extractUrls(step.output);
                      if (results.length === 0) return null;
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderUrlResults(results)}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {xSocialMediaSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">X Social Media Results</h2>
                  <div className="space-y-4">
                    {xSocialMediaSteps.map((step) => {
                      const results = extractUrls(step.output);
                      if (results.length === 0) return null;
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderUrlResults(results)}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {fileTypeSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">File Type Results</h2>
                  <div className="space-y-4">
                    {fileTypeSteps.map((step) => {
                      const results = extractUrls(step.output);
                      if (results.length === 0) return null;
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderUrlResults(results)}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {googleDockingSteps.length > 0 && (
                <div>
                  <h2 className="text-xl font-bold mb-4 text-gray-900">Google Docking Results</h2>
                  <div className="space-y-4">
                    {googleDockingSteps.map((step) => {
                      const results = extractUrls(step.output);
                      if (results.length === 0) return null;
                      return (
                        <div key={step.id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                          <p className="text-sm font-medium text-gray-700 mb-3">
                            Query: <span className="text-gray-900">{step.search_parameter}</span>
                          </p>
                          {renderUrlResults(results)}
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {socialMediaSteps.length === 0 &&
                xSocialMediaSteps.length === 0 &&
                fileTypeSteps.length === 0 &&
                googleDockingSteps.length === 0 && (
                  <div className="text-center py-12 text-gray-500">
                    <p>No results available.</p>
                  </div>
                )}
            </div>
          )}

          {/* Keywords Tab */}
          {activeTab === 'keywords' && (
            <div className="space-y-6">
              {Object.keys(keywordsByCategory).length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {Object.entries(keywordsByCategory).map(([category, keywords]) => (
                    <div
                      key={category}
                      className="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-5 border border-gray-200 hover:shadow-lg transition-shadow"
                    >
                      <h3 className="font-bold text-gray-900 mb-3 capitalize">
                        {category.replace(/_/g, ' ')}
                      </h3>
                      <div className="flex flex-wrap gap-2 mb-3">
                        {keywords.map((keyword, index) => (
                          <span
                            key={index}
                            className="px-2.5 py-1 bg-blue-100 text-blue-800 text-xs rounded-md font-medium hover:bg-blue-200 transition-colors cursor-default"
                          >
                            {keyword}
                          </span>
                        ))}
                      </div>
                      <div className="pt-3 border-t border-gray-300">
                        <p className="text-xs text-gray-600">
                          <span className="font-semibold">{keywords.length}</span> unique keyword
                          {keywords.length !== 1 ? 's' : ''}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-12 text-gray-500">
                  <p>No keywords found in this pipeline.</p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

