import type { ScjCase } from '../types';

interface ScjRowProps {
  scjCase: ScjCase;
  index: number;
}

export default function ScjRow({ scjCase, index }: ScjRowProps) {
  return (
    <tr key={scjCase.id || index} className="hover:bg-gray-50">
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
        {scjCase.id_expediente}
      </td>
      <td className="px-6 py-4 text-sm text-gray-900">{scjCase.no_sentencia}</td>
      <td className="px-6 py-4 text-sm text-gray-900">{scjCase.involucrados}</td>
      <td className="px-6 py-4 text-sm text-gray-500 truncate max-w-xs">
        {scjCase.desc_tribunal}
      </td>
      <td className="px-6 py-4 text-sm text-gray-900">
        <a
          href={scjCase.url_blob}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-600 hover:underline hover:text-blue-800 transition-colors"
        >
          {scjCase?.url_blob?.split('/').pop() || 'No URL'}
        </a>
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
        {scjCase.fecha_fallo}
      </td>
    </tr>
  );
}

