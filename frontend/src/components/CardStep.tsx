import type { DynamicPipelineStep, GoogleDockingResult } from '../types';
import CardUrlResults from './CardUrlResults';

// Helper to extract URLs from step output
const extractUrls = (output: any): GoogleDockingResult[] => {
  if (!output) return [];
  if (Array.isArray(output)) {
    return output.filter((item) => item && (item.link || item.url));
  }
  return [];
};

// Helper to get domain type display name
const getDomainTypeDisplay = (domainType: string): string => {
  return domainType
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
};

interface CardStepProps {
  step: DynamicPipelineStep;
  showDetails?: boolean;
}

export default function CardStep({ step, showDetails = true }: CardStepProps) {
  const results = extractUrls(step.output);
  const hasResults = results.length > 0;

  return (
    <div
      key={step.id}
      className={`p-4 border rounded-lg transition-all ${
        step.success
          ? 'border-green-300 bg-green-50/50 hover:bg-green-50'
          : 'border-red-300 bg-red-50/50 hover:bg-red-50'
      }`}
    >
      <div className="flex items-start justify-between mb-2">
        <div className="flex items-center gap-2 flex-1">
          <span className="font-semibold text-gray-900 text-sm">
            {getDomainTypeDisplay(step.domain_type)}
          </span>
          <span
            className={`px-2 py-1 text-xs rounded-full ${
              step.success
                ? 'bg-green-200 text-green-800'
                : 'bg-red-200 text-red-800'
            }`}
          >
            {step.success ? '✓ Success' : '✗ Failed'}
          </span>
          {step.depth !== undefined && (
            <span className="px-2 py-1 bg-gray-200 text-gray-700 text-xs rounded">
              Depth {step.depth}
            </span>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <p className="text-sm text-gray-700">
          <span className="font-medium">Query:</span>{' '}
          <span className="text-gray-900">{step.search_parameter}</span>
        </p>

        {step.error && (
          <div className="p-2 bg-red-100 border border-red-300 rounded text-sm text-red-800">
            <span className="font-medium">Error:</span> {step.error}
          </div>
        )}

        {showDetails && step.keywords && step.keywords.length > 0 && (
          <div>
            <p className="text-xs font-medium text-gray-700 mb-1">Keywords:</p>
            <div className="flex flex-wrap gap-1">
              {step.keywords.map((keyword, idx) => (
                <span
                  key={idx}
                  className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded"
                >
                  {keyword}
                </span>
              ))}
            </div>
          </div>
        )}

        {showDetails &&
          step.keywords_per_category &&
          Object.keys(step.keywords_per_category).length > 0 && (
            <div>
              <p className="text-xs font-medium text-gray-700 mb-1">Keywords by Category:</p>
              <div className="space-y-1">
                {Object.entries(step.keywords_per_category).map(([category, keywords]) => (
                  <div key={category}>
                    <span className="text-xs font-medium text-gray-600 capitalize">
                      {category.replace(/_/g, ' ')}:
                    </span>{' '}
                    <div className="inline-flex flex-wrap gap-1">
                      {keywords.map((kw, idx) => (
                        <span
                          key={idx}
                          className="px-1.5 py-0.5 bg-gray-100 text-gray-700 text-xs rounded"
                        >
                          {kw}
                        </span>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

        {hasResults && showDetails && <CardUrlResults results={results} />}
      </div>
    </div>
  );
}

