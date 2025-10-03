'use client'

import React from 'react'
import {
  Box,
  Grid,
  Paper,
  Typography,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Chip,
  LinearProgress,
  Alert,
  CircularProgress,
  Divider,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  Category as CategoryIcon,
  People as PeopleIcon,
  Star as StarIcon,
  Timeline as TimelineIcon,
} from '@mui/icons-material'

import { Skill, Category, SkillStats } from '@/types'

interface SkillAnalyticsProps {
  stats: SkillStats | null
  skills: Skill[]
  categories: Category[]
  loading?: boolean
}

const SkillAnalytics: React.FC<SkillAnalyticsProps> = ({
  stats,
  skills,
  categories,
  loading = false,
}) => {
  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  if (!stats) {
    return (
      <Alert severity="info">
        No analytics data available. Skills and categories need to be created first.
      </Alert>
    )
  }

  const topSkills = stats.most_popular_skills || []
  const skillsByCategory = stats.skills_by_category || []
  const recentSkills = stats.recent_skills || []

  return (
    <Box>
      {/* Overview Cards */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <CategoryIcon color="primary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="h4" component="div">
                    {stats.total_skills}
                  </Typography>
                  <Typography color="text.secondary">
                    Total Skills
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <CategoryIcon color="secondary" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="h4" component="div">
                    {stats.total_categories}
                  </Typography>
                  <Typography color="text.secondary">
                    Categories
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <PeopleIcon color="success" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="h4" component="div">
                    {skills.reduce((sum, skill) => sum + (skill.employee_count || 0), 0)}
                  </Typography>
                  <Typography color="text.secondary">
                    Total Assignments
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <TrendingUpIcon color="warning" sx={{ mr: 2 }} />
                <Box>
                  <Typography variant="h4" component="div">
                    {topSkills.length > 0 ? topSkills[0].employee_count || 0 : 0}
                  </Typography>
                  <Typography color="text.secondary">
                    Most Popular
                  </Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Most Popular Skills */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, height: '100%' }}>
            <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
              <StarIcon sx={{ mr: 1 }} />
              Most Popular Skills
            </Typography>
            <Divider sx={{ mb: 2 }} />
            
            {topSkills.length === 0 ? (
              <Typography color="text.secondary">
                No skill usage data available
              </Typography>
            ) : (
              <List>
                {topSkills.slice(0, 10).map((skill, index) => (
                  <ListItem key={skill.id} sx={{ px: 0 }}>
                    <ListItemIcon>
                      <Typography variant="body2" color="text.secondary">
                        #{index + 1}
                      </Typography>
                    </ListItemIcon>
                    <ListItemText
                      primary={skill.name}
                      secondary={`${skill.employee_count || 0} employees`}
                    />
                    <Box sx={{ width: 100, ml: 2 }}>
                      <LinearProgress
                        variant="determinate"
                        value={Math.min(
                          ((skill.employee_count || 0) / (topSkills[0]?.employee_count || 1)) * 100,
                          100
                        )}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                    </Box>
                  </ListItem>
                ))}
              </List>
            )}
          </Paper>
        </Grid>

        {/* Skills by Category */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, height: '100%' }}>
            <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
              <CategoryIcon sx={{ mr: 1 }} />
              Skills by Category
            </Typography>
            <Divider sx={{ mb: 2 }} />
            
            {skillsByCategory.length === 0 ? (
              <Typography color="text.secondary">
                No category data available
              </Typography>
            ) : (
              <List>
                {skillsByCategory.map((item, index) => (
                  <ListItem key={index} sx={{ px: 0 }}>
                    <ListItemText
                      primary={item.category}
                      secondary={`${item.count} skills`}
                    />
                    <Box sx={{ width: 100, ml: 2 }}>
                      <LinearProgress
                        variant="determinate"
                        value={Math.min(
                          (item.count / Math.max(...skillsByCategory.map(c => c.count))) * 100,
                          100
                        )}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                    </Box>
                  </ListItem>
                ))}
              </List>
            )}
          </Paper>
        </Grid>

        {/* Recent Skills */}
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
              <TimelineIcon sx={{ mr: 1 }} />
              Recently Added Skills
            </Typography>
            <Divider sx={{ mb: 2 }} />
            
            {recentSkills.length === 0 ? (
              <Typography color="text.secondary">
                No recent skills data available
              </Typography>
            ) : (
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                {recentSkills.slice(0, 20).map(skill => (
                  <Chip
                    key={skill.id}
                    label={skill.name}
                    variant="outlined"
                    size="small"
                    icon={<CategoryIcon />}
                  />
                ))}
              </Box>
            )}
          </Paper>
        </Grid>

        {/* Category Distribution */}
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Category Distribution
            </Typography>
            <Divider sx={{ mb: 2 }} />
            
            {categories.length === 0 ? (
              <Typography color="text.secondary">
                No categories available
              </Typography>
            ) : (
              <Grid container spacing={2}>
                {categories.map(category => {
                  const categorySkills = skills.filter(skill =>
                    skill.categories?.some(cat => cat.id === category.id)
                  )
                  const percentage = skills.length > 0 
                    ? (categorySkills.length / skills.length) * 100 
                    : 0

                  return (
                    <Grid item xs={12} sm={6} md={4} key={category.id}>
                      <Card variant="outlined">
                        <CardContent>
                          <Typography variant="h6" gutterBottom>
                            {category.name}
                          </Typography>
                          <Typography color="text.secondary" gutterBottom>
                            {categorySkills.length} skills
                          </Typography>
                          <LinearProgress
                            variant="determinate"
                            value={percentage}
                            sx={{ height: 8, borderRadius: 4 }}
                          />
                          <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                            {percentage.toFixed(1)}% of all skills
                          </Typography>
                        </CardContent>
                      </Card>
                    </Grid>
                  )
                })}
              </Grid>
            )}
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}

export default SkillAnalytics
