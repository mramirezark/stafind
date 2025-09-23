# StaffFind Authentication System

This document describes the comprehensive authentication and authorization system implemented in StaffFind.

## 🔐 Overview

StaffFind implements a secure authentication system with:
- **JWT-based authentication** with refresh tokens
- **Role-based access control (RBAC)** with multiple user roles
- **Secure password hashing** using bcrypt
- **Session management** with token revocation
- **Frontend authentication context** with automatic token handling

## 🏗️ Architecture

### Backend Components

#### 1. Database Schema
- **Roles Table**: Defines available roles in the system (created first in V1)
- **Users Table**: Stores user credentials and profile information (references roles)
- **User Roles Table**: Many-to-many relationship between users and roles
- **User Sessions Table**: Tracks active sessions and tokens

#### 2. Authentication Models
- **User Model**: Complete user profile with role information
- **Role Model**: Role definitions with descriptions
- **Session Model**: Token management and expiration tracking

#### 3. Security Components
- **JWT Token Management**: Secure token generation and validation
- **Password Hashing**: bcrypt with configurable cost
- **Session Management**: Token storage and revocation
- **Middleware**: Route protection and role-based access

#### 4. API Endpoints
- **Public Routes**: Login, register, refresh token
- **Protected Routes**: User profile, protected resources
- **Admin Routes**: User management, role management

### Frontend Components

#### 1. Authentication Context
- **AuthProvider**: Manages authentication state
- **useAuth Hook**: Provides authentication functions
- **Token Management**: Automatic token storage and refresh

#### 2. Authentication Forms
- **LoginForm**: User login with validation
- **RegisterForm**: User registration with role selection
- **AuthWrapper**: Protects routes and manages auth state

#### 3. User Interface
- **Navigation**: User info display and logout
- **Role Indicators**: Visual role badges and permissions
- **Protected Components**: Role-based UI rendering

## 👥 User Roles

### Default Roles

1. **Admin**
   - Full system access
   - User management
   - Role management
   - All CRUD operations

2. **HR Manager**
   - Access to all employee data
   - Job request management
   - Engineer matching
   - User profile management

3. **Hiring Manager**
   - Job request creation and management
   - Engineer search and matching
   - Limited user access

4. **Employee**
   - View own profile
   - Basic system access
   - Limited functionality

### Role Permissions

| Feature | Admin | HR Manager | Hiring Manager | Employee |
|---------|-------|------------|----------------|----------|
| User Management | ✅ | ✅ | ❌ | ❌ |
| Role Management | ✅ | ❌ | ❌ | ❌ |
| Engineer CRUD | ✅ | ✅ | ✅ | ❌ |
| Job Requests | ✅ | ✅ | ✅ | ❌ |
| Search Engineers | ✅ | ✅ | ✅ | ❌ |
| View Own Profile | ✅ | ✅ | ✅ | ✅ |
| System Settings | ✅ | ❌ | ❌ | ❌ |

## 🔧 API Endpoints

### Authentication Endpoints

#### Public Endpoints
```http
POST /api/v1/auth/login
POST /api/v1/auth/register
POST /api/v1/auth/refresh
```

#### Protected Endpoints
```http
GET    /api/v1/auth/profile
PUT    /api/v1/auth/profile
POST   /api/v1/auth/change-password
POST   /api/v1/auth/logout
```

#### Admin Endpoints
```http
GET    /api/v1/admin/users
GET    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id
DELETE /api/v1/admin/users/:id

GET    /api/v1/roles
GET    /api/v1/roles/:id
POST   /api/v1/admin/roles
PUT    /api/v1/admin/roles/:id
DELETE /api/v1/admin/roles/:id
```

### Request/Response Examples

#### Login Request
```json
POST /api/v1/auth/login
{
  "email": "admin@stafind.com",
  "password": "admin123"
}
```

#### Login Response
```json
{
  "user": {
    "id": 1,
    "email": "admin@stafind.com",
    "first_name": "System",
    "last_name": "Administrator",
    "role": {
      "id": 1,
      "name": "admin",
      "description": "System administrator with full access"
    },
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Register Request
```json
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role_id": 4
}
```

## 🔒 Security Features

### Password Security
- **bcrypt Hashing**: Industry-standard password hashing
- **Configurable Cost**: Adjustable hashing complexity
- **Password Validation**: Minimum length and complexity requirements

### Token Security
- **JWT Tokens**: Secure, stateless authentication
- **Token Expiration**: 24-hour token lifetime
- **Token Refresh**: Automatic token renewal
- **Session Storage**: Server-side session tracking
- **Token Revocation**: Secure logout with token invalidation

### Access Control
- **Route Protection**: Middleware-based route protection
- **Role Verification**: Server-side role checking
- **Permission Validation**: Granular permission system
- **CSRF Protection**: Built-in CSRF protection

### Data Protection
- **Password Exclusion**: Passwords never returned in API responses
- **Sensitive Data**: Secure handling of authentication data
- **Input Validation**: Comprehensive request validation
- **Error Handling**: Secure error messages without information leakage

## 🚀 Usage Guide

### Backend Setup

1. **Database Migration**
   ```bash
   # Run migrations to create auth tables (V1 creates users and roles first)
   make migrate
   ```

2. **Default Admin User**
   - Email: `admin@stafind.com`
   - Password: `admin123`
   - Role: `admin`

3. **Environment Variables**
   ```env
   JWT_SECRET=your-secret-key-change-in-production
   DATABASE_URL=postgres://user:pass@localhost:5432/stafind
   ```

### Frontend Integration

1. **Authentication Context**
   ```tsx
   import { useAuth } from '@/lib/auth'
   
   function MyComponent() {
     const { user, isAuthenticated, login, logout } = useAuth()
     
     if (!isAuthenticated) {
       return <LoginForm />
     }
     
     return <div>Welcome, {user.first_name}!</div>
   }
   ```

2. **Protected Routes**
   ```tsx
   import { AuthWrapper } from '@/components/AuthWrapper'
   
   function App() {
     return (
       <AuthWrapper>
         <MainApplication />
       </AuthWrapper>
     )
   }
   ```

3. **Role-Based Rendering**
   ```tsx
   function AdminPanel() {
     const { isAdmin } = useAuth()
     
     if (!isAdmin()) {
       return <div>Access denied</div>
     }
     
     return <AdminContent />
   }
   ```

### API Integration

1. **Authenticated Requests**
   ```javascript
   // Token is automatically included in requests
   const response = await api.get('/api/v1/engineers')
   ```

2. **Role Checking**
   ```javascript
   // Server automatically validates user roles
   const response = await api.get('/api/v1/admin/users')
   ```

## 🛠️ Development

### Adding New Roles

1. **Database**
   ```sql
   INSERT INTO roles (name, description) VALUES 
   ('new_role', 'Description of new role');
   ```

2. **Backend**
   ```go
   // Add role constants
   const NewRole = "new_role"
   
   // Add permission checks
   func (u *User) HasNewRole() bool {
       return u.HasRole(NewRole)
   }
   ```

3. **Frontend**
   ```tsx
   // Add role checking
   const { hasRole } = useAuth()
   
   if (hasRole('new_role')) {
     // Show role-specific content
   }
   ```

### Custom Permissions

1. **Service Layer**
   ```go
   func (s *userService) CustomAction(userID int) error {
       user, err := s.userRepo.GetUserByID(userID)
       if err != nil {
           return err
       }
       
       if !user.HasRole("required_role") {
           return NewValidationError("Insufficient permissions")
       }
       
       // Perform action
       return nil
   }
   ```

2. **Handler Layer**
   ```go
   func (h *Handlers) CustomEndpoint(c *fiber.Ctx) error {
       user, err := middleware.GetCurrentUser(c)
       if err != nil {
           return handleServiceError(c, err)
       }
       
       if !middleware.HasRole(c, "required_role") {
           return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
               "error": "Insufficient permissions",
           })
       }
       
       // Handle request
       return c.JSON(fiber.Map{"success": true})
   }
   ```

## 🔍 Testing

### Backend Tests

1. **Unit Tests**
   ```bash
   go test ./internal/auth/...
   go test ./internal/services/...
   ```

2. **Integration Tests**
   ```bash
   go test ./cmd/server/...
   ```

### Frontend Tests

1. **Component Tests**
   ```bash
   npm test -- --testPathPattern=auth
   ```

2. **E2E Tests**
   ```bash
   npm run test:e2e
   ```

### API Testing

1. **Authentication Flow**
   ```bash
   # Login
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@stafind.com","password":"admin123"}'
   
   # Use token
   curl -X GET http://localhost:8080/api/v1/auth/profile \
     -H "Authorization: Bearer YOUR_TOKEN"
   ```

## 📋 Security Checklist

- [ ] JWT secret key is strong and unique
- [ ] Password hashing uses bcrypt with appropriate cost
- [ ] Tokens have reasonable expiration times
- [ ] Session management includes token revocation
- [ ] All routes are properly protected
- [ ] Role-based access control is implemented
- [ ] Input validation is comprehensive
- [ ] Error messages don't leak sensitive information
- [ ] CORS is properly configured
- [ ] HTTPS is enforced in production

## 🚨 Security Considerations

### Production Deployment

1. **Environment Variables**
   - Use strong, unique JWT secrets
   - Configure secure database connections
   - Set appropriate CORS origins

2. **Database Security**
   - Use connection encryption
   - Implement proper access controls
   - Regular security updates

3. **Application Security**
   - Enable HTTPS only
   - Implement rate limiting
   - Regular security audits

### Monitoring

1. **Authentication Metrics**
   - Failed login attempts
   - Token refresh rates
   - Session durations

2. **Security Alerts**
   - Unusual access patterns
   - Multiple failed logins
   - Token anomalies

## 📚 Additional Resources

- [JWT Best Practices](https://tools.ietf.org/html/rfc7519)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [bcrypt Security](https://en.wikipedia.org/wiki/Bcrypt)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Material-UI Authentication](https://mui.com/material-ui/getting-started/usage/)

## 🤝 Contributing

When contributing to the authentication system:

1. Follow security best practices
2. Add comprehensive tests
3. Update documentation
4. Consider backward compatibility
5. Review with security team

## 📞 Support

For authentication-related issues:
- Check the logs for detailed error messages
- Verify database connectivity
- Ensure proper environment configuration
- Review token expiration settings
