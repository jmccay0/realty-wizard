import axios from 'axios';
import type { Project, Property, Disclosure, ContractTerms, Deadline, Document, ProjectSummary } from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // If error is 401 and we haven't retried yet, try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) {
          // No refresh token, redirect to login
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');
          window.location.href = '/login';
          return Promise.reject(error);
        }

        // Try to refresh the token
        const response = await axios.post(`${API_BASE_URL}/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token } = response.data;
        localStorage.setItem('access_token', access_token);

        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

// Projects
export const listProjects = () => api.get<Project[]>('/projects');
export const getProject = (id: string) => api.get<Project>(`/projects/${id}`);
export const createProject = (data: Partial<Project>) => api.post<Project>('/projects', data);
export const updateProject = (id: string, data: Partial<Project>) => api.put<Project>(`/projects/${id}`, data);
export const getProjectSummary = (id: string) => api.get<ProjectSummary>(`/projects/${id}/summary`);

// Property
export const getProperty = (projectId: string) => api.get<Property>(`/projects/${projectId}/property`);
export const createProperty = (projectId: string, data: Partial<Property>) =>
  api.post<Property>(`/projects/${projectId}/property`, data);
export const updateProperty = (projectId: string, data: Partial<Property>) =>
  api.put<Property>(`/projects/${projectId}/property`, data);

// Disclosure
export const getDisclosure = (projectId: string) => api.get<Disclosure>(`/projects/${projectId}/disclosure`);
export const createDisclosure = (projectId: string, data: Partial<Disclosure>) =>
  api.post<Disclosure>(`/projects/${projectId}/disclosure`, data);
export const updateDisclosure = (projectId: string, data: Partial<Disclosure>) =>
  api.put<Disclosure>(`/projects/${projectId}/disclosure`, data);

// Contract
export const getContract = (projectId: string) => api.get<ContractTerms>(`/projects/${projectId}/contract`);
export const createContract = (projectId: string, data: Partial<ContractTerms>) =>
  api.post<ContractTerms>(`/projects/${projectId}/contract`, data);

// Deadlines
export const listDeadlines = (projectId: string) => api.get<Deadline[]>(`/projects/${projectId}/deadlines`);
export const updateDeadline = (id: string, data: Partial<Deadline>) =>
  api.put<Deadline>(`/deadlines/${id}`, data);

// Documents
export const listDocuments = (projectId: string) => api.get<Document[]>(`/projects/${projectId}/documents`);
export const generateDocument = (projectId: string, type: string) =>
  api.post<{ document: Document; pdf: number[] }>(`/projects/${projectId}/documents/generate`, { type });
export const downloadDocument = (documentId: string) =>
  api.get(`/documents/${documentId}/download`, { responseType: 'blob' });
