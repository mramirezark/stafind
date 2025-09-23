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
} from '@mui/icons-material'
import { 
  Dashboard, 
  JobRequestForm, 
  EmployeeManagement, 
  Navigation, 
  AuthWrapper 
} from '@/components'

const actions = [
  { icon: <DashboardIcon />, name: 'Dashboard', key: 'dashboard' },
  { icon: <WorkIcon />, name: 'Job Request', key: 'job-request' },
  { icon: <PersonAddIcon />, name: 'Add Employee', key: 'employee' },
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
      case 'job-request':
        return <JobRequestForm />
      case 'employee':
        return <EmployeeManagement />
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
                {activeView === 'job-request' && 'Create Job Request'}
                {activeView === 'employee' && 'Employee Management'}
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {activeView === 'dashboard' && 'Overview of your engineering team and recent activity'}
                {activeView === 'job-request' && 'Submit a new job request with skill requirements'}
                {activeView === 'employee' && 'Add and manage employee profiles with integrated search and filtering'}
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
