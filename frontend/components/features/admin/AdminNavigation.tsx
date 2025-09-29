'use client'

import { useState } from 'react'
import {
  Box,
  Tabs,
  Tab,
  Typography,
  Paper,
  useTheme,
} from '@mui/material'
import {
  People as UsersIcon,
  VpnKey as APIKeysIcon,
  AdminPanelSettings as AdminIcon,
} from '@mui/icons-material'

interface AdminNavigationProps {
  activeTab: number
  onTabChange: (tab: number) => void
}

const adminTabs = [
  { 
    label: 'Users', 
    value: 0, 
    icon: <UsersIcon />,
    description: 'Manage user accounts and permissions'
  },
  { 
    label: 'API Keys', 
    value: 1, 
    icon: <APIKeysIcon />,
    description: 'Manage API keys for external services'
  },
]

export function AdminNavigation({ activeTab, onTabChange }: AdminNavigationProps) {
  const theme = useTheme()

  return (
    <Paper elevation={1} sx={{ mb: 3 }}>
      <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
          <AdminIcon color="primary" />
          <Typography variant="h6" component="h1">
            Admin Panel
          </Typography>
        </Box>
        <Typography variant="body2" color="text.secondary">
          Manage users, roles, and API keys
        </Typography>
      </Box>
      
      <Tabs
        value={activeTab}
        onChange={(_, newValue) => onTabChange(newValue)}
        variant="fullWidth"
        sx={{
          '& .MuiTab-root': {
            textTransform: 'none',
            minHeight: 64,
            gap: 1,
          },
        }}
      >
        {adminTabs.map((tab) => (
          <Tab
            key={tab.value}
            label={
              <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 0.5 }}>
                {tab.icon}
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  {tab.label}
                </Typography>
                <Typography variant="caption" color="text.secondary" sx={{ textAlign: 'center' }}>
                  {tab.description}
                </Typography>
              </Box>
            }
            value={tab.value}
          />
        ))}
      </Tabs>
    </Paper>
  )
}
