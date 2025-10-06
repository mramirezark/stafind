'use client'

import React, { useState, useEffect } from 'react'
import {
  Box,
  Tabs,
  Tab,
  Paper,
  Typography,
  Alert,
  CircularProgress,
  Fab,
  Tooltip,
} from '@mui/material'
import {
  Add as AddIcon,
  Analytics as AnalyticsIcon,
  Category as CategoryIcon,
  List as ListIcon,
} from '@mui/icons-material'

import { Skill, Category, SkillStats } from '@/types'
import { endpoints } from '@/lib/api'
import SkillList from './SkillList'
import SkillForm from './SkillForm'
import CategoryManagement from './CategoryManagement'
import SkillAnalytics from './SkillAnalytics'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel({ children, value, index, ...other }: TabPanelProps) {
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`skill-tabpanel-${index}`}
      aria-labelledby={`skill-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  )
}

const SkillManagement: React.FC = () => {
  const [activeTab, setActiveTab] = useState(0)
  const [skills, setSkills] = useState<Skill[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [stats, setStats] = useState<SkillStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formOpen, setFormOpen] = useState(false)
  const [editingSkill, setEditingSkill] = useState<Skill | null>(null)

  const loadSkills = async () => {
    try {
      setLoading(true)
      const response = await endpoints.skills.list()
      setSkills(response.data || [])
    } catch (err) {
      setError('Failed to load skills')
      console.error('Error loading skills:', err)
    } finally {
      setLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const response = await endpoints.categories.list()
      setCategories(response.data || [])
    } catch (err) {
      console.error('Error loading categories:', err)
    }
  }

  const loadStats = async () => {
    try {
      const response = await endpoints.skills.stats()
      setStats(response.data)
    } catch (err) {
      console.error('Error loading stats:', err)
    }
  }

  useEffect(() => {
    loadSkills()
    loadCategories()
    loadStats()
  }, [])

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue)
  }

  const handleCreateSkill = () => {
    setEditingSkill(null)
    setFormOpen(true)
  }

  const handleEditSkill = (skill: Skill) => {
    setEditingSkill(skill)
    setFormOpen(true)
  }

  const handleFormClose = () => {
    setFormOpen(false)
    setEditingSkill(null)
  }

  const handleFormSubmit = async (skillData: any) => {
    try {
      if (editingSkill) {
        await endpoints.skills.update(editingSkill.id, skillData)
      } else {
        await endpoints.skills.create(skillData)
      }
      await loadSkills()
      await loadStats()
      setFormOpen(false)
      setEditingSkill(null)
    } catch (err) {
      console.error('Error saving skill:', err)
      throw err
    }
  }

  const handleDeleteSkill = async (id: number) => {
    try {
      await endpoints.skills.delete(id)
      await loadSkills()
      await loadStats()
    } catch (err) {
      console.error('Error deleting skill:', err)
      throw err
    }
  }

  const handleRefresh = () => {
    loadSkills()
    loadCategories()
    loadStats()
  }

  if (loading && skills.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Skill Management
        </Typography>
        <Tabs value={activeTab} onChange={handleTabChange} aria-label="skill management tabs">
          <Tab
            icon={<ListIcon />}
            label="Skills"
            id="skill-tab-0"
            aria-controls="skill-tabpanel-0"
          />
          <Tab
            icon={<CategoryIcon />}
            label="Categories"
            id="skill-tab-1"
            aria-controls="skill-tabpanel-1"
          />
          <Tab
            icon={<AnalyticsIcon />}
            label="Analytics"
            id="skill-tab-2"
            aria-controls="skill-tabpanel-2"
          />
        </Tabs>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <TabPanel value={activeTab} index={0}>
        <SkillList
          skills={skills}
          categories={categories}
          onEdit={handleEditSkill}
          onDelete={handleDeleteSkill}
          onRefresh={handleRefresh}
          onCreate={handleCreateSkill}
          loading={loading}
        />
      </TabPanel>

      <TabPanel value={activeTab} index={1}>
        <CategoryManagement
          categories={categories}
          onRefresh={loadCategories}
        />
      </TabPanel>

      <TabPanel value={activeTab} index={2}>
        <SkillAnalytics
          stats={stats}
          skills={skills}
          categories={categories}
          loading={loading}
        />
      </TabPanel>


      {/* Floating Action Button for creating skills */}
      {activeTab === 0 && (
        <Tooltip title="Add New Skill">
          <Fab
            color="primary"
            aria-label="add skill"
            sx={{
              position: 'fixed',
              bottom: 16,
              right: 16,
              zIndex: 1000,
              boxShadow: '0 4px 8px rgba(0,0,0,0.3)',
              '&:hover': {
                boxShadow: '0 6px 12px rgba(0,0,0,0.4)',
              }
            }}
            onClick={handleCreateSkill}
          >
            <AddIcon />
          </Fab>
        </Tooltip>
      )}

      {/* Skill Form Dialog */}
      <SkillForm
        open={formOpen}
        onClose={handleFormClose}
        onSubmit={handleFormSubmit}
        initialData={editingSkill}
        categories={categories}
        title={editingSkill ? 'Edit Skill' : 'Create New Skill'}
      />
    </Box>
  )
}

export default SkillManagement
