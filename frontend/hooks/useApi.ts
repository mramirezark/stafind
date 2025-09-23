/**
 * Custom hooks for API data management
 * 
 * These hooks provide a consistent way to manage API state,
 * loading, and error handling across the application.
 */

import { useState, useEffect, useCallback } from 'react'
import { 
  employeeService, 
  jobRequestService, 
  authService, 
  skillService, 
  dashboardService, 
  searchService 
} from '@/services/api'
import { Employee, JobRequest, Skill, User, SearchCriteria, DashboardStats } from '@/types'

// ============================================================================
// GENERIC API HOOK
// ============================================================================

interface UseApiState<T> {
  data: T | null
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

function useApi<T>(
  apiCall: () => Promise<T>,
  dependencies: any[] = []
): UseApiState<T> {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await apiCall()
      setData(result)
    } catch (err: any) {
      setError(err.message)
      console.error('API Error:', err)
    } finally {
      setLoading(false)
    }
  }, dependencies)

  useEffect(() => {
    fetchData()
  }, [fetchData])

  return {
    data,
    loading,
    error,
    refetch: fetchData
  }
}

// ============================================================================
// EMPLOYEE HOOKS
// ============================================================================

export function useEmployees() {
  return useApi(() => employeeService.getEmployees())
}

export function useEmployee(id: number) {
  return useApi(() => employeeService.getEmployee(id), [id])
}

export function useCreateEmployee() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const createEmployee = useCallback(async (employeeData: Partial<Employee>) => {
    try {
      setLoading(true)
      setError(null)
      const result = await employeeService.createEmployee(employeeData)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    createEmployee,
    loading,
    error
  }
}

export function useUpdateEmployee() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const updateEmployee = useCallback(async (id: number, employeeData: Partial<Employee>) => {
    try {
      setLoading(true)
      setError(null)
      const result = await employeeService.updateEmployee(id, employeeData)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    updateEmployee,
    loading,
    error
  }
}

export function useDeleteEmployee() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const deleteEmployee = useCallback(async (id: number) => {
    try {
      setLoading(true)
      setError(null)
      await employeeService.deleteEmployee(id)
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    deleteEmployee,
    loading,
    error
  }
}

// ============================================================================
// JOB REQUEST HOOKS
// ============================================================================

export function useJobRequests() {
  return useApi(() => jobRequestService.getJobRequests())
}

export function useJobRequest(id: number) {
  return useApi(() => jobRequestService.getJobRequest(id), [id])
}

export function useCreateJobRequest() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const createJobRequest = useCallback(async (jobRequestData: Partial<JobRequest>) => {
    try {
      setLoading(true)
      setError(null)
      const result = await jobRequestService.createJobRequest(jobRequestData)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    createJobRequest,
    loading,
    error
  }
}

// ============================================================================
// SKILL HOOKS
// ============================================================================

export function useSkills() {
  return useApi(() => skillService.getSkills())
}

export function useCreateSkill() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const createSkill = useCallback(async (skillData: Partial<Skill>) => {
    try {
      setLoading(true)
      setError(null)
      const result = await skillService.createSkill(skillData)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    createSkill,
    loading,
    error
  }
}

// ============================================================================
// SEARCH HOOKS
// ============================================================================

export function useSearchEmployees() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [results, setResults] = useState<Employee[]>([])

  const searchEmployees = useCallback(async (searchCriteria: SearchCriteria) => {
    try {
      setLoading(true)
      setError(null)
      const result = await searchService.searchEmployees(searchCriteria)
      setResults(result)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    searchEmployees,
    results,
    loading,
    error
  }
}

// ============================================================================
// DASHBOARD HOOKS
// ============================================================================

export function useDashboardData() {
  return useApi(() => dashboardService.getDashboardData())
}

export function useDashboardStats() {
  return useApi(() => dashboardService.getDashboardStats())
}

// ============================================================================
// AUTH HOOKS
// ============================================================================

export function useProfile() {
  return useApi(() => authService.getProfile())
}

export function useUpdateProfile() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const updateProfile = useCallback(async (userData: Partial<User>) => {
    try {
      setLoading(true)
      setError(null)
      const result = await authService.updateProfile(userData)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    updateProfile,
    loading,
    error
  }
}

// ============================================================================
// UTILITY HOOKS
// ============================================================================

export function useApiCache() {
  const clearCache = useCallback((pattern?: string) => {
    // Clear cache for all services
    employeeService.clearCache(pattern)
    jobRequestService.clearCache(pattern)
    authService.clearCache(pattern)
    skillService.clearCache(pattern)
    dashboardService.clearCache(pattern)
    searchService.clearCache(pattern)
  }, [])

  return {
    clearCache
  }
}
