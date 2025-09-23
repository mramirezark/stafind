# Service Architecture Documentation

## Overview

The API services have been refactored into domain-specific services following the Single Responsibility Principle. Each service handles a specific domain of the application, making the codebase more maintainable, testable, and scalable.

## Architecture

### 🏗️ **Service Structure**

```
services/api/
├── baseService.ts           # Base class with common functionality
├── employeeService.ts       # Employee domain operations
├── jobRequestService.ts     # Job request domain operations
├── authService.ts          # Authentication domain operations
├── skillService.ts         # Skill domain operations
├── dashboardService.ts     # Dashboard domain operations
├── searchService.ts        # Search domain operations
├── index.ts               # Centralized exports
├── apiService.ts          # Legacy service (backward compatibility)
└── README.md              # General API documentation
```

### 🎯 **Domain Services**

| Service | Domain | Responsibilities |
|---------|--------|------------------|
| `EmployeeService` | Employees | CRUD operations, skills management, search |
| `JobRequestService` | Job Requests | CRUD operations, matching, status management |
| `AuthService` | Authentication | Login, registration, profile management, password reset |
| `SkillService` | Skills | CRUD operations, categories, search, analytics |
| `DashboardService` | Dashboard | Statistics, analytics, recent data, trends |
| `SearchService` | Search | Employee search, suggestions, filters, analytics |

## Base Service Class

### `BaseApiService`

All domain services extend from `BaseApiService` which provides:

- **HTTP Request Handling**: Generic request method with error handling
- **Caching**: Built-in caching with TTL and invalidation
- **Error Management**: Consistent error formatting and handling
- **Cache Management**: Domain-specific cache clearing

```typescript
export abstract class BaseApiService {
  protected cache = new Map<string, { data: any; timestamp: number }>()
  protected readonly CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  protected async request<T>(method: 'GET' | 'POST' | 'PUT' | 'DELETE', url: string, data?: any, useCache: boolean = true): Promise<T>
  protected handleError(error: any): Error
  public clearCache(pattern?: string): void
  protected clearDomainCache(): void
  protected abstract getDomainName(): string
}
```

## Domain Services

### 1. Employee Service

**File**: `employeeService.ts`

**Responsibilities**:
- Employee CRUD operations
- Employee skills management
- Employee search functionality

**Key Methods**:
```typescript
// CRUD Operations
getEmployees(): Promise<Employee[]>
getEmployee(id: number): Promise<Employee>
createEmployee(data: Partial<Employee>): Promise<Employee>
updateEmployee(id: number, data: Partial<Employee>): Promise<Employee>
deleteEmployee(id: number): Promise<void>

// Skills Management
getEmployeeSkills(employeeId: number): Promise<any[]>
addEmployeeSkill(employeeId: number, skillData: any): Promise<any>
updateEmployeeSkill(employeeId: number, skillId: number, skillData: any): Promise<any>
removeEmployeeSkill(employeeId: number, skillId: number): Promise<void>

// Search
searchEmployees(searchCriteria: SearchCriteria): Promise<Employee[]>
```

### 2. Job Request Service

**File**: `jobRequestService.ts`

**Responsibilities**:
- Job request CRUD operations
- Matching functionality
- Status management

**Key Methods**:
```typescript
// CRUD Operations
getJobRequests(): Promise<JobRequest[]>
getJobRequest(id: number): Promise<JobRequest>
createJobRequest(data: Partial<JobRequest>): Promise<JobRequest>
updateJobRequest(id: number, data: Partial<JobRequest>): Promise<JobRequest>
deleteJobRequest(id: number): Promise<void>

// Matching
getJobRequestMatches(id: number): Promise<any[]>
createMatch(jobRequestId: number, employeeId: number, matchData: any): Promise<any>
deleteMatch(jobRequestId: number, matchId: number): Promise<void>

// Status Management
updateJobRequestStatus(id: number, status: string): Promise<JobRequest>
getJobRequestsByStatus(status: string): Promise<JobRequest[]>
getJobRequestsByDepartment(department: string): Promise<JobRequest[]>
```

### 3. Authentication Service

**File**: `authService.ts`

**Responsibilities**:
- User authentication
- Profile management
- Password operations

**Key Methods**:
```typescript
// Authentication
login(username: string, password: string): Promise<{ user: User; token: string }>
register(userData: RegisterData): Promise<{ user: User; token: string }>
logout(): Promise<void>
refreshToken(): Promise<{ token: string }>

// Profile Management
getProfile(): Promise<User>
updateProfile(userData: Partial<User>): Promise<User>
changePassword(currentPassword: string, newPassword: string): Promise<void>

// Password Reset
requestPasswordReset(email: string): Promise<void>
resetPassword(token: string, newPassword: string): Promise<void>

// Email Verification
sendEmailVerification(): Promise<void>
verifyEmail(token: string): Promise<void>
```

### 4. Skill Service

**File**: `skillService.ts`

**Responsibilities**:
- Skill CRUD operations
- Category management
- Search and analytics

**Key Methods**:
```typescript
// CRUD Operations
getSkills(): Promise<Skill[]>
getSkill(id: number): Promise<Skill>
createSkill(data: Partial<Skill>): Promise<Skill>
updateSkill(id: number, data: Partial<Skill>): Promise<Skill>
deleteSkill(id: number): Promise<void>

// Categories
getSkillsByCategory(category: string): Promise<Skill[]>
getSkillCategories(): Promise<string[]>
createSkillCategory(categoryData: { name: string; description?: string }): Promise<any>

// Search & Analytics
searchSkills(query: string): Promise<Skill[]>
getPopularSkills(limit?: number): Promise<Skill[]>
getTrendingSkills(limit?: number): Promise<Skill[]>
```

### 5. Dashboard Service

**File**: `dashboardService.ts`

**Responsibilities**:
- Dashboard statistics
- Analytics and trends
- Recent data

**Key Methods**:
```typescript
// Statistics
getDashboardStats(): Promise<DashboardStats>
getDashboardData(): Promise<{ stats: DashboardStats; recentRequests: JobRequest[]; employees: Employee[] }>

// Recent Data
getRecentJobRequests(limit?: number): Promise<JobRequest[]>
getRecentEmployees(limit?: number): Promise<Employee[]>
getRecentMatches(limit?: number): Promise<any[]>

// Analytics
getDepartmentStats(): Promise<any[]>
getSkillDemandStats(): Promise<any[]>
getMatchingSuccessRate(): Promise<{ success_rate: number; total_matches: number }>
getMonthlyTrends(months?: number): Promise<any[]>

// Activity
getActivityFeed(limit?: number): Promise<any[]>
getUserActivity(userId: number, limit?: number): Promise<any[]>
```

### 6. Search Service

**File**: `searchService.ts`

**Responsibilities**:
- Employee search functionality
- Search suggestions
- Search analytics

**Key Methods**:
```typescript
// Basic Search
searchEmployees(searchCriteria: SearchCriteria): Promise<Employee[]>
searchEmployeesBySkills(skills: string[], operator?: 'AND' | 'OR'): Promise<Employee[]>
searchEmployeesByDepartment(department: string): Promise<Employee[]>
searchEmployeesByLocation(location: string): Promise<Employee[]>
searchEmployeesByExperience(level: string): Promise<Employee[]>

// Advanced Search
advancedSearch(criteria: AdvancedSearchCriteria): Promise<Employee[]>
searchWithFilters(filters: SearchFilters): Promise<PaginatedResponse<Employee>>

// Suggestions
getSearchSuggestions(query: string): Promise<string[]>
getSkillSuggestions(query: string): Promise<string[]>
getDepartmentSuggestions(query: string): Promise<string[]>

// Analytics
getSearchHistory(): Promise<any[]>
saveSearchQuery(query: string, results_count: number): Promise<void>
getPopularSearches(limit?: number): Promise<any[]>
```

## Usage Examples

### Direct Service Usage

```typescript
import { employeeService, authService, dashboardService } from '@/services/api'

// Employee operations
const employees = await employeeService.getEmployees()
const newEmployee = await employeeService.createEmployee(employeeData)

// Authentication
const { user, token } = await authService.login(username, password)

// Dashboard
const stats = await dashboardService.getDashboardStats()
```

### Using Custom Hooks

```typescript
import { useEmployees, useDashboardData, useSearchEmployees } from '@/hooks/useApi'

function MyComponent() {
  const { data: employees, loading, error } = useEmployees()
  const { data: dashboardData } = useDashboardData()
  const { searchEmployees, results } = useSearchEmployees()
  
  // Component logic...
}
```

### Service Instances

```typescript
import { services } from '@/services/api'

// Access all services
const { employee, auth, dashboard, search } = services

// Use services
const employees = await services.employee.getEmployees()
const stats = await services.dashboard.getDashboardStats()
```

## Benefits

### ✅ **Single Responsibility Principle**
Each service handles one domain, making code easier to understand and maintain.

### ✅ **Better Organization**
Related functionality is grouped together, improving code discoverability.

### ✅ **Easier Testing**
Domain services can be tested independently with focused test suites.

### ✅ **Scalability**
New domains can be added without affecting existing services.

### ✅ **Maintainability**
Changes to one domain don't affect others, reducing risk of bugs.

### ✅ **Reusability**
Services can be reused across different parts of the application.

## Migration from Monolithic Service

### Before (Monolithic)
```typescript
import { apiService } from '@/services/api/apiService'

// All operations in one service
const employees = await apiService.getEmployees()
const stats = await apiService.getDashboardStats()
const skills = await apiService.getSkills()
```

### After (Domain-Specific)
```typescript
import { employeeService, dashboardService, skillService } from '@/services/api'

// Domain-specific services
const employees = await employeeService.getEmployees()
const stats = await dashboardService.getDashboardStats()
const skills = await skillService.getSkills()
```

## Best Practices

### 1. Use Domain Services
Always use the appropriate domain service for your operations.

### 2. Leverage Hooks
Use custom hooks for React components instead of direct service calls.

### 3. Handle Errors
Always handle loading and error states in your components.

### 4. Cache Management
Let services handle caching automatically, use manual clearing sparingly.

### 5. Type Safety
Use TypeScript interfaces for better type safety and IntelliSense.

## Future Enhancements

- **Service Composition**: Combine multiple services for complex operations
- **Service Middleware**: Add logging, monitoring, or validation middleware
- **Service Versioning**: Support for API versioning
- **Service Mocking**: Easy mocking for testing
- **Service Documentation**: Auto-generated API documentation

## Conclusion

The domain-specific service architecture provides a clean, maintainable, and scalable foundation for the application. Each service is focused on its specific domain, making the codebase easier to understand, test, and extend.
