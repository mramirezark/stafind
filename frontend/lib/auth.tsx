'use client'

import React, { createContext, useContext, useEffect, useState } from 'react'
import { api, endpoints } from './api'
import { User, AuthContextType, RegisterData } from '@/types'

interface LoginResponse {
  user: User
  token: string
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // Check if user is authenticated
  const isAuthenticated = !!user && !!token

  // Initialize auth state from localStorage
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const storedToken = localStorage.getItem('authToken')
        const storedUser = localStorage.getItem('authUser')

        if (storedToken && storedUser) {
          setToken(storedToken)
          setUser(JSON.parse(storedUser))
          
          // Set authorization header
          api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`
          
          // Verify token is still valid by fetching profile
          try {
            const response = await endpoints.auth.getProfile()
            setUser(response.data)
            localStorage.setItem('authUser', JSON.stringify(response.data))
          } catch (error) {
            // Token is invalid, clear auth state
            clearAuthState()
          }
        }
      } catch (error) {
        console.error('Auth initialization error:', error)
        clearAuthState()
      } finally {
        setIsLoading(false)
      }
    }

    initializeAuth()
  }, [])

  const clearAuthState = () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem('authToken')
    localStorage.removeItem('authUser')
    delete api.defaults.headers.common['Authorization']
  }

  const login = async (username: string, password: string) => {
    try {
      const response = await endpoints.auth.login({ username, password })
      const { user: userData, token: authToken } = response.data

      setUser(userData)
      setToken(authToken)
      
      localStorage.setItem('authToken', authToken)
      localStorage.setItem('authUser', JSON.stringify(userData))
      
      // Set authorization header
      api.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Login failed')
    }
  }

  const register = async (userData: RegisterData) => {
    try {
      const response = await endpoints.auth.register(userData)
      // After successful registration, automatically log in
      await login(userData.email, userData.password)
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Registration failed')
    }
  }

  const logout = async () => {
    try {
      if (token) {
        await endpoints.auth.logout()
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuthState()
    }
  }

  const updateProfile = async (userData: Partial<User>) => {
    try {
      const response = await endpoints.auth.updateProfile(userData)
      setUser(response.data)
      localStorage.setItem('authUser', JSON.stringify(response.data))
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Profile update failed')
    }
  }

  const changePassword = async (currentPassword: string, newPassword: string) => {
    try {
      await endpoints.auth.changePassword({ current_password: currentPassword, new_password: newPassword })
    } catch (error: any) {
      throw new Error(error.response?.data?.error || 'Password change failed')
    }
  }

  const hasRole = (roleName: string): boolean => {
    if (!user) return false
    
    // Check primary role
    if (user.role?.name === roleName) return true
    
    // Check additional roles
    if (user.roles) {
      return user.roles.some(role => role.name === roleName)
    }
    
    return false
  }

  const isAdmin = (): boolean => hasRole('admin')
  const isHRManager = (): boolean => hasRole('hr_manager')
  const isHiringManager = (): boolean => hasRole('hiring_manager')

  const value: AuthContextType = {
    user,
    token,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    updateProfile,
    changePassword,
    hasRole,
    isAdmin,
    isHRManager,
    isHiringManager,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

// Export types for use in components
export type { User, RegisterData, LoginResponse }
