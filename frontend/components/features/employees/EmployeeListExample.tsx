'use client'

import { useState } from 'react'
import { useEmployees } from '@/hooks/useEmployees'
import { useAuth } from '@/hooks/useAuth'
import { EmployeeCard, EmployeeFormDialog } from './'
import { Button } from '@mui/material'

/**
 * Example of a modularized EmployeeList component
 * This shows how to use the custom hooks and services
 */
export const EmployeeListExample = () => {
  const {
    employees,
    loading,
    error,
    createEmployee,
    updateEmployee,
    deleteEmployee,
    refresh,
  } = useEmployees()

  const { isAdmin, isHRManager } = useAuth()

  const [formOpen, setFormOpen] = useState(false)
  const [editingEmployee, setEditingEmployee] = useState<any>(null)

  const handleCreateEmployee = async (data: any) => {
    try {
      await createEmployee(data)
      setFormOpen(false)
    } catch (error) {
      console.error('Failed to create employee:', error)
    }
  }

  const handleUpdateEmployee = async (data: any) => {
    try {
      await updateEmployee(editingEmployee.id, data)
      setFormOpen(false)
      setEditingEmployee(null)
    } catch (error) {
      console.error('Failed to update employee:', error)
    }
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

  const handleEditEmployee = (employee: any) => {
    setEditingEmployee(employee)
    setFormOpen(true)
  }

  const canManageEmployees = isAdmin() || isHRManager()

  if (loading && employees.length === 0) {
    return <div>Loading employees...</div>
  }

  if (error) {
    return (
      <div>
        <p>Error: {error}</p>
        <Button onClick={refresh}>Retry</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Employees ({employees.length})</h1>
        {canManageEmployees && (
          <Button onClick={() => setFormOpen(true)}>
            Add Employee
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {employees.map((employee) => (
          <EmployeeCard
            key={employee.id}
            employee={employee}
            onEdit={canManageEmployees ? handleEditEmployee : undefined}
            onDelete={canManageEmployees ? handleDeleteEmployee : undefined}
            viewMode="grid"
          />
        ))}
      </div>

      {formOpen && (
        <EmployeeFormDialog
          open={formOpen}
          onClose={() => {
            setFormOpen(false)
            setEditingEmployee(null)
          }}
          employee={editingEmployee}
          onSubmit={editingEmployee ? handleUpdateEmployee : handleCreateEmployee}
          skills={[]}
          skillsLoading={false}
          loading={false}
          error={null}
        />
      )}
    </div>
  )
}

/**
 * Benefits of this modularized approach:
 * 
 * 1. **Separation of Concerns**: 
 *    - useEmployees hook handles all employee data logic
 *    - useAuth hook handles authentication logic
 *    - Component only handles UI logic
 * 
 * 2. **Reusability**: 
 *    - useEmployees can be used in other components
 *    - EmployeeCard can be reused anywhere
 *    - Services can be used in different contexts
 * 
 * 3. **Testability**: 
 *    - Hooks can be tested independently
 *    - Services can be mocked easily
 *    - Components are easier to test
 * 
 * 4. **Type Safety**: 
 *    - Full TypeScript support throughout
 *    - Compile-time error checking
 *    - Better IDE support and autocomplete
 * 
 * 5. **Maintainability**: 
 *    - Changes to API structure only affect services
 *    - UI changes don't affect business logic
 *    - Easy to add new features
 */
