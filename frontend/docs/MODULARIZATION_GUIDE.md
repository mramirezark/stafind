# React/Next.js Modularization Guide for StaffFind

## 🏗️ Recommended Folder Structure

```
frontend/
├── app/                          # Next.js App Router
│   ├── (auth)/                   # Route groups
│   │   ├── login/
│   │   └── register/
│   ├── dashboard/
│   ├── employees/
│   └── layout.tsx
├── components/                   # Reusable UI components
│   ├── ui/                      # Basic UI components
│   │   ├── Button/
│   │   │   ├── Button.tsx
│   │   │   ├── Button.test.tsx
│   │   │   └── index.ts
│   │   ├── Card/
│   │   ├── Modal/
│   │   └── index.ts
│   ├── forms/                   # Form components
│   │   ├── EmployeeForm/
│   │   ├── LoginForm/
│   │   └── index.ts
│   ├── layout/                  # Layout components
│   │   ├── Header/
│   │   ├── Sidebar/
│   │   └── index.ts
│   └── features/                # Feature-specific components
│       ├── auth/
│       ├── employees/
│       └── dashboard/
├── hooks/                       # Custom React hooks
│   ├── useAuth.ts
│   ├── useApi.ts
│   ├── useEmployees.ts
│   └── index.ts
├── services/                    # Business logic & API services
│   ├── api/                     # API layer
│   │   ├── client.ts           # Axios instance
│   │   ├── endpoints.ts        # API endpoints
│   │   └── types.ts            # API types
│   ├── auth/
│   │   ├── authService.ts
│   │   └── types.ts
│   ├── employees/
│   │   ├── employeeService.ts
│   │   └── types.ts
│   └── index.ts
├── store/                       # State management (if using Redux/Zustand)
│   ├── slices/
│   ├── store.ts
│   └── index.ts
├── utils/                       # Utility functions
│   ├── constants.ts
│   ├── helpers.ts
│   ├── validators.ts
│   └── index.ts
├── types/                       # Global TypeScript types
│   ├── api.ts
│   ├── auth.ts
│   ├── employee.ts
│   └── index.ts
└── lib/                         # External library configurations
    ├── auth.tsx                 # Auth context (existing)
    ├── api.ts                   # API client (existing)
    └── index.ts
```

## 🎯 **Modularization Principles**

### 1. **API Layer Modularization**

#### **Current vs. Recommended Approach**

**Current (Good):**
```typescript
// lib/api.ts
export const endpoints = {
  auth: {
    login: (data) => api.post('/api/v1/auth/login', data),
  }
}
```

**Recommended (Better):**
```typescript
// services/api/client.ts
import axios from 'axios'

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  timeout: 10000,
})

// services/api/endpoints.ts
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    REGISTER: '/api/v1/auth/register',
    REFRESH: '/api/v1/auth/refresh',
  },
  EMPLOYEES: {
    LIST: '/api/v1/employees',
    CREATE: '/api/v1/employees',
    UPDATE: (id: number) => `/api/v1/employees/${id}`,
    DELETE: (id: number) => `/api/v1/employees/${id}`,
  }
} as const

// services/auth/authService.ts
import { apiClient } from '../api/client'
import { API_ENDPOINTS } from '../api/endpoints'
import { LoginRequest, LoginResponse } from '../api/types'

export class AuthService {
  static async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.LOGIN, credentials)
    return response.data
  }

  static async register(userData: RegisterRequest): Promise<RegisterResponse> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.REGISTER, userData)
    return response.data
  }

  static async refreshToken(): Promise<{ token: string }> {
    const response = await apiClient.post(API_ENDPOINTS.AUTH.REFRESH)
    return response.data
  }
}
```

### 2. **Custom Hooks for API Calls**

```typescript
// hooks/useEmployees.ts
import { useState, useEffect } from 'react'
import { EmployeeService } from '@/services/employees/employeeService'
import { Employee } from '@/types/employee'

export const useEmployees = () => {
  const [employees, setEmployees] = useState<Employee[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchEmployees = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await EmployeeService.getAll()
      setEmployees(data)
    } catch (err: any) {
      setError(err.message || 'Failed to fetch employees')
    } finally {
      setLoading(false)
    }
  }

  const createEmployee = async (employeeData: CreateEmployeeRequest) => {
    try {
      setLoading(true)
      const newEmployee = await EmployeeService.create(employeeData)
      setEmployees(prev => [...prev, newEmployee])
      return newEmployee
    } catch (err: any) {
      setError(err.message || 'Failed to create employee')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const updateEmployee = async (id: number, employeeData: UpdateEmployeeRequest) => {
    try {
      setLoading(true)
      const updatedEmployee = await EmployeeService.update(id, employeeData)
      setEmployees(prev => prev.map(emp => emp.id === id ? updatedEmployee : emp))
      return updatedEmployee
    } catch (err: any) {
      setError(err.message || 'Failed to update employee')
      throw err
    } finally {
      setLoading(false)
    }
  }

  const deleteEmployee = async (id: number) => {
    try {
      setLoading(true)
      await EmployeeService.delete(id)
      setEmployees(prev => prev.filter(emp => emp.id !== id))
    } catch (err: any) {
      setError(err.message || 'Failed to delete employee')
      throw err
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchEmployees()
  }, [])

  return {
    employees,
    loading,
    error,
    fetchEmployees,
    createEmployee,
    updateEmployee,
    deleteEmployee,
  }
}
```

### 3. **Component Modularization**

```typescript
// components/ui/Button/Button.tsx
import { ButtonHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/utils/helpers'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', loading, children, ...props }, ref) => {
    return (
      <button
        className={cn(
          'inline-flex items-center justify-center rounded-md font-medium transition-colors',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
          'disabled:pointer-events-none disabled:opacity-50',
          {
            'bg-primary text-primary-foreground hover:bg-primary/90': variant === 'primary',
            'bg-secondary text-secondary-foreground hover:bg-secondary/80': variant === 'secondary',
            'border border-input hover:bg-accent': variant === 'outline',
            'hover:bg-accent hover:text-accent-foreground': variant === 'ghost',
          },
          {
            'h-8 px-3 text-sm': size === 'sm',
            'h-10 px-4 py-2': size === 'md',
            'h-12 px-8 text-lg': size === 'lg',
          },
          className
        )}
        ref={ref}
        disabled={loading}
        {...props}
      >
        {loading && <Spinner className="mr-2 h-4 w-4" />}
        {children}
      </button>
    )
  }
)

Button.displayName = 'Button'
```

```typescript
// components/features/employees/EmployeeList/EmployeeList.tsx
import { useEmployees } from '@/hooks/useEmployees'
import { EmployeeCard } from '../EmployeeCard'
import { EmployeeForm } from '../EmployeeForm'
import { Button } from '@/components/ui/Button'
import { useState } from 'react'

export const EmployeeList = () => {
  const {
    employees,
    loading,
    error,
    createEmployee,
    updateEmployee,
    deleteEmployee,
  } = useEmployees()

  const [formOpen, setFormOpen] = useState(false)
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null)

  if (loading && employees.length === 0) {
    return <EmployeeListSkeleton />
  }

  if (error) {
    return <EmployeeListError error={error} onRetry={() => fetchEmployees()} />
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Employees</h1>
        <Button onClick={() => setFormOpen(true)}>
          Add Employee
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {employees.map((employee) => (
          <EmployeeCard
            key={employee.id}
            employee={employee}
            onEdit={setEditingEmployee}
            onDelete={deleteEmployee}
          />
        ))}
      </div>

      <EmployeeForm
        open={formOpen}
        onClose={() => setFormOpen(false)}
        employee={editingEmployee}
        onSubmit={editingEmployee ? updateEmployee : createEmployee}
      />
    </div>
  )
}
```

### 4. **Service Layer Pattern**

```typescript
// services/employees/employeeService.ts
import { apiClient } from '../api/client'
import { API_ENDPOINTS } from '../api/endpoints'
import { Employee, CreateEmployeeRequest, UpdateEmployeeRequest } from '@/types/employee'

export class EmployeeService {
  static async getAll(): Promise<Employee[]> {
    const response = await apiClient.get(API_ENDPOINTS.EMPLOYEES.LIST)
    return response.data
  }

  static async getById(id: number): Promise<Employee> {
    const response = await apiClient.get(API_ENDPOINTS.EMPLOYEES.UPDATE(id))
    return response.data
  }

  static async create(data: CreateEmployeeRequest): Promise<Employee> {
    const response = await apiClient.post(API_ENDPOINTS.EMPLOYEES.CREATE, data)
    return response.data
  }

  static async update(id: number, data: UpdateEmployeeRequest): Promise<Employee> {
    const response = await apiClient.put(API_ENDPOINTS.EMPLOYEES.UPDATE(id), data)
    return response.data
  }

  static async delete(id: number): Promise<void> {
    await apiClient.delete(API_ENDPOINTS.EMPLOYEES.DELETE(id))
  }

  static async search(criteria: SearchCriteria): Promise<Employee[]> {
    const response = await apiClient.post(API_ENDPOINTS.EMPLOYEES.SEARCH, criteria)
    return response.data
  }
}
```

### 5. **Type Safety & Validation**

```typescript
// types/employee.ts
export interface Employee {
  id: number
  name: string
  email: string
  department: string
  level: 'junior' | 'mid' | 'senior'
  location: string
  bio: string
  skills: Skill[]
  created_at: string
  updated_at: string
}

export interface CreateEmployeeRequest {
  name: string
  email: string
  department: string
  level: string
  location: string
  bio: string
  skills: EmployeeSkillRequest[]
}

export interface UpdateEmployeeRequest extends Partial<CreateEmployeeRequest> {
  id: number
}

// utils/validators.ts
import { z } from 'zod'

export const employeeSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email('Invalid email address'),
  department: z.string().min(1, 'Department is required'),
  level: z.enum(['junior', 'mid', 'senior']),
  location: z.string().optional(),
  bio: z.string().optional(),
  skills: z.array(z.string()).min(1, 'At least one skill is required'),
})

export type EmployeeFormData = z.infer<typeof employeeSchema>
```

## 🚀 **Benefits of This Approach**

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Reusability**: Components and hooks can be easily reused
3. **Testability**: Each module can be tested independently
4. **Maintainability**: Changes are isolated to specific modules
5. **Type Safety**: Full TypeScript support throughout
6. **Scalability**: Easy to add new features without affecting existing code

## 📝 **Migration Strategy**

1. **Start with API Layer**: Move API calls to services
2. **Create Custom Hooks**: Extract component logic to hooks
3. **Modularize Components**: Break down large components
4. **Add Type Definitions**: Create comprehensive type system
5. **Implement Validation**: Add form validation with Zod
6. **Add Error Boundaries**: Implement error handling

This modular approach will make your StaffFind application much more maintainable and scalable!
