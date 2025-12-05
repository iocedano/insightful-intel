import type { GoogleDockingResult } from '../types';

interface DockingRowProps {
  docking: GoogleDockingResult;
  index: number;
}

export default function DockingRow({ docking, index }: DockingRowProps) {
  console.log('docking', docking);
  return (
    <tr key={`${docking.url}-${index}`} className="hover:bg-gray-50">
      <td className="px-6 py-4 text-sm text-gray-900 truncate max-w-md">
        {docking.title}
      </td>
      <td className="px-6 py-4 text-sm text-gray-500 truncate max-w-xs">
        {docking.description || docking.snippet}
      </td>
      <td className="px-6 py-4 text-sm">
        <a
          href={docking.link}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-600 hover:text-blue-800 truncate max-w-xs block"
        >
          {docking.link}
        </a>
      </td>
    </tr>
  );
}

