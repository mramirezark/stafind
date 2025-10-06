# Local Storage Setup for StaffFind

This document explains how to use local storage for file uploads in the StaffFind application.

## 🗂️ **Directory Structure**

```
stafind/
├── backend/
│   ├── uploads/              # Local file storage directory
│   │   ├── resumes/          # Resume files
│   │   └── temp/             # Temporary files
│   └── internal/handlers/
├── frontend/
│   └── lib/
│       └── storage.ts        # Storage configuration
└── .gitignore               # Excludes uploads/ directory
```

## ⚙️ **Configuration**

### Backend Configuration
The backend automatically creates the `uploads/` directory when the server starts.

### Frontend Configuration
The frontend is configured to use local storage by default:

```typescript
// frontend/lib/storage.ts
const getStorageConfig = (): StorageConfig => {
  return {
    provider: 'local',
    maxFileSize: 10 * 1024 * 1024, // 10MB
    allowedTypes: [
      'application/pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/msword',
    ],
  }
}
```

## 🚀 **API Endpoints**

### File Upload
- **POST** `/api/v1/upload` - Upload single file
- **POST** `/api/v1/upload/multiple` - Upload multiple files
- **DELETE** `/api/v1/upload/:filename` - Delete file
- **GET** `/api/v1/upload/` - List all files
- **GET** `/uploads/:filename` - Serve file

### Example Usage

#### Upload Single File
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@resume.pdf"
```

#### Upload Multiple Files
```bash
curl -X POST http://localhost:8080/api/v1/upload/multiple \
  -F "files=@resume1.pdf" \
  -F "files=@resume2.docx"
```

#### List Files
```bash
curl http://localhost:8080/api/v1/upload/
```

#### Delete File
```bash
curl -X DELETE http://localhost:8080/api/v1/upload/1234567890_abc123.pdf
```

## 📁 **File Storage Details**

### File Naming
Files are stored with unique names to prevent conflicts:
- Format: `{timestamp}_{random_string}.{extension}`
- Example: `1703123456_abc123.pdf`

### File Validation
- **Size Limit**: 10MB per file
- **Allowed Types**: PDF, DOCX, DOC
- **Security**: Directory traversal protection

### Directory Permissions
- **Upload Directory**: `0755` (readable by all, writable by owner)
- **Files**: `0644` (readable by all, writable by owner)

## 🔧 **Development Setup**

### 1. Start the Backend Server
```bash
cd backend
go run cmd/server/main.go
```

### 2. Start the Frontend
```bash
cd frontend
npm run dev
```

### 3. Test File Upload
```bash
cd backend
go run test_upload.go
```

## 🛡️ **Security Considerations**

### Current Implementation
- ✅ File type validation
- ✅ File size limits
- ✅ Directory traversal protection
- ✅ Unique filename generation

### Production Recommendations
- 🔒 Add authentication to upload endpoints
- 🔒 Implement file virus scanning
- 🔒 Add rate limiting
- 🔒 Use HTTPS for file uploads
- 🔒 Implement file cleanup policies

## 📊 **Monitoring**

### File Storage Usage
```bash
# Check uploads directory size
du -sh backend/uploads/

# List all files
ls -la backend/uploads/
```

### Logs
File upload operations are logged by the Fiber logger middleware.

## 🔄 **Migration to Cloud Storage**

When ready to move to cloud storage:

1. **Update Storage Configuration**:
   ```typescript
   // frontend/lib/storage.ts
   const getStorageConfig = (): StorageConfig => {
     return {
       provider: 'aws-s3', // or 'google-cloud', 'azure-blob'
       bucket: 'your-bucket-name',
       region: 'us-east-1',
       // ... other cloud-specific config
     }
   }
   ```

2. **Migrate Existing Files**:
   - Upload existing files to cloud storage
   - Update database references
   - Remove local files

3. **Update Backend**:
   - Replace local file handlers with cloud storage handlers
   - Update file serving logic

## 🐛 **Troubleshooting**

### Common Issues

1. **Upload Directory Not Created**
   - Check server logs for permission errors
   - Ensure the backend process has write permissions

2. **File Upload Fails**
   - Check file size (must be < 10MB)
   - Verify file type is allowed
   - Check network connectivity

3. **Files Not Accessible**
   - Verify the file exists in `backend/uploads/`
   - Check file permissions
   - Ensure the server is running

### Debug Commands
```bash
# Check uploads directory
ls -la backend/uploads/

# Check file permissions
ls -la backend/uploads/resumes/

# Test file upload
curl -v -X POST http://localhost:8080/api/v1/upload -F "file=@test.pdf"
```

## 📝 **Notes**

- Files are stored locally and will be lost if the server is restarted
- This setup is suitable for development and testing
- For production, consider using cloud storage for better reliability and scalability
- The `uploads/` directory is excluded from git to prevent committing large files
