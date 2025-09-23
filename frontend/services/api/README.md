# API Service Architecture

This directory contains the centralized API service layer that serves as the single point of truth for all API calls in the application.

## Architecture Overview

### 🎯 **Single Point of Truth**
All API calls are centralized in `apiService.ts`, making it easy to:
- Modify API endpoints in one place
- Add consistent error handling
- Implement caching strategies
- Add authentication headers
- Monitor API usage

### 🏗️ **Service Structure**

```
services/api/
├── apiService.ts          # Main API service class
├── README.md             # This documentation
└── (legacy files)        # Old API files (can be removed)
```

## Key Features

### ✅ **Centralized API Calls**
- All API endpoints defined in one place
- Consistent request/response handling
- Easy to modify endpoints without touching components

### ✅ **Built-in Caching**
- Automatic caching for GET requests (5-minute TTL)
- Cache invalidation on data mutations
- Manual cache clearing capabilities

### ✅ **Error Handling**
- Consistent error formatting
- Network error detection
- Server error handling
- Automatic token refresh on 401 errors

### ✅ **Type Safety**
- Full TypeScript support
- Type-safe API responses
- IntelliSense support

## Usage Examples

### Basic API Calls
```typescript
import { apiService } from '@/services/api/apiService'

// Get all employees
const employees = await apiService.getEmployees()

// Create new employee
const newEmployee = await apiService.createEmployee(employeeData)

// Search employees
const results = await apiService.searchEmployees(searchCriteria)
```

### Using Custom Hooks
```typescript
import { useEmployees, useDashboardData } from '@/hooks/useApi'

function MyComponent() {
  const { data: employees, loading, error, refetch } = useEmployees()
  const { data: dashboardData } = useDashboardData()
  
  // Component logic...
}
```

## API Service Methods

### Authentication
- `login(username, password)` - User login
- `register(userData)` - User registration
- `logout()` - User logout
- `getProfile()` - Get user profile
- `updateProfile(userData)` - Update user profile
- `changePassword(current, new)` - Change password

### Employees
- `getEmployees()` - Get all employees
- `getEmployee(id)` - Get employee by ID
- `createEmployee(data)` - Create new employee
- `updateEmployee(id, data)` - Update employee
- `deleteEmployee(id)` - Delete employee

### Job Requests
- `getJobRequests()` - Get all job requests
- `getJobRequest(id)` - Get job request by ID
- `createJobRequest(data)` - Create new job request
- `updateJobRequest(id, data)` - Update job request
- `deleteJobRequest(id)` - Delete job request
- `getJobRequestMatches(id)` - Get matches for job request

### Skills
- `getSkills()` - Get all skills
- `createSkill(data)` - Create new skill
- `updateSkill(id, data)` - Update skill
- `deleteSkill(id)` - Delete skill

### Search
- `searchEmployees(criteria)` - Search employees by criteria

### Dashboard
- `getDashboardStats()` - Get dashboard statistics
- `getDashboardData()` - Get complete dashboard data

## Custom Hooks

### Data Fetching Hooks
```typescript
// Basic data fetching
const { data, loading, error, refetch } = useEmployees()
const { data: skills } = useSkills()
const { data: dashboardData } = useDashboardData()

// Mutation hooks
const { createEmployee, loading, error } = useCreateEmployee()
const { updateEmployee, loading, error } = useUpdateEmployee()
const { deleteEmployee, loading, error } = useDeleteEmployee()
```

### Search Hooks
```typescript
const { searchEmployees, results, loading, error } = useSearchEmployees()
```

## Caching Strategy

### Automatic Caching
- GET requests are automatically cached for 5 minutes
- Cache is invalidated when related data is modified
- Manual cache clearing available

### Cache Management
```typescript
import { apiService } from '@/services/api/apiService'

// Clear all cache
apiService.clearCache()

// Clear specific cache pattern
apiService.clearCache('employees')
```

## Error Handling

### Consistent Error Format
All API errors are converted to standard Error objects with descriptive messages:
- Network errors: "Network error: Unable to connect to server"
- Server errors: "404: Resource not found"
- Validation errors: "400: Invalid input data"

### Error Recovery
- Automatic token refresh on 401 errors
- Cache invalidation on authentication changes
- Graceful degradation for network issues

## Migration Benefits

### Before (Scattered API Calls)
```typescript
// In Dashboard component
const [employeesResponse, jobRequestsResponse] = await Promise.all([
  api.get('/employees'),
  api.get('/job-requests')
])

// In SearchEmployees component
const response = await endpoints.search.employees(searchCriteria)
```

### After (Centralized API Service)
```typescript
// In Dashboard component
const { data: dashboardData } = useDashboardData()

// In SearchEmployees component
const { searchEmployees, results } = useSearchEmployees()
```

## Best Practices

### 1. Use Custom Hooks
Always use the provided custom hooks instead of calling `apiService` directly in components.

### 2. Handle Loading States
```typescript
const { data, loading, error } = useEmployees()

if (loading) return <LoadingSpinner />
if (error) return <ErrorMessage error={error} />
return <DataComponent data={data} />
```

### 3. Cache Management
- Let the service handle automatic caching
- Use manual cache clearing sparingly
- Clear cache after mutations

### 4. Error Handling
- Always handle loading and error states
- Provide meaningful error messages to users
- Log errors for debugging

## Future Enhancements

- **Request Deduplication**: Prevent duplicate requests
- **Optimistic Updates**: Update UI before server response
- **Background Sync**: Sync data in background
- **Offline Support**: Cache data for offline usage
- **Request Queuing**: Queue requests when offline

## Migration Complete ✅

All components now use the centralized API service:
- ✅ Single point of truth for all API calls
- ✅ Consistent error handling and loading states
- ✅ Built-in caching and optimization
- ✅ Type-safe API interactions
- ✅ Easy to maintain and modify

The API layer is now much more maintainable and follows React best practices!
