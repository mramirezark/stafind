'use client'

import { useState, useEffect } from 'react'
import {
  Box,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
  Pagination,
  Tooltip,
  Snackbar,
  FormControlLabel,
  Switch,
} from '@mui/material'
import {
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  VpnKey as KeyIcon,
  ContentCopy as CopyIcon,
  Refresh as RotateIcon,
  Block as DeactivateIcon,
  CheckCircle as ActiveIcon,
} from '@mui/icons-material'
import { APIKey, CreateAPIKeyRequest } from '@/types'
import { apiKeyService } from '@/services/api'

export function APIKeyManagement() {
  const [apiKeys, setApiKeys] = useState<APIKey[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedKey, setSelectedKey] = useState<APIKey | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [formData, setFormData] = useState<CreateAPIKeyRequest>({
    service_name: '',
    description: '',
    expires_at: '',
  })
  const [newKeyDialog, setNewKeyDialog] = useState(false)
  const [newKey, setNewKey] = useState<string>('')

  const limit = 10

  useEffect(() => {
    loadAPIKeys()
  }, [page])

  const loadAPIKeys = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await apiKeyService.getAPIKeys(limit, (page - 1) * limit)
      console.log('API Keys response received:', response)
      
      // Ensure response has the expected structure
      if (response && Array.isArray(response.data)) {
        setApiKeys(response.data)
        if (response.pagination) {
          setTotalPages(Math.ceil(response.pagination.total / limit))
        } else {
          setTotalPages(1)
        }
      } else {
        console.error('Invalid API keys response structure:', response)
        setApiKeys([])
        setTotalPages(1)
        setError('Invalid response format from server')
      }
    } catch (err) {
      setError('Failed to load API keys')
      console.error('Error loading API keys:', err)
      setApiKeys([])
      setTotalPages(1)
    } finally {
      setLoading(false)
    }
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, key: APIKey) => {
    setAnchorEl(event.currentTarget)
    setSelectedKey(key)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedKey(null)
  }

  const handleCreateKey = () => {
    setFormData({
      service_name: '',
      description: '',
      expires_at: '',
    })
    setDialogOpen(true)
    handleMenuClose()
  }


  const handleDeactivateKey = async () => {
    if (!selectedKey) return
    try {
      await apiKeyService.deactivateAPIKey(selectedKey.id)
      setSuccess('API key deactivated successfully')
      await loadAPIKeys()
    } catch (err) {
      setError('Failed to deactivate API key')
    }
    handleMenuClose()
  }

  const handleRotateKey = async () => {
    if (!selectedKey) return
    try {
      const response = await apiKeyService.rotateAPIKey(selectedKey.id)
      setNewKey(response.key_preview) // Assuming the response includes the new key
      setNewKeyDialog(true)
      setSuccess('API key rotated successfully')
      await loadAPIKeys()
    } catch (err) {
      setError('Failed to rotate API key')
    }
    handleMenuClose()
  }

  const handleSubmit = async () => {
    try {
      const response = await apiKeyService.createAPIKey(formData)
      setNewKey(response.key_preview) // Assuming the response includes the new key
      setNewKeyDialog(true)
      setSuccess('API key created successfully')
      setDialogOpen(false)
      await loadAPIKeys()
    } catch (err) {
      setError('Failed to create API key')
    }
  }

  const handleCopyKey = (keyPreview: string) => {
    navigator.clipboard.writeText(keyPreview)
    setSuccess('API key copied to clipboard')
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString()
  }

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const isExpired = (expiresAt?: string) => {
    if (!expiresAt) return false
    return new Date(expiresAt) < new Date()
  }

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h5" component="h2">
          API Key Management
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleCreateKey}
        >
          Create API Key
        </Button>
      </Box>

      {error && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }} 
          onClose={() => setError(null)}
          action={
            <Button color="inherit" size="small" onClick={loadAPIKeys}>
              Retry
            </Button>
          }
        >
          {error}
        </Alert>
      )}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Service</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Key Preview</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Expires</TableCell>
              <TableCell>Last Used</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {apiKeys.length === 0 && !loading ? (
              <TableRow>
                <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    No API keys found
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              apiKeys.map((key) => (
              <TableRow key={key.id}>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <KeyIcon color="action" />
                    <Typography variant="body2" fontWeight="medium">
                      {key.service_name}
                    </Typography>
                  </Box>
                </TableCell>
                <TableCell>{key.description}</TableCell>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="body2" fontFamily="monospace">
                      {key.key_preview}
                    </Typography>
                    <Tooltip title="Copy key">
                      <IconButton
                        size="small"
                        onClick={() => handleCopyKey(key.key_preview)}
                      >
                        <CopyIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip
                    label={key.is_active ? 'Active' : 'Inactive'}
                    color={key.is_active ? 'success' : 'default'}
                    size="small"
                    icon={key.is_active ? <ActiveIcon /> : <DeactivateIcon />}
                  />
                </TableCell>
                <TableCell>
                  {key.expires_at ? (
                    <Typography
                      variant="body2"
                      color={isExpired(key.expires_at) ? 'error' : 'text.primary'}
                    >
                      {formatDate(key.expires_at)}
                    </Typography>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      Never
                    </Typography>
                  )}
                </TableCell>
                <TableCell>
                  {key.last_used_at ? (
                    <Typography variant="body2">
                      {formatDateTime(key.last_used_at)}
                    </Typography>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      Never
                    </Typography>
                  )}
                </TableCell>
                <TableCell align="right">
                  <IconButton
                    onClick={(e) => handleMenuOpen(e, key)}
                    size="small"
                  >
                    <MoreVertIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      {totalPages > 1 && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
          <Pagination
            count={totalPages}
            page={page}
            onChange={(_, newPage) => setPage(newPage)}
            color="primary"
          />
        </Box>
      )}

      {/* API Key Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleRotateKey}>
          <RotateIcon fontSize="small" sx={{ mr: 1 }} />
          Rotate
        </MenuItem>
        <MenuItem onClick={handleDeactivateKey}>
          <DeactivateIcon fontSize="small" sx={{ mr: 1 }} />
          Deactivate
        </MenuItem>
      </Menu>

      {/* API Key Dialog */}
      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          Create API Key
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Service Name"
              value={formData.service_name}
              onChange={(e) => setFormData({ ...formData, service_name: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Description"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
            <TextField
              label="Expires At"
              type="datetime-local"
              value={formData.expires_at}
              onChange={(e) => setFormData({ ...formData, expires_at: e.target.value })}
              fullWidth
              InputLabelProps={{ shrink: true }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">
            Create
          </Button>
        </DialogActions>
      </Dialog>

      {/* New Key Dialog */}
      <Dialog open={newKeyDialog} onClose={() => setNewKeyDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>New API Key</DialogTitle>
        <DialogContent>
          <Alert severity="warning" sx={{ mb: 2 }}>
            This is the only time you'll see this API key. Make sure to copy it and store it securely.
          </Alert>
          <TextField
            label="API Key"
            value={newKey}
            fullWidth
            multiline
            rows={3}
            InputProps={{
              readOnly: true,
              style: { fontFamily: 'monospace' }
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setNewKeyDialog(false)}>Close</Button>
          <Button
            onClick={() => handleCopyKey(newKey)}
            variant="contained"
            startIcon={<CopyIcon />}
          >
            Copy Key
          </Button>
        </DialogActions>
      </Dialog>

      {/* Success Snackbar */}
      <Snackbar
        open={Boolean(success)}
        autoHideDuration={6000}
        onClose={() => setSuccess(null)}
        message={success}
      />
    </Box>
  )
}
