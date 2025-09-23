/**
 * API Services Index
 * 
 * Centralized exports for all domain-specific API services.
 * Provides both class exports and singleton instances for maximum flexibility.
 * 
 * @example
 * ```typescript
 * // Direct service usage
 * import { employeeService, authService } from '@/services/api'
 * 
 * // Services object usage
 * import { services } from '@/services/api'
 * const employees = await services.employee.getEmployees()
 * 
 * // Class usage (for testing or custom instances)
 * import { EmployeeService } from '@/services/api'
 * const customService = new EmployeeService()
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

// Job Request Service
export { JobRequestService, jobRequestService } from './jobRequestService'

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
// SERVICES OBJECT - LAZY LOADED FOR PERFORMANCE
// ============================================================================

/**
 * Centralized services object for easy access to all domain services.
 * Uses lazy loading to improve initial bundle size and performance.
 * 
 * @example
 * ```typescript
 * import { services } from '@/services/api'
 * 
 * // All services are available
 * const employees = await services.employee.getEmployees()
 * const stats = await services.dashboard.getDashboardStats()
 * const skills = await services.skill.getSkills()
 * ```
 */
export const services = {
  get employee() {
    return require('./employeeService').employeeService
  },
  get jobRequest() {
    return require('./jobRequestService').jobRequestService
  },
  get auth() {
    return require('./authService').authService
  },
  get skill() {
    return require('./skillService').skillService
  },
  get dashboard() {
    return require('./dashboardService').dashboardService
  },
  get search() {
    return require('./searchService').searchService
  },
} as const

// ============================================================================
// TYPES
// ============================================================================

/**
 * Type for the services object
 */
export type Services = typeof services

/**
 * Type for individual service names
 */
export type ServiceName = keyof Services

/**
 * Type for service instances
 */
export type ServiceInstance = Services[ServiceName]

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

/**
 * Get a service by name
 * 
 * @param name - The service name
 * @returns The service instance
 * 
 * @example
 * ```typescript
 * import { getService } from '@/services/api'
 * 
 * const employeeService = getService('employee')
 * const employees = await employeeService.getEmployees()
 * ```
 */
export function getService<T extends ServiceName>(name: T): Services[T] {
  return services[name]
}

/**
 * Check if a service name is valid
 * 
 * @param name - The service name to check
 * @returns True if the service name is valid
 * 
 * @example
 * ```typescript
 * import { isValidService } from '@/services/api'
 * 
 * if (isValidService('employee')) {
 *   // Safe to use the service
 * }
 * ```
 */
export function isValidService(name: string): name is ServiceName {
  return name in services
}

/**
 * Get all available service names
 * 
 * @returns Array of all service names
 * 
 * @example
 * ```typescript
 * import { getServiceNames } from '@/services/api'
 * 
 * const allServices = getServiceNames()
 * console.log('Available services:', allServices)
 * ```
 */
export function getServiceNames(): ServiceName[] {
  return Object.keys(services) as ServiceName[]
}
