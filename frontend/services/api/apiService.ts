/**
 * Centralized API Service
 * 
 * This is the single point of truth for all API calls in the application.
 * All components should use this service instead of making direct API calls.
 */

import { api } from '@/lib/api'
import { 
  Employee, 
  Skill, 
  User, 
  RegisterData, 
  SearchCriteria,
  DashboardStats,
  ApiResponse,
  PaginatedResponse 
} from '@/types'

// ============================================================================
// API SERVICE CLASS
// ============================================================================

class ApiService {
  private cache = new Map<string, { data: any; timestamp: number }>()
  private readonly CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  /**
   * Generic method to handle API calls with caching
   */
  private async request<T>(
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
  private handleError(error: any): Error {
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
      const keysToDelete: string[] = []
      this.cache.forEach((_, key) => {
        if (key.includes(pattern)) {
          keysToDelete.push(key)
        }
      })
      keysToDelete.forEach(key => this.cache.delete(key))
    } else {
      this.cache.clear()
    }
  }

  // ============================================================================
  // AUTHENTICATION METHODS
  // ============================================================================

  async login(username: string, password: string): Promise<{ user: User; token: string }> {
    return this.request('POST', '/api/v1/auth/login', { username, password }, false)
  }

  async register(userData: RegisterData): Promise<{ user: User; token: string }> {
    return this.request('POST', '/api/v1/auth/register', userData, false)
  }

  async logout(): Promise<void> {
    await this.request('POST', '/api/v1/auth/logout', undefined, false)
    this.clearCache()
  }

  async getProfile(): Promise<User> {
    return this.request('GET', '/api/v1/auth/profile')
  }

  async updateProfile(userData: Partial<User>): Promise<User> {
    const result = await this.request('PUT', '/api/v1/auth/profile', userData, false)
    this.clearCache('auth')
    return result as User
  }

  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    await this.request('POST', '/api/v1/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword
    }, false)
  }

  // ============================================================================
  // EMPLOYEE METHODS
  // ============================================================================

  async getEmployees(): Promise<Employee[]> {
    return this.request('GET', '/api/v1/employees')
  }

  async getEmployee(id: number): Promise<Employee> {
    return this.request('GET', `/api/v1/employees/${id}`)
  }

  async createEmployee(employeeData: Partial<Employee>): Promise<Employee> {
    const result = await this.request('POST', '/api/v1/employees', employeeData, false)
    this.clearCache('employees')
    return result as Employee
  }

  async updateEmployee(id: number, employeeData: Partial<Employee>): Promise<Employee> {
    const result = await this.request('PUT', `/api/v1/employees/${id}`, employeeData, false)
    this.clearCache('employees')
    return result as Employee
  }

  async deleteEmployee(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/employees/${id}`, undefined, false)
    this.clearCache('employees')
  }



  // ============================================================================
  // SKILL METHODS
  // ============================================================================

  async getSkills(): Promise<Skill[]> {
    return this.request('GET', '/api/v1/skills')
  }

  async createSkill(skillData: Partial<Skill>): Promise<Skill> {
    const result = await this.request('POST', '/api/v1/skills', skillData, false)
    this.clearCache('skills')
    return result as Skill
  }

  async updateSkill(id: number, skillData: Partial<Skill>): Promise<Skill> {
    const result = await this.request('PUT', `/api/v1/skills/${id}`, skillData, false)
    this.clearCache('skills')
    return result as Skill
  }

  async deleteSkill(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/skills/${id}`, undefined, false)
    this.clearCache('skills')
  }

  // ============================================================================
  // SEARCH METHODS
  // ============================================================================

  async searchEmployees(searchCriteria: SearchCriteria): Promise<Employee[]> {
    return this.request('POST', '/api/v1/search', searchCriteria, false)
  }

  // ============================================================================
  // DASHBOARD METHODS
  // ============================================================================

  async getDashboardStats(): Promise<DashboardStats> {
    return this.request('GET', '/api/v1/dashboard/stats')
  }

  async getDashboardData(): Promise<{
    stats: DashboardStats
    employees: Employee[]
  }> {
    // Fetch all dashboard data in parallel
    const [stats, employees] = await Promise.all([
      this.getDashboardStats(),
      this.getEmployees()
    ])

    return {
      stats,
      employees
    }
  }

  // ============================================================================
  // ROLES METHODS
  // ============================================================================

  async getRoles(): Promise<any[]> {
    return this.request('GET', '/api/v1/roles')
  }

  // ============================================================================
  // ADMIN METHODS
  // ============================================================================

  async getUsers(params?: { page?: number; limit?: number }): Promise<PaginatedResponse<User>> {
    const queryParams = params ? `?${new URLSearchParams(params as any).toString()}` : ''
    return this.request('GET', `/api/v1/admin/users${queryParams}`)
  }

  async updateUser(id: number, userData: Partial<User>): Promise<User> {
    const result = await this.request('PUT', `/api/v1/admin/users/${id}`, userData, false)
    this.clearCache('users')
    return result as User
  }

  async deleteUser(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/admin/users/${id}`, undefined, false)
    this.clearCache('users')
  }
}

// ============================================================================
// SINGLETON INSTANCE
// ============================================================================

// Export a singleton instance
export const apiService = new ApiService()

// Export the class for testing purposes
export { ApiService }
