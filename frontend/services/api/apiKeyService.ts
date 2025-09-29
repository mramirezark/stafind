import { API_ENDPOINTS } from './endpoints'
import { ADMIN_ENDPOINTS } from './adminEndpoints'
import { apiClient } from './client'
import { APIKey, CreateAPIKeyRequest, UpdateAPIKeyRequest, ApiResponse, PaginatedResponse } from '@/types'

export class APIKeyService {
  /**
   * Get all API keys with pagination
   */
  async getAPIKeys(limit: number = 10, offset: number = 0): Promise<PaginatedResponse<APIKey>> {
    try {
      const response = await apiClient.get(ADMIN_ENDPOINTS.API_KEYS.LIST, {
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
      } else if (response.data && response.data.api_keys && Array.isArray(response.data.api_keys)) {
        // API keys in nested property
        return {
          data: response.data.api_keys,
          pagination: response.data.pagination || {
            page: Math.floor(offset / limit) + 1,
            limit,
            total: response.data.api_keys.length,
            totalPages: Math.ceil(response.data.api_keys.length / limit)
          },
          success: true
        }
      } else {
        console.warn('Unexpected API keys response structure:', response.data)
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
      console.error('Failed to fetch API keys:', error)
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
   * Get a specific API key by ID
   */
  async getAPIKey(id: number): Promise<APIKey> {
    try {
      const response = await apiClient.get(ADMIN_ENDPOINTS.API_KEYS.GET(id))
      return response.data
    } catch (error) {
      console.error('Failed to fetch API key:', error)
      throw error
    }
  }

  /**
   * Create a new API key
   */
  async createAPIKey(data: CreateAPIKeyRequest): Promise<APIKey> {
    try {
      const response = await apiClient.post(ADMIN_ENDPOINTS.API_KEYS.CREATE, data)
      return response.data
    } catch (error) {
      console.error('Failed to create API key:', error)
      throw error
    }
  }


  /**
   * Deactivate an API key
   */
  async deactivateAPIKey(id: number): Promise<void> {
    try {
      await apiClient.post(ADMIN_ENDPOINTS.API_KEYS.DEACTIVATE(id))
    } catch (error) {
      console.error('Failed to deactivate API key:', error)
      throw error
    }
  }

  /**
   * Rotate an API key (create new, deactivate old)
   */
  async rotateAPIKey(id: number): Promise<APIKey> {
    try {
      const response = await apiClient.post(ADMIN_ENDPOINTS.API_KEYS.ROTATE(id))
      return response.data
    } catch (error) {
      console.error('Failed to rotate API key:', error)
      throw error
    }
  }

  /**
   * Validate an API key
   */
  async validateAPIKey(apiKey: string): Promise<{ valid: boolean; service_name?: string; description?: string; is_active?: boolean }> {
    try {
      const response = await apiClient.post(API_ENDPOINTS.API_KEYS.VALIDATE, { api_key: apiKey })
      return response.data
    } catch (error) {
      console.error('Failed to validate API key:', error)
      throw error
    }
  }
}

export const apiKeyService = new APIKeyService()
