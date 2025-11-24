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

            return (
              <div
                key={pipeline.id}
                className="border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-shadow cursor-pointer"
                onClick={() => setSelectedPipeline(pipeline.id)}
              >
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-gray-900">Pipeline {pipeline.config.query}</h3>
                  <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                    Depth {pipeline.config.max_depth}
                  </span>
                  <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                    Domains {pipeline.config?.available_domains?.join(', ')}
                  </span>
                  <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
                    Query {pipeline.config?.query}
                  </span>
                </div>

                <div className="grid grid-cols-3 gap-2 mb-4">
                  <div className="text-center">
                    <p className="text-2xl font-bold text-gray-900">{pipeline.total_steps}</p>
                    <p className="text-xs text-gray-500">Steps</p>
                  </div>
                  <div className="text-center">
                    <p className="text-2xl font-bold text-green-600">{pipeline.successful_steps}</p>
                    <p className="text-xs text-gray-500">Success</p>
                  </div>
                  <div className="text-center">
                    <p className="text-2xl font-bold text-red-600">{pipeline.failed_steps}</p>
                    <p className="text-xs text-gray-500">Failed</p>
                  </div>
                </div>

                {categoryCount > 0 && (
                  <div className="mt-4 pt-4 border-t border-gray-200">
                    <p className="text-sm font-medium text-gray-700 mb-2">
                      Keywords: {totalKeywords} across {categoryCount} categor{categoryCount !== 1 ? 'ies' : 'y'}
                    </p>
                    <div className="flex flex-wrap gap-1">
                      {Object.keys(keywordsByCategory).slice(0, 3).map((category) => (
                        <span
                          key={category}
                          className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                        >
                          {category.replace(/_/g, ' ')}
                        </span>
                      ))}
                      {categoryCount > 3 && (
                        <span className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded">
                          +{categoryCount - 3} more
                        </span>
                      )}
                    </div>
                  </div>
                )}

                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    console.log('viewing pipeline details', pipeline);
                    setSelectedPipeline(pipeline.id);
                  }}
                  className="mt-4 w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                >
                  View Details
                </button>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

