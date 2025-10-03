'use client'

import React, { useState } from 'react'
import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Tooltip,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
  Stack,
  Chip,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material'
import { LoadingButton } from '@mui/lab'

import { Category, CategoryFormData } from '@/types'
import { endpoints } from '@/lib/api'

interface CategoryManagementProps {
  categories: Category[]
  onRefresh: () => void
}

const CategoryManagement: React.FC<CategoryManagementProps> = ({
  categories,
  onRefresh,
}) => {
  const [formOpen, setFormOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [formData, setFormData] = useState<CategoryFormData>({
    name: '',
    description: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [deleteConfirm, setDeleteConfirm] = useState<Category | null>(null)

  const handleCreate = () => {
    setEditingCategory(null)
    setFormData({ name: '', description: '' })
    setFormOpen(true)
    setError(null)
  }

  const handleEdit = (category: Category) => {
    setEditingCategory(category)
    setFormData({
      name: category.name,
      description: category.description || '',
    })
    setFormOpen(true)
    setError(null)
  }

  const handleFormClose = () => {
    setFormOpen(false)
    setEditingCategory(null)
    setFormData({ name: '', description: '' })
    setError(null)
  }

  const handleFormSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    
    if (!formData.name.trim()) {
      setError('Category name is required')
      return
    }

    try {
      setLoading(true)
      setError(null)

      if (editingCategory) {
        await endpoints.categories.update(editingCategory.id, formData)
      } else {
        await endpoints.categories.create(formData)
      }

      await onRefresh()
      handleFormClose()
    } catch (err: any) {
      setError(err.message || 'Failed to save category')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (category: Category) => {
    try {
      setLoading(true)
      await endpoints.categories.delete(category.id)
      await onRefresh()
      setDeleteConfirm(null)
    } catch (err: any) {
      setError(err.message || 'Failed to delete category')
    } finally {
      setLoading(false)
    }
  }

  const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({ ...prev, name: event.target.value }))
  }

  const handleDescriptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({ ...prev, description: event.target.value }))
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          Skill Categories ({categories.length})
        </Typography>
        <Stack direction="row" spacing={1}>
          <Tooltip title="Refresh">
            <IconButton onClick={onRefresh} disabled={loading}>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={handleCreate}
          >
            Add Category
          </Button>
        </Stack>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Categories Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell align="right">Created</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={4} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : categories.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} align="center" sx={{ py: 4 }}>
                    <Typography variant="body2" color="text.secondary">
                      No categories found. Create your first category to get started.
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                categories.map(category => (
                  <TableRow key={category.id} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {category.name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="text.secondary">
                        {category.description || 'No description'}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Typography variant="caption" color="text.secondary">
                        {category.created_at
                          ? new Date(category.created_at).toLocaleDateString()
                          : '-'}
                      </Typography>
                    </TableCell>
                    <TableCell align="center">
                      <Stack direction="row" spacing={1} justifyContent="center">
                        <Tooltip title="Edit">
                          <IconButton
                            size="small"
                            onClick={() => handleEdit(category)}
                          >
                            <EditIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Delete">
                          <IconButton
                            size="small"
                            color="error"
                            onClick={() => setDeleteConfirm(category)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Tooltip>
                      </Stack>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Category Form Dialog */}
      <Dialog
        open={formOpen}
        onClose={handleFormClose}
        maxWidth="sm"
        fullWidth
        PaperProps={{
          component: 'form',
          onSubmit: handleFormSubmit,
        }}
      >
        <DialogTitle>
          <Typography variant="h6">
            {editingCategory ? 'Edit Category' : 'Create New Category'}
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
              label="Category Name"
              value={formData.name}
              onChange={handleNameChange}
              placeholder="e.g., Programming Languages, Frameworks, Tools"
              disabled={loading}
              error={!formData.name.trim() && formData.name !== ''}
              helperText={
                !formData.name.trim() && formData.name !== ''
                  ? 'Category name is required'
                  : 'Enter a descriptive name for the category'
              }
            />

            <TextField
              fullWidth
              multiline
              rows={3}
              label="Description"
              value={formData.description}
              onChange={handleDescriptionChange}
              placeholder="Optional description of what skills belong in this category"
              disabled={loading}
            />
          </Stack>
        </DialogContent>

        <DialogActions>
          <Button
            onClick={handleFormClose}
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
            {editingCategory ? 'Update Category' : 'Create Category'}
          </LoadingButton>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      {deleteConfirm && (
        <Dialog
          open={Boolean(deleteConfirm)}
          onClose={() => setDeleteConfirm(null)}
          maxWidth="sm"
          fullWidth
        >
          <DialogTitle>
            <Typography variant="h6" color="error">
              Delete Category
            </Typography>
          </DialogTitle>

          <DialogContent>
            <Typography>
              Are you sure you want to delete the category "{deleteConfirm.name}"? 
              This action cannot be undone and may affect skills associated with this category.
            </Typography>
          </DialogContent>

          <DialogActions>
            <Button
              onClick={() => setDeleteConfirm(null)}
              disabled={loading}
            >
              Cancel
            </Button>
            <LoadingButton
              onClick={() => handleDelete(deleteConfirm)}
              color="error"
              variant="contained"
              loading={loading}
            >
              Delete
            </LoadingButton>
          </DialogActions>
        </Dialog>
      )}
    </Box>
  )
}

export default CategoryManagement
