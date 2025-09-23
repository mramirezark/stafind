/**
 * Search API Service
 * 
 * Handles all search-related API calls
 */

import { BaseApiService } from './baseService'
import { Employee, SearchCriteria } from '@/types'

export class SearchService extends BaseApiService {
  protected getDomainName(): string {
    return 'search'
  }

  // ============================================================================
  // EMPLOYEE SEARCH OPERATIONS
  // ============================================================================

  /**
   * Search employees by criteria
   */
  async searchEmployees(searchCriteria: SearchCriteria): Promise<Employee[]> {
    return this.request('POST', '/api/v1/search', searchCriteria, false)
  }

  /**
   * Search employees by skills
   */
  async searchEmployeesBySkills(skills: string[], operator: 'AND' | 'OR' = 'AND'): Promise<Employee[]> {
    return this.request('POST', '/api/v1/search/skills', {
      skills,
      operator
    }, false)
  }

  /**
   * Search employees by department
   */
  async searchEmployeesByDepartment(department: string): Promise<Employee[]> {
    return this.request('GET', `/api/v1/search/department/${encodeURIComponent(department)}`)
  }

  /**
   * Search employees by location
   */
  async searchEmployeesByLocation(location: string): Promise<Employee[]> {
    return this.request('GET', `/api/v1/search/location/${encodeURIComponent(location)}`)
  }

  /**
   * Search employees by experience level
   */
  async searchEmployeesByExperience(level: string): Promise<Employee[]> {
    return this.request('GET', `/api/v1/search/experience/${encodeURIComponent(level)}`)
  }

  // ============================================================================
  // ADVANCED SEARCH OPERATIONS
  // ============================================================================

  /**
   * Advanced search with multiple criteria
   */
  async advancedSearch(criteria: {
    skills?: string[]
    department?: string
    location?: string
    experience_level?: string
    availability?: boolean
    min_experience?: number
    max_experience?: number
  }): Promise<Employee[]> {
    return this.request('POST', '/api/v1/search/advanced', criteria, false)
  }

  /**
   * Search with filters
   */
  async searchWithFilters(filters: {
    query?: string
    skills?: string[]
    department?: string
    location?: string
    experience_level?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
    page?: number
    limit?: number
  }): Promise<{
    employees: Employee[]
    total: number
    page: number
    limit: number
  }> {
    return this.request('POST', '/api/v1/search/filters', filters, false)
  }

  // ============================================================================
  // SEARCH SUGGESTIONS
  // ============================================================================

  /**
   * Get search suggestions
   */
  async getSearchSuggestions(query: string): Promise<string[]> {
    return this.request('GET', `/api/v1/search/suggestions?q=${encodeURIComponent(query)}`)
  }

  /**
   * Get skill suggestions
   */
  async getSkillSuggestions(query: string): Promise<string[]> {
    return this.request('GET', `/api/v1/search/skill-suggestions?q=${encodeURIComponent(query)}`)
  }

  /**
   * Get department suggestions
   */
  async getDepartmentSuggestions(query: string): Promise<string[]> {
    return this.request('GET', `/api/v1/search/department-suggestions?q=${encodeURIComponent(query)}`)
  }

  // ============================================================================
  // SEARCH ANALYTICS
  // ============================================================================

  /**
   * Get search history
   */
  async getSearchHistory(): Promise<any[]> {
    return this.request('GET', '/api/v1/search/history')
  }

  /**
   * Save search query
   */
  async saveSearchQuery(query: string, results_count: number): Promise<void> {
    await this.request('POST', '/api/v1/search/save', {
      query,
      results_count
    }, false)
  }

  /**
   * Get popular searches
   */
  async getPopularSearches(limit: number = 10): Promise<any[]> {
    return this.request('GET', `/api/v1/search/popular?limit=${limit}`)
  }
}

// Export singleton instance
export const searchService = new SearchService()
