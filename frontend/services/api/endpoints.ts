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

  // User management endpoints (admin only)
  USERS: {
    LIST: '/api/v1/users',
    CREATE: '/api/v1/users',
    GET: (id: number) => `/api/v1/users/${id}`,
    UPDATE: (id: number) => `/api/v1/users/${id}`,
    DELETE: (id: number) => `/api/v1/users/${id}`,
  },

  // Role management endpoints (admin only)
  ROLES: {
    LIST: '/api/v1/roles',
    CREATE: '/api/v1/roles',
    GET: (id: number) => `/api/v1/roles/${id}`,
    UPDATE: (id: number) => `/api/v1/roles/${id}`,
    DELETE: (id: number) => `/api/v1/roles/${id}`,
  },

  // Health check
  HEALTH: '/health',
} as const

// Type for endpoint values
export type ApiEndpoint = typeof API_ENDPOINTS[keyof typeof API_ENDPOINTS]
export type ApiEndpointFunction = (id: number) => string
