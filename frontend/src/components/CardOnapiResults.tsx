import type { Entity } from '../types';

interface CardOnapiResultsProps {
  results: Entity[];
}

export default function CardOnapiResults({ results }: CardOnapiResultsProps) {
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
}

