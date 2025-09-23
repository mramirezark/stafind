# File Tracking System for StaffFind

This document explains the comprehensive file tracking system implemented for managing uploaded files with database persistence.

## 🗄️ **Database Schema**

### Uploaded Files Table
```sql
CREATE TABLE uploaded_files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,                    -- Unique stored filename
    original_filename VARCHAR(255) NOT NULL,           -- Original upload filename
    file_path VARCHAR(500) NOT NULL,                   -- Full disk path
    file_size BIGINT NOT NULL,                         -- Size in bytes
    content_type VARCHAR(100) NOT NULL,                -- MIME type
    file_hash VARCHAR(64),                             -- SHA-256 hash for deduplication
    upload_type VARCHAR(50) DEFAULT 'resume',          -- Type: resume, document, image, etc.
    status VARCHAR(20) DEFAULT 'active',               -- Status: active, deleted, archived
    uploaded_by INTEGER REFERENCES users(id),          -- User who uploaded
    employee_id INTEGER REFERENCES employees(id),      -- Associated employee
    metadata JSONB,                                     -- Additional metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL                          -- Soft delete timestamp
);
```

### Key Features
- ✅ **Unique Filenames**: Prevents conflicts with timestamp + random string
- ✅ **File Deduplication**: SHA-256 hash prevents duplicate uploads
- ✅ **Soft Deletes**: Files marked as deleted, not physically removed
- ✅ **Metadata Storage**: JSONB field for flexible metadata
- ✅ **User Tracking**: Links files to uploaders and employees
- ✅ **Status Management**: Track file lifecycle states

## 🏗️ **Architecture**

### Backend Components

**1. Models (`models/uploaded_file.go`)**
```go
type UploadedFile struct {
    ID               int       `json:"id"`
    Filename         string    `json:"filename"`
    OriginalFilename string    `json:"original_filename"`
    FilePath         string    `json:"file_path"`
    FileSize         int64     `json:"file_size"`
    ContentType      string    `json:"content_type"`
    FileHash         *string   `json:"file_hash,omitempty"`
    UploadType       string    `json:"upload_type"`
    Status           string    `json:"status"`
    UploadedBy       *int      `json:"uploaded_by,omitempty"`
    EmployeeID       *int      `json:"employee_id,omitempty"`
    Metadata         *string   `json:"metadata,omitempty"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}
```

**2. Repository (`repositories/uploaded_file_repository.go`)**
- ✅ **CRUD Operations**: Create, Read, Update, Delete
- ✅ **Advanced Filtering**: By type, status, user, employee
- ✅ **Pagination**: Efficient large dataset handling
- ✅ **Search**: Full-text search on filenames
- ✅ **Statistics**: File counts, sizes, types
- ✅ **Cleanup**: Remove old deleted files

**3. Service (`services/uploaded_file_service.go`)**
- ✅ **Business Logic**: File validation, hash calculation
- ✅ **Duplicate Detection**: Prevent duplicate uploads
- ✅ **File Integrity**: Verify files exist on disk
- ✅ **Cleanup Operations**: Remove orphaned files
- ✅ **Statistics**: Comprehensive file analytics

**4. Handlers (`handlers/file_upload_handlers.go`)**
- ✅ **REST API**: Complete CRUD endpoints
- ✅ **File Upload**: Single and multiple file uploads
- ✅ **File Management**: Update, delete, list files
- ✅ **Statistics**: File upload analytics
- ✅ **Error Handling**: Comprehensive error responses

## 🚀 **API Endpoints**

### File Upload
```bash
# Upload single file
POST /api/v1/upload
Content-Type: multipart/form-data
Body: file (binary)

# Upload multiple files
POST /api/v1/upload/multiple
Content-Type: multipart/form-data
Body: files[] (binary array)
```

### File Management
```bash
# List files with pagination and filtering
GET /api/v1/files/?page=1&page_size=20&upload_type=resume&status=active&search=john

# Get file by ID
GET /api/v1/files/{id}

# Update file metadata
PUT /api/v1/files/{id}
Body: { "status": "archived", "employee_id": 123 }

# Soft delete file
DELETE /api/v1/files/{id}

# Get file statistics
GET /api/v1/files/stats
```

### File Access
```bash
# Serve file directly
GET /uploads/{filename}

# Legacy file listing (filesystem only)
GET /api/v1/upload/
```

## 📊 **File Statistics**

### Available Metrics
```json
{
  "total_files": 150,
  "total_size": 52428800,
  "files_by_type": {
    "resume": 120,
    "document": 25,
    "image": 5
  },
  "files_by_status": {
    "active": 140,
    "deleted": 8,
    "archived": 2
  },
  "recent_uploads": [...]
}
```

### Query Parameters
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20)
- `upload_type`: Filter by type (resume, document, image)
- `status`: Filter by status (active, deleted, archived)
- `search`: Search in filenames
- `sort_by`: Sort field (created_at, file_size, filename)
- `sort_order`: Sort direction (asc, desc)

## 🔧 **Setup and Usage**

### 1. Database Migration
```bash
# Run the migration to create the table
cd backend
flyway migrate
```

### 2. Start the Server
```bash
cd backend
go run cmd/server/main.go
```

### 3. Test File Upload
```bash
cd backend
go run test_file_tracking.go
```

### 4. Frontend Integration
The frontend `ResumeUpload` component automatically uses the new tracking system:

```typescript
// File upload now includes database tracking
const uploadResponse = await fetch('/api/v1/upload', {
  method: 'POST',
  body: formData,
})

// Response includes database ID and metadata
const result = await uploadResponse.json()
console.log('File ID:', result.file.id)
console.log('File Hash:', result.file.fileHash)
```

## 🛡️ **Security Features**

### File Validation
- ✅ **Size Limits**: 10MB maximum per file
- ✅ **Type Validation**: Only PDF, DOCX, DOC allowed
- ✅ **Hash Verification**: SHA-256 integrity checking
- ✅ **Path Security**: Directory traversal protection

### Access Control
- ✅ **User Tracking**: Link uploads to users
- ✅ **Employee Association**: Link files to employees
- ✅ **Soft Deletes**: Preserve audit trail
- ✅ **Metadata Encryption**: Secure sensitive data

## 📈 **Performance Optimizations**

### Database Indexes
```sql
-- Performance indexes
CREATE INDEX idx_uploaded_files_filename ON uploaded_files(filename);
CREATE INDEX idx_uploaded_files_upload_type ON uploaded_files(upload_type);
CREATE INDEX idx_uploaded_files_status ON uploaded_files(status);
CREATE INDEX idx_uploaded_files_uploaded_by ON uploaded_files(uploaded_by);
CREATE INDEX idx_uploaded_files_employee_id ON uploaded_files(employee_id);
CREATE INDEX idx_uploaded_files_created_at ON uploaded_files(created_at);
CREATE INDEX idx_uploaded_files_file_hash ON uploaded_files(file_hash);
```

### Caching Strategy
- ✅ **File Hash Caching**: Prevent duplicate calculations
- ✅ **Metadata Caching**: Cache frequently accessed data
- ✅ **Statistics Caching**: Cache file statistics
- ✅ **Pagination Caching**: Cache paginated results

## 🔄 **File Lifecycle Management**

### Upload Process
1. **File Validation**: Check size, type, security
2. **Hash Calculation**: Generate SHA-256 hash
3. **Duplicate Check**: Verify file not already uploaded
4. **Physical Storage**: Save file to disk
5. **Database Record**: Create tracking record
6. **Response**: Return file metadata and ID

### Update Process
1. **Validation**: Check file exists and user permissions
2. **Metadata Update**: Update database record
3. **File Integrity**: Verify physical file still exists
4. **Response**: Return updated file data

### Delete Process
1. **Soft Delete**: Mark as deleted in database
2. **Physical Removal**: Delete file from disk
3. **Audit Trail**: Preserve deletion timestamp
4. **Cleanup**: Schedule for permanent removal

## 🧹 **Maintenance Operations**

### Cleanup Old Files
```go
// Remove files deleted more than 30 days ago
rowsAffected, err := fileService.CleanupOldFiles(30)
```

### File Integrity Check
```go
// Verify all database files exist on disk
files, err := fileService.ListFiles(filters)
for _, file := range files {
    if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
        // Mark as deleted or handle missing file
    }
}
```

### Statistics Monitoring
```go
// Get comprehensive file statistics
stats, err := fileService.GetFileStats()
fmt.Printf("Total files: %d, Total size: %d bytes\n", 
    stats.TotalFiles, stats.TotalSize)
```

## 🐛 **Troubleshooting**

### Common Issues

**1. File Upload Fails**
- Check file size (must be < 10MB)
- Verify file type is allowed
- Ensure uploads directory exists and is writable

**2. Database Errors**
- Verify database connection
- Check if migration was run
- Ensure proper permissions

**3. File Not Found**
- Check if file exists in uploads directory
- Verify filename in database matches disk
- Check file permissions

### Debug Commands
```bash
# Check database table
psql -d stafind -c "SELECT COUNT(*) FROM uploaded_files;"

# Check uploads directory
ls -la backend/uploads/

# Test file upload
curl -X POST http://localhost:8080/api/v1/upload -F "file=@test.pdf"

# List files
curl http://localhost:8080/api/v1/files/

# Get statistics
curl http://localhost:8080/api/v1/files/stats
```

## 📝 **Future Enhancements**

### Planned Features
- 🔄 **File Versioning**: Track file versions and changes
- 🔍 **Advanced Search**: Full-text search in file contents
- 📊 **Analytics Dashboard**: Visual file usage statistics
- 🔐 **Access Control**: Role-based file access permissions
- ☁️ **Cloud Storage**: Integration with AWS S3, Google Cloud
- 🔄 **Sync**: Real-time file synchronization
- 📱 **Mobile API**: Mobile-optimized file management

### Integration Points
- **Resume Parsing**: Link parsed data to file records
- **Employee Management**: Associate files with employee profiles
- **Job Requests**: Link files to job applications
- **Audit Logging**: Track all file operations
- **Backup System**: Automated file backup and recovery

## 🎯 **Benefits**

### For Developers
- ✅ **Centralized Management**: Single source of truth for files
- ✅ **Type Safety**: Strong typing with Go structs
- ✅ **Error Handling**: Comprehensive error management
- ✅ **Testing**: Easy to test and mock
- ✅ **Documentation**: Self-documenting API

### For Users
- ✅ **File History**: Track all uploaded files
- ✅ **Search**: Find files quickly
- ✅ **Organization**: Categorize and tag files
- ✅ **Security**: Secure file handling
- ✅ **Performance**: Fast file operations

### For Administrators
- ✅ **Monitoring**: Track file usage and storage
- ✅ **Cleanup**: Automated maintenance operations
- ✅ **Audit**: Complete file operation history
- ✅ **Statistics**: Comprehensive usage analytics
- ✅ **Control**: Manage file access and permissions

The file tracking system provides a robust, scalable solution for managing uploaded files with full database persistence, comprehensive metadata tracking, and advanced management capabilities. 🚀
