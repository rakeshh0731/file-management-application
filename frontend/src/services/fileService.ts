import api from './api';
import { File as FileType } from '../types/file';

const API_ORIGIN = new URL(api.defaults.baseURL!).origin;

export interface FilterParams {
  search?: string;
  file_type?: string;
  size_min?: number;
  size_max?: number;
  uploaded_after?: string;
  uploaded_before?: string;
}

export const fileService = {
  async uploadFile(file: File): Promise<FileType> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await api.post(`/files/`, formData);
    return response.data;
  },

  async getFiles(filters: FilterParams): Promise<FileType[]> {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== null && value !== undefined && value !== '') {
        params.append(key, String(value));
      }
    });
    const response = await api.get(`/files/`, { params });
    return response.data;
  },

  async deleteFile(id: string): Promise<void> {
    await api.delete(`/files/${id}/`);
  },

  async downloadFile(fileUrl: string, filename: string): Promise<void> {
    try {
      // The fileUrl from the backend is a path (e.g., /uploads/...).
      // We need to construct the full URL to the backend server.
      const fullUrl = `${API_ORIGIN}${fileUrl}`;
      const response = await api.get(fullUrl, {
        responseType: 'blob',
      });
      
      // Create a blob URL and trigger download
      const blob = new Blob([response.data]);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Download error:', error);
      throw new Error('Failed to download file');
    }
  },
}; 