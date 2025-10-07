/**
 * Custom hooks for API data management
 * 
 * These hooks provide a consistent way to manage API state,
 * loading, and error handling across the application.
 */

import { useState, useEffect, useCallback } from 'react'
import { 
  employeeService, 
  authService, 
  skillService, 
  dashboardService, 
  searchService 
} from '@/services/api'
import { aiAgentService, AIAgentRequest, AIAgentResponse, SkillExtractionResponse, CreateAIAgentRequest } from '@/services/ai/aiAgentService'
import { Employee, Skill, User, SearchCriteria, DashboardStats, DashboardMetrics, TopSuggestedEmployee, SkillDemandStats } from '@/types'

// ============================================================================
// GENERIC API HOOK
// ============================================================================

interface UseApiState<T> {
  data: T | null
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useApi<T>(
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

export function useDashboardMetrics() {
  return useApi(() => dashboardService.getDashboardMetrics())
}

export function useTopSuggestedEmployees(limit: number = 5) {
  return useApi(() => dashboardService.getTopSuggestedEmployees(limit), [limit])
}

export function useSkillDemandStats() {
  return useApi(() => dashboardService.getSkillDemandStats())
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
// AI AGENT HOOKS
// ============================================================================

export function useAIAgentRequests(limit: number = 50, offset: number = 0) {
  return useApi(() => aiAgentService.getRequests(limit, offset), [limit, offset])
}

export function useAIAgentRequest(id: number) {
  return useApi(() => {
    if (id <= 0) {
      return Promise.resolve(null);
    }
    return aiAgentService.getRequest(id);
  }, [id])
}

export function useProcessAIAgentRequest() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const processRequest = useCallback(async (id: number) => {
    try {
      setLoading(true)
      setError(null)
      const result = await aiAgentService.processRequest(id)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    processRequest,
    loading,
    error
  }
}

export function useCreateAndProcessAIAgentRequest() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const processNewRequest = useCallback(async (data: CreateAIAgentRequest) => {
    try {
      setLoading(true)
      setError(null)
      const result = await aiAgentService.processNewRequest(data)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    processNewRequest,
    loading,
    error
  }
}

export function useExtractSkills() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const extractSkills = useCallback(async (text: string) => {
    try {
      setLoading(true)
      setError(null)
      const result = await aiAgentService.extractSkills(text)
      return result
    } catch (err: any) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }, [])

  return {
    extractSkills,
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
    authService.clearCache(pattern)
    skillService.clearCache(pattern)
    dashboardService.clearCache(pattern)
    searchService.clearCache(pattern)
    aiAgentService.clearCache(pattern)
  }, [])

  return {
    clearCache
  }
}
