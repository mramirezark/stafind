'use client'

import { TextField, InputAdornment, IconButton } from '@mui/material'
import { Search as SearchIcon, Clear as ClearIcon } from '@mui/icons-material'
import { SearchBarProps } from './interfaces'

export function EmployeeSearchBar({ searchQuery, onSearchChange }: SearchBarProps) {
  return (
    <TextField
      fullWidth
      placeholder="Search employees..."
      value={searchQuery}
      onChange={(e) => onSearchChange(e.target.value)}
      InputProps={{
        startAdornment: (
          <InputAdornment position="start">
            <SearchIcon />
          </InputAdornment>
        ),
        endAdornment: searchQuery && (
          <InputAdornment position="end">
            <IconButton size="small" onClick={() => onSearchChange('')}>
              <ClearIcon />
            </IconButton>
          </InputAdornment>
        ),
      }}
    />
  )
}
