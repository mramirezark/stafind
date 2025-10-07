'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
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
  Badge,
  Divider,
} from '@mui/material'
import {
  Dashboard as DashboardIcon,
  PersonAdd as PersonAddIcon,
  SmartToy as AIAgentIcon,
  AdminPanelSettings as AdminIcon,
  Menu as MenuIcon,
  Close as CloseIcon,
  AccountCircle as AccountIcon,
  Logout as LogoutIcon,
  Settings as SettingsIcon,
  Psychology as SkillsIcon,
  Notifications as NotificationsIcon,
  Search as SearchIcon,
  TrendingUp as TrendingUpIcon,
} from '@mui/icons-material'
import { useAuth } from '@/lib/auth'
import { NavigationProps } from '@/types'

const mainNavigationItems = [
  { 
    key: 'dashboard', 
    label: 'Dashboard', 
    icon: <DashboardIcon />, 
    color: '#1e3a8a',
    gradient: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 100%)'
  },
  { 
    key: 'employee', 
    label: 'Employees', 
    icon: <PersonAddIcon />, 
    color: '#10b981',
    gradient: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)'
  },
  { 
    key: 'skills', 
    label: 'Skills', 
    icon: <SkillsIcon />, 
    color: '#f59e0b',
    gradient: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)'
  },
  { 
    key: 'ai-agent', 
    label: 'AI Agent', 
    icon: <AIAgentIcon />, 
    color: '#06b6d4',
    gradient: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)'
  },
]

const adminNavigationItems = [
  { 
    key: 'admin', 
    label: 'Admin', 
    icon: <AdminIcon />, 
    color: '#ef4444',
    gradient: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)'
  },
]

export function Navigation({ activeView, onViewChange }: NavigationProps) {
  const [mobileOpen, setMobileOpen] = useState(false)
  const [userMenuAnchor, setUserMenuAnchor] = useState<null | HTMLElement>(null)
  const [notificationsAnchor, setNotificationsAnchor] = useState<null | HTMLElement>(null)
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  const { user, logout, isAdmin } = useAuth()

  // Get navigation items based on user permissions
  const filteredMainNavigationItems = mainNavigationItems
  const filteredAdminNavigationItems = isAdmin() ? adminNavigationItems : []

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

  const handleNotificationsOpen = (event: React.MouseEvent<HTMLElement>) => {
    setNotificationsAnchor(event.currentTarget)
  }

  const handleNotificationsClose = () => {
    setNotificationsAnchor(null)
  }

  const handleLogout = async () => {
    await logout()
    handleUserMenuClose()
  }

  // Mock notifications data
  const notifications = [
    {
      id: 1,
      title: 'New Employee Added',
      message: 'John Doe has been added to the Engineering team',
      time: '2 minutes ago',
      type: 'success',
      read: false
    },
    {
      id: 2,
      title: 'AI Agent Request',
      message: 'New skill extraction request from Sarah Wilson',
      time: '15 minutes ago',
      type: 'info',
      read: false
    },
    {
      id: 3,
      title: 'System Update',
      message: 'Database maintenance completed successfully',
      time: '1 hour ago',
      type: 'warning',
      read: true
    }
  ]

  const unreadCount = notifications.filter(n => !n.read).length

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
    <Box sx={{ width: 280, height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <Box sx={{ 
        p: 3, 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'space-between',
        background: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 100%)',
        color: 'white'
      }}>
        <Box display="flex" alignItems="center" gap={2}>
          <Avatar sx={{ 
            width: 40, 
            height: 40, 
            bgcolor: 'rgba(255, 255, 255, 0.2)',
            backdropFilter: 'blur(10px)'
          }}>
            <TrendingUpIcon />
          </Avatar>
          <Box>
            <Typography variant="h6" component="div" sx={{ fontWeight: 700 }}>
              Staff Manager
            </Typography>
            <Typography variant="caption" sx={{ opacity: 0.8 }}>
              Employee Management
            </Typography>
          </Box>
        </Box>
        {isMobile && (
          <IconButton onClick={handleDrawerToggle} sx={{ color: 'white' }}>
            <CloseIcon />
          </IconButton>
        )}
      </Box>

      {/* Navigation Items */}
      <Box sx={{ flex: 1, p: 2 }}>
        <List sx={{ p: 0 }}>
          {filteredMainNavigationItems.map((item, index) => (
            <motion.div
              key={item.key}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.1, duration: 0.3 }}
            >
              <ListItem disablePadding sx={{ mb: 1 }}>
                <ListItemButton
                  selected={activeView === item.key}
                  onClick={() => handleNavClick(item.key)}
                  sx={{
                    borderRadius: 2,
                    py: 1.5,
                    px: 2,
                    transition: 'all 0.2s ease-in-out',
                    '&.Mui-selected': {
                      background: item.gradient,
                      color: 'white',
                      boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                      transform: 'translateX(4px)',
                      '&:hover': {
                        background: item.gradient,
                        transform: 'translateX(6px)',
                      },
                      '& .MuiListItemIcon-root': {
                        color: 'white',
                      },
                    },
                    '&:hover': {
                      backgroundColor: 'rgba(99, 102, 241, 0.08)',
                      transform: 'translateX(2px)',
                    },
                  }}
                >
                  <ListItemIcon sx={{ 
                    color: activeView === item.key ? 'white' : item.color,
                    minWidth: 40,
                    transition: 'all 0.2s ease-in-out'
                  }}>
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText 
                    primary={item.label} 
                    primaryTypographyProps={{
                      fontWeight: activeView === item.key ? 600 : 500,
                      fontSize: '0.95rem'
                    }}
                  />
                  {activeView === item.key && (
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={{ scale: 1 }}
                      transition={{ delay: 0.1 }}
                    >
                      <Box
                        sx={{
                          width: 6,
                          height: 6,
                          borderRadius: '50%',
                          bgcolor: 'white',
                          ml: 1
                        }}
                      />
                    </motion.div>
                  )}
                </ListItemButton>
              </ListItem>
            </motion.div>
          ))}
        </List>
      </Box>

      {/* User Section */}
      <Box sx={{ p: 2, borderTop: '1px solid rgba(0, 0, 0, 0.1)' }}>
        <Box 
          sx={{ 
            p: 2, 
            borderRadius: 2, 
            bgcolor: 'rgba(99, 102, 241, 0.05)',
            border: '1px solid rgba(99, 102, 241, 0.1)'
          }}
        >
          <Box display="flex" alignItems="center" gap={2}>
            <Avatar sx={{ 
              width: 40, 
              height: 40, 
              bgcolor: 'rgba(99, 102, 241, 0.1)',
              color: theme.palette.primary.main
            }}>
              {getUserInitials()}
            </Avatar>
            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                {user?.first_name} {user?.last_name}
              </Typography>
              <Typography variant="caption" color="text.secondary" noWrap>
                @{user?.username}
              </Typography>
            </Box>
          </Box>
          {user?.role && (
            <Chip
              label={user.role.name.replace('_', ' ').toUpperCase()}
              color={getRoleColor(user.role.name) as any}
              size="small"
              sx={{ mt: 1, fontSize: '0.75rem' }}
            />
          )}
        </Box>
      </Box>
    </Box>
  )

  return (
    <>
      {/* Left Sidebar - Desktop Only */}
      {!isMobile && (
          <Drawer
            variant="permanent"
            sx={{
              width: 280,
              flexShrink: 0,
              '& .MuiDrawer-paper': {
                width: 280,
                boxSizing: 'border-box',
                background: 'linear-gradient(180deg, #ffffff 0%, #f8fafc 100%)',
                borderRight: '1px solid rgba(226, 232, 240, 0.3)',
                boxShadow: '4px 0 24px rgba(0, 0, 0, 0.06)',
                backdropFilter: 'blur(10px)',
              },
            }}
          >
          {drawer}
        </Drawer>
      )}

      {/* Mobile Drawer */}
      <Drawer
        variant="temporary"
        anchor="left"
        open={mobileOpen}
        onClose={handleDrawerToggle}
        ModalProps={{
          keepMounted: true,
        }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: 280 },
        }}
      >
        {drawer}
      </Drawer>

      <AppBar 
        position="sticky" 
        elevation={0}
        sx={{
          background: 'rgba(255, 255, 255, 0.98)',
          backdropFilter: 'blur(20px)',
          borderBottom: '1px solid rgba(226, 232, 240, 0.2)',
          ml: { xs: 0, md: '280px' }, // Offset for desktop sidebar
          width: { xs: '100%', md: 'calc(100% - 280px)' },
          boxShadow: '0 1px 3px rgba(0, 0, 0, 0.05), 0 1px 2px rgba(0, 0, 0, 0.1)',
          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        }}
      >
        <Toolbar sx={{ px: { xs: 2, md: 4 } }}>
          {isMobile && (
            <IconButton
              color="inherit"
              aria-label="open drawer"
              edge="start"
              onClick={handleDrawerToggle}
              sx={{ 
                mr: 2,
                color: 'text.primary',
                bgcolor: 'rgba(99, 102, 241, 0.1)',
                '&:hover': {
                  bgcolor: 'rgba(99, 102, 241, 0.2)',
                }
              }}
            >
              <MenuIcon />
            </IconButton>
          )}
          
          <Box sx={{ flexGrow: 1 }} />

          {!isMobile && (
            <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
              {filteredAdminNavigationItems.map((item) => (
                <motion.div
                  key={item.key}
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <Button
                    color="inherit"
                    startIcon={item.icon}
                    onClick={() => handleNavClick(item.key)}
                    variant={activeView === item.key ? 'contained' : 'text'}
                    sx={{
                      borderRadius: 2,
                      px: 2,
                      py: 1,
                      textTransform: 'none',
                      fontWeight: 600,
                      ...(activeView === item.key ? {
                        background: item.gradient,
                        color: 'white',
                        boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                        '&:hover': {
                          background: item.gradient,
                          boxShadow: '0 6px 16px rgba(0, 0, 0, 0.2)',
                        }
                      } : {
                        color: 'text.primary',
                        '&:hover': {
                          bgcolor: 'rgba(99, 102, 241, 0.08)',
                        },
                      }),
                    }}
                  >
                    {item.label}
                  </Button>
                </motion.div>
              ))}
              
              {/* Notifications */}
              <IconButton
                onClick={handleNotificationsOpen}
                sx={{ 
                  color: 'text.primary',
                  bgcolor: 'rgba(99, 102, 241, 0.1)',
                  '&:hover': {
                    bgcolor: 'rgba(99, 102, 241, 0.2)',
                  }
                }}
              >
                <Badge badgeContent={unreadCount} color="error">
                  <NotificationsIcon />
                </Badge>
              </IconButton>
              
              {/* User Menu */}
              <Box sx={{ ml: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                {user?.role && (
                  <Chip
                    label={user.role.name.replace('_', ' ').toUpperCase()}
                    color={getRoleColor(user.role.name) as any}
                    size="small"
                    sx={{ 
                      fontWeight: 600,
                      fontSize: '0.75rem'
                    }}
                  />
                )}
                <IconButton
                  onClick={handleUserMenuOpen}
                  aria-label="user menu"
                  sx={{
                    '&:hover': {
                      bgcolor: 'rgba(99, 102, 241, 0.1)',
                    }
                  }}
                >
                  <Avatar sx={{ 
                    width: 36, 
                    height: 36, 
                    background: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 100%)',
                    color: 'white',
                    fontWeight: 600
                  }}>
                    {getUserInitials()}
                  </Avatar>
                </IconButton>
              </Box>
            </Box>
          )}
        </Toolbar>
      </AppBar>

      {/* Notifications Menu */}
      <Menu
        anchorEl={notificationsAnchor}
        open={Boolean(notificationsAnchor)}
        onClose={handleNotificationsClose}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        PaperProps={{
          sx: { width: 360, maxHeight: 400 }
        }}
      >
        <Box sx={{ p: 2, borderBottom: '1px solid rgba(0, 0, 0, 0.1)' }}>
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            Notifications
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {unreadCount} unread notifications
          </Typography>
        </Box>
        
        <Box sx={{ maxHeight: 300, overflow: 'auto' }}>
          {notifications.length > 0 ? (
            notifications.map((notification) => (
              <MenuItem 
                key={notification.id}
                sx={{ 
                  py: 2,
                  borderBottom: '1px solid rgba(0, 0, 0, 0.05)',
                  bgcolor: notification.read ? 'transparent' : 'rgba(99, 102, 241, 0.05)'
                }}
              >
                <Box sx={{ width: '100%' }}>
                  <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={1}>
                    <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                      {notification.title}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {notification.time}
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="text.secondary">
                    {notification.message}
                  </Typography>
                  <Box display="flex" justifyContent="flex-end" mt={1}>
                    <Chip
                      label={notification.type}
                      size="small"
                      color={
                        notification.type === 'success' ? 'success' :
                        notification.type === 'warning' ? 'warning' :
                        notification.type === 'error' ? 'error' : 'info'
                      }
                      sx={{ fontSize: '0.75rem', height: 20 }}
                    />
                  </Box>
                </Box>
              </MenuItem>
            ))
          ) : (
            <Box sx={{ p: 3, textAlign: 'center' }}>
              <Typography color="text.secondary">
                No notifications
              </Typography>
            </Box>
          )}
        </Box>
        
        <Box sx={{ p: 1, borderTop: '1px solid rgba(0, 0, 0, 0.1)' }}>
          <Button 
            fullWidth 
            variant="text" 
            size="small"
            onClick={handleNotificationsClose}
          >
            Mark all as read
          </Button>
        </Box>
      </Menu>

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

    </>
  )
}
