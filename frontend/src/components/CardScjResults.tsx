import type { ScjCase } from '../types';

interface CardScjResultsProps {
  results: ScjCase[];
}

export default function CardScjResults({ results }: CardScjResultsProps) {
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
                {/* {result.url_blob && ( */}
                  <div>
                    <span className="text-gray-500">URL:</span>{' '}
                    <a
                      href={result.url_blob}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:underline hover:text-blue-800 transition-colors"
                    >
                      {result.url_blob.split('/').pop() || 'No URL'}
                    </a>
                  </div>
                {/* )} */}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

