// API Endpoints Configuration
export const API_ENDPOINTS = {
  // Authentication endpoints
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    REGISTER: '/api/v1/auth/register',
    LOGOUT: '/api/v1/auth/logout',
    REFRESH: '/api/v1/auth/refresh',
    PROFILE: '/api/v1/auth/profile',
  },

  // Employee endpoints
  EMPLOYEES: {
    LIST: '/api/v1/employees',
    CREATE: '/api/v1/employees',
    GET: (id: number) => `/api/v1/employees/${id}`,
    UPDATE: (id: number) => `/api/v1/employees/${id}`,
    DELETE: (id: number) => `/api/v1/employees/${id}`,
    SEARCH: '/api/v1/search',
  },

  // Job Request endpoints
  JOB_REQUESTS: {
    LIST: '/api/v1/job-requests',
    CREATE: '/api/v1/job-requests',
    GET: (id: number) => `/api/v1/job-requests/${id}`,
    UPDATE: (id: number) => `/api/v1/job-requests/${id}`,
    DELETE: (id: number) => `/api/v1/job-requests/${id}`,
  },

  // Skills endpoints
  SKILLS: {
    LIST: '/api/v1/skills',
    CREATE: '/api/v1/skills',
    GET: (id: number) => `/api/v1/skills/${id}`,
    UPDATE: (id: number) => `/api/v1/skills/${id}`,
    DELETE: (id: number) => `/api/v1/skills/${id}`,
  },

  // Role management endpoints (public/protected)
  ROLES: {
    LIST: '/api/v1/roles',
    GET: (id: number) => `/api/v1/roles/${id}`,
  },

  // API Key validation (public)
  API_KEYS: {
    VALIDATE: '/api-keys/validate',
  },

  // Health check
  HEALTH: '/health',
} as const

// Type for endpoint values
export type ApiEndpoint = typeof API_ENDPOINTS[keyof typeof API_ENDPOINTS]
export type ApiEndpointFunction = (id: number) => string
