'use client'

import { useState, useMemo } from 'react'
import {
  Box,
  Card,
  CardContent,
  Button,
  Typography,
  Alert,
  CircularProgress,
  Grid,
} from '@mui/material'
import { Add as AddIcon, CloudUpload as UploadIcon } from '@mui/icons-material'
import { useEmployees, useSkills, useCreateEmployee, useUpdateEmployee, useDeleteEmployee } from '@/hooks/useApi'
import { useAuth } from '@/hooks/useAuth'
import { SearchFilters, EmployeeFormData } from './interfaces'
import { Employee } from '@/types'
import { EmployeeSearchBar } from './EmployeeSearchBar'
import { EmployeeFilters } from './EmployeeFilters'
import { ViewControls } from './ViewControls'
import { EmployeeList } from './EmployeeList'
import { EmployeeFormDialog } from './EmployeeFormDialog'
import { ResumeUpload } from './ResumeUpload'

export function EmployeeManagement() {
  // Hooks
  const { data: employees, loading, error, refetch } = useEmployees()
  const { createEmployee } = useCreateEmployee()
  const { updateEmployee } = useUpdateEmployee()
  const { deleteEmployee } = useDeleteEmployee()
  const { data: skills, loading: skillsLoading } = useSkills()
  const { isAdmin, isHRManager } = useAuth()

  // State
  const [filters, setFilters] = useState<SearchFilters>({
    searchQuery: '',
    departmentFilter: '',
    levelFilter: '',
    skillFilter: [],
  })
  const [viewMode, setViewMode] = useState<'grid' | 'list' | 'table'>('grid')
  const [formOpen, setFormOpen] = useState(false)
  const [resumeUploadOpen, setResumeUploadOpen] = useState(false)
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null)
  const [formData, setFormData] = useState<EmployeeFormData>({
    name: '',
    email: '',
    department: '',
    level: '',
    location: '',
    bio: '',
    skills: [],
  })
  const [formError, setFormError] = useState<string | null>(null)
  const [formLoading, setFormLoading] = useState(false)

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
    })
    setFormError(null)
    setFormOpen(true)
  }

  const handleBulkImport = async (employees: any[]) => {
    try {
      // Convert parsed data to employee format
      const employeeData = employees.map(emp => ({
        name: emp.name,
        email: emp.email,
        department: 'Engineering', // Default department
        level: emp.experience || 'Mid-level',
        location: emp.location || '',
        bio: emp.bio || '',
        skills: emp.skills.map((skill: string) => ({
          name: skill,
          category: 'Technical',
          proficiency_level: 3,
          years_experience: 1,
        }))
      }))

      // Create employees in bulk
      for (const emp of employeeData) {
        await createEmployee(emp)
      }

      // Refresh the employee list
      refetch()
    } catch (error) {
      console.error('Bulk import failed:', error)
      throw error
    }
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
    })
    setFormError(null)
    setFormOpen(true)
  }

  const handleDeleteEmployee = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this employee?')) {
      try {
        await deleteEmployee(id)
      } catch (error) {
        console.error('Failed to delete employee:', error)
      }
    }
  }

  const handleFormSubmit = async (data: EmployeeFormData) => {
    if (!data.name || !data.email || !data.department || !data.level) {
      setFormError('Please fill in all required fields')
      return
    }

    try {
      setFormLoading(true)
      setFormError(null)

      // Convert skills from string array to EmployeeSkill array
      const employeeData = {
        ...data,
        skills: data.skills.map(skillName => {
          const skill = skills?.find(s => s.name === skillName)
          return {
            id: skill?.id || 0,
            name: skillName,
            category: skill?.category || 'Other',
            proficiency_level: 3, // Default proficiency level
            years_experience: 1, // Default years of experience
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
              variant="outlined"
              startIcon={<UploadIcon />}
              onClick={() => setResumeUploadOpen(true)}
              size="large"
            >
              Bulk Import
            </Button>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleAddEmployee}
              size="large"
            >
              Add Employee
            </Button>
          </Box>
        )}
      </Box>

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

      {/* Resume Upload Dialog */}
      <ResumeUpload
        open={resumeUploadOpen}
        onClose={() => setResumeUploadOpen(false)}
        onBulkImport={handleBulkImport}
      />
    </Box>
  )
}