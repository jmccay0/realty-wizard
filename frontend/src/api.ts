import axios from 'axios';
import type { Project, Property, Disclosure, ContractTerms, Deadline, Document, ProjectSummary } from './types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

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
