'use client'

import React from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
  LinearProgress,
  Alert,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  IconButton,
  Chip,
  CircularProgress,
  Paper,
} from '@mui/material'
import {
  CloudUpload as UploadIcon,
  Delete as DeleteIcon,
  CheckCircle as CheckIcon,
  Error as ErrorIcon,
  Description as FileIcon,
} from '@mui/icons-material'
import { useResumeUpload } from './hooks/useResumeUpload'
import { ResumeUploadProps } from './interfaces/resumeUpload'

export function ResumeUpload({ onBulkImport, onClose, open }: ResumeUploadProps) {
  const {
    files,
    isProcessing,
    overallProgress,
    addFiles,
    removeFile,
    processFiles,
    clearFiles,
  } = useResumeUpload()

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = event.target.files
    if (selectedFiles) {
      addFiles(Array.from(selectedFiles))
    }
  }

  const handleDragOver = (event: React.DragEvent) => {
    event.preventDefault()
  }

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault()
    const droppedFiles = event.dataTransfer.files
    if (droppedFiles) {
      addFiles(Array.from(droppedFiles))
    }
  }

  const handleProcessFiles = async () => {
    try {
      const results = await processFiles()
      if (results.length > 0) {
        await onBulkImport(results)
        clearFiles()
        onClose()
      }
    } catch (error) {
      console.error('Failed to process files:', error)
    }
  }

  const handleClose = () => {
    clearFiles()
    onClose()
  }

  const renderStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <FileIcon />
      case 'success':
        return <CheckIcon color="success" />
      case 'error':
        return <ErrorIcon color="error" />
      case 'processing':
        return <LinearProgress />
      default:
        return <FileIcon />
    }
  }

  const successfulFiles = files.filter(f => f.status === 'success')
  const hasErrors = files.some(f => f.status === 'error')

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        Bulk Import Resumes
        <IconButton
          aria-label="close"
          onClick={handleClose}
          sx={{
            position: 'absolute',
            right: 8,
            top: 8,
          }}
        >
          <DeleteIcon />
        </IconButton>
      </DialogTitle>
      
      <DialogContent>
        <Box sx={{ mt: 2 }}>
          {/* Upload Area */}
          <Paper
            onDragOver={handleDragOver}
            onDrop={handleDrop}
            onClick={() => document.getElementById('file-input')?.click()}
            sx={{
              p: 4,
              textAlign: 'center',
              border: '2px dashed',
              borderColor: 'grey.300',
              backgroundColor: 'background.paper',
              cursor: 'pointer',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                borderColor: 'primary.main',
                backgroundColor: 'action.hover',
              },
            }}
          >
            <input
              id="file-input"
              type="file"
              multiple
              accept=".pdf,.doc,.docx"
              onChange={handleFileSelect}
              style={{ display: 'none' }}
            />
            <UploadIcon sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Drag & drop resumes here
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              or click to select files
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Supports PDF, DOC, DOCX files
            </Typography>
          </Paper>

          {/* File List */}
          {files.length > 0 && (
            <Box sx={{ mt: 3 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6">
                  Files ({files.length})
                </Typography>
                <Button
                  variant="contained"
                  onClick={handleProcessFiles}
                  disabled={isProcessing || files.length === 0}
                  startIcon={<UploadIcon />}
                >
                  {isProcessing ? 'Processing...' : 'Process Files'}
                </Button>
              </Box>

              {/* Overall Progress */}
              {isProcessing && (
                <Box sx={{ mb: 2 }}>
                  <Typography variant="body2" gutterBottom>
                    Processing files... {Math.round(overallProgress)}%
                  </Typography>
                  <LinearProgress variant="determinate" value={overallProgress} />
                </Box>
              )}

              <List>
                {files.map((file) => (
                  <Box key={file.id}>
                    <ListItem>
                      <Box sx={{ mr: 2 }}>
                        {renderStatusIcon(file.status)}
                      </Box>
                      <ListItemText
                        primary={file.file.name}
                        secondary={
                          <Box>
                            <Typography variant="caption" display="block">
                              {(file.file.size / 1024).toFixed(2)} KB
                            </Typography>
                            {file.status === 'processing' && (
                              <LinearProgress
                                variant="determinate"
                                value={file.progress}
                                sx={{ mt: 1 }}
                              />
                            )}
                            {file.status === 'error' && file.error && (
                              <Typography variant="caption" color="error">
                                {file.error}
                              </Typography>
                            )}
                            {file.status === 'success' && file.parsedData && (
                              <Box sx={{ mt: 1 }}>
                                <Typography variant="caption" display="block">
                                  {file.parsedData.name} - {file.parsedData.email}
                                </Typography>
                                <Box sx={{ mt: 0.5 }}>
                                  {file.parsedData.skills.slice(0, 3).map((skill, index) => (
                                    <Chip
                                      key={index}
                                      label={skill}
                                      size="small"
                                      sx={{ mr: 0.5, mb: 0.5 }}
                                    />
                                  ))}
                                  {file.parsedData.skills.length > 3 && (
                                    <Chip
                                      label={`+${file.parsedData.skills.length - 3} more`}
                                      size="small"
                                      variant="outlined"
                                    />
                                  )}
                                </Box>
                              </Box>
                            )}
                          </Box>
                        }
                      />
                      <IconButton
                        onClick={() => removeFile(file.id)}
                        disabled={isProcessing}
                        color="error"
                      >
                        <DeleteIcon />
                      </IconButton>
                    </ListItem>
                  </Box>
                ))}
              </List>

              {/* Success/Error Messages */}
              {successfulFiles.length > 0 && (
                <Alert severity="success" sx={{ mt: 2 }}>
                  {successfulFiles.length} files processed successfully
                </Alert>
              )}

              {hasErrors && (
                <Alert severity="error" sx={{ mt: 2 }}>
                  Some files failed to process. Please check the errors above.
                </Alert>
              )}
            </Box>
          )}
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={handleClose}>
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleProcessFiles}
          disabled={isProcessing || successfulFiles.length === 0}
          startIcon={isProcessing ? <CircularProgress size={16} /> : <UploadIcon />}
        >
          {isProcessing ? 'Processing...' : `Import ${successfulFiles.length} Employees`}
        </Button>
      </DialogActions>
    </Dialog>
  )
}