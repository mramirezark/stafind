'use client'

import { useState, useEffect } from 'react'
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  CircularProgress,
  Alert,
} from '@mui/material'
import {
  People as PeopleIcon,
  Work as WorkIcon,
  TrendingUp as TrendingUpIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material'
import { useDashboardMetrics, useTopSuggestedEmployees, useSkillDemandStats } from '@/hooks/useApi'
import { DashboardStats, TopSuggestedEmployee, SkillDemandStats } from '@/types'

export function Dashboard() {
  const { data: dashboardMetrics, loading, error, refetch } = useDashboardMetrics()
  const { data: topEmployees, loading: topEmployeesLoading } = useTopSuggestedEmployees(5)
  const { data: skillDemand, loading: skillDemandLoading } = useSkillDemandStats()

  const stats = dashboardMetrics?.stats || null
  const mostRequestedSkills = dashboardMetrics?.most_requested_skills || skillDemand || []
  const topSuggestedEmployees = dashboardMetrics?.top_suggested_employees || topEmployees || []
  const recentRequests = dashboardMetrics?.recent_requests || []


  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case 'high':
        return 'error'
      case 'medium':
        return 'warning'
      case 'low':
        return 'success'
      default:
        return 'default'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'open':
        return 'success'
      case 'in-progress':
        return 'warning'
      case 'closed':
        return 'default'
      default:
        return 'default'
    }
  }

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    )
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    )
  }

  return (
    <Grid container spacing={3}>
      {/* Stats Cards */}
      <Grid item xs={12} sm={6} md={3}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center">
              <PeopleIcon color="primary" sx={{ mr: 2, fontSize: 40 }} />
              <Box>
                <Typography color="textSecondary" gutterBottom>
                  Total Employees
                </Typography>
                <Typography variant="h4">
                  {stats?.totalEmployees || 0}
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>

      <Grid item xs={12} sm={6} md={3}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center">
              <WorkIcon color="secondary" sx={{ mr: 2, fontSize: 40 }} />
              <Box>
                <Typography color="textSecondary" gutterBottom>
                  Total Requests
                </Typography>
                <Typography variant="h4">
                  {stats?.totalRequests || 0}
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>

      <Grid item xs={12} sm={6} md={3}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center">
              <TrendingUpIcon color="success" sx={{ mr: 2, fontSize: 40 }} />
              <Box>
                <Typography color="textSecondary" gutterBottom>
                  Completed
                </Typography>
                <Typography variant="h4">
                  {stats?.completedRequests || 0}
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>

      <Grid item xs={12} sm={6} md={3}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center">
              <ScheduleIcon color="warning" sx={{ mr: 2, fontSize: 40 }} />
              <Box>
                <Typography color="textSecondary" gutterBottom>
                  Pending
                </Typography>
                <Typography variant="h4">
                  {stats?.pendingRequests || 0}
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>

      {/* Most Requested Skills */}
      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Most Requested Skills
          </Typography>
          {skillDemandLoading ? (
            <Box display="flex" justifyContent="center" p={2}>
              <CircularProgress />
            </Box>
          ) : mostRequestedSkills.length > 0 ? (
            <Box>
              {mostRequestedSkills.slice(0, 5).map((skill, index) => (
                <Box key={skill.skill_name} display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                  <Box display="flex" alignItems="center">
                    <Typography variant="body2" sx={{ mr: 1, fontWeight: 'bold' }}>
                      {index + 1}.
                    </Typography>
                    <Typography variant="body2">
                      {skill.skill_name}
                    </Typography>
                    <Chip 
                      label={skill.category} 
                      size="small" 
                      variant="outlined" 
                      sx={{ ml: 1 }}
                    />
                  </Box>
                  <Typography variant="body2" color="text.secondary">
                    {skill.count} requests
                  </Typography>
                </Box>
              ))}
            </Box>
          ) : (
            <Typography color="textSecondary">
              No skill data available
            </Typography>
          )}
        </Paper>
      </Grid>

      {/* Top Suggested Employees */}
      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Top Suggested Employees
          </Typography>
          {topEmployeesLoading ? (
            <Box display="flex" justifyContent="center" p={2}>
              <CircularProgress />
            </Box>
          ) : topSuggestedEmployees.length > 0 ? (
            <Box>
              {topSuggestedEmployees.slice(0, 5).map((employee, index) => (
                <Box key={employee.employee_id} display="flex" justifyContent="space-between" alignItems="center" mb={1}>
                  <Box display="flex" alignItems="center">
                    <Typography variant="body2" sx={{ mr: 1, fontWeight: 'bold' }}>
                      {index + 1}.
                    </Typography>
                    <Box>
                      <Typography variant="body2" fontWeight="medium">
                        {employee.employee_name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {employee.department} • {employee.level}
                      </Typography>
                    </Box>
                  </Box>
                  <Box textAlign="right">
                    <Typography variant="body2" color="text.secondary">
                      {employee.match_count} matches
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Avg: {employee.avg_match_score.toFixed(1)}
                    </Typography>
                  </Box>
                </Box>
              ))}
            </Box>
          ) : (
            <Typography color="textSecondary">
              No employee data available
            </Typography>
          )}
        </Paper>
      </Grid>

      {/* Recent AI Agent Requests */}
      <Grid item xs={12}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Recent AI Agent Requests
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>User</TableCell>
                  <TableCell>Message</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Created</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {recentRequests.length > 0 ? (
                  recentRequests.map((request) => (
                    <TableRow key={request.id}>
                      <TableCell>{request.user_name}</TableCell>
                      <TableCell>
                        <Typography variant="body2" noWrap sx={{ maxWidth: 200 }}>
                          {request.message_text || 'No message text'}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={request.status}
                          color={getStatusColor(request.status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        {new Date(request.created_at).toLocaleDateString()}
                      </TableCell>
                    </TableRow>
                  ))
                ) : (
                  <TableRow>
                    <TableCell colSpan={4} align="center">
                      <Typography color="textSecondary">
                        No recent requests found
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      </Grid>
    </Grid>
  )
}
