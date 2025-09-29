export interface ResumeFile {
  id: string
  file: File
  status: 'pending' | 'processing' | 'success' | 'error'
  progress: number
  parsedData?: ParsedResumeData
  error?: string
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

export interface ResumeUploadProps {
  open: boolean
  onClose: () => void
  onBulkImport: (employees: ParsedResumeData[]) => Promise<void>
}
