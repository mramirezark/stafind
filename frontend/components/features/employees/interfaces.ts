import { Employee, Skill } from '@/types'

// ============================================================================
// EMPLOYEE MANAGEMENT INTERFACES
// ============================================================================

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

// ============================================================================
// SEARCH & FILTER INTERFACES
// ============================================================================

export interface SearchFilters {
  searchQuery: string
  departmentFilter: string
  levelFilter: string
  skillFilter: string[]
}

export interface SearchBarProps {
  searchQuery: string
  onSearchChange: (query: string) => void
}

export interface FilterControlsProps {
  filters: SearchFilters
  onFiltersChange: (filters: Partial<SearchFilters>) => void
  onClearFilters: () => void
  skills: Skill[]
  skillsLoading: boolean
}

// ============================================================================
// EMPLOYEE CARD INTERFACES
// ============================================================================

export interface EmployeeCardProps {
  employee: Employee
  viewMode: 'grid' | 'list' | 'table'
  onEdit?: (employee: Employee) => void
  onDelete?: (id: number) => void
}

// ============================================================================
// EMPLOYEE FORM INTERFACES
// ============================================================================

export interface EmployeeFormDialogProps {
  open: boolean
  onClose: () => void
  onSubmit: (data: EmployeeFormData) => Promise<void>
  employee?: Employee | null
  skills: Skill[]
  skillsLoading: boolean
  loading: boolean
  error: string | null
}

export interface EmployeeFormProps {
  formData: EmployeeFormData
  onFormDataChange: (data: Partial<EmployeeFormData>) => void
  skills: Skill[]
  skillsLoading: boolean
  error: string | null
}

// ============================================================================
// VIEW CONTROLS INTERFACES
// ============================================================================

export interface ViewControlsProps {
  viewMode: 'grid' | 'list' | 'table'
  onViewModeChange: (mode: 'grid' | 'list' | 'table') => void
  onClearFilters: () => void
}

// ============================================================================
// EMPLOYEE LIST INTERFACES
// ============================================================================

export interface EmployeeListProps {
  employees: Employee[]
  viewMode: 'grid' | 'list' | 'table'
  onEdit?: (employee: Employee) => void
  onDelete?: (id: number) => void
}

export interface EmployeeTableProps {
  employees: Employee[]
  onEdit?: (employee: Employee) => void
  onDelete?: (id: number) => void
  page: number
  rowsPerPage: number
  onPageChange: (event: unknown, newPage: number) => void
  onRowsPerPageChange: (event: React.ChangeEvent<HTMLInputElement>) => void
}

// ============================================================================
// MAIN COMPONENT INTERFACES
// ============================================================================

export interface EmployeeManagementProps {
  // Optional props for customization
  showSearch?: boolean
  showFilters?: boolean
  allowCreate?: boolean
  allowEdit?: boolean
  allowDelete?: boolean
  defaultViewMode?: 'grid' | 'list' | 'table'
}
