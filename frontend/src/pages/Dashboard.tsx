import { useState, useEffect } from 'react';
import { api } from '../api';
import type { DynamicPipelineResult, DynamicPipelineStep } from '../types';
import PipelineDetails from '../components/PipelineDetails';



export default function Dashboard() {
  const [pipelines, setPipelines] = useState<DynamicPipelineResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedPipeline, setSelectedPipeline] = useState<string | null>(null);
  const [selectedPipelineSteps, setSelectedPipelineSteps] = useState<DynamicPipelineStep[]>([]);

  useEffect(() => {
    setLoading(true);
    setError(null);
    setSelectedPipelineSteps([]);
      api.getPipelines(0, 50).then((response) => {
        if (response.success && response.data) {
          setPipelines(response.data as DynamicPipelineResult[]);
        }
      }).catch((err: any) => {
        setError(err.message || 'Failed to load pipelines');
      }).finally(() => {
        setLoading(false);
      });
  }, []);

  console.log('selectedPipeline changed', selectedPipeline);

  useEffect(() => {
    console.log('selectedPipeline changed', selectedPipeline);
    if (selectedPipeline) {
      console.log('selectedPipeline', selectedPipeline);
      api.getPipelineSteps(selectedPipeline as string).then((response) => {
        if (response.success && response.steps) {
          setSelectedPipelineSteps(response.steps as DynamicPipelineStep[]);
        }
      }).catch((err: any) => {
        console.error('Error loading pipeline steps:', err);
        setSelectedPipelineSteps([]);
      });
    }
  }, [selectedPipeline]); 


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

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="text-center py-8">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading pipelines...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        </div>
      </div>
    );
  }

  if (pipelines.length === 0) {
    return (
      <div className="max-w-7xl mx-auto p-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h1 className="text-3xl font-bold mb-6">Dashboard</h1>
          <div className="text-center py-12">
            <p className="text-gray-600 text-lg">No pipelines found</p>
            <p className="text-gray-500 mt-2">Create a pipeline to see it here</p>
          </div>
        </div>
      </div>
    );
  }

  // If a pipeline is selected, show its details
  if (selectedPipeline && selectedPipelineSteps.length > 0) {
    const pipeline = pipelines.find((p) => p.id === selectedPipeline);
    if (!pipeline) {
      setSelectedPipeline(null);
      return null;
    }

    return (
      <PipelineDetails
        pipeline={pipeline}
        steps={selectedPipelineSteps}
        onBack={() => setSelectedPipeline(null)}
        showBackButton={true}
      />
    );
  }

  // Main dashboard view - list of pipelines
  return (
    <div className="max-w-7xl mx-auto p-6 space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h1 className="text-3xl font-bold mb-6">Dashboard</h1>
        <p className="text-gray-600 mb-6">
          View and analyze your saved pipelines. Click on a pipeline to see detailed steps grouped by depth.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {pipelines.map((pipeline) => {
            const keywordsByCategory = aggregateKeywordsByCategory(pipeline.steps || []);
            const categoryCount = Object.keys(keywordsByCategory).length;
            const totalKeywords = Object.values(keywordsByCategory).reduce(
              (sum, keywords) => sum + keywords.length,
              0
            );

            // Format timestamps
            const formatDate = (date: Date | string) => {
              const d = typeof date === 'string' ? new Date(date) : date;
              return d.toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                year: 'numeric',
              });
            };

            const formatTime = (date: Date | string) => {
              const d = typeof date === 'string' ? new Date(date) : date;
              return d.toLocaleTimeString('en-US', {
                hour: '2-digit',
                minute: '2-digit',
              });
            };


            // Calculate execution duration
            const getExecutionDuration = () => {
              if (!pipeline.created_at || !pipeline.updated_at) return null;
              const created = typeof pipeline.created_at === 'string' 
                ? new Date(pipeline.created_at) 
                : pipeline.created_at;
              const updated = typeof pipeline.updated_at === 'string' 
                ? new Date(pipeline.updated_at) 
                : pipeline.updated_at;
              const durationMs = updated.getTime() - created.getTime();
              const durationSec = Math.floor(durationMs / 1000);
              const durationMin = Math.floor(durationSec / 60);
              
              if (durationMin > 0) {
                return `${durationMin}m ${durationSec % 60}s`;
              }
              return `${durationSec}s`;
            };

            const executionDuration = getExecutionDuration();
            const successRate = pipeline.total_steps > 0 
              ? Math.round((pipeline.successful_steps / pipeline.total_steps) * 100) 
              : 0;

            return (
              <div
                key={pipeline.id}
                className="border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-all duration-200 cursor-pointer bg-white"
                onClick={() => setSelectedPipeline(pipeline.id)}
              >
                {/* Header Section */}
                <div className="mb-4">
                  <div className="flex items-start justify-between mb-2">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-gray-900 mb-1">
                        {pipeline.config.query || 'Untitled Pipeline'}
                      </h3>
                      <p className="text-xs text-gray-500">
                        Pipeline ID: {pipeline.id.substring(0, 8)}...
                      </p>
                    </div>
                    <div className="flex flex-col items-end gap-1">
                      <span className={`px-2 py-1 text-xs font-medium rounded ${
                        successRate >= 80 
                          ? 'bg-green-100 text-green-800' 
                          : successRate >= 50 
                          ? 'bg-yellow-100 text-yellow-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {successRate}% Success
                      </span>
                    </div>
                  </div>

                  {/* Configuration Badges */}
                  <div className="flex flex-wrap gap-2 mb-3">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
                      <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clipRule="evenodd" />
                      </svg>
                      Depth: {pipeline.config.max_depth}
                    </span>
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-50 text-purple-700">
                      <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                        <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                        <path fillRule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clipRule="evenodd" />
                      </svg>
                      {pipeline.config?.available_domains?.length || 0} Domains
                    </span>
                    {pipeline.config.skip_duplicates && (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-50 text-gray-700">
                        Skip Duplicates
                      </span>
                    )}
                  </div>
                </div>

                {/* Statistics Grid */}
                <div className="grid grid-cols-3 gap-3 mb-4 p-3 bg-gray-50 rounded-lg">
                  <div className="text-center">
                    <p className="text-2xl font-bold text-gray-900">{pipeline.total_steps}</p>
                    <p className="text-xs text-gray-600 font-medium">Total Steps</p>
                  </div>
                  <div className="text-center border-x border-gray-200">
                    <p className="text-2xl font-bold text-green-600">{pipeline.successful_steps}</p>
                    <p className="text-xs text-gray-600 font-medium">Successful</p>
                  </div>
                  <div className="text-center">
                    <p className="text-2xl font-bold text-red-600">{pipeline.failed_steps}</p>
                    <p className="text-xs text-gray-600 font-medium">Failed</p>
                  </div>
                </div>

                {/* Progress Indicator */}
                <div className="mb-4">
                  <div className="flex items-center justify-between text-xs text-gray-600 mb-1">
                    <span>Max Depth Reached</span>
                    <span>{pipeline.max_depth_reached} / {pipeline.config.max_depth}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full transition-all"
                      style={{
                        width: `${(pipeline.max_depth_reached / pipeline.config.max_depth) * 100}%`,
                      }}
                    ></div>
                  </div>
                </div>

                {/* Keywords Section */}
                {categoryCount > 0 && (
                  <div className="mt-4 pt-4 border-t border-gray-200">
                    <div className="flex items-center justify-between mb-2">
                      <p className="text-sm font-semibold text-gray-700">
                        Extracted Keywords
                      </p>
                      <span className="text-xs text-gray-500">
                        {totalKeywords} keywords
                      </span>
                    </div>
                    <div className="flex flex-wrap gap-1.5 mb-2">
                      {Object.keys(keywordsByCategory).slice(0, 4).map((category) => (
                        <span
                          key={category}
                          className="inline-flex items-center px-2 py-1 bg-indigo-50 text-indigo-700 text-xs font-medium rounded-md"
                        >
                          {category.replace(/_/g, ' ')}
                          <span className="ml-1 text-indigo-500">
                            ({keywordsByCategory[category].length})
                          </span>
                        </span>
                      ))}
                      {categoryCount > 4 && (
                        <span className="inline-flex items-center px-2 py-1 bg-gray-100 text-gray-600 text-xs font-medium rounded-md">
                          +{categoryCount - 4} more
                        </span>
                      )}
                    </div>
                  </div>
                )}

                {/* Timestamps Section */}
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <div className="grid grid-cols-2 gap-3 text-xs">
                    <div>
                      <p className="text-gray-500 font-medium mb-1">Created</p>
                      <p className="text-gray-900 font-semibold">
                        {pipeline.created_at ? formatDate(pipeline.created_at) : 'N/A'}
                      </p>
                      <p className="text-gray-500">
                        {pipeline.created_at ? formatTime(pipeline.created_at) : ''}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-500 font-medium mb-1">Last Updated</p>
                      <p className="text-gray-900 font-semibold">
                        {pipeline.updated_at ? formatDate(pipeline.updated_at) : 'N/A'}
                      </p>
                      <p className="text-gray-500">
                        {pipeline.updated_at ? formatTime(pipeline.updated_at) : ''}
                      </p>
                    </div>
                  </div>
                  {executionDuration && (
                    <div className="mt-3 pt-3 border-t border-gray-100">
                      <div className="flex items-center text-xs text-gray-600">
                        <svg className="w-4 h-4 mr-1.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span className="font-medium">Execution Time:</span>
                        <span className="ml-1">{executionDuration}</span>
                      </div>
                    </div>
                  )}
                </div>

                {/* Action Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedPipeline(pipeline.id);
                  }}
                  className="mt-4 w-full px-4 py-2.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium text-sm transition-colors duration-200 flex items-center justify-center"
                >
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  View Full Details
                </button>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

