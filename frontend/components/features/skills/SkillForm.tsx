'use client'

import React, { useState, useEffect } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Box,
  Typography,
  Alert,
  CircularProgress,
  Autocomplete,
  Stack,
} from '@mui/material'
import { LoadingButton } from '@mui/lab'

import { Skill, Category, SkillFormData } from '@/types'
import { endpoints } from '@/lib/api'

interface SkillFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: SkillFormData) => Promise<void>
  initialData?: Skill | null
  categories: Category[]
  title?: string
}

const SkillForm: React.FC<SkillFormProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
  categories,
  title = 'Skill Form',
}) => {
  const [formData, setFormData] = useState<SkillFormData>({
    name: '',
    categories: [],
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedCategories, setSelectedCategories] = useState<Category[]>([])

  useEffect(() => {
    if (initialData) {
      setFormData({
        name: initialData.name,
        categories: initialData.categories?.map(cat => cat.id) || [],
      })
      setSelectedCategories(initialData.categories || [])
    } else {
      setFormData({
        name: '',
        categories: [],
      })
      setSelectedCategories([])
    }
    setError(null)
  }, [initialData, open])

  const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({ ...prev, name: event.target.value }))
  }

  const handleCategoryChange = (event: any, newValue: Category[]) => {
    setSelectedCategories(newValue)
    setFormData(prev => ({
      ...prev,
      categories: newValue.map(cat => cat.id),
    }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    
    if (!formData.name.trim()) {
      setError('Skill name is required')
      return
    }

    try {
      setLoading(true)
      setError(null)
      await onSubmit(formData)
    } catch (err: any) {
      setError(err.message || 'Failed to save skill')
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    if (!loading) {
      onClose()
    }
  }

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        component: 'form',
        onSubmit: handleSubmit,
      }}
    >
      <DialogTitle>
        <Typography variant="h6" component="div">
          {title}
        </Typography>
      </DialogTitle>

      <DialogContent>
        <Stack spacing={3} sx={{ mt: 1 }}>
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          <TextField
            autoFocus
            required
            fullWidth
            label="Skill Name"
            value={formData.name}
            onChange={handleNameChange}
            placeholder="e.g., React, Python, Machine Learning"
            disabled={loading}
            error={!formData.name.trim() && formData.name !== ''}
            helperText={
              !formData.name.trim() && formData.name !== ''
                ? 'Skill name is required'
                : 'Enter a descriptive name for the skill'
            }
          />

          <Autocomplete
            multiple
            options={categories}
            getOptionLabel={(option) => option.name}
            value={selectedCategories}
            onChange={handleCategoryChange}
            disabled={loading}
            renderTags={(value, getTagProps) =>
              value.map((option, index) => (
                <Chip
                  variant="outlined"
                  label={option.name}
                  {...getTagProps({ index })}
                  key={option.id}
                />
              ))
            }
            renderInput={(params) => (
              <TextField
                {...params}
                label="Categories"
                placeholder="Select or search categories"
                helperText="Choose relevant categories for this skill"
              />
            )}
            renderOption={(props, option) => (
              <Box component="li" {...props}>
                <Box>
                  <Typography variant="body2">{option.name}</Typography>
                  {option.description && (
                    <Typography variant="caption" color="text.secondary">
                      {option.description}
                    </Typography>
                  )}
                </Box>
              </Box>
            )}
          />

          {selectedCategories.length > 0 && (
            <Box>
              <Typography variant="body2" color="text.secondary" gutterBottom>
                Selected Categories:
              </Typography>
              <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                {selectedCategories.map(category => (
                  <Chip
                    key={category.id}
                    label={category.name}
                    size="small"
                    color="primary"
                    variant="outlined"
                  />
                ))}
              </Box>
            </Box>
          )}
        </Stack>
      </DialogContent>

      <DialogActions>
        <Button
          onClick={handleClose}
          disabled={loading}
        >
          Cancel
        </Button>
        <LoadingButton
          type="submit"
          variant="contained"
          loading={loading}
          disabled={!formData.name.trim()}
        >
          {initialData ? 'Update Skill' : 'Create Skill'}
        </LoadingButton>
      </DialogActions>
    </Dialog>
  )
}

export default SkillForm
