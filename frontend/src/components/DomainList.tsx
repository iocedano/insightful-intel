import { useState, useEffect } from 'react';
import { api } from '../api';
import { type Entity, type ScjCase, type Register, type PgrNews, type GoogleDockingResult, type DomainType, DOMAIN_TYPE_MAP} from '../types';
import OnapiRow from './OnapiRow';
import ScjRow from './ScjRow';
import DgiiRow from './DgiiRow';
import PgrRow from './PgrRow';
import DockingRow from './DockingRow';
import { getHeaders } from '../helpers/common';

interface DomainListProps {
  domain: DomainType;
  title: string;
}

export default function DomainList({ domain, title }: DomainListProps) {
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [limit] = useState(10);
  const [count, setCount] = useState(0);

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      let response: any;
      switch (domain) {
        case DOMAIN_TYPE_MAP.ONAPI:
          response = await api.getOnapi(offset, limit);
          break;
        case DOMAIN_TYPE_MAP.SCJ:
          response = await api.getScj(offset, limit);
          break;
        case DOMAIN_TYPE_MAP.DGII:
          response = await api.getDgii(offset, limit);
          break;
        case DOMAIN_TYPE_MAP.PGR:
          response = await api.getPgr(offset, limit);
          break;
        case DOMAIN_TYPE_MAP.GOOGLE_DOCKING:
          response = await api.getDocking(offset, limit);
          break;
      }
      setData(response.data || []);
      setCount(response.count || 0);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [offset, domain]);

  const renderRow = (item: any, index: number) => {
    switch (domain) {
      case DOMAIN_TYPE_MAP.ONAPI:
        return <OnapiRow key={(item as Entity).id || index} entity={item as Entity} index={index} />;
      case DOMAIN_TYPE_MAP.SCJ:
        return <ScjRow key={(item as ScjCase).id || index} scjCase={item as ScjCase} index={index} />;
      case DOMAIN_TYPE_MAP.DGII:
        return <DgiiRow key={(item as Register).id || index} register={item as Register} index={index} />;
      case DOMAIN_TYPE_MAP.PGR:
        return <PgrRow key={(item as PgrNews).id || index} news={item as PgrNews} index={index} />;
      case DOMAIN_TYPE_MAP.GOOGLE_DOCKING: case DOMAIN_TYPE_MAP.SOCIAL_MEDIA: case DOMAIN_TYPE_MAP.FILE_TYPE: case DOMAIN_TYPE_MAP.X_SOCIAL_MEDIA:
        return <DockingRow key={(item as GoogleDockingResult).id || index} docking={item as GoogleDockingResult} index={index} />;
      default:
        return null;
    }
  };


  return (
    <div className="max-w-7xl mx-auto p-6">
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold">{title}</h2>
          <span className="text-sm text-gray-500">Total: {count}</span>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {loading ? (
          <div className="text-center py-8">
            <p className="text-gray-500">Loading...</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    {getHeaders(domain).map((header) => (
                      <th
                        key={header}
                        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                      >
                        {header}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {data.length === 0 ? (
                    <tr>
                      <td
                        colSpan={getHeaders(domain).length}
                        className="px-6 py-4 text-center text-gray-500"
                      >
                        No data available
                      </td>
                    </tr>
                  ) : (
                    data.map((item, index) => renderRow(item, index))
                  )}
                </tbody>
              </table>
            </div>

            <div className="mt-4 flex items-center justify-between">
              <button
                onClick={() => setOffset(Math.max(0, offset - limit))}
                disabled={offset === 0 || loading}
                className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <span className="text-sm text-gray-500">
                Showing {offset + 1} to {Math.min(offset + limit, count)} of {count}
              </span>
              <button
                onClick={() => setOffset(offset + limit)}
                disabled={offset + limit >= count || loading}
                className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}



