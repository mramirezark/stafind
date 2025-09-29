'use client'

import { Button as MuiButton, ButtonProps as MuiButtonProps } from '@mui/material'
import { forwardRef } from 'react'

export interface ButtonProps extends Omit<MuiButtonProps, 'variant'> {
  variant?: 'primary' | 'secondary' | 'outlined' | 'text' | 'contained'
  loading?: boolean
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ variant = 'contained', loading = false, children, disabled, ...props }, ref) => {
    return (
      <MuiButton
        ref={ref}
        variant={variant === 'primary' ? 'contained' : variant === 'secondary' ? 'outlined' : variant}
        disabled={disabled || loading}
        {...props}
      >
        {loading ? 'Loading...' : children}
      </MuiButton>
    )
  }
)

Button.displayName = 'Button'
