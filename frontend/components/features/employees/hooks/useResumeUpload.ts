import { useState, useCallback } from 'react'
import { uploadService, ParsedResumeData } from '@/services/upload'
import { ResumeFile } from '../interfaces/resumeUpload'

export function useResumeUpload() {
  const [files, setFiles] = useState<ResumeFile[]>([])
  const [isProcessing, setIsProcessing] = useState(false)
  const [overallProgress, setOverallProgress] = useState(0)

  const addFiles = useCallback((newFiles: File[]) => {
    const resumeFiles: ResumeFile[] = newFiles.map(file => ({
      id: Math.random().toString(36).substr(2, 9),
      file,
      status: 'pending',
      progress: 0,
    }))
    setFiles(prev => [...prev, ...resumeFiles])
  }, [])

  const removeFile = useCallback((id: string) => {
    setFiles(prev => prev.filter(file => file.id !== id))
  }, [])

  const processFiles = useCallback(async (): Promise<ParsedResumeData[]> => {
    setIsProcessing(true)
    setOverallProgress(0)
    
    const results: ParsedResumeData[] = []
    
    for (let i = 0; i < files.length; i++) {
      const file = files[i]
      
      // Update file status to processing
      setFiles(prev => prev.map(f => 
        f.id === file.id ? { ...f, status: 'processing' as const } : f
      ))
      
      try {
        const parsedData = await uploadService.parseResume(file.file)
        
        // Update file status to success
        setFiles(prev => prev.map(f => 
          f.id === file.id 
            ? { ...f, status: 'success' as const, parsedData, progress: 100 }
            : f
        ))
        
        results.push(parsedData)
      } catch (error) {
        // Update file status to error
        setFiles(prev => prev.map(f => 
          f.id === file.id 
            ? { 
                ...f, 
                status: 'error' as const, 
                error: error instanceof Error ? error.message : 'Unknown error'
              }
            : f
        ))
      }
      
      // Update overall progress
      const progress = ((i + 1) / files.length) * 100
      setOverallProgress(progress)
    }
    
    setIsProcessing(false)
    return results
  }, [files])

  const clearFiles = useCallback(() => {
    setFiles([])
    setOverallProgress(0)
    setIsProcessing(false)
  }, [])

  return {
    files,
    isProcessing,
    overallProgress,
    addFiles,
    removeFile,
    processFiles,
    clearFiles,
  }
}
