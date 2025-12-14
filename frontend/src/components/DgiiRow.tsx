import type { Register } from '../types';

interface DgiiRowProps {
  register: Register;
  index: number;
}

export default function DgiiRow({ register, index }: DgiiRowProps) {
  return (
    <tr key={register.id || index} className="hover:bg-gray-50">
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
        {register.rnc}
      </td>
      <td className="px-6 py-4 text-sm text-gray-900">{register.razon_social}</td>
      <td className="px-6 py-4 text-sm text-gray-500">{register.nombre_comercial}</td>
      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
        {register.estado}
      </td>
    </tr>
  );
}

