import React, { useState } from 'react';
import { fileService, FilterParams } from '../services/fileService';
import { File as FileType } from '../types/file';
import { DocumentIcon, TrashIcon, ArrowDownTrayIcon, MagnifyingGlassIcon, FunnelIcon } from '@heroicons/react/24/outline';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

interface FilterState {
  search: string;
  file_type: string;
  size_min: string;
  size_max: string;
  uploaded_after: string;
  uploaded_before: string;
}

const initialFilterState: FilterState = {
  search: '',
  file_type: '',
  size_min: '',
  size_max: '',
  uploaded_after: '',
  uploaded_before: '',
};

export const FileList: React.FC = () => {
  const queryClient = useQueryClient();
  const [inputFilters, setInputFilters] = useState<FilterState>(initialFilterState);
  const [submittedFilters, setSubmittedFilters] = useState<FilterState>(initialFilterState);
  const [showAdvanced, setShowAdvanced] = useState(false);

  // Query for fetching files
  const { data: files, isLoading, error } = useQuery({
    queryKey: ['files', submittedFilters],
    queryFn: () => {
      const apiParams: FilterParams = {
        ...submittedFilters,
        size_min: submittedFilters.size_min ? parseInt(submittedFilters.size_min, 10) * 1024 : undefined,
        size_max: submittedFilters.size_max ? parseInt(submittedFilters.size_max, 10) * 1024 : undefined,
      };
      return fileService.getFiles(apiParams);
    },
  });

  // Mutation for deleting files
  const deleteMutation = useMutation({
    mutationFn: fileService.deleteFile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['files'] });
    },
  });

  // Mutation for downloading files
  const downloadMutation = useMutation({
    mutationFn: ({ fileUrl, filename }: { fileUrl: string; filename: string }) =>
      fileService.downloadFile(fileUrl, filename),
  });

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
    } catch (err) {
      console.error('Delete error:', err);
    }
  };

  const handleDownload = async (fileUrl: string, filename: string) => {
    try {
      await downloadMutation.mutateAsync({ fileUrl, filename });
    } catch (err) {
      console.error('Download error:', err);
    }
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setInputFilters(prev => ({ ...prev, [name]: value }));
  };

  const handleFilterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmittedFilters(inputFilters);
  };

  const clearFilters = () => {
    setInputFilters(initialFilterState);
    setSubmittedFilters(initialFilterState);
  };

  // Determine button states
  // Search is enabled only if the inputs have changed since the last search.
  const isDirty = JSON.stringify(inputFilters) !== JSON.stringify(submittedFilters);
  // Clear is enabled if any filter input has a value.
  const hasInput = JSON.stringify(inputFilters) !== JSON.stringify(initialFilterState);

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-4 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-3">
            <div className="h-8 bg-gray-200 rounded"></div>
            <div className="h-8 bg-gray-200 rounded"></div>
            <div className="h-8 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <div className="bg-red-50 border-l-4 border-red-400 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-red-400"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-red-700">Failed to load files. Please try again.</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-gray-900">File Vault</h2>
      </div>

      <form onSubmit={handleFilterSubmit} className="bg-gray-50 p-4 rounded-lg border border-gray-200 mb-6 space-y-4">
        {/* Top row for primary actions */}
        <div className="flex flex-col md:flex-row md:justify-between md:items-end gap-4">
          {/* Search input - takes up remaining space */}
          <div className="flex-grow">
            <label htmlFor="search" className="sr-only">Filename</label>
            <input
              type="text"
              name="search"
              id="search"
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="Search by filename..."
              value={inputFilters.search}
              onChange={handleFilterChange}
            />
          </div>

          {/* Action buttons */}
          <div className="flex items-center space-x-2">
            <button
              type="button"
              onClick={() => setShowAdvanced(!showAdvanced)}
              className="inline-flex items-center justify-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
            >
              <FunnelIcon className="h-5 w-5 mr-2" />
              Filters
            </button>
            <button
              type="button"
              onClick={clearFilters}
              disabled={!hasInput}
              className="inline-flex justify-center items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-50 disabled:text-gray-400 disabled:cursor-not-allowed"
            >
              Clear
            </button>
            <button
              type="submit"
              disabled={!isDirty}
              className="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:bg-primary-300 disabled:cursor-not-allowed"
            >
              <MagnifyingGlassIcon className="h-5 w-5 mr-2" />
              Search
            </button>
          </div>
        </div>

        {showAdvanced && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 border-t border-gray-200 pt-4">
            <div>
              <label htmlFor="file_type" className="block text-sm font-medium text-gray-700">File Type</label>
              <input type="text" name="file_type" id="file_type" value={inputFilters.file_type} onChange={handleFilterChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm" placeholder="e.g., image/jpeg" />
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div>
                <label htmlFor="size_min" className="block text-sm font-medium text-gray-700">Min Size (KB)</label>
                <input type="number" name="size_min" id="size_min" value={inputFilters.size_min} onChange={handleFilterChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm" placeholder="0" />
              </div>
              <div>
                <label htmlFor="size_max" className="block text-sm font-medium text-gray-700">Max Size (KB)</label>
                <input type="number" name="size_max" id="size_max" value={inputFilters.size_max} onChange={handleFilterChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm" placeholder="10240" />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <div>
                <label htmlFor="uploaded_after" className="block text-sm font-medium text-gray-700">Uploaded After</label>
                <input type="date" name="uploaded_after" id="uploaded_after" value={inputFilters.uploaded_after} onChange={handleFilterChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm" />
              </div>
              <div>
                <label htmlFor="uploaded_before" className="block text-sm font-medium text-gray-700">Uploaded Before</label>
                <input type="date" name="uploaded_before" id="uploaded_before" value={inputFilters.uploaded_before} onChange={handleFilterChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm" />
              </div>
            </div>
          </div>
        )}
      </form>

      {!files || files.length === 0 ? (
        <div className="text-center py-12">
          <DocumentIcon className="mx-auto h-12 w-12 text-gray-400" />
          <h3 className="mt-2 text-sm font-medium text-gray-900">No files</h3>
          <p className="mt-1 text-sm text-gray-500">
            Get started by uploading a file
          </p>
        </div>
      ) : (
        <div className="mt-6 flow-root">
          <ul className="-my-5 divide-y divide-gray-200">
            {files.map((file) => (
              <li key={file.id} className="py-4">
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <DocumentIcon className="h-8 w-8 text-gray-400" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {file.original_filename}
                    </p>
                    <p className="text-sm text-gray-500">
                      {file.file_type} • {(file.size / 1024).toFixed(2)} KB
                    </p>
                    <p className="text-sm text-gray-500">
                      Uploaded {new Date(file.uploaded_at).toLocaleString()}
                    </p>
                  </div>
                  <div className="flex space-x-2">
                    <button
                      onClick={() => handleDownload(file.file, file.original_filename)}
                      disabled={downloadMutation.isPending}
                      className="inline-flex items-center px-3 py-2 border border-transparent shadow-sm text-sm leading-4 font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                    >
                      <ArrowDownTrayIcon className="h-4 w-4 mr-1" />
                      Download
                    </button>
                    <button
                      onClick={() => handleDelete(file.id)}
                      disabled={deleteMutation.isPending}
                      className="inline-flex items-center px-3 py-2 border border-transparent shadow-sm text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                    >
                      <TrashIcon className="h-4 w-4 mr-1" />
                      Delete
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}; 