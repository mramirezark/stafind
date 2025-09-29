/**
 * TypeScript type definitions for StaffFind application
 * 
 * This file contains all interfaces and types used throughout the application.
 * Types are organized by domain for better maintainability and discoverability.
 */

// ============================================================================
// CORE DOMAIN TYPES
// ============================================================================

/**
 * Employee-related types
 */
export interface Employee {
  id: number
  name: string
  email: string
  department: string
  level: string
  location: string
  bio: string
  current_project: string | null
  skills: EmployeeSkill[]
  created_at: string
  updated_at: string
}

export interface EmployeeSkill {
  id: number
  name: string
  category: string
  proficiency_level?: number
  years_experience?: number
}

export interface EmployeeFormData {
  name: string
  email: string
  department: string
  level: string
  location: string
  bio: string
  current_project: string | null
  skills: string[]
}


/**
 * Skill-related types
 */
export interface Skill {
  id: number
  name: string
  category: string
}

/**
 * Match-related types (for AI agent)
 */
export interface Match {
  id: number
  employee_id: number
  match_score: number
  matching_skills: string[]
  notes?: string
  created_at: string
  employee?: Employee
}

/**
 * Search-related types
 */
export interface SearchCriteria {
  required_skills: string[]
  preferred_skills: string[]
  department: string
  experience_level: string
  location: string
}

export interface SearchRequest {
  required_skills: string[]
  preferred_skills: string[]
  department?: string
  experience_level?: string
  location?: string
  min_match_score?: number
}

// ============================================================================
// AUTHENTICATION TYPES
// ============================================================================

export interface User {
  id: number
  username: string
  email: string
  first_name: string
  last_name: string
  role_id?: number
  role?: Role
  is_active: boolean
  last_login?: string
  created_at: string
  updated_at: string
  roles?: Role[]
}

export interface Role {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface APIKey {
  id: number
  service_name: string
  description: string
  key_hash: string
  key_preview: string
  is_active: boolean
  expires_at?: string
  last_used_at?: string
  created_at: string
  updated_at: string
}

export interface CreateAPIKeyRequest {
  service_name: string
  description: string
  expires_at?: string
}

export interface UpdateAPIKeyRequest {
  service_name?: string
  description?: string
  is_active?: boolean
  expires_at?: string
}

export interface CreateUserRequest extends RegisterData {
  role_id?: number
}

export interface UpdateUserRequest {
  username?: string
  email?: string
  first_name?: string
  last_name?: string
  role_id?: number
  is_active?: boolean
}

export interface RegisterData {
  username: string
  email: string
  password: string
  first_name: string
  last_name: string
  role_id?: number
}

export interface AuthContextType {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (userData: RegisterData) => Promise<void>
  logout: () => void
  updateProfile: (userData: Partial<User>) => Promise<void>
  changePassword: (currentPassword: string, newPassword: string) => Promise<void>
  hasRole: (roleName: string) => boolean
  isAdmin: () => boolean
  isHRManager: () => boolean
  isHiringManager: () => boolean
}

// ============================================================================
// COMPONENT PROPS TYPES
// ============================================================================

/**
 * Navigation component props
 */
export interface NavigationProps {
  activeView: string
  onViewChange: (view: string) => void
}

/**
 * Employee component props
 */
export interface EmployeeCardProps {
  employee: Employee
  onEdit: (employee: Employee) => void
  onDelete: (id: number) => void
}

export interface EmployeeFormProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: EmployeeFormData) => Promise<void>
  initialData?: EmployeeFormData
  availableSkills: Skill[]
  loading?: boolean
  title?: string
}

/**
 * Authentication component props
 */
export interface AuthCardProps {
  onSuccess?: () => void
}

export interface LoginFormProps {
  onSuccess?: () => void
}

export interface RegisterFormProps {
  onSuccess?: () => void
  onToggleMode?: () => void
}

export interface AuthWrapperProps {
  children: React.ReactNode
}

export interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

// ============================================================================
// DASHBOARD TYPES
// ============================================================================

export interface DashboardStats {
  totalEmployees: number
  totalRequests: number
  completedRequests: number
  pendingRequests: number
}

export interface TopSuggestedEmployee {
  employee_id: number
  employee_name: string
  employee_email: string
  department: string
  level: string
  location: string
  current_project: string | null
  match_count: number
  avg_match_score: number
}

export interface SkillDemandStats {
  skill_name: string
  count: number
  category: string
}

export interface AIAgentRequest {
  id: number
  teams_message_id: string
  channel_id: string
  user_id: string
  user_name: string
  message_text: string
  attachment_url?: string
  extracted_text?: string
  extracted_skills?: string[]
  status: string
  error?: string
  created_at: string
  processed_at?: string
}

export interface DashboardMetrics {
  stats: DashboardStats
  most_requested_skills: SkillDemandStats[]
  top_suggested_employees: TopSuggestedEmployee[]
  recent_requests: AIAgentRequest[]
}

// ============================================================================
// FORM OPTIONS TYPES
// ============================================================================

export interface SelectOption {
  value: string
  label: string
}

export interface ExperienceLevel extends SelectOption {
  value: 'junior' | 'mid' | 'senior'
}

export interface Priority extends SelectOption {
  value: 'low' | 'medium' | 'high' | 'urgent'
}

// ============================================================================
// API RESPONSE TYPES
// ============================================================================

export interface ApiResponse<T = any> {
  data: T
  message?: string
  success: boolean
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

// ============================================================================
// ERROR TYPES
// ============================================================================

export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
}

// ============================================================================
// UTILITY TYPES
// ============================================================================

export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

export type FormMode = 'create' | 'edit' | 'view'

export type ViewType = 'dashboard' | 'employee' | 'ai-agent' | 'admin'

// ============================================================================
// CONSTANTS
// ============================================================================

export const EXPERIENCE_LEVELS: ExperienceLevel[] = [
  { value: 'junior', label: 'Junior (0-2 years)' },
  { value: 'mid', label: 'Mid-level (3-5 years)' },
  { value: 'senior', label: 'Senior (6+ years)' },
]

export const PRIORITIES: Priority[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'urgent', label: 'Urgent' },
]

export const DEPARTMENTS = [
  'Engineering',
  'Data Science',
  'Product',
  'Design',
  'Marketing',
  'Sales',
  'Operations',
] as const

export type Department = typeof DEPARTMENTS[number]

// ============================================================================
// TYPE GUARDS
// ============================================================================

export const isEmployee = (obj: any): obj is Employee => {
  return obj && typeof obj.id === 'number' && typeof obj.name === 'string'
}


export const isSkill = (obj: any): obj is Skill => {
  return obj && typeof obj.id === 'number' && typeof obj.name === 'string'
}

// ============================================================================
// RE-EXPORTS FOR CONVENIENCE
// ============================================================================

// Re-export commonly used types for easier imports
export type {
  Employee as EmployeeType,
  Skill as SkillType,
  User as UserType,
  Match as MatchType,
}
