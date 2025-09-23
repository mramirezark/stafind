import { useState, useEffect, useCallback } from 'react'
import { AuthService, LoginRequest, RegisterRequest, User, ChangePasswordRequest, UpdateProfileRequest } from '@/services/auth/authService'

interface UseAuthReturn {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (userData: RegisterRequest) => Promise<void>
  logout: () => void
  refreshToken: () => Promise<void>
  updateProfile: (profileData: UpdateProfileRequest) => Promise<void>
  changePassword: (passwordData: ChangePasswordRequest) => Promise<void>
  hasRole: (roleName: string) => boolean
  isAdmin: () => boolean
  isHRManager: () => boolean
  isHiringManager: () => boolean
}

/**
 * Custom hook for authentication management
 */
export const useAuth = (): UseAuthReturn => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  /**
   * Initialize auth state from localStorage
   */
  useEffect(() => {
    const initializeAuth = () => {
      try {
        const storedToken = localStorage.getItem('authToken')
        const storedUser = localStorage.getItem('authUser')

        if (storedToken && storedUser) {
          setToken(storedToken)
          setUser(JSON.parse(storedUser))
        }
      } catch (error) {
        console.error('Failed to initialize auth:', error)
        // Clear invalid data
        localStorage.removeItem('authToken')
        localStorage.removeItem('authUser')
      } finally {
        setIsLoading(false)
      }
    }

    initializeAuth()
  }, [])

  /**
   * Login user
   */
  const login = useCallback(async (credentials: LoginRequest): Promise<void> => {
    try {
      setIsLoading(true)
      const response = await AuthService.login(credentials)
      
      setUser(response.user)
      setToken(response.token)
      
      // Store in localStorage
      localStorage.setItem('authToken', response.token)
      localStorage.setItem('authUser', JSON.stringify(response.user))
    } catch (error: any) {
      throw new Error(error.response?.data?.error || error.message || 'Login failed')
    } finally {
      setIsLoading(false)
    }
  }, [])

  /**
   * Register new user
   */
  const register = useCallback(async (userData: RegisterRequest): Promise<void> => {
    try {
      setIsLoading(true)
      const newUser = await AuthService.register(userData)
      
      // Auto-login after successful registration
      await login({
        username: userData.username,
        password: userData.password,
      })
    } catch (error: any) {
      throw new Error(error.response?.data?.error || error.message || 'Registration failed')
    } finally {
      setIsLoading(false)
    }
  }, [login])

  /**
   * Logout user
   */
  const logout = useCallback(() => {
    try {
      // Call logout API (optional, for server-side session cleanup)
      AuthService.logout().catch(console.error)
    } catch (error) {
      console.error('Logout API call failed:', error)
    } finally {
      // Clear local state
      setUser(null)
      setToken(null)
      
      // Clear localStorage
      localStorage.removeItem('authToken')
      localStorage.removeItem('authUser')
    }
  }, [])

  /**
   * Refresh authentication token
   */
  const refreshToken = useCallback(async (): Promise<void> => {
    try {
      const response = await AuthService.refreshToken()
      setToken(response.token)
      localStorage.setItem('authToken', response.token)
    } catch (error: any) {
      console.error('Token refresh failed:', error)
      // If refresh fails, logout user
      logout()
      throw new Error('Session expired. Please login again.')
    }
  }, [logout])

  /**
   * Update user profile
   */
  const updateProfile = useCallback(async (profileData: UpdateProfileRequest): Promise<void> => {
    try {
      const updatedUser = await AuthService.updateProfile(profileData)
      setUser(updatedUser)
      localStorage.setItem('authUser', JSON.stringify(updatedUser))
    } catch (error: any) {
      throw new Error(error.response?.data?.error || error.message || 'Failed to update profile')
    }
  }, [])

  /**
   * Change user password
   */
  const changePassword = useCallback(async (passwordData: ChangePasswordRequest): Promise<void> => {
    try {
      await AuthService.changePassword(passwordData)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || error.message || 'Failed to change password')
    }
  }, [])

  /**
   * Check if user has specific role
   */
  const hasRole = useCallback((roleName: string): boolean => {
    if (!user) return false
    return AuthService.hasRole(user, roleName)
  }, [user])

  /**
   * Check if user is admin
   */
  const isAdmin = useCallback((): boolean => {
    if (!user) return false
    return AuthService.isAdmin(user)
  }, [user])

  /**
   * Check if user is HR manager
   */
  const isHRManager = useCallback((): boolean => {
    if (!user) return false
    return AuthService.isHRManager(user)
  }, [user])

  /**
   * Check if user is hiring manager
   */
  const isHiringManager = useCallback((): boolean => {
    if (!user) return false
    return AuthService.isHiringManager(user)
  }, [user])

  return {
    user,
    token,
    isAuthenticated: !!user && !!token,
    isLoading,
    login,
    register,
    logout,
    refreshToken,
    updateProfile,
    changePassword,
    hasRole,
    isAdmin,
    isHRManager,
    isHiringManager,
  }
}
