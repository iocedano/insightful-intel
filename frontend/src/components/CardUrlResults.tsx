import type { GoogleDockingResult } from '../types';

interface CardUrlResultsProps {
  results: GoogleDockingResult[];
}

export default function CardUrlResults({ results }: CardUrlResultsProps) {
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
}

