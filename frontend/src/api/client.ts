import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { toast } from 'react-hot-toast';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000, // 30 seconds
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    // Add request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add auth token if available
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Add response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<{ message?: string }>) => {
        const errorMessage = error.response?.data?.message || 'An error occurred';
        toast.error(errorMessage);
        return Promise.reject(error);
      }
    );
  }

  // Health check
  async healthCheck(): Promise<{ status: string }> {
    const response = await this.client.get('/health');
    return response.data;
  }

  // Get supported formats
  async getSupportedFormats(): Promise<{
    [key: string]: {
      input: string[];
      output: string[];
    };
  }> {
    const response = await this.client.get('/formats');
    return response.data.data;
  }

  // Convert file
  async convertFile(
    file: File,
    targetFormat: string,
    options: Record<string, any> = {}
  ): Promise<{ downloadUrl: string }> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('format', targetFormat);

    // Add additional options
    Object.entries(options).forEach(([key, value]) => {
      formData.append(key, value);
    });

    const response = await this.client.post('/convert', formData, {
      responseType: 'blob',
    });

    // Create a download URL for the converted file
    const url = window.URL.createObjectURL(new Blob([response.data]));
    return { downloadUrl: url };
  }

  // Generic request method
  async request<T = any>(config: AxiosRequestConfig): Promise<AxiosResponse<T>> {
    return this.client.request<T>(config);
  }
}

export const apiClient = new ApiClient();
export default apiClient;
