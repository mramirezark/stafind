/**
 * Domain-specific type exports
 * 
 * This file provides domain-specific exports for better organization
 * and to avoid circular dependencies.
 */

// Re-export all types from the main index file
export * from './index'

// Domain-specific groupings for easier imports
export type {
  // Employee domain
  Employee,
  EmployeeSkill,
  EmployeeFormData,
  EmployeeCardProps,
  EmployeeFormProps,
} from './index'

export type {
  // Job Request domain
  JobRequest,
  JobRequestFormData,
} from './index'

export type {
  // Authentication domain
  User,
  Role,
  RegisterData,
  AuthContextType,
  AuthCardProps,
  LoginFormProps,
  RegisterFormProps,
  AuthWrapperProps,
} from './index'

export type {
  // UI/Component domain
  NavigationProps,
  TabPanelProps,
  DashboardStats,
} from './index'

export type {
  // Search domain
  SearchCriteria,
  SearchRequest,
  Match,
} from './index'

export type {
  // Utility domain
  Skill,
  SelectOption,
  ExperienceLevel,
  Priority,
  Department,
  LoadingState,
  FormMode,
  ViewType,
} from './index'
