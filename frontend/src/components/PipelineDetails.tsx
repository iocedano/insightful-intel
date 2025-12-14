import { useState, useMemo, memo } from 'react';
import type { DynamicPipelineStep, DynamicPipelineResult, GoogleDockingResult, Register, ScjCase, Entity, PgrNews } from '../types';
import { DOMAIN_TYPE_MAP } from '../types';
import CardUrlResults from './CardUrlResults';
import CardDgiiResults from './CardDgiiResults';
import CardPgrResults from './CardPgrResults';
import CardScjResults from './CardScjResults';
import CardOnapiResults from './CardOnapiResults';
import CardStep from './CardStep';

interface PipelineDetailsProps {
  pipeline: DynamicPipelineResult;
  steps: DynamicPipelineStep[];
  onBack?: () => void;
  showBackButton?: boolean;
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

const filterStepsByDomainType = (steps: DynamicPipelineStep[], domainType: string) => {
  return steps.filter((step) => step.domain_type === domainType &&  step.output !== null && step.output.length > 0);
};

// Helper to extract URLs from step output
const extractUrls = (output: any): GoogleDockingResult[] => {
  if (!output) return [];
  if (Array.isArray(output)) {
    return output.filter((item) => item && (item.link || item.url));
  }
  return [];
};



function PipelineDetails(props: PipelineDetailsProps) {
  const [activeTab, setActiveTab] = useState<TabType>('overview');
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedDepths, setExpandedDepths] = useState<Set<number>>(new Set([0]));
  const steps = props.steps;
  const pipeline = props.pipeline;
  const showBackButton = props.showBackButton;
  const onBack = props.onBack;

  const stepsByDepth = useMemo(() => groupStepsByDepth(steps), [steps]);
  const keywordsByCategory = useMemo(() => aggregateKeywordsByCategory(steps), [steps]);
  const depthLevels = useMemo(
    () => Object.keys(stepsByDepth).map(Number).sort((a, b) => a - b),
    [stepsByDepth]
  );

  // Compute a mapping of domain type to filtered steps (by output) for quick lookup
  const domainPerStep = useMemo(() => {
    console.log('domainPerStep', steps.length);
    return {
      [DOMAIN_TYPE_MAP.SOCIAL_MEDIA]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.SOCIAL_MEDIA.toUpperCase()),
      [DOMAIN_TYPE_MAP.GOOGLE_DOCKING]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.GOOGLE_DOCKING.toUpperCase()),
      [DOMAIN_TYPE_MAP.FILE_TYPE]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.FILE_TYPE.toUpperCase()),
      [DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA.toUpperCase()),
      [DOMAIN_TYPE_MAP.DGII]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.DGII.toUpperCase()),
      [DOMAIN_TYPE_MAP.PGR]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.PGR.toUpperCase()),
      [DOMAIN_TYPE_MAP.SCJ]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.SCJ.toUpperCase()),
      [DOMAIN_TYPE_MAP.ONAPI]: filterStepsByDomainType(steps, DOMAIN_TYPE_MAP.ONAPI.toUpperCase()),
    };
  }, [steps.length]);

  // Filter steps based on search query
  const filteredSteps = useMemo(() => {
    if (!searchQuery.trim()) return steps;
    const query = searchQuery.toLowerCase();
    return steps.filter((step) => {
      return (
        step.search_parameter?.toLowerCase().includes(query) ||
        step.domain_type?.toLowerCase().includes(query) ||
        step.keywords?.some((kw) => kw.toLowerCase().includes(query)) ||
        !!step.error
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


  const tabs: { id: TabType; label: string; count?: number }[] = [
    { id: 'overview', label: 'Overview' },
    { id: 'timeline', label: 'Steps', count: steps.length },
    { id: 'results', label: 'Results' },
    { id: 'keywords', label: 'Keywords', count: Object.keys(keywordsByCategory).length },
  ];

  const dockingSteps = domainPerStep[DOMAIN_TYPE_MAP.GOOGLE_DOCKING];
  const dgiiSteps = domainPerStep[DOMAIN_TYPE_MAP.DGII];
  const pgrSteps = domainPerStep[DOMAIN_TYPE_MAP.PGR];
  const scjSteps = domainPerStep[DOMAIN_TYPE_MAP.SCJ];
  const onapiSteps = domainPerStep[DOMAIN_TYPE_MAP.ONAPI];
  const socialMediaSteps = domainPerStep[DOMAIN_TYPE_MAP.SOCIAL_MEDIA];
  const xSocialMediaSteps = domainPerStep[DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA];
  const fileTypeSteps = domainPerStep[DOMAIN_TYPE_MAP.FILE_TYPE];
  const googleDockingSteps = domainPerStep[DOMAIN_TYPE_MAP.GOOGLE_DOCKING];

  return (
    <div className="max-w-7xl mx-auto p-4 md:p-6 space-y-6">
      <div className="bg-white rounded-lg shadow-lg">
        {/* Header */}
        <div className="border-b border-gray-200 p-4 md:p-6">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl md:text-3xl font-bold text-gray-900">Pipeline Details</h1>
              <p className="text-sm text-gray-500 mt-1">Execution ID: {pipeline.id}</p>
              <p className="text-sm text-gray-500 mt-1">Created At: {new Date(pipeline.created_at).toLocaleString()}</p>
              <p className="text-sm text-gray-500 mt-1">Updated At: {new Date(pipeline.updated_at).toLocaleString()}</p>
              <p className="text-sm text-gray-500 mt-1">Duration: {Math.round((new Date(pipeline.updated_at).getTime() - new Date(pipeline.created_at).getTime()) / 60000)} minutes</p>
              <p className="text-sm text-gray-500 mt-1">Query: {pipeline.config?.query}</p>
              <p className="text-sm text-gray-500 mt-1">Domains: {pipeline.config?.available_domains?.join(', ')}</p>
              <p className="text-sm text-gray-500 mt-1">Max Depth: {pipeline.config?.max_depth}</p>
            </div>
            {showBackButton && (
              <button
                onClick={onBack}
                className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors font-medium"
              >
                ← Back to Dashboard
              </button>
            )}
          </div>

          {/* Summary Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4 mt-6">
            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg border border-blue-200">
              <p className="text-xs md:text-sm text-blue-700 font-medium mb-1">Total Steps</p>
              <p className="text-xl md:text-2xl font-bold text-blue-900">{pipeline?.total_steps || pipeline?.steps?.length}</p>
            </div>
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-lg border border-green-200">
              <p className="text-xs md:text-sm text-green-700 font-medium mb-1">Successful</p>
              <p className="text-xl md:text-2xl font-bold text-green-900">
              {pipeline?.successful_steps || pipeline?.steps?.filter(s => s.success).length}
              </p>
              <p className="text-xs text-green-600 mt-1">
                {pipeline?.successful_steps && pipeline?.total_steps > 0
                  ? Math.round((pipeline?.successful_steps / pipeline?.total_steps) * 100)
                  : 0}%
              </p>
            </div>
            <div className="bg-gradient-to-br from-red-50 to-red-100 p-4 rounded-lg border border-red-200">
              <p className="text-xs md:text-sm text-red-700 font-medium mb-1">Failed</p>
              <p className="text-xl md:text-2xl font-bold text-red-900">
              {pipeline?.failed_steps || pipeline?.steps?.filter(s => !s.success).length || 0}
              </p>
              <p className="text-xs text-red-600 mt-1">
                {pipeline?.failed_steps && pipeline?.total_steps > 0
                  ? Math.round((pipeline?.failed_steps / pipeline?.total_steps) * 100)
                  : 0}%
              </p>
            </div>
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-lg border border-purple-200">
              <p className="text-xs md:text-sm text-purple-700 font-medium mb-1">Max Depth</p>
              <p className="text-xl md:text-2xl font-bold text-purple-900">
              {pipeline?.max_depth_reached || 0}
              </p>
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
                  {dockingSteps.length > 0 && (
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                      <h3 className="font-semibold text-blue-900 mb-2">Google Docking</h3>
                      <p className="text-2xl font-bold text-blue-700">{dockingSteps.length}</p>
                      <p className="text-sm text-blue-600">steps with results</p>
                    </div>
                  )}
                  {dgiiSteps.length > 0 && (
                    <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                      <h3 className="font-semibold text-green-900 mb-2">Dirección General de Impuestos Internos (DGII)</h3>
                      <p className="text-2xl font-bold text-green-700">{dgiiSteps.length}</p>
                      <p className="text-sm text-green-600">steps with results</p>
                    </div>
                  )}
                  {pgrSteps.length > 0 && (
                    <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
                      <h3 className="font-semibold text-purple-900 mb-2">Ministerio de Procuraduría General de la República (PGR)</h3>
                      <p className="text-2xl font-bold text-purple-700">{pgrSteps.length}</p>
                      <p className="text-sm text-purple-600">steps with results</p>
                    </div>
                  )}
                  {scjSteps.length > 0 && (
                    <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
                      <h3 className="font-semibold text-orange-900 mb-2">Sala de Casación de la Suprema Corte de Justicia (SCJ)</h3>
                      <p className="text-2xl font-bold text-orange-700">{scjSteps.length}</p>
                      <p className="text-sm text-orange-600">steps with results</p>
                    </div>
                  )}
                  {onapiSteps.length > 0 && (
                    <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                      <h3 className="font-semibold text-red-900 mb-2">Oficina Nacional de la Propiedad Industrial (ONAPI)</h3>
                      <p className="text-2xl font-bold text-red-700">{onapiSteps.length}</p>
                      <p className="text-sm text-red-600">steps with results</p>
                    </div>
                  )}
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
                          {depthSteps.map((step) => (
                            <CardStep key={step.id} step={step} />
                          ))}
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
                          <CardUrlResults results={results} />
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
                          <CardDgiiResults results={step.output as Register[]} />
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
                          <CardPgrResults results={step.output as PgrNews[]} />
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
                          <CardScjResults results={step.output as ScjCase[]} />
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
                          <CardOnapiResults results={step.output as Entity[]} />
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
                          <CardUrlResults results={results} />
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
                          <CardUrlResults results={results} />
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
                          <CardUrlResults results={results} />
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
                          <CardUrlResults results={results} />
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

function checkIfEqual(prevProps: PipelineDetailsProps, nextProps: PipelineDetailsProps) {
  return prevProps.pipeline.id === nextProps.pipeline.id && prevProps.steps.length === nextProps.steps.length;
}

export default memo(PipelineDetails, checkIfEqual);