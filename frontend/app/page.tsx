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
        return <Dashboard key="dashboard" />
      case 'employee':
        return <EmployeeManagement key="employee" />
      case 'skills':
        return <SkillManagement key="skills" />
      case 'ai-agent':
        return <AIAgentManagement key="ai-agent" />
      case 'admin':
        return <AdminDashboard key="admin" />
      default:
        return <Dashboard key="default" />
    }
  }

  return (
    <>
      {/* Navigation */}
      <Navigation activeView={activeView} onViewChange={setActiveView} />

      {/* Main Content */}
      <Box sx={{ 
        ml: { xs: 0, md: '280px' }, // Offset for desktop sidebar
        mt: { xs: 8, md: 0 }, // Offset for mobile app bar
        minHeight: '100vh',
        backgroundColor: 'background.default'
      }}>
        <Container maxWidth="xl" sx={{ py: 4 }}>
          {/* Modern Page Header */}
          <Box sx={{ mb: 4 }}>
            <Typography 
              variant="h3" 
              component="h1" 
              sx={{ 
                fontWeight: 700,
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                backgroundClip: 'text',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                mb: 1
              }}
            >
              {activeView === 'dashboard' && 'Dashboard'}
              {activeView === 'employee' && 'Employee Management'}
              {activeView === 'skills' && 'Skill Management'}
              {activeView === 'ai-agent' && 'AI Agent Management'}
              {activeView === 'admin' && 'Admin Panel'}
            </Typography>
            <Typography 
              variant="h6" 
              color="text.secondary"
              sx={{ 
                fontWeight: 400,
                maxWidth: '600px',
                lineHeight: 1.6
              }}
            >
              {activeView === 'dashboard' && 'Comprehensive overview of your engineering team performance, project metrics, and key insights'}
              {activeView === 'employee' && 'Streamlined employee lifecycle management with advanced search, filtering, and profile management capabilities'}
              {activeView === 'skills' && 'Centralized skills inventory with analytics, category management, and competency tracking'}
              {activeView === 'ai-agent' && 'Intelligent automation for employee matching, skill extraction, and request management'}
              {activeView === 'admin' && 'System administration, user management, and security controls for your organization'}
            </Typography>
          </Box>

          <Grid container spacing={4}>

          {/* Content */}
          <Grid item xs={12}>
            <Box className="fade-in">
              {renderActiveView()}
            </Box>
          </Grid>
         </Grid>
       </Container>

       {/* Modern Speed Dial */}
       <SpeedDial
         ariaLabel="Quick actions"
         sx={{ 
           position: 'fixed', 
           bottom: 24, 
           right: 24,
           '& .MuiFab-primary': {
             background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
             boxShadow: '0 8px 32px rgba(102, 126, 234, 0.3)',
             '&:hover': {
               background: 'linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%)',
               transform: 'scale(1.05)',
               boxShadow: '0 12px 40px rgba(102, 126, 234, 0.4)',
             },
             transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
           }
         }}
         icon={<SpeedDialIcon />}
       >
         {actions.map((action) => (
           <SpeedDialAction
             key={action.key}
             icon={action.icon}
             tooltipTitle={action.name}
             onClick={() => handleActionClick(action.key)}
             sx={{
               '& .MuiFab-primary': {
                 background: 'rgba(255, 255, 255, 0.9)',
                 color: '#4a5568',
                 boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
                 '&:hover': {
                   background: 'rgba(255, 255, 255, 1)',
                   transform: 'scale(1.1)',
                   boxShadow: '0 8px 20px rgba(0, 0, 0, 0.2)',
                 },
                 transition: 'all 0.2s cubic-bezier(0.4, 0, 0.2, 1)'
               }
             }}
           />
         ))}
       </SpeedDial>
      </Box>
    </>
  )
}

export default function HomePage() {
  return (
    <AuthWrapper>
      <MainApp />
    </AuthWrapper>
  )
}
