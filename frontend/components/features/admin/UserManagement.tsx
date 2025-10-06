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
  FormControl,
  InputLabel,
  Select,
  Switch,
  FormControlLabel,
  Alert,
  CircularProgress,
  Pagination,
  Tooltip,
} from '@mui/material'
import {
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Person as PersonIcon,
  AdminPanelSettings as AdminIcon,
} from '@mui/icons-material'
import { User, Role, CreateUserRequest, UpdateUserRequest } from '@/types'
import { userManagementService } from '@/services/api'

export function UserManagement() {
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [rolesLoading, setRolesLoading] = useState(true)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create')
  const [formData, setFormData] = useState<CreateUserRequest>({
    username: '',
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    role_id: undefined,
  })

  const limit = 10

  useEffect(() => {
    loadUsers()
    loadRoles()
  }, [page])

  const loadUsers = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await userManagementService.getUsers(limit, (page - 1) * limit)
      console.log('Users response received:', response)
      
      // Ensure response has the expected structure
      if (response && Array.isArray(response.data)) {
        setUsers(response.data)
        if (response.pagination) {
          setTotalPages(Math.ceil(response.pagination.total / limit))
        } else {
          setTotalPages(1)
        }
      } else {
        console.error('Invalid users response structure:', response)
        setUsers([])
        setTotalPages(1)
        setError('Invalid response format from server')
      }
    } catch (err) {
      setError('Failed to load users')
      console.error('Error loading users:', err)
      setUsers([])
      setTotalPages(1)
    } finally {
      setLoading(false)
    }
  }

  const loadRoles = async () => {
    try {
      setRolesLoading(true)
      const rolesData = await userManagementService.getRoles()
      console.log('Roles data received:', rolesData)
      // Ensure rolesData is an array
      if (Array.isArray(rolesData)) {
        setRoles(rolesData)
      } else {
        console.error('Roles data is not an array:', rolesData)
        setRoles([])
      }
    } catch (err) {
      console.error('Error loading roles:', err)
      setRoles([]) // Set empty array on error
    } finally {
      setRolesLoading(false)
    }
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, user: User) => {
    setAnchorEl(event.currentTarget)
    setSelectedUser(user)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedUser(null)
  }

  const handleCreateUser = () => {
    setDialogMode('create')
    setFormData({
      username: '',
      email: '',
      password: '',
      first_name: '',
      last_name: '',
      role_id: undefined,
    })
    setDialogOpen(true)
    handleMenuClose()
  }

  const handleEditUser = () => {
    if (!selectedUser) return
    setDialogMode('edit')
    setFormData({
      username: selectedUser.username,
      email: selectedUser.email,
      password: '', // Not needed for edit, but required by type
      first_name: selectedUser.first_name,
      last_name: selectedUser.last_name,
      role_id: selectedUser.role_id,
    })
    setDialogOpen(true)
    handleMenuClose()
  }

  const handleDeleteUser = async () => {
    if (!selectedUser) return
    try {
      await userManagementService.deleteUser(selectedUser.id)
      await loadUsers()
    } catch (err) {
      setError('Failed to delete user')
    }
    handleMenuClose()
  }

  const handleSubmit = async () => {
    try {
      if (dialogMode === 'create') {
        await userManagementService.createUser(formData)
      } else if (selectedUser) {
        const updateData: UpdateUserRequest = {
          username: formData.username,
          email: formData.email,
          first_name: formData.first_name,
          last_name: formData.last_name,
          role_id: formData.role_id,
        }
        await userManagementService.updateUser(selectedUser.id, updateData)
      }
      setDialogOpen(false)
      await loadUsers()
    } catch (err) {
      setError(`Failed to ${dialogMode} user`)
    }
  }

  const getRoleColor = (roleName?: string) => {
    switch (roleName) {
      case 'admin':
        return 'error'
      case 'hr_manager':
        return 'warning'
      case 'hiring_manager':
        return 'info'
      default:
        return 'default'
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString()
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
          User Management
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleCreateUser}
        >
          Add User
        </Button>
      </Box>

      {error && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }} 
          onClose={() => setError(null)}
          action={
            <Button color="inherit" size="small" onClick={loadUsers}>
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
              <TableCell>User</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Role</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Created</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.length === 0 && !loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                  <Typography variant="body2" color="text.secondary">
                    No users found
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <PersonIcon color="action" />
                    <Box>
                      <Typography variant="body2" fontWeight="medium">
                        {user.first_name} {user.last_name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        @{user.username}
                      </Typography>
                    </Box>
                  </Box>
                </TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  {user.role ? (
                    <Chip
                      label={user.role.name.replace('_', ' ').toUpperCase()}
                      color={getRoleColor(user.role.name) as any}
                      size="small"
                    />
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      No role
                    </Typography>
                  )}
                </TableCell>
                <TableCell>
                  <Chip
                    label={user.is_active ? 'Active' : 'Inactive'}
                    color={user.is_active ? 'success' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>{formatDate(user.created_at)}</TableCell>
                <TableCell align="right">
                  <IconButton
                    onClick={(e) => handleMenuOpen(e, user)}
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

      {/* User Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleEditUser}>
          <EditIcon fontSize="small" sx={{ mr: 1 }} />
          Edit
        </MenuItem>
        <MenuItem onClick={handleDeleteUser} sx={{ color: 'error.main' }}>
          <DeleteIcon fontSize="small" sx={{ mr: 1 }} />
          Delete
        </MenuItem>
      </Menu>

      {/* User Dialog */}
      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {dialogMode === 'create' ? 'Create User' : 'Edit User'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Username"
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              fullWidth
              required
            />
            {dialogMode === 'create' && (
              <TextField
                label="Password"
                type="password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                fullWidth
                required
              />
            )}
            <TextField
              label="First Name"
              value={formData.first_name}
              onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Last Name"
              value={formData.last_name}
              onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
              fullWidth
              required
            />
            <FormControl fullWidth>
              <InputLabel>Role</InputLabel>
              <Select
                value={formData.role_id || ''}
                onChange={(e) => setFormData({ ...formData, role_id: Number(e.target.value) || undefined })}
                label="Role"
                disabled={rolesLoading}
              >
                {rolesLoading ? (
                  <MenuItem disabled>Loading roles...</MenuItem>
                ) : Array.isArray(roles) && roles.length > 0 ? (
                  roles.map((role) => (
                    <MenuItem key={role.id} value={role.id}>
                      {role.name.replace('_', ' ').toUpperCase()}
                    </MenuItem>
                  ))
                ) : (
                  <MenuItem disabled>No roles available</MenuItem>
                )}
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">
            {dialogMode === 'create' ? 'Create' : 'Update'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
