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
      <td className="px-6 py-4 text-sm text-gray-500 truncate max-w-xs">
        {scjCase.desc_materia}
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
        {scjCase.fecha_fallo}
      </td>
    </tr>
  );
}

