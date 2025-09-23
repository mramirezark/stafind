'use client'

import { useState } from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  Tabs,
  Tab,
  Divider,
} from '@mui/material'
import { Login as LoginIcon, PersonAdd as PersonAddIcon } from '@mui/icons-material'
import { LoginForm } from './LoginForm'
import { RegisterForm } from './RegisterForm'
import { AuthCardProps, TabPanelProps } from '@/types'

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`auth-tabpanel-${index}`}
      aria-labelledby={`auth-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ pt: 3 }}>
          {children}
        </Box>
      )}
    </div>
  )
}

export function AuthCard({ onSuccess }: AuthCardProps) {
  const [tabValue, setTabValue] = useState(0)

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  const handleToggleMode = () => {
    setTabValue(tabValue === 0 ? 1 : 0)
  }

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        p: 2,
      }}
    >
      <Card sx={{ maxWidth: 500, width: '100%', boxShadow: 3 }}>
        <CardContent sx={{ p: 4 }}>
          <Box textAlign="center" mb={3}>
            <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
              StaffFind
            </Typography>
            <Typography variant="body1" color="text.secondary">
              Welcome to StaffFind - Your Employee Management Platform
            </Typography>
          </Box>

          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs value={tabValue} onChange={handleTabChange} aria-label="auth tabs">
              <Tab
                icon={<LoginIcon />}
                label="Sign In"
                id="auth-tab-0"
                aria-controls="auth-tabpanel-0"
              />
              <Tab
                icon={<PersonAddIcon />}
                label="Sign Up"
                id="auth-tab-1"
                aria-controls="auth-tabpanel-1"
              />
            </Tabs>
          </Box>

          <TabPanel value={tabValue} index={0}>
            <LoginForm onSuccess={onSuccess} />
          </TabPanel>

          <TabPanel value={tabValue} index={1}>
            <RegisterForm onSuccess={onSuccess} onToggleMode={handleToggleMode} />
          </TabPanel>
        </CardContent>
      </Card>
    </Box>
  )
}
