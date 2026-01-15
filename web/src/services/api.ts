import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Project {
  id: string;
  name: string;
  created_at: string;
}

export interface Deployment {
  id: string;
  project_id: string;
  status: 'queued' | 'building' | 'running' | 'failed' | 'stopped';
  container_id: string;
  port: number;
  image_tag: string;
  created_at: string;
  updated_at: string;
}

export const projectApi = {
  create: async (name: string): Promise<Project> => {
    const response = await api.post<Project>('/projects', { name });
    return response.data;
  },
};

export const deploymentApi = {
  create: async (projectId: string, code: string): Promise<Deployment> => {
    const response = await api.post<Deployment>('/deployments', {
      project_id: projectId,
      code,
    });
    return response.data;
  },
};

export default api;
