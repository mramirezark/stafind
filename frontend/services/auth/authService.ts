import { apiClient } from '../api/client'
import { API_ENDPOINTS } from '../api/endpoints'

// Auth types
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  user: User
  token: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  first_name: string
  last_name: string
  role_id?: number
}

export interface User {
  id: number
  username: string
  email: string
  first_name: string
  last_name: string
  role_id?: number
  role?: Role
  is_active: boolean
  last_login?: string
  created_at: string
  updated_at: string
  roles?: Role[]
}

export interface Role {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface UpdateProfileRequest {
  username?: string
  email?: string
  first_name?: string
  last_name?: string
}

/**
 * Authentication Service
 * Handles all authentication-related API calls
 */
export class AuthService {
  /**
   * Login user with username and password
   */
  static async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.LOGIN, credentials)
    return response.data
  }

  /**
   * Register a new user
   */
  static async register(userData: RegisterRequest): Promise<User> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.REGISTER, userData)
    return response.data
  }

  /**
   * Logout current user
   */
  static async logout(): Promise<{ message: string }> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.LOGOUT)
    return response.data
  }

  /**
   * Refresh authentication token
   */
  static async refreshToken(): Promise<{ token: string }> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.REFRESH)
    return response.data
  }

  /**
   * Get current user profile
   */
  static async getProfile(): Promise<User> {
    const response = await apiClient.get(API_ENDPOINTS.AUTH.PROFILE)
    return response.data
  }

  /**
   * Update user profile
   */
  static async updateProfile(profileData: UpdateProfileRequest): Promise<User> {
    const response = await apiClient.put(API_ENDPOINTS.AUTH.PROFILE, profileData)
    return response.data
  }

  /**
   * Change user password
   */
  static async changePassword(passwordData: ChangePasswordRequest): Promise<{ message: string }> {
    const response = await apiClient.put('/api/v1/auth/change-password', passwordData)
    return response.data
  }

  /**
   * Check if user has specific role
   */
  static hasRole(user: User, roleName: string): boolean {
    if (user.role?.name === roleName) return true
    if (user.roles?.some(role => role.name === roleName)) return true
    return false
  }

  /**
   * Check if user is admin
   */
  static isAdmin(user: User): boolean {
    return this.hasRole(user, 'admin')
  }

  /**
   * Check if user is HR manager
   */
  static isHRManager(user: User): boolean {
    return this.hasRole(user, 'hr_manager')
  }

  /**
   * Check if user is hiring manager
   */
  static isHiringManager(user: User): boolean {
    return this.hasRole(user, 'hiring_manager')
  }
}
