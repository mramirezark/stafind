'use client'

import { useState } from 'react'
import {
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material'
import { PersonAdd as PersonAddIcon } from '@mui/icons-material'
import { useAuth } from '@/lib/auth'
import { RegisterFormProps } from '@/types'

export function RegisterForm({ onSuccess, onToggleMode }: RegisterFormProps) {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    first_name: '',
    last_name: '',
    role_id: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { register } = useAuth()

  const handleInputChange = (field: string) => (event: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: event.target.value
    }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()

    // Validation
    if (!formData.username || !formData.email || !formData.password || !formData.first_name || !formData.last_name) {
      setError('Please fill in all required fields')
      return
    }

    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters long')
      return
    }

    try {
      setLoading(true)
      setError(null)

      const registerData = {
        username: formData.username,
        email: formData.email,
        password: formData.password,
        first_name: formData.first_name,
        last_name: formData.last_name,
        role_id: formData.role_id ? parseInt(formData.role_id) : undefined,
      }

      await register(registerData)
      onSuccess?.()
    } catch (err: any) {
      setError(err.message || 'Registration failed')
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
        value={formData.username}
        onChange={handleInputChange('username')}
        margin="normal"
        required
        autoComplete="username"
        helperText="Choose a unique username for login"
      />

      <TextField
        fullWidth
        label="Email"
        type="email"
        value={formData.email}
        onChange={handleInputChange('email')}
        margin="normal"
        required
        autoComplete="email"
      />

      <TextField
        fullWidth
        label="First Name"
        value={formData.first_name}
        onChange={handleInputChange('first_name')}
        margin="normal"
        required
        autoComplete="given-name"
      />

      <TextField
        fullWidth
        label="Last Name"
        value={formData.last_name}
        onChange={handleInputChange('last_name')}
        margin="normal"
        required
        autoComplete="family-name"
      />

      <TextField
        fullWidth
        label="Password"
        type="password"
        value={formData.password}
        onChange={handleInputChange('password')}
        margin="normal"
        required
        autoComplete="new-password"
      />

      <TextField
        fullWidth
        label="Confirm Password"
        type="password"
        value={formData.confirmPassword}
        onChange={handleInputChange('confirmPassword')}
        margin="normal"
        required
        autoComplete="new-password"
      />

      <FormControl fullWidth margin="normal">
        <InputLabel>Role (Optional)</InputLabel>
        <Select
          value={formData.role_id}
          onChange={handleInputChange('role_id')}
          label="Role (Optional)"
        >
          <MenuItem value="">No specific role</MenuItem>
          <MenuItem value="2">HR Manager</MenuItem>
          <MenuItem value="3">Hiring Manager</MenuItem>
        </Select>
      </FormControl>

      <Button
        type="submit"
        fullWidth
        variant="contained"
        startIcon={loading ? <CircularProgress size={20} /> : <PersonAddIcon />}
        disabled={loading}
        sx={{ mt: 2, mb: 2 }}
      >
        {loading ? 'Creating Account...' : 'Create Account'}
      </Button>

      {onToggleMode && (
        <Box textAlign="center">
          <Typography variant="body2" color="text.secondary">
            Already have an account?{' '}
            <Button variant="text" onClick={onToggleMode} sx={{ p: 0, minWidth: 'auto' }}>
              Sign In
            </Button>
          </Typography>
        </Box>
      )}
    </Box>
  )
}
