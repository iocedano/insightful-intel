import { useState, useEffect, useRef } from 'react';
import { api, createPipelineStream } from '../api';
import type { DynamicPipelineStep } from '../types';

export default function Pipeline() {
  const [query, setQuery] = useState('');
  const [depth, setDepth] = useState(3);
  const [skipDuplicates, setSkipDuplicates] = useState(true);
  const [streaming, setStreaming] = useState(false);
  const [steps, setSteps] = useState<DynamicPipelineStep[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [summary, setSummary] = useState<any>(null);
  const [pipelineId, setPipelineId] = useState<string | null>(null);
  const pollingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const handleExecute = async () => {
    if (!query.trim()) {
      setError('Please enter a search query');
      return;
    }

    setLoading(true);
    setError(null);
    setSteps([]);
    setSummary(null);
    
    // Clear any existing polling interval
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    setPipelineId(null);

    try {
      if (streaming) {
        // Use streaming
        const closeStream = createPipelineStream(
          query,
          depth,
          skipDuplicates,
          (data) => {
            if (data.step) {
              setSteps((prev) => [...prev, data.step]);
            }
            if (data.step?.domain_type === 'SUMMARY') {
              setSummary(data.step.output);
            }
          },
          () => {
            setLoading(false);
          },
          (err) => {
            setError(err.message);
            setLoading(false);
          }
        );

        // Store close function for cleanup
        return () => closeStream();
      } else {
        // Use regular API call - returns execution_id immediately
        const response = await api.executeDynamicPipeline(query, depth, skipDuplicates, false);
        
        // Extract execution ID from response (pipeline runs in background)
        if (response.execution_id) {
          setPipelineId(response.execution_id);
        }
        
        setLoading(false);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to execute pipeline');
      setLoading(false);
    }
  };

  // Poll for pipeline steps every 15 seconds
  useEffect(() => {
    if (!pipelineId) {
      // Clear any existing interval if pipelineId is cleared
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
        pollingIntervalRef.current = null;
      }
      return;
    }

    // Function to fetch and update steps
    const fetchSteps = async () => {
      try {
        const response = await api.getPipelineSteps(pipelineId);
        if (response.success && response.steps) {
          // Map the steps to the expected format
          const mappedSteps = response.steps.map((step: any) => ({
            id: step.id,
            domain_type: step.domain_type || step.name,
            search_parameter: step.search_parameter,
            success: step.success,
            error: step.error,
            output: step.output,
            keywords_per_category: step.keywords_per_category,
            depth: step.depth || 0,
            category: step.category,
            keywords: step.keywords,
          }));
          setSteps(mappedSteps);
        }
      } catch (err) {
        console.error('Error fetching pipeline steps:', err);
        // Don't show error to user for polling failures, just log it
      }
    };

    // Fetch immediately
    fetchSteps();

    // Set up polling interval (15 seconds = 15000 milliseconds)
    pollingIntervalRef.current = setInterval(fetchSteps, 15000);

    // Cleanup function
    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
        pollingIntervalRef.current = null;
      }
    };
  }, [pipelineId]);

  const handleSave = async () => {
    if (steps.length === 0) {
      setError('No pipeline results to save');
      return;
    }

    try {
      const pipelineData = {
        steps: steps,
        total_steps: steps.length,
        successful_steps: steps.filter(s => s.success).length,
        failed_steps: steps.filter(s => !s.success).length,
        max_depth_reached: Math.max(...steps.map(s => s.depth || 0)),
        config: {
          max_depth: depth,
          max_concurrent_steps: 10,
          delay_between_steps: 2,
          skip_duplicates: skipDuplicates,
        },
      };

      await api.savePipeline(pipelineData);
      alert('Pipeline saved successfully!');
    } catch (err: any) {
      setError(err.message || 'Failed to save pipeline');
    }
  };

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-4">Dynamic Pipeline</h2>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Search Query
            </label>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Enter search query (e.g., Novasco)"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Max Depth
              </label>
              <input
                type="number"
                min="1"
                max="10"
                value={depth}
                onChange={(e) => setDepth(parseInt(e.target.value) || 3)}
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div className="flex items-center space-x-2 pt-8">
              <input
                type="checkbox"
                id="skipDuplicates"
                checked={skipDuplicates}
                onChange={(e) => setSkipDuplicates(e.target.checked)}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label htmlFor="skipDuplicates" className="text-sm font-medium text-gray-700">
                Skip Duplicates
              </label>
            </div>

            <div className="flex items-center space-x-2 pt-8">
              <input
                type="checkbox"
                id="streaming"
                checked={streaming}
                onChange={(e) => setStreaming(e.target.checked)}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label htmlFor="streaming" className="text-sm font-medium text-gray-700">
                Stream Results
              </label>
            </div>
          </div>

          <div className="flex space-x-4">
            <button
              onClick={handleExecute}
              disabled={loading}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {loading ? 'Executing...' : 'Execute Pipeline'}
            </button>
            {steps.length > 0 && (
              <button
                onClick={handleSave}
                className="px-6 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
              >
                Save Pipeline
              </button>
            )}
          </div>
        </div>

        {error && (
          <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {summary && (
          <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-md">
            <h3 className="font-semibold text-blue-900 mb-2">Summary</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span className="text-blue-700">Total Steps:</span>
                <span className="ml-2 font-semibold">{summary.total_steps || steps.length}</span>
              </div>
              <div>
                <span className="text-blue-700">Successful:</span>
                <span className="ml-2 font-semibold text-green-600">
                  {summary.successful_steps || steps.filter(s => s.success).length}
                </span>
              </div>
              <div>
                <span className="text-blue-700">Failed:</span>
                <span className="ml-2 font-semibold text-red-600">
                  {summary.failed_steps || steps.filter(s => !s.success).length}
                </span>
              </div>
              <div>
                <span className="text-blue-700">Max Depth:</span>
                <span className="ml-2 font-semibold">{summary.max_depth_reached || 0}</span>
              </div>
            </div>
          </div>
        )}
      </div>

      {steps.length > 0 && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-xl font-bold mb-4">Pipeline Steps ({steps.length})</h3>
          <div className="space-y-4 max-h-96 overflow-y-auto">
            {steps.map((step, index) => (
              <div
                key={index}
                className={`p-4 border rounded-md ${
                  step.success
                    ? 'border-green-200 bg-green-50'
                    : 'border-red-200 bg-red-50'
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-2 mb-2">
                      <span className="font-semibold text-gray-900">
                        {step.domain_type}
                      </span>
                      <span className="text-sm text-gray-500">
                        Depth: {step.depth || 0}
                      </span>
                      {step.success ? (
                        <span className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded">
                          Success
                        </span>
                      ) : (
                        <span className="px-2 py-1 bg-red-100 text-red-800 text-xs rounded">
                          Failed
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-600 mb-2">
                      Query: <span className="font-medium">{step.search_parameter}</span>
                    </p>
                    {step.error && (
                      <p className="text-sm text-red-600">Error: {step.error}</p>
                    )}
                    {step.keywords_per_category && (
                      <div className="mt-2">
                        <p className="text-xs font-medium text-gray-700 mb-1">Keywords:</p>
                        <div className="flex flex-wrap gap-1">
                          {Object.entries(step.keywords_per_category).map(([cat, keywords]) =>
                            (keywords as string[]).map((kw, i) => (
                              <span
                                key={`${cat}-${i}`}
                                className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
                              >
                                {kw}
                              </span>
                            ))
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}



