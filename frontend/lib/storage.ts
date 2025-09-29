/**
 * Storage Configuration
 * 
 * This file handles file storage configuration for different environments
 * and storage providers (local, AWS S3, Google Cloud, etc.)
 */

export interface StorageConfig {
  provider: 'local' | 'aws-s3' | 'google-cloud' | 'azure-blob'
  bucket?: string
  region?: string
  accessKeyId?: string
  secretAccessKey?: string
  endpoint?: string
  maxFileSize: number
  allowedTypes: string[]
}

export interface UploadResult {
  url: string
  key: string
  size: number
  contentType: string
}

export interface StorageService {
  uploadFile(file: File, path: string): Promise<UploadResult>
  deleteFile(key: string): Promise<void>
  getFileUrl(key: string): string
}

// Default configuration
const defaultConfig: StorageConfig = {
  provider: 'local',
  maxFileSize: 10 * 1024 * 1024, // 10MB
  allowedTypes: [
    'application/pdf',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/msword',
    'image/jpeg',
    'image/png',
  ],
}

// Environment-based configuration
const getStorageConfig = (): StorageConfig => {
  // Force local storage for now
  return {
    ...defaultConfig,
    provider: 'local',
    maxFileSize: 10 * 1024 * 1024, // 10MB
    allowedTypes: [
      'application/pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/msword',
    ],
  }
}

// Local Storage Service (for development)
class LocalStorageService implements StorageService {
  private baseUrl: string

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  }

  async uploadFile(file: File, path: string): Promise<UploadResult> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('path', path)

    const response = await fetch(`${this.baseUrl}/api/v1/upload`, {
      method: 'POST',
      body: formData,
    })

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.statusText}`)
    }

    const result = await response.json()
    return {
      url: result.url,
      key: result.key,
      size: file.size,
      contentType: file.type,
    }
  }

  async deleteFile(key: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/api/v1/upload/${key}`, {
      method: 'DELETE',
    })

    if (!response.ok) {
      throw new Error(`Delete failed: ${response.statusText}`)
    }
  }

  getFileUrl(key: string): string {
    return `${this.baseUrl}/uploads/${key}`
  }
}

// AWS S3 Storage Service
class S3StorageService implements StorageService {
  private config: StorageConfig
  private s3: any // AWS S3 client

  constructor(config: StorageConfig) {
    this.config = config
    // Initialize AWS S3 client here
    // this.s3 = new AWS.S3({ ... })
  }

  async uploadFile(file: File, path: string): Promise<UploadResult> {
    // Implement S3 upload
    throw new Error('S3 upload not implemented yet')
  }

  async deleteFile(key: string): Promise<void> {
    // Implement S3 delete
    throw new Error('S3 delete not implemented yet')
  }

  getFileUrl(key: string): string {
    return `https://${this.config.bucket}.s3.${this.config.region}.amazonaws.com/${key}`
  }
}

// Google Cloud Storage Service
class GCSStorageService implements StorageService {
  private config: StorageConfig

  constructor(config: StorageConfig) {
    this.config = config
  }

  async uploadFile(file: File, path: string): Promise<UploadResult> {
    // Implement GCS upload
    throw new Error('GCS upload not implemented yet')
  }

  async deleteFile(key: string): Promise<void> {
    // Implement GCS delete
    throw new Error('GCS delete not implemented yet')
  }

  getFileUrl(key: string): string {
    return `https://storage.googleapis.com/${this.config.bucket}/${key}`
  }
}

// Azure Blob Storage Service
class AzureStorageService implements StorageService {
  private config: StorageConfig

  constructor(config: StorageConfig) {
    this.config = config
  }

  async uploadFile(file: File, path: string): Promise<UploadResult> {
    // Implement Azure upload
    throw new Error('Azure upload not implemented yet')
  }

  async deleteFile(key: string): Promise<void> {
    // Implement Azure delete
    throw new Error('Azure delete not implemented yet')
  }

  getFileUrl(key: string): string {
    return `https://${this.config.bucket}.blob.core.windows.net/${key}`
  }
}

// Storage Service Factory
export const createStorageService = (): StorageService => {
  const config = getStorageConfig()

  switch (config.provider) {
    case 'aws-s3':
      return new S3StorageService(config)
    case 'google-cloud':
      return new GCSStorageService(config)
    case 'azure-blob':
      return new AzureStorageService(config)
    case 'local':
    default:
      return new LocalStorageService()
  }
}

// Export singleton instance
export const storageService = createStorageService()

// Utility functions
export const validateFile = (file: File, config: StorageConfig = defaultConfig): { valid: boolean; error?: string } => {
  // Check file size
  if (file.size > config.maxFileSize) {
    return {
      valid: false,
      error: `File size must be less than ${config.maxFileSize / 1024 / 1024}MB`
    }
  }

  // Check file type
  if (!config.allowedTypes.includes(file.type)) {
    return {
      valid: false,
      error: `File type ${file.type} is not allowed. Allowed types: ${config.allowedTypes.join(', ')}`
    }
  }

  return { valid: true }
}

export const generateFilePath = (originalName: string, prefix: string = 'resumes'): string => {
  const timestamp = Date.now()
  const randomString = Math.random().toString(36).substring(2, 8)
  const extension = originalName.split('.').pop()
  return `${prefix}/${timestamp}-${randomString}.${extension}`
}

// Resume-specific storage functions
export const uploadResume = async (file: File): Promise<UploadResult> => {
  const validation = validateFile(file)
  if (!validation.valid) {
    throw new Error(validation.error)
  }

  const path = generateFilePath(file.name, 'resumes')
  return storageService.uploadFile(file, path)
}

export const deleteResume = async (key: string): Promise<void> => {
  return storageService.deleteFile(key)
}

export const getResumeUrl = (key: string): string => {
  return storageService.getFileUrl(key)
}
