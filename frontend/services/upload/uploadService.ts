import { apiClient } from '../api/client'

export interface FileUploadResponse {
  id: number
  filename: string
  originalName: string
  size: number
  contentType: string
  uploadType: string
  status: string
  url: string
  createdAt: string
}

export interface FileUploadRequest {
  file: File
  uploadType?: string
  employeeId?: number
}

export interface ParsedResumeData {
  name: string
  email: string
  phone?: string
  skills: string[]
  experience: string
  location?: string
  bio?: string
}

class UploadService {
  async uploadFile(request: FileUploadRequest): Promise<FileUploadResponse> {
    const formData = new FormData()
    formData.append('file', request.file)
    
    if (request.uploadType) {
      formData.append('upload_type', request.uploadType)
    }
    
    if (request.employeeId) {
      formData.append('employee_id', request.employeeId.toString())
    }

    const response = await apiClient.post('/api/v1/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    return response.data.file
  }

  async uploadMultipleFiles(files: File[], uploadType?: string): Promise<FileUploadResponse[]> {
    const formData = new FormData()
    
    files.forEach(file => {
      formData.append('files', file)
    })
    
    if (uploadType) {
      formData.append('upload_type', uploadType)
    }

    const response = await apiClient.post('/api/v1/upload/multiple', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    return response.data.uploadedFiles || []
  }

  async parseResume(file: File): Promise<ParsedResumeData> {
    // First upload the file
    const uploadResponse = await this.uploadFile({
      file,
      uploadType: 'resume'
    })

    // Simulate processing delay
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Mock data - replace with actual resume parsing logic
    const mockData: ParsedResumeData = {
      name: `John Doe ${Math.random().toString(36).substr(2, 5)}`,
      email: `john.doe${Math.random().toString(36).substr(2, 5)}@example.com`,
      skills: ['JavaScript', 'React', 'Node.js', 'TypeScript', 'Python'],
      experience: 'Senior',
      location: 'San Francisco, CA',
    }
    
    return mockData
  }

  async getFileStats() {
    const response = await apiClient.get('/api/v1/files/stats')
    return response.data
  }

  async listFiles(params?: {
    page?: number
    pageSize?: number
    uploadType?: string
    status?: string
    search?: string
  }) {
    const response = await apiClient.get('/api/v1/files/', { params })
    return response.data
  }

  async deleteFile(fileId: number) {
    const response = await apiClient.delete(`/api/v1/files/${fileId}`)
    return response.data
  }
}

export const uploadService = new UploadService()
