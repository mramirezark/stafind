'use client'

import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Alert,
  CircularProgress,
} from '@mui/material'
import { EmployeeFormDialogProps } from './interfaces'
import { EmployeeForm } from './EmployeeForm'

export function EmployeeFormDialog({
  open,
  onClose,
  onSubmit,
  employee,
  skills,
  skillsLoading,
  loading,
  error
}: EmployeeFormDialogProps) {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {employee ? 'Edit Employee' : 'Add New Employee'}
      </DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <EmployeeForm
          formData={{
            name: employee?.name || '',
            email: employee?.email || '',
            department: employee?.department || '',
            level: employee?.level || '',
            location: employee?.location || '',
            bio: employee?.bio || '',
            skills: employee?.skills?.map(skill => skill.name) || [],
            current_project: employee?.current_project || '',
          }}
          onFormDataChange={() => {}} // This will be handled by the parent
          skills={skills}
          skillsLoading={skillsLoading}
          error={error}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          onClick={() => onSubmit({
            name: employee?.name || '',
            email: employee?.email || '',
            department: employee?.department || '',
            level: employee?.level || '',
            location: employee?.location || '',
            bio: employee?.bio || '',
            skills: employee?.skills?.map(skill => skill.name) || [],
            current_project: employee?.current_project || '',
          })}
          variant="contained"
          disabled={loading}
          startIcon={loading ? <CircularProgress size={16} /> : undefined}
        >
          {loading ? 'Saving...' : employee ? 'Update' : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
