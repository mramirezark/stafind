import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('authToken')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API endpoints
export const endpoints = {
  // Authentication
  auth: {
    login: (data: { username: string; password: string }) => api.post('/api/v1/auth/login', data),
    register: (data: any) => api.post('/api/v1/auth/register', data),
    logout: () => api.post('/api/v1/auth/logout'),
    refresh: () => api.post('/api/v1/auth/refresh'),
    getProfile: () => api.get('/api/v1/auth/profile'),
    updateProfile: (data: any) => api.put('/api/v1/auth/profile', data),
    changePassword: (data: { current_password: string; new_password: string }) => 
      api.post('/api/v1/auth/change-password', data),
  },

  // Employees
  employees: {
    list: () => api.get('/api/v1/employees'),
    get: (id: number) => api.get(`/api/v1/employees/${id}`),
    create: (data: any) => api.post('/api/v1/employees', data),
    update: (id: number, data: any) => api.put(`/api/v1/employees/${id}`, data),
    delete: (id: number) => api.delete(`/api/v1/employees/${id}`),
  },
  
  
  // Skills
  skills: {
    list: () => api.get('/api/v1/skills'),
    get: (id: number) => api.get(`/api/v1/skills/${id}`),
    create: (data: any) => api.post('/api/v1/skills', data),
    update: (id: number, data: any) => api.put(`/api/v1/skills/${id}`, data),
    delete: (id: number) => api.delete(`/api/v1/skills/${id}`),
    search: (query: string) => api.get(`/api/v1/skills/search?q=${encodeURIComponent(query)}`),
    popular: (limit?: number) => api.get(`/api/v1/skills/popular${limit ? `?limit=${limit}` : ''}`),
    withCount: () => api.get('/api/v1/skills/with-count'),
    stats: () => api.get('/api/v1/skills/stats'),
    byCategory: (category: string) => api.get(`/api/v1/skills/category/${encodeURIComponent(category)}`),
    byEmployee: (employeeId: number) => api.get(`/api/v1/skills/employee/${employeeId}`),
    batch: {
      create: (data: any[]) => api.post('/api/v1/skills/batch', data),
      update: (data: any[]) => api.put('/api/v1/skills/batch', data),
      delete: (ids: number[]) => api.delete('/api/v1/skills/batch', { data: { ids } }),
    },
  },
  
  // Categories
  categories: {
    list: () => api.get('/api/v1/skills/categories'),
    create: (data: any) => api.post('/api/v1/skills/categories', data),
    update: (id: number, data: any) => api.put(`/api/v1/skills/categories/${id}`, data),
    delete: (id: number) => api.delete(`/api/v1/skills/categories/${id}`),
  },
  
  // Search
  search: {
    employees: (data: any) => api.post('/api/v1/search', data),
  },

  // Roles
  roles: {
    list: () => api.get('/api/v1/roles'),
    get: (id: number) => api.get(`/api/v1/roles/${id}`),
  },

  // AI Agent
  aiAgent: {
    getRequests: (params?: { limit?: number; offset?: number }) => 
      api.get('/api/v1/ai-agent/requests', { params }),
    getRequest: (id: number) => api.get(`/api/v1/ai-agent/requests/${id}`),
    getResponse: (id: number) => api.get(`/api/v1/ai-agent/responses/${id}`),
    processRequest: (id: number) => api.post(`/api/v1/ai-agent/requests/${id}/process`),
    process: (data: any) => api.post('/api/v1/ai-agent/process', data),
    extractSkills: (data: any) => api.post('/api/v1/ai-agent/extract-skills', data),
  },

  // Admin endpoints
  admin: {
    users: {
      list: (params?: { page?: number; limit?: number }) => 
        api.get('/api/v1/admin/users', { params }),
      get: (id: number) => api.get(`/api/v1/admin/users/${id}`),
      update: (id: number, data: any) => api.put(`/api/v1/admin/users/${id}`, data),
      delete: (id: number) => api.delete(`/api/v1/admin/users/${id}`),
    },
    roles: {
      create: (data: any) => api.post('/api/v1/admin/roles', data),
      update: (id: number, data: any) => api.put(`/api/v1/admin/roles/${id}`, data),
      delete: (id: number) => api.delete(`/api/v1/admin/roles/${id}`),
    },
  },
}

export default api
