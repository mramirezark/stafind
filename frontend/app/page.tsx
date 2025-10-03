'use client'

import { useState } from 'react'
import {
  Container,
  Grid,
  Paper,
  Typography,
  Box,
  Fab,
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon,
} from '@mui/material'
import {
  Dashboard as DashboardIcon,
  PersonAdd as PersonAddIcon,
  Work as WorkIcon,
  SmartToy as AIAgentIcon,
  AdminPanelSettings as AdminIcon,
  Psychology as SkillsIcon,
} from '@mui/icons-material'
import { 
  Dashboard, 
  EmployeeManagement, 
  Navigation, 
  AuthWrapper,
  AIAgentManagement,
  SkillExtractionTool,
  AdminDashboard,
  SkillManagement
} from '@/components'

const actions = [
  { icon: <DashboardIcon />, name: 'Dashboard', key: 'dashboard' },
  { icon: <PersonAddIcon />, name: 'Add Employee', key: 'employee' },
  { icon: <SkillsIcon />, name: 'Skills', key: 'skills' },
  { icon: <AIAgentIcon />, name: 'AI Agent', key: 'ai-agent' },
]

function MainApp() {
  const [activeView, setActiveView] = useState('dashboard')

  const handleActionClick = (key: string) => {
    setActiveView(key)
  }

  const renderActiveView = () => {
    switch (activeView) {
      case 'dashboard':
        return <Dashboard />
      case 'employee':
        return <EmployeeManagement />
      case 'skills':
        return <SkillManagement />
      case 'ai-agent':
        return <AIAgentManagement />
      case 'admin':
        return <AdminDashboard />
      default:
        return <Dashboard />
    }
  }

  return (
    <Box sx={{ minHeight: '100vh', backgroundColor: 'background.default' }}>
      {/* Navigation */}
      <Navigation activeView={activeView} onViewChange={setActiveView} />

      {/* Main Content */}
      <Container maxWidth="xl" sx={{ py: 3 }}>
        <Grid container spacing={3}>
          {/* Page Header */}
          <Grid item xs={12}>
            <Paper elevation={2} sx={{ p: 3 }}>
              <Typography variant="h4" component="h1" gutterBottom>
                {activeView === 'dashboard' && 'Dashboard'}
                {activeView === 'employee' && 'Employee Management'}
                {activeView === 'skills' && 'Skill Management'}
                {activeView === 'ai-agent' && 'AI Agent Management'}
                {activeView === 'admin' && 'Admin Panel'}
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {activeView === 'dashboard' && 'Overview of your engineering team and recent activity'}
                {activeView === 'employee' && 'Add and manage employee profiles with integrated search and filtering'}
                {activeView === 'skills' && 'Manage skills, categories, and analyze skill usage across your organization'}
                {activeView === 'ai-agent' && 'Monitor AI agent requests and manage automated employee matching'}
                {activeView === 'admin' && 'Manage users, roles, and API keys for system administration'}
              </Typography>
            </Paper>
          </Grid>

          {/* Content */}
          <Grid item xs={12}>
            <Box className="fade-in">
              {renderActiveView()}
            </Box>
          </Grid>
        </Grid>
      </Container>

      {/* Speed Dial for Quick Actions */}
      <SpeedDial
        ariaLabel="Quick actions"
        sx={{ position: 'fixed', bottom: 16, right: 16 }}
        icon={<SpeedDialIcon />}
      >
        {actions.map((action) => (
          <SpeedDialAction
            key={action.key}
            icon={action.icon}
            tooltipTitle={action.name}
            onClick={() => handleActionClick(action.key)}
          />
        ))}
      </SpeedDial>
    </Box>
  )
}

export default function HomePage() {
  return (
    <AuthWrapper>
      <MainApp />
    </AuthWrapper>
  )
}
