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

// Service Marketplace Types
export interface Service {
  id: string;
  name: string;
  category: string;
  description: string;
  typical_timeline?: string;
  estimated_cost_min?: number;
  estimated_cost_max?: number;
  created_at: string;
  updated_at: string;
}

export interface Provider {
  id: string;
  service_id: string;
  business_name: string;
  contact_name?: string;
  email: string;
  phone: string;
  address?: string;
  city?: string;
  state?: string;
  zip?: string;
  bio?: string;
  years_experience?: number;
  license_number?: string;
  insurance_verified: boolean;
  availability_status: string;
  rating_average: number;
  rating_count: number;
  created_at: string;
  updated_at: string;
}

export interface ServiceRequest {
  id: string;
  project_id?: string;
  user_email: string;
  user_name: string;
  user_phone?: string;
  service_id: string;
  provider_id?: string;
  property_address?: string;
  requested_date?: string;
  preferred_time?: string;
  status: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface ProviderReview {
  id: string;
  provider_id: string;
  service_request_id: string;
  user_email: string;
  rating: number;
  review_text?: string;
  created_at: string;
}

// Service Marketplace API
export const listServices = (category?: string) => {
  const url = category ? `/services?category=${category}` : '/services';
  return api.get<Service[]>(url);
};

export const getService = (id: string) => api.get<Service>(`/services/${id}`);

export const listProviders = (serviceId: string) =>
  api.get<Provider[]>(`/services/${serviceId}/providers`);

export const getProvider = (id: string) => api.get<Provider>(`/providers/${id}`);

export const listProviderReviews = (providerId: string) =>
  api.get<ProviderReview[]>(`/providers/${providerId}/reviews`);

export const createProviderReview = (providerId: string, data: Partial<ProviderReview>) =>
  api.post<ProviderReview>(`/providers/${providerId}/reviews`, data);

export const createServiceRequest = (data: Partial<ServiceRequest>) =>
  api.post<ServiceRequest>('/service-requests', data);

export const listServiceRequests = (userEmail: string) =>
  api.get<ServiceRequest[]>(`/service-requests?user_email=${userEmail}`);

export const getServiceRequest = (id: string) =>
  api.get<ServiceRequest>(`/service-requests/${id}`);

export const updateServiceRequest = (id: string, data: Partial<ServiceRequest>) =>
  api.put<ServiceRequest>(`/service-requests/${id}`, data);
