'use client'

import { TextField, TextFieldProps } from '@mui/material'
import { forwardRef } from 'react'

export interface FormFieldProps extends Omit<TextFieldProps, 'variant'> {
  variant?: 'outlined' | 'filled' | 'standard'
  error?: boolean
  helperText?: string
}

export const FormField = forwardRef<HTMLInputElement, FormFieldProps>(
  ({ variant = 'outlined', error = false, helperText, ...props }, ref) => {
    return (
      <TextField
        ref={ref}
        variant={variant}
        error={error}
        helperText={helperText}
        fullWidth
        {...props}
      />
    )
  }
)

FormField.displayName = 'FormField'
