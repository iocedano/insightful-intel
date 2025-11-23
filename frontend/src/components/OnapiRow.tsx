import type { Entity } from '../types';

interface OnapiRowProps {
  entity: Entity;
  index: number;
}

export default function OnapiRow({ entity, index }: OnapiRowProps) {
  return (
    <tr key={entity.id || index} className="hover:bg-gray-50">
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
        {entity.serie_expediente}-{entity.numero_expediente}
      </td>
      <td className="px-6 py-4 text-sm text-gray-900">{entity.tipo}</td>
      <td className="px-6 py-4 text-sm text-gray-900">{entity.subtipo}</td>
      <td className="px-6 py-4 text-sm text-gray-500 truncate max-w-xs">
        {entity.texto}
      </td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
        {entity.titular}
      </td>
    </tr>
  );
}

