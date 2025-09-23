# Components Directory Structure

This directory follows React best practices for component organization, using a feature-based approach with shared components.

## Directory Structure

```
components/
├── index.ts                    # Main export file
├── README.md                   # This documentation
├── features/                   # Feature-based components
│   ├── employees/             # Employee management feature
│   │   ├── index.ts
│   │   ├── EmployeeCard.tsx
│   │   ├── EmployeeForm.tsx
│   │   ├── EmployeeList.tsx
│   │   ├── EmployeeManagement.tsx
│   │   └── EmployeeListExample.tsx
│   ├── job-requests/          # Job request feature
│   │   ├── index.ts
│   │   └── JobRequestForm.tsx
│   ├── search/                # Search feature
│   │   ├── index.ts
│   │   └── SearchEmployees.tsx
│   └── dashboard/             # Dashboard feature
│       ├── index.ts
│       └── Dashboard.tsx
├── shared/                    # Shared components
│   ├── layout/               # Layout components
│   │   ├── index.ts
│   │   ├── Navigation.tsx
│   │   └── AuthWrapper.tsx
│   ├── ui/                   # Reusable UI components
│   │   ├── index.ts
│   │   └── Button.tsx
│   └── forms/                # Form components
│       ├── index.ts
│       └── FormField.tsx
└── auth/                     # Authentication components
    ├── index.ts
    ├── AuthCard.tsx
    ├── LoginForm.tsx
    └── RegisterForm.tsx
```

## Organization Principles

### 1. Feature-Based Organization
Components are organized by business features rather than technical concerns:
- **employees/**: All employee-related components
- **job-requests/**: Job request management components
- **search/**: Search functionality components
- **dashboard/**: Dashboard and analytics components

### 2. Shared Components
Common components used across multiple features:
- **layout/**: Navigation, wrappers, and layout components
- **ui/**: Reusable UI elements (buttons, inputs, etc.)
- **forms/**: Form-specific components and utilities

### 3. Authentication
Authentication components are kept separate as they're used across the entire application.

## Usage Patterns

### Importing Components

```typescript
// Import from main components index (recommended)
import { EmployeeCard, Dashboard, Navigation } from '@/components'

// Import from specific feature
import { EmployeeCard } from '@/components/features/employees'

// Import from shared components
import { Button } from '@/components/shared/ui'
```

### Creating New Components

1. **Feature Components**: Place in appropriate feature directory
2. **Shared Components**: Place in `shared/` directory based on type
3. **Auth Components**: Place in `auth/` directory

### Component Naming Conventions

- **PascalCase**: All component files use PascalCase
- **Descriptive Names**: Component names clearly indicate their purpose
- **Consistent Suffixes**: Use consistent suffixes like `Form`, `Card`, `List`

## Best Practices

### 1. Single Responsibility
Each component should have a single, well-defined responsibility.

### 2. Composition over Inheritance
Prefer composition and props over complex inheritance patterns.

### 3. Consistent Props Interface
Use TypeScript interfaces for all component props, defined in `/types`.

### 4. Index Files
Each directory should have an `index.ts` file for clean exports.

### 5. Documentation
Complex components should include JSDoc comments explaining their purpose and usage.

## File Structure Guidelines

### Component Files
- One component per file
- Component name should match file name
- Export as named export, not default export

### Index Files
- Export all public components
- Use re-exports for clean API
- Group related exports logically

### Type Definitions
- All types defined in `/types` directory
- Import types from centralized location
- Use consistent naming conventions

## Migration Notes

This structure replaces the previous flat organization:
- ✅ Removed duplicate components
- ✅ Organized by feature domains
- ✅ Created shared component categories
- ✅ Added proper index files
- ✅ Updated import paths throughout application

## Future Considerations

- **Storybook Integration**: Consider adding Storybook for component documentation
- **Testing**: Add component-specific test files
- **Lazy Loading**: Implement lazy loading for feature components
- **Theme Integration**: Ensure all components work with the design system
