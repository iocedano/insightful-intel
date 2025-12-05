import { useState } from 'react';
import { api } from '../api';
import type { DomainSearchResult } from '../types';
import DomainOutput from '../components/DomainOutput';

export default function Search() {
  const [query, setQuery] = useState('');
  const [domain, setDomain] = useState('');
  const [results, setResults] = useState<DomainSearchResult | DomainSearchResult[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    console.log('handleSearch', query, domain);
    if (!query.trim()) {
      setError('Please enter a search query');
      return;
    }

    setLoading(true);
    setError(null);
    setResults(null);

    try {
      const data = await api.search(query, domain || undefined);
      setResults(data);
      console.log('data', data);
    } catch (err: any) {
      console.error('err', err);
      setError(err.message || 'Search failed');
      setResults(null);
    } finally {
      setLoading(false);
    }
  };

  const renderResult = (result: DomainSearchResult, index?: number) => (
    <div
      key={result.id || index}
      className={`p-4 border rounded-md ${
        result.success
          ? 'border-green-200 bg-green-50'
          : 'border-red-200 bg-red-50'
      }`}
    >
      <div className="flex items-start justify-between mb-2">
        <div>
          <span className="font-semibold text-gray-900">{result.name}</span>
          {result.success ? (
            <span className="ml-2 px-2 py-1 bg-green-100 text-green-800 text-xs rounded">
              Success
            </span>
          ) : (
            <span className="ml-2 px-2 py-1 bg-red-100 text-red-800 text-xs rounded">
              Failed
            </span>
          )}
        </div>
      </div>
      <p className="text-sm text-gray-600 mb-2">
        Query: <span className="font-medium">{result.search_parameter}</span>
      </p>
      {result.error && (
        <p className="text-sm text-red-600 mb-2">Error: {result.error}</p>
      )}
      {result.keywords_per_category && (
        <div className="mt-2">
          <p className="text-xs font-medium text-gray-700 mb-1">Keywords:</p>
          <div className="flex flex-wrap gap-1">
            {Object.entries(result.keywords_per_category).map(([cat, keywords]) =>
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
      {result && 'output' in result ? <DomainOutput result={result as DomainSearchResult} /> : null}
    </div>
  );

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold mb-4">Domain Search</h2>

        <div className="space-y-4 mb-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Search Query
            </label>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              placeholder="Enter search query"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Domain (Optional)
            </label>
            <select
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">All Domains</option>
              <option value="onapi">ONAPI</option>
              <option value="scj">SCJ</option>
              <option value="dgii">DGII</option>
              <option value="pgr">PGR</option>
              <option value="docking">Google Docking</option>
              <option value="social_media">Social Media</option>
              <option value="file_type">File Type</option>
              <option value="x_social_media">X Social Media</option>
            </select>
          </div>

          <button
            onClick={handleSearch}
            disabled={loading}
            className="w-full px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {loading ? 'Searching...' : 'Search'}
          </button>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {results && (
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">
              Results {results && 'output' in results && Array.isArray(results.output) ? `(${results.output.length})` : ''}
            </h3>
            {Array.isArray(results) ? (
              results.map((result, index) => renderResult(result, index))
            ) : (
              results && 'output' in results ? renderResult(results as DomainSearchResult) : null
            )}
          </div>
        )}
      </div>
    </div>
  );
}



