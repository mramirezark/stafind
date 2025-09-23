/**
 * Base API Service Class
 * 
 * Contains common functionality shared across all domain services:
 * - HTTP request handling
 * - Caching
 * - Error handling
 * - Authentication
 */

import { api } from '@/lib/api'

// ============================================================================
// BASE SERVICE CLASS
// ============================================================================

export abstract class BaseApiService {
  protected cache = new Map<string, { data: any; timestamp: number }>()
  protected readonly CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  /**
   * Generic method to handle API calls with caching
   */
  protected async request<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    url: string,
    data?: any,
    useCache: boolean = true
  ): Promise<T> {
    const cacheKey = `${method}:${url}:${JSON.stringify(data || {})}`
    
    // Check cache for GET requests
    if (method === 'GET' && useCache) {
      const cached = this.cache.get(cacheKey)
      if (cached && Date.now() - cached.timestamp < this.CACHE_DURATION) {
        return cached.data
      }
    }

    try {
      let response
      switch (method) {
        case 'GET':
          response = await api.get(url)
          break
        case 'POST':
          response = await api.post(url, data)
          break
        case 'PUT':
          response = await api.put(url, data)
          break
        case 'DELETE':
          response = await api.delete(url)
          break
      }

      // Cache successful GET requests
      if (method === 'GET' && useCache) {
        this.cache.set(cacheKey, {
          data: response.data,
          timestamp: Date.now()
        })
      }

      return response.data
    } catch (error: any) {
      console.error(`API Error [${method} ${url}]:`, error)
      throw this.handleError(error)
    }
  }

  /**
   * Handle API errors consistently
   */
  protected handleError(error: any): Error {
    if (error.response) {
      // Server responded with error status
      const message = error.response.data?.message || error.response.data?.error || 'Server error'
      const status = error.response.status
      return new Error(`${status}: ${message}`)
    } else if (error.request) {
      // Request was made but no response received
      return new Error('Network error: Unable to connect to server')
    } else {
      // Something else happened
      return new Error(error.message || 'An unexpected error occurred')
    }
  }

  /**
   * Clear cache for specific pattern or all cache
   */
  public clearCache(pattern?: string): void {
    if (pattern) {
      for (const key of Array.from(this.cache.keys())) {
        if (key.includes(pattern)) {
          this.cache.delete(key)
        }
      }
    } else {
      this.cache.clear()
    }
  }

  /**
   * Clear cache for this service's domain
   */
  protected clearDomainCache(): void {
    this.clearCache(this.getDomainName())
  }

  /**
   * Get domain name for cache clearing
   * Override in subclasses to provide domain-specific cache clearing
   */
  protected abstract getDomainName(): string
}
