'use client'

import { useState } from 'react'
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  IconButton,
  useMediaQuery,
  useTheme,
  Menu,
  MenuItem,
  Avatar,
  Chip,
} from '@mui/material'
import {
  Dashboard as DashboardIcon,
  PersonAdd as PersonAddIcon,
  Work as WorkIcon,
  Menu as MenuIcon,
  Close as CloseIcon,
  AccountCircle as AccountIcon,
  Logout as LogoutIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material'
import { useAuth } from '@/lib/auth'
import { NavigationProps } from '@/types'

const navigationItems = [
  { key: 'dashboard', label: 'Dashboard', icon: <DashboardIcon /> },
  { key: 'job-request', label: 'Job Request', icon: <WorkIcon /> },
  { key: 'employee', label: 'Employees', icon: <PersonAddIcon /> },
]

export function Navigation({ activeView, onViewChange }: NavigationProps) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [userMenuAnchor, setUserMenuAnchor] = useState<null | HTMLElement>(null)
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  const { user, logout } = useAuth()

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen)
  }

  const handleNavClick = (view: string) => {
    onViewChange(view)
    if (isMobile) {
      setMobileOpen(false)
    }
  }

  const handleUserMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setUserMenuAnchor(event.currentTarget)
  }

  const handleUserMenuClose = () => {
    setUserMenuAnchor(null)
  }

  const handleLogout = async () => {
    await logout()
    handleUserMenuClose()
  }

  const getUserInitials = () => {
    if (!user) return 'U'
    return `${user.first_name.charAt(0)}${user.last_name.charAt(0)}`.toUpperCase()
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

  const drawer = (
    <Box sx={{ width: 250 }}>
      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <Typography variant="h6" component="div">
          StaffFind
        </Typography>
        {isMobile && (
          <IconButton onClick={handleDrawerToggle}>
            <CloseIcon />
          </IconButton>
        )}
      </Box>
      <List>
        {navigationItems.map((item) => (
          <ListItem key={item.key} disablePadding>
            <ListItemButton
              selected={activeView === item.key}
              onClick={() => handleNavClick(item.key)}
              sx={{
                '&.Mui-selected': {
                  backgroundColor: theme.palette.primary.light,
                  color: 'white',
                  '&:hover': {
                    backgroundColor: theme.palette.primary.main,
                  },
                },
              }}
            >
              <ListItemIcon sx={{ color: activeView === item.key ? 'white' : 'inherit' }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  )

  return (
    <>
      <AppBar position="sticky" elevation={1}>
        <Toolbar>
          {isMobile && (
            <IconButton
              color="inherit"
              aria-label="open drawer"
              edge="start"
              onClick={handleDrawerToggle}
              sx={{ mr: 2 }}
            >
              <MenuIcon />
            </IconButton>
          )}
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            StaffFind
          </Typography>
          {!isMobile && (
            <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
              {navigationItems.map((item) => (
                <Button
                  key={item.key}
                  color="inherit"
                  startIcon={item.icon}
                  onClick={() => handleNavClick(item.key)}
                  variant={activeView === item.key ? 'outlined' : 'text'}
                  sx={{
                    borderColor: 'rgba(255, 255, 255, 0.23)',
                    color: 'white',
                    '&:hover': {
                      borderColor: 'white',
                    },
                  }}
                >
                  {item.label}
                </Button>
              ))}
              
              {/* User Menu */}
              <Box sx={{ ml: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                {user?.role && (
                  <Chip
                    label={user.role.name.replace('_', ' ').toUpperCase()}
                    color={getRoleColor(user.role.name) as any}
                    size="small"
                    variant="outlined"
                    sx={{ color: 'white', borderColor: 'rgba(255, 255, 255, 0.5)' }}
                  />
                )}
                <IconButton
                  color="inherit"
                  onClick={handleUserMenuOpen}
                  aria-label="user menu"
                >
                  <Avatar sx={{ width: 32, height: 32, bgcolor: 'rgba(255, 255, 255, 0.2)' }}>
                    {getUserInitials()}
                  </Avatar>
                </IconButton>
              </Box>
            </Box>
          )}
        </Toolbar>
      </AppBar>

      {/* User Menu */}
      <Menu
        anchorEl={userMenuAnchor}
        open={Boolean(userMenuAnchor)}
        onClose={handleUserMenuClose}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
      >
        <MenuItem disabled>
          <Box>
            <Typography variant="subtitle2">
              {user?.first_name} {user?.last_name}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              @{user?.username}
            </Typography>
          </Box>
        </MenuItem>
        <MenuItem onClick={handleUserMenuClose}>
          <ListItemIcon>
            <SettingsIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Settings</ListItemText>
        </MenuItem>
        <MenuItem onClick={handleLogout}>
          <ListItemIcon>
            <LogoutIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Logout</ListItemText>
        </MenuItem>
      </Menu>

      {/* Mobile Drawer */}
      <Drawer
        variant="temporary"
        anchor="left"
        open={mobileOpen}
        onClose={handleDrawerToggle}
        ModalProps={{
          keepMounted: true, // Better open performance on mobile.
        }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: 250 },
        }}
      >
        {drawer}
      </Drawer>
    </>
  )
}
