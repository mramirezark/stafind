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
  
  // Job Requests
  jobRequests: {
    list: () => api.get('/api/v1/job-requests'),
    get: (id: number) => api.get(`/api/v1/job-requests/${id}`),
    create: (data: any) => api.post('/api/v1/job-requests', data),
    getMatches: (id: number) => api.get(`/api/v1/job-requests/${id}/matches`),
  },
  
  // Skills
  skills: {
    list: () => api.get('/api/v1/skills'),
    create: (data: any) => api.post('/api/v1/skills', data),
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
