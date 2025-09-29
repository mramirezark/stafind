/**
 * OPTIMIZED API Services Index
 * 
 * This is an alternative optimized version using dynamic imports for maximum performance.
 * Use this version if you want the best possible bundle splitting and lazy loading.
 * 
 * @example
 * ```typescript
 * // Dynamic service loading
 * import { getService } from '@/services/api/optimizedIndex'
 * 
 * const employeeService = await getService('employee')
 * const employees = await employeeService.getEmployees()
 * ```
 */

// ============================================================================
// BASE SERVICE
// ============================================================================

export { BaseApiService } from './baseService'

// ============================================================================
// DOMAIN SERVICES - CLASSES AND SINGLETONS
// ============================================================================

// Employee Service
export { EmployeeService, employeeService } from './employeeService'


// Authentication Service
export { AuthService, authService } from './authService'

// Skill Service
export { SkillService, skillService } from './skillService'

// Dashboard Service
export { DashboardService, dashboardService } from './dashboardService'

// Search Service
export { SearchService, searchService } from './searchService'

// ============================================================================
// LEGACY API SERVICE (for backward compatibility)
// ============================================================================

export { ApiService, apiService } from './apiService'

// ============================================================================
// TYPES
// ============================================================================

export type ServiceName = 'employee' | 'auth' | 'skill' | 'dashboard' | 'search'

export type ServiceMap = {
  employee: typeof import('./employeeService').employeeService
  auth: typeof import('./authService').authService
  skill: typeof import('./skillService').skillService
  dashboard: typeof import('./dashboardService').dashboardService
  search: typeof import('./searchService').searchService
}

// ============================================================================
// DYNAMIC SERVICE LOADING
// ============================================================================

/**
 * Service loader cache to avoid repeated dynamic imports
 */
const serviceCache = new Map<ServiceName, Promise<any>>()

/**
 * Get a service dynamically with caching
 * 
 * @param name - The service name
 * @returns Promise that resolves to the service instance
 * 
 * @example
 * ```typescript
 * import { getService } from '@/services/api/optimizedIndex'
 * 
 * const employeeService = await getService('employee')
 * const employees = await employeeService.getEmployees()
 * ```
 */
export async function getService<T extends ServiceName>(name: T): Promise<ServiceMap[T]> {
  if (serviceCache.has(name)) {
    return serviceCache.get(name)! as Promise<ServiceMap[T]>
  }

  const servicePromise = (async () => {
    switch (name) {
      case 'employee':
        return (await import('./employeeService')).employeeService
      case 'auth':
        return (await import('./authService')).authService
      case 'skill':
        return (await import('./skillService')).skillService
      case 'dashboard':
        return (await import('./dashboardService')).dashboardService
      case 'search':
        return (await import('./searchService')).searchService
      default:
        throw new Error(`Unknown service: ${name}`)
    }
  })() as Promise<ServiceMap[T]>

  serviceCache.set(name, servicePromise)
  return servicePromise
}

/**
 * Preload all services for better performance
 * 
 * @returns Promise that resolves when all services are loaded
 * 
 * @example
 * ```typescript
 * import { preloadServices } from '@/services/api/optimizedIndex'
 * 
 * // Preload all services
 * await preloadServices()
 * ```
 */
export async function preloadServices(): Promise<void> {
  const serviceNames: ServiceName[] = ['employee', 'auth', 'skill', 'dashboard', 'search']
  await Promise.all(serviceNames.map(name => getService(name)))
}

/**
 * Clear service cache
 * 
 * @example
 * ```typescript
 * import { clearServiceCache } from '@/services/api/optimizedIndex'
 * 
 * // Clear cache to force reload
 * clearServiceCache()
 * ```
 */
export function clearServiceCache(): void {
  serviceCache.clear()
}

/**
 * Check if a service is cached
 * 
 * @param name - The service name
 * @returns True if the service is cached
 */
export function isServiceCached(name: ServiceName): boolean {
  return serviceCache.has(name)
}

/**
 * Get all cached service names
 * 
 * @returns Array of cached service names
 */
export function getCachedServices(): ServiceName[] {
  return Array.from(serviceCache.keys())
}
