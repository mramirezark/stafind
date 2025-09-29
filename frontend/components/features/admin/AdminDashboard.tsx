'use client'

import { useState } from 'react'
import { Box, Container } from '@mui/material'
import { AdminNavigation } from './AdminNavigation'
import { UserManagement } from './UserManagement'
import { APIKeyManagement } from './APIKeyManagement'

export function AdminDashboard() {
  const [activeTab, setActiveTab] = useState(0)

  const handleTabChange = (tab: number) => {
    setActiveTab(tab)
  }

  const renderTabContent = () => {
    switch (activeTab) {
      case 0:
        return <UserManagement />
      case 1:
        return <APIKeyManagement />
      default:
        return <UserManagement />
    }
  }

  return (
    <Container maxWidth="xl" sx={{ py: 3 }}>
      <AdminNavigation activeTab={activeTab} onTabChange={handleTabChange} />
      {renderTabContent()}
    </Container>
  )
}
