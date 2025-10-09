'use client'

import { useState, useMemo, useEffect } from 'react'
import {
  Box,
  Card,
  CardContent,
  Button,
  Typography,
  Alert,
  CircularProgress,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material'
import { Add as AddIcon } from '@mui/icons-material'
import { useEmployees, useSkills, useCreateEmployee, useUpdateEmployee, useDeleteEmployee } from '@/hooks/useApi'
import { useAuth } from '@/hooks/useAuth'
import { SearchFilters, EmployeeFormData } from './interfaces'
import { Employee } from '@/types'
import { EmployeeSearchBar } from './EmployeeSearchBar'
import { EmployeeFilters } from './EmployeeFilters'
import { ViewControls } from './ViewControls'
import { EmployeeList } from './EmployeeList'
import { EmployeeFormDialog } from './EmployeeFormDialog'

export function EmployeeManagement() {
  // Hooks
  const { data: employees, loading, error, refetch } = useEmployees()
  const { createEmployee } = useCreateEmployee()
  const { updateEmployee } = useUpdateEmployee()
  const { deleteEmployee } = useDeleteEmployee()
  const { data: skills, loading: skillsLoading } = useSkills()
  const { isAdmin, isHRManager } = useAuth()

  // Debug: Log when employees data changes
  useEffect(() => {
    if (employees) {
      console.log('Employees data changed:', {
        count: employees.length,
        employees: employees.map(emp => ({ id: emp.id, name: emp.name }))
      })
    }
  }, [employees])

  // State
  const [filters, setFilters] = useState<SearchFilters>({
    searchQuery: '',
    departmentFilter: '',
    levelFilter: '',
    skillFilter: [],
  })
  const [viewMode, setViewMode] = useState<'grid' | 'list' | 'table'>('table')
  const [formOpen, setFormOpen] = useState(false)
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null)
  const [formData, setFormData] = useState<EmployeeFormData>({
    name: '',
    email: '',
    department: '',
    level: '',
    location: '',
    bio: '',
    skills: [],
    current_project: '',
  })
  const [formError, setFormError] = useState<string | null>(null)
  const [formLoading, setFormLoading] = useState(false)
  const [getLatestLoading, setGetLatestLoading] = useState(false)
  const [getAllLoading, setGetAllLoading] = useState(false)
  const [getLatestSuccess, setGetLatestSuccess] = useState(false)
  const [getAllSuccess, setGetAllSuccess] = useState(false)
  const [webhookError, setWebhookError] = useState<string | null>(null)
  
  // Delete confirmation dialog state
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [employeeToDelete, setEmployeeToDelete] = useState<Employee | null>(null)
  const [deleteLoading, setDeleteLoading] = useState(false)

  // Permissions
  const canManageEmployees = isAdmin() || isHRManager()

  // Filtered employees
  const filteredEmployees = useMemo(() => {
    if (!employees) return []

    return employees.filter(employee => {
      const matchesSearch = !filters.searchQuery || 
        employee.name.toLowerCase().includes(filters.searchQuery.toLowerCase()) ||
        employee.email.toLowerCase().includes(filters.searchQuery.toLowerCase()) ||
        employee.bio?.toLowerCase().includes(filters.searchQuery.toLowerCase())

      const matchesDepartment = !filters.departmentFilter || employee.department === filters.departmentFilter
      const matchesLevel = !filters.levelFilter || employee.level === filters.levelFilter
      const matchesSkills = filters.skillFilter.length === 0 || 
        filters.skillFilter.some((skill: string) => employee.skills?.some(empSkill => empSkill.name === skill))

      return matchesSearch && matchesDepartment && matchesLevel && matchesSkills
    })
  }, [employees, filters])

  // Handlers
  const handleFiltersChange = (newFilters: Partial<SearchFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }))
  }

  const handleClearFilters = () => {
    setFilters({
      searchQuery: '',
      departmentFilter: '',
      levelFilter: '',
      skillFilter: [],
    })
  }

  const handleViewModeChange = (mode: 'grid' | 'list' | 'table') => {
    setViewMode(mode)
  }

  const handleAddEmployee = () => {
    setEditingEmployee(null)
    setFormData({
      name: '',
      email: '',
      department: '',
      level: '',
      location: '',
      bio: '',
      skills: [],
      current_project: '',
    })
    setFormError(null)
    setFormOpen(true)
  }


  const handleEditEmployee = (employee: Employee) => {
    setEditingEmployee(employee)
    setFormData({
      name: employee.name,
      email: employee.email,
      department: employee.department,
      level: employee.level,
      location: employee.location,
      bio: employee.bio || '',
      skills: employee.skills?.map(skill => skill.name) || [],
      current_project: employee.current_project || '',
    })
    setFormError(null)
    setFormOpen(true)
  }

  const handleDeleteEmployee = (id: number) => {
    const employee = employees?.find(emp => emp.id === id)
    if (employee) {
      setEmployeeToDelete(employee)
      setDeleteDialogOpen(true)
    }
  }

  const confirmDeleteEmployee = async () => {
    if (!employeeToDelete) return

    setDeleteLoading(true)
    try {
      await deleteEmployee(employeeToDelete.id)
      setDeleteDialogOpen(false)
      setEmployeeToDelete(null)
      // Refresh the employee list after successful deletion
      refetch()
    } catch (error) {
      console.error('Failed to delete employee:', error)
      // You could add a toast notification here for better UX
    } finally {
      setDeleteLoading(false)
    }
  }

  const cancelDeleteEmployee = () => {
    setDeleteDialogOpen(false)
    setEmployeeToDelete(null)
  }

  const handleFormSubmit = async (data: EmployeeFormData) => {
    if (!data.name || !data.email || !data.department || !data.level) {
      setFormError('Please fill in all required fields')
      return
    }

    try {
      setFormLoading(true)
      setFormError(null)

      // Convert skills from string array to EmployeeSkill format
      const employeeData: Partial<Employee> = {
        name: data.name,
        email: data.email,
        department: data.department,
        level: data.level,
        location: data.location,
        bio: data.bio,
        current_project: data.current_project,
        skills: data.skills.map(skillName => {
          const skill = skills?.find(s => s.name === skillName)
          return {
            id: skill?.id || 0,
            name: skillName,
            category: skill?.categories?.[0]?.name || 'Other',
            proficiency_level: 3,
            years_experience: 1.0
          }
        })
      }

      if (editingEmployee) {
        await updateEmployee(editingEmployee.id, employeeData)
      } else {
        await createEmployee(employeeData)
      }

      setFormOpen(false)
      setEditingEmployee(null)
    } catch (err: any) {
      setFormError(err.message || 'Failed to save employee')
    } finally {
      setFormLoading(false)
    }
  }

  const handleFormClose = () => {
    setFormOpen(false)
    setEditingEmployee(null)
    setFormError(null)
  }

  const callWebhook = async (includeDate = true, buttonType: 'latest' | 'all' = 'latest') => {
    // Set loading state based on button type
    if (buttonType === 'latest') {
      setGetLatestLoading(true)
      setGetLatestSuccess(false)
    } else {
      setGetAllLoading(true)
      setGetAllSuccess(false)
    }
    setWebhookError(null)
    
    try {
      const webhookUrl = process.env.NEXT_PUBLIC_WEBHOOK_URL
      if (!webhookUrl) {
        throw new Error('Webhook URL not configured')
      }
      
      // Build request body based on whether to include date filter
      const requestBody: any = {
        source: 'employee-management'
      }
      
      if (includeDate) {
        let date = new Date()
        date.setDate(date.getDate() - 1)
        date.setHours(0, 0, 0, 0) // Set to beginning of the day (00:00:00)
        requestBody.timestamp = date.toISOString()
      }
      
      const response = await fetch(webhookUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody)
      })
      
      if (!response.ok) {
        throw new Error(`Webhook failed with status: ${response.status}`)
      }
      
      // Set success state based on button type
      if (buttonType === 'latest') {
        setGetLatestSuccess(true)
        setTimeout(() => setGetLatestSuccess(false), 3000)
      } else {
        setGetAllSuccess(true)
        setTimeout(() => setGetAllSuccess(false), 3000)
      }
    } catch (err) {
      setWebhookError(err instanceof Error ? err.message : 'Failed to call webhook')
    } finally {
      // Clear loading state based on button type
      if (buttonType === 'latest') {
        setGetLatestLoading(false)
      } else {
        setGetAllLoading(false)
      }
    }
  }

  // Loading state
  if (loading && !employees) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  // Error state
  if (error) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
        <Button onClick={refetch} variant="outlined">
          Retry
        </Button>
      </Box>
    )
  }

  return (
    <Box>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h4" gutterBottom>
            Employee Management
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {filteredEmployees.length} of {employees?.length || 0} employees
          </Typography>
        </Box>
        {canManageEmployees && (
          <Box display="flex" gap={2}>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleAddEmployee}
              size="large"
            >
              Add Employee
            </Button>
            <Button
              variant="contained"
              color={getLatestSuccess ? "success" : "primary"}
              onClick={() => callWebhook(true, 'latest')}
              disabled={getLatestLoading}
              size="large"
            >
              {getLatestLoading ? 'Calling...' : getLatestSuccess ? 'Get Latest Called!' : 'Get Latest'}
            </Button>
            <Button
              variant="contained"
              color={getAllSuccess ? "success" : "primary"}
              onClick={() => callWebhook(false, 'all')}
              disabled={getAllLoading}
              size="large"
            >
              {getAllLoading ? 'Calling...' : getAllSuccess ? 'Get All Called!' : 'Get All'}
            </Button>
          </Box>
        )}
      </Box>

      {/* Webhook Error Alert */}
      {webhookError && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setWebhookError(null)}>
          {webhookError}
        </Alert>
      )}

      {/* Search and Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            {/* Search */}
            <Grid item xs={12} md={4}>
              <EmployeeSearchBar
                searchQuery={filters.searchQuery}
                onSearchChange={(query) => handleFiltersChange({ searchQuery: query })}
              />
            </Grid>

            {/* Filters */}
            <Grid item xs={12} md={7}>
              <EmployeeFilters
                filters={filters}
                onFiltersChange={handleFiltersChange}
                onClearFilters={handleClearFilters}
                skills={skills || []}
                skillsLoading={skillsLoading}
              />
            </Grid>

            {/* View Controls */}
            <Grid item xs={12} md={1}>
              <Box display="flex" gap={1}>
                <ViewControls
                  viewMode={viewMode}
                  onViewModeChange={handleViewModeChange}
                  onClearFilters={handleClearFilters}
                />
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Employee List */}
      <EmployeeList
        employees={filteredEmployees}
        viewMode={viewMode}
        onEdit={canManageEmployees ? handleEditEmployee : undefined}
        onDelete={canManageEmployees ? handleDeleteEmployee : undefined}
      />

      {/* Employee Form Dialog */}
      <EmployeeFormDialog
        open={formOpen}
        onClose={handleFormClose}
        onSubmit={handleFormSubmit}
        employee={editingEmployee}
        skills={skills || []}
        skillsLoading={skillsLoading}
        loading={formLoading}
        error={formError}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={cancelDeleteEmployee}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Delete Employee
        </DialogTitle>
        <DialogContent>
          <Typography variant="body1">
            Are you sure you want to delete{' '}
            <strong>{employeeToDelete?.name}</strong>?
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            This action cannot be undone. All employee data will be permanently removed.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={cancelDeleteEmployee}
            disabled={deleteLoading}
          >
            Cancel
          </Button>
          <Button
            onClick={confirmDeleteEmployee}
            color="error"
            variant="contained"
            disabled={deleteLoading}
            startIcon={deleteLoading ? <CircularProgress size={20} /> : null}
          >
            {deleteLoading ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>

    </Box>
  )
}