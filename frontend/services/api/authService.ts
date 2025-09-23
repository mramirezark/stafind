/**
 * Authentication API Service
 * 
 * Handles all authentication-related API calls
 */

import { BaseApiService } from './baseService'
import { User, RegisterData } from '@/types'

export class AuthService extends BaseApiService {
  protected getDomainName(): string {
    return 'auth'
  }

  // ============================================================================
  // AUTHENTICATION OPERATIONS
  // ============================================================================

  /**
   * User login
   */
  async login(username: string, password: string): Promise<{ user: User; token: string }> {
    const result = await this.request('POST', '/api/v1/auth/login', { username, password }, false)
    // Clear all cache on login
    this.clearCache()
    return result
  }

  /**
   * User registration
   */
  async register(userData: RegisterData): Promise<{ user: User; token: string }> {
    const result = await this.request('POST', '/api/v1/auth/register', userData, false)
    // Clear all cache on registration
    this.clearCache()
    return result
  }

  /**
   * User logout
   */
  async logout(): Promise<void> {
    await this.request('POST', '/api/v1/auth/logout', undefined, false)
    // Clear all cache on logout
    this.clearCache()
  }

  /**
   * Refresh authentication token
   */
  async refreshToken(): Promise<{ token: string }> {
    return this.request('POST', '/api/v1/auth/refresh', undefined, false)
  }

  // ============================================================================
  // USER PROFILE OPERATIONS
  // ============================================================================

  /**
   * Get user profile
   */
  async getProfile(): Promise<User> {
    return this.request('GET', '/api/v1/auth/profile')
  }

  /**
   * Update user profile
   */
  async updateProfile(userData: Partial<User>): Promise<User> {
    const result = await this.request('PUT', '/api/v1/auth/profile', userData, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Change user password
   */
  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    await this.request('POST', '/api/v1/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword
    }, false)
  }

  // ============================================================================
  // PASSWORD RESET OPERATIONS
  // ============================================================================

  /**
   * Request password reset
   */
  async requestPasswordReset(email: string): Promise<void> {
    await this.request('POST', '/api/v1/auth/forgot-password', { email }, false)
  }

  /**
   * Reset password with token
   */
  async resetPassword(token: string, newPassword: string): Promise<void> {
    await this.request('POST', '/api/v1/auth/reset-password', {
      token,
      new_password: newPassword
    }, false)
  }

  // ============================================================================
  // EMAIL VERIFICATION OPERATIONS
  // ============================================================================

  /**
   * Send email verification
   */
  async sendEmailVerification(): Promise<void> {
    await this.request('POST', '/api/v1/auth/send-verification', undefined, false)
  }

  /**
   * Verify email with token
   */
  async verifyEmail(token: string): Promise<void> {
    await this.request('POST', '/api/v1/auth/verify-email', { token }, false)
  }
}

// Export singleton instance
export const authService = new AuthService()
