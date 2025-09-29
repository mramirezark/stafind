// Admin API Endpoints Configuration
export const ADMIN_ENDPOINTS = {
  // User management endpoints (admin only)
  USERS: {
    LIST: '/api/v1/admin/users',
    CREATE: '/api/v1/admin/users',
    GET: (id: number) => `/api/v1/admin/users/${id}`,
    UPDATE: (id: number) => `/api/v1/admin/users/${id}`,
    DELETE: (id: number) => `/api/v1/admin/users/${id}`,
  },

  // Role management endpoints (admin only)
  ROLES: {
    CREATE: '/api/v1/admin/roles',
    UPDATE: (id: number) => `/api/v1/admin/roles/${id}`,
    DELETE: (id: number) => `/api/v1/admin/roles/${id}`,
  },

  // API Key management endpoints (admin only)
  API_KEYS: {
    LIST: '/api/v1/admin/api-keys',
    CREATE: '/api/v1/admin/api-keys',
    GET: (id: number) => `/api/v1/admin/api-keys/${id}`,
    DEACTIVATE: (id: number) => `/api/v1/admin/api-keys/${id}/deactivate`,
    ROTATE: (id: number) => `/api/v1/admin/api-keys/${id}/rotate`,
  },
} as const

// Type for admin endpoint values
export type AdminEndpoint = typeof ADMIN_ENDPOINTS[keyof typeof ADMIN_ENDPOINTS]
export type AdminEndpointFunction = (id: number) => string
