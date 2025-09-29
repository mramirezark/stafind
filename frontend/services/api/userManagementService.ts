import { API_ENDPOINTS } from './endpoints'
import { ADMIN_ENDPOINTS } from './adminEndpoints'
import { apiClient } from './client'
import { User, Role, RegisterData, ApiResponse, PaginatedResponse, CreateUserRequest, UpdateUserRequest } from '@/types'

export class UserManagementService {
  /**
   * Get all users with pagination
   */
  async getUsers(limit: number = 10, offset: number = 0): Promise<PaginatedResponse<User>> {
    try {
      const response = await apiClient.get(ADMIN_ENDPOINTS.USERS.LIST, {
        params: { limit, offset }
      })
      
      // Handle different response structures
      if (response.data && Array.isArray(response.data.data) && response.data.pagination) {
        // Standard paginated response
        return response.data
      } else if (response.data && Array.isArray(response.data)) {
        // Direct array response - create pagination wrapper
        return {
          data: response.data,
          pagination: {
            page: Math.floor(offset / limit) + 1,
            limit,
            total: response.data.length,
            totalPages: Math.ceil(response.data.length / limit)
          },
          success: true
        }
      } else if (response.data && response.data.users && Array.isArray(response.data.users)) {
        // Users in nested property
        return {
          data: response.data.users,
          pagination: response.data.pagination || {
            page: Math.floor(offset / limit) + 1,
            limit,
            total: response.data.users.length,
            totalPages: Math.ceil(response.data.users.length / limit)
          },
          success: true
        }
      } else {
        console.warn('Unexpected users response structure:', response.data)
        return {
          data: [],
          pagination: {
            page: 1,
            limit,
            total: 0,
            totalPages: 0
          },
          success: false
        }
      }
    } catch (error) {
      console.error('Failed to fetch users:', error)
      // Return empty response on error
      return {
        data: [],
        pagination: {
          page: 1,
          limit,
          total: 0,
          totalPages: 0
        },
        success: false
      }
    }
  }

  /**
   * Get a specific user by ID
   */
  async getUser(id: number): Promise<User> {
    try {
      const response = await apiClient.get(ADMIN_ENDPOINTS.USERS.GET(id))
      return response.data
    } catch (error) {
      console.error('Failed to fetch user:', error)
      throw error
    }
  }

  /**
   * Create a new user
   */
  async createUser(data: CreateUserRequest): Promise<User> {
    try {
      const response = await apiClient.post(ADMIN_ENDPOINTS.USERS.CREATE, data)
      return response.data
    } catch (error) {
      console.error('Failed to create user:', error)
      throw error
    }
  }

  /**
   * Update an existing user
   */
  async updateUser(id: number, data: UpdateUserRequest): Promise<User> {
    try {
      const response = await apiClient.put(ADMIN_ENDPOINTS.USERS.UPDATE(id), data)
      return response.data
    } catch (error) {
      console.error('Failed to update user:', error)
      throw error
    }
  }

  /**
   * Delete a user
   */
  async deleteUser(id: number): Promise<void> {
    try {
      await apiClient.delete(ADMIN_ENDPOINTS.USERS.DELETE(id))
    } catch (error) {
      console.error('Failed to delete user:', error)
      throw error
    }
  }

  /**
   * Get all roles
   */
  async getRoles(): Promise<Role[]> {
    try {
      const response = await apiClient.get(API_ENDPOINTS.ROLES.LIST)
      // Handle different response structures
      if (Array.isArray(response.data)) {
        return response.data
      } else if (response.data && Array.isArray(response.data.data)) {
        return response.data.data
      } else if (response.data && Array.isArray(response.data.roles)) {
        return response.data.roles
      } else {
        console.warn('Unexpected roles response structure:', response.data)
        return []
      }
    } catch (error) {
      console.error('Failed to fetch roles:', error)
      // Return empty array on error
      return []
    }
  }

  /**
   * Create a new role
   */
  async createRole(data: { name: string; description: string }): Promise<Role> {
    const response = await apiClient.post(ADMIN_ENDPOINTS.ROLES.CREATE, data)
    return response.data
  }

  /**
   * Update an existing role
   */
  async updateRole(id: number, data: { name: string; description: string }): Promise<Role> {
    const response = await apiClient.put(ADMIN_ENDPOINTS.ROLES.UPDATE(id), data)
    return response.data
  }

  /**
   * Delete a role
   */
  async deleteRole(id: number): Promise<void> {
    await apiClient.delete(ADMIN_ENDPOINTS.ROLES.DELETE(id))
  }
}

export const userManagementService = new UserManagementService()
