import type { Register } from '../types';

interface CardDgiiResultsProps {
  results: Register[];
}

export default function CardDgiiResults({ results }: CardDgiiResultsProps) {
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
}

