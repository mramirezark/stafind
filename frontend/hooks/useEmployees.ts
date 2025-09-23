import { useState, useEffect, useCallback } from 'react'
import { employeeService } from '@/services/api'
import { Employee, SearchCriteria } from '@/types'

interface UseEmployeesReturn {
  employees: Employee[]
  loading: boolean
  error: string | null
  fetchEmployees: () => Promise<void>
  createEmployee: (data: Partial<Employee>) => Promise<Employee>
  updateEmployee: (id: number, data: Partial<Employee>) => Promise<Employee>
  deleteEmployee: (id: number) => Promise<void>
  searchEmployees: (criteria: SearchCriteria) => Promise<Employee[]>
  refresh: () => Promise<void>
  clearError: () => void
}

/**
 * Custom hook for managing employee data and operations
 */
export const useEmployees = (): UseEmployeesReturn => {
  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  /**
   * Fetch all employees
   */
  const fetchEmployees = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await employeeService.getEmployees()
      setEmployees(data)
    } catch (err: any) {
      setError(err.message || 'Failed to fetch employees')
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * Create a new employee
   */
  const createEmployee = useCallback(async (employeeData: Partial<Employee>): Promise<Employee> => {
    try {
      setLoading(true)
      setError(null)
      const newEmployee = await employeeService.createEmployee(employeeData)
      setEmployees(prev => [...prev, newEmployee])
      return newEmployee
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to create employee'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * Update an existing employee
   */
  const updateEmployee = useCallback(async (id: number, employeeData: Partial<Employee>): Promise<Employee> => {
    try {
      setLoading(true)
      setError(null)
      const updatedEmployee = await employeeService.updateEmployee(id, employeeData)
      setEmployees(prev => prev.map(emp => emp.id === id ? updatedEmployee : emp))
      return updatedEmployee
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to update employee'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * Delete an employee
   */
  const deleteEmployee = useCallback(async (id: number): Promise<void> => {
    try {
      setLoading(true)
      setError(null)
      await employeeService.deleteEmployee(id)
      setEmployees(prev => prev.filter(emp => emp.id !== id))
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to delete employee'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * Search employees based on criteria
   */
  const searchEmployees = useCallback(async (criteria: SearchCriteria): Promise<Employee[]> => {
    try {
      setLoading(true)
      setError(null)
      const results = await employeeService.searchEmployees(criteria)
      return results
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to search employees'
      setError(errorMessage)
      throw new Error(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [])

  /**
   * Refresh employee data
   */
  const refresh = useCallback(async () => {
    await fetchEmployees()
  }, [fetchEmployees])

  /**
   * Clear error state
   */
  const clearError = useCallback(() => {
    setError(null)
  }, [])

  // Fetch employees on mount
  useEffect(() => {
    fetchEmployees()
  }, [fetchEmployees])

  return {
    employees,
    loading,
    error,
    fetchEmployees,
    createEmployee,
    updateEmployee,
    deleteEmployee,
    searchEmployees,
    refresh,
    clearError,
  }
}
