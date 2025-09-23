# Types Directory

This directory contains all TypeScript type definitions for the StaffFind application.

## File Organization

### `index.ts`
The main types file containing all interfaces, types, and constants organized by domain:

- **Core Domain Types**: Employee, JobRequest, Skill, Match, SearchCriteria
- **Authentication Types**: User, Role, AuthContextType, RegisterData
- **Component Props Types**: NavigationProps, EmployeeCardProps, etc.
- **Dashboard Types**: DashboardStats
- **Form Options Types**: SelectOption, ExperienceLevel, Priority
- **API Response Types**: ApiResponse, PaginatedResponse
- **Error Types**: ApiError
- **Utility Types**: LoadingState, FormMode, ViewType
- **Constants**: EXPERIENCE_LEVELS, PRIORITIES, DEPARTMENTS
- **Type Guards**: isEmployee, isJobRequest, isSkill

### `domains.ts`
Domain-specific type exports for better organization and to avoid circular dependencies. Provides grouped exports by domain (Employee, JobRequest, Authentication, etc.).

## Usage

### Basic Import
```typescript
import { Employee, JobRequest, User } from '@/types'
```

### Domain-specific Import
```typescript
import { Employee, EmployeeFormData } from '@/types/domains'
```

### Constants Import
```typescript
import { EXPERIENCE_LEVELS, DEPARTMENTS, PRIORITIES } from '@/types'
```

## Best Practices

1. **Single Source of Truth**: All types are defined in this directory and imported where needed
2. **Domain Organization**: Types are grouped by business domain for better maintainability
3. **Consistent Naming**: All interfaces use PascalCase, types use descriptive names
4. **Type Guards**: Include type guards for runtime type checking
5. **Constants**: Shared constants are exported to avoid duplication
6. **Documentation**: All types include JSDoc comments where helpful

## Adding New Types

When adding new types:

1. Add the type definition to `index.ts` in the appropriate domain section
2. Export it in `domains.ts` if it belongs to a specific domain
3. Update this README if adding new domain categories
4. Ensure the type follows the established naming conventions

## Type Organization

- **Core Types**: Business entities (Employee, JobRequest, etc.)
- **Props Types**: Component prop interfaces
- **Form Types**: Form data interfaces
- **API Types**: Request/response interfaces
- **Utility Types**: Helper types and enums
- **Constants**: Shared constant values
