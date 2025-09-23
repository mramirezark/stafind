'use client'

import { useState } from 'react'
import {
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
} from '@mui/material'
import { Login as LoginIcon } from '@mui/icons-material'
import { useAuth } from '@/lib/auth'
import { LoginFormProps } from '@/types'

export function LoginForm({ onSuccess }: LoginFormProps) {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { login } = useAuth()

  const handleInputChange = (field: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({
      ...prev,
      [field]: event.target.value
    }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    
    if (!formData.username || !formData.password) {
      setError('Please fill in all fields')
      return
    }

    try {
      setLoading(true)
      setError(null)
      await login(formData.username, formData.password)
      onSuccess?.()
    } catch (err: any) {
      setError(err.message || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%', maxWidth: 400 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <TextField
        fullWidth
        label="Username"
        type="text"
        value={formData.username}
        onChange={handleInputChange('username')}
        margin="normal"
        required
        autoComplete="username"
        autoFocus
      />

      <TextField
        fullWidth
        label="Password"
        type="password"
        value={formData.password}
        onChange={handleInputChange('password')}
        margin="normal"
        required
        autoComplete="current-password"
      />

      <Button
        type="submit"
        fullWidth
        variant="contained"
        startIcon={loading ? <CircularProgress size={20} /> : <LoginIcon />}
        disabled={loading}
        sx={{ mt: 2, mb: 2 }}
      >
        {loading ? 'Signing In...' : 'Sign In'}
      </Button>

      <Box textAlign="center">
        <Typography variant="caption" color="text.secondary">
          Demo: username: admin, password: admin123
        </Typography>
      </Box>
    </Box>
  )
}
