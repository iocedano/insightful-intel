import type { PgrNews } from '../types';

interface PgrRowProps {
  news: PgrNews;
  index: number;
}

export default function PgrRow({ news, index }: PgrRowProps) {
  return (
    <tr key={news.id || index} className="hover:bg-gray-50">
      <td className="px-6 py-4 text-sm text-gray-900 truncate max-w-md">
        {news.title}
      </td>
      <td className="px-6 py-4 text-sm">
        <a
          href={news.url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-600 hover:text-blue-800 truncate max-w-xs block"
        >
          {news.url}
        </a>
      </td>
    </tr>
  );
}

