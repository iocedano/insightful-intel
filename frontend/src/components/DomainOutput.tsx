import { type DomainSearchResult, type Entity, type ScjCase, type PgrNews, type GoogleDockingResult, type DomainType, DOMAIN_TYPE_MAP, type Register  } from '../types';
import OnapiRow from './OnapiRow';
import ScjRow from './ScjRow';
import DgiiRow from './DgiiRow';
import PgrRow from './PgrRow';
import DockingRow from './DockingRow';
import { getHeaders } from '../helpers/common';

interface DomainOutputProps {
  result: DomainSearchResult;
}

export default function DomainOutput({ result }: DomainOutputProps) {
  console.log('result', result);
  if (!result.output || !Array.isArray(result.output) || result.output.length === 0) {
    console.log('result.output', result.output);
    return null;
  }

  // Normalize domain_type to handle both uppercase constants and lowercase strings
  const domainType = result.name.toLowerCase();

  const renderRow = (item: Entity | ScjCase | PgrNews | GoogleDockingResult | Register, index: number) => {
    console.log('item', item);
    switch (domainType) {
      case DOMAIN_TYPE_MAP.ONAPI:
        return <OnapiRow key={(item as Entity).id || index} entity={item as Entity} index={index} />;
      case DOMAIN_TYPE_MAP.DGII:
        return <DgiiRow key={(item as Register).id || index} register={item as Register} index={index} />;
      case DOMAIN_TYPE_MAP.SCJ:
        return <ScjRow key={(item as ScjCase).id || index} scjCase={item as ScjCase} index={index} />;
      case DOMAIN_TYPE_MAP.PGR:
        return <PgrRow key={(item as PgrNews).id || index} news={item as PgrNews} index={index} />;
      case DOMAIN_TYPE_MAP.GOOGLE_DOCKING: case DOMAIN_TYPE_MAP.SOCIAL_MEDIA: case DOMAIN_TYPE_MAP.FILE_TYPE: case DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA:
        return <DockingRow key={(item as GoogleDockingResult).id || index} docking={item as GoogleDockingResult} index={index} />;
      default:
        return null;
    }
  };

  const headers = getHeaders(domainType as DomainType);
  if (headers.length === 0) {
    console.log('headers', headers);
    return null;
  }

  return (
    <div className="mt-2 overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-100">
          <tr>
            {headers.map((header) => (
              <th
                key={header}
                className="px-4 py-2 text-left text-xs font-medium text-gray-700 uppercase tracking-wider"
              >
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {result.output.map((item, index) => renderRow(item, index))}
        </tbody>
      </table>
    </div>
  );
}

