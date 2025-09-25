import React, { useEffect, useState } from 'react';
import { getTools, type DomainTool } from '../api/client';

export const ToolListDemo: React.FC = () => {
  const [tools, setTools] = useState<DomainTool[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTools = async () => {
      try {
        setLoading(true);
        // The API client is fully typed with modern async/await!
        const response = await getTools({
          query: {
            limit: 10,
            offset: 0,
          },
        });

        // The response has the data property
        if (response.data) {
          // Response is typed as Record<string, DomainTool[]>
          const toolsArray = Object.values(response.data).flat();
          setTools(toolsArray);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch tools');
      } finally {
        setLoading(false);
      }
    };

    fetchTools();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2">Loading tools...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4 m-4">
        <div className="text-red-800">
          <strong>Error:</strong> {error}
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">Tool Inventory</h2>

      {tools.length === 0 || tools[0] == null ? (
        <div className="text-center py-8 text-gray-500">
          No tools found. Try adding some tools first.
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {tools.map((tool) => (
            <div key={tool.id} className="bg-white rounded-lg shadow-md p-6 border border-gray-200">
              <div className="flex justify-between items-start mb-4">
                <h3 className="text-lg font-semibold text-gray-900">{tool.name}</h3>
                <StatusBadge status={tool.status} />
              </div>

              <div className="space-y-2 text-sm text-gray-600">
                <div>
                  <strong>ID:</strong> {tool.id}
                </div>
                {tool.current_user_id && (
                  <div>
                    <strong>Checked out to:</strong> {tool.current_user_id}
                  </div>
                )}
                {tool.last_checked_out_at && (
                  <div>
                    <strong>Last checked out:</strong>{' '}
                    {new Date(tool.last_checked_out_at).toLocaleDateString()}
                  </div>
                )}
                <div>
                  <strong>Created:</strong> {new Date(tool.created_at ?? '').toLocaleDateString()}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

const StatusBadge: React.FC<{ status?: string }> = ({ status }) => {
  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'IN_OFFICE':
        return 'bg-green-100 text-green-800';
      case 'CHECKED_OUT':
        return 'bg-yellow-100 text-yellow-800';
      case 'MAINTENANCE':
        return 'bg-orange-100 text-orange-800';
      case 'LOST':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusText = (status?: string) => {
    switch (status) {
      case 'IN_OFFICE':
        return 'In Office';
      case 'CHECKED_OUT':
        return 'Checked Out';
      case 'MAINTENANCE':
        return 'Maintenance';
      case 'LOST':
        return 'Lost';
      default:
        return 'Unknown';
    }
  };

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(status)}`}
    >
      {getStatusText(status)}
    </span>
  );
};
