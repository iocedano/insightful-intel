import { useState, useEffect, useRef } from 'react';
import { api } from '../api';
import type { DynamicPipelineResult } from '../types';
import PipelineDetails from '../components/PipelineDetails';

export default function Pipeline() {
  const [query, setQuery] = useState('');
  const [depth, setDepth] = useState(3);
  const [skipDuplicates, setSkipDuplicates] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pipelineId, setPipelineId] = useState<string | null>(null);
  const [pipelineResult, setPipelineResult] = useState<DynamicPipelineResult | null>(null);
  const pollingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const handleExecute = async () => {
    if (!query.trim()) {
      setError('Please enter a search query');
      return;
    }

    // Clear any existing polling interval
    if (pollingIntervalRef.current) {
      clearInterval(pollingIntervalRef.current);
      pollingIntervalRef.current = null;
    }
    setPipelineId(null);


    api.executeDynamicPipeline(query, depth, skipDuplicates, false).then((response) => {
      setPipelineId(response.execution_id);
      setLoading(true);
      setError(null);
    }).catch((err) => {
      console.error('Error executing pipeline:', err);  
      setError('Failed to execute pipeline');
      setPipelineResult(null);
      setLoading(false);
    });
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
    const fetchPipelineResult = async () => {
      api.getPipelineByID(pipelineId).then((response) => {
        if (response.success && response.data) {
          setPipelineResult(response.data as DynamicPipelineResult);
        } else {
          setError('Failed to fetch pipeline result');
        }
      }).catch((err) => {
        setLoading(false);
        setError('Failed to fetch pipeline result');
        console.error('Error fetching pipeline result:', err);
      })
    };

    // Fetch immediately
    fetchPipelineResult();

    // Set up polling interval (5 seconds = 5000 milliseconds)
    pollingIntervalRef.current = setInterval(fetchPipelineResult, 5000);

    // Cleanup function
    return () => {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
        pollingIntervalRef.current = null;
      }
    };
  }, [pipelineId]);

  useEffect(() => {
    if (pipelineResult?.total_steps && pipelineResult?.total_steps > 0) {
      if (pollingIntervalRef.current) {
        clearInterval(pollingIntervalRef.current);
        pollingIntervalRef.current = null;
        setLoading(false);
        setError(null);
      }
    }
  }, [pipelineResult?.total_steps || 0]);


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
          </div>

          <div className="flex space-x-4">
            <button
              onClick={handleExecute}
              disabled={loading}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {loading ? 'Executing...' : 'Execute Pipeline'}
            </button>
          </div>
        </div>

        {error && (
          <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {/* {pipelineResult && (
          <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-md">
            <h3 className="font-semibold text-blue-900 mb-2">Summary</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span className="text-blue-700">Total Steps:</span>
                <span className="ml-2 font-semibold">{pipelineResult?.total_steps || 0}</span>
              </div>
              <div>
                <span className="text-blue-700">Successful:</span>
                <span className="ml-2 font-semibold text-green-600">
                  {pipelineResult?.successful_steps || pipelineResult?.steps?.filter(s => s.success).length}
                </span>
              </div>
              <div>
                <span className="text-blue-700">Failed:</span>
                <span className="ml-2 font-semibold text-red-600">
                  {pipelineResult?.failed_steps || pipelineResult?.steps?.filter(s => !s.success).length}
                </span>
              </div>
              <div>
                <span className="text-blue-700">Max Depth:</span>
                <span className="ml-2 font-semibold">{pipelineResult?.max_depth_reached || 0}</span>
              </div>
            </div>
          </div>
        )} */}

        {loading && (
          <div className="mt-4 p-4 bg-gray-50 border border-gray-200 rounded-md">
            <p className="text-gray-700">Loading...</p>
          </div>
        )}
      </div>

      {pipelineResult && pipelineResult.steps?.length && pipelineResult.steps.length > 0 && (
        <PipelineDetails
          pipeline={pipelineResult}
          steps={pipelineResult?.steps || []}
          showBackButton={false}
        />
      )}
    </div>
  );
}
