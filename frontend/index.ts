/**
 * Main frontend exports
 * 
 * This file provides a centralized way to import
 * components, hooks, services, and utilities.
 */

// Component exports
export * from './components'

// Hook exports
export * from './hooks/useAuth'
export * from './hooks/useEmployees'

// Service exports
export * from './services/api/client'
export * from './services/api/endpoints'
export * from './services/auth/authService'
export * from './services/employees/employeeService'

// Utility exports
export * from './utils/helpers'

// Type exports (excluding EmployeeCardProps to avoid conflict)
export type {
  Employee,
  Skill,
  User,
  SearchCriteria,
  DashboardStats,
  Match,
  AuthContextType,
  NavigationProps,
  EmployeeFormProps,
  AuthCardProps,
  LoginFormProps,
  RegisterFormProps,
  AuthWrapperProps,
  TabPanelProps,
  SelectOption,
  ExperienceLevel,
  Priority,
  ApiResponse,
  PaginatedResponse,
  ApiError,
  LoadingState,
  FormMode,
  ViewType,
  Department,
  EmployeeType,
  SkillType,
  UserType,
  MatchType,
} from './types'
