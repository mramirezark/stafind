'use client'

import { useState } from 'react'
import { Box, CircularProgress, Typography } from '@mui/material'
import { useAuth } from '@/lib/auth'
import { LoginForm, RegisterForm } from '@/components/auth'
import { AuthWrapperProps } from '@/types'

export function AuthWrapper({ children }: AuthWrapperProps) {
  const { isAuthenticated, isLoading } = useAuth()
  const [isLoginMode, setIsLoginMode] = useState(true)

  // Show loading spinner while checking authentication
  if (isLoading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="100vh"
        sx={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
      >
        <CircularProgress size={60} sx={{ color: 'white' }} />
      </Box>
    )
  }

  // Show login/register forms if not authenticated
  if (!isAuthenticated) {
    return isLoginMode ? (
      <LoginForm />
    ) : (
      <RegisterForm />
    )
  }

  // Show main application if authenticated
  return <>{children}</>
}
