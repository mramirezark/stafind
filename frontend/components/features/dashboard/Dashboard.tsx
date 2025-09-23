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
import { useDashboardData } from '@/hooks/useApi'
import { DashboardStats, JobRequest } from '@/types'

export function Dashboard() {
  const { data: dashboardData, loading, error, refetch } = useDashboardData()

  const stats = dashboardData?.stats || null
  const recentRequests = dashboardData?.recentRequests || []

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
                  Job Requests
                </Typography>
                <Typography variant="h4">
                  {stats?.totalJobRequests || 0}
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
                  Active Matches
                </Typography>
                <Typography variant="h4">
                  {stats?.activeMatches || 0}
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
                  This Week
                </Typography>
                <Typography variant="h4">
                  {stats?.recentRequests || 0}
                </Typography>
              </Box>
            </Box>
          </CardContent>
        </Card>
      </Grid>

      {/* Recent Job Requests */}
      <Grid item xs={12}>
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Recent Job Requests
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Title</TableCell>
                  <TableCell>Department</TableCell>
                  <TableCell>Priority</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Created</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {recentRequests.length > 0 ? (
                  recentRequests.map((request) => (
                    <TableRow key={request.id}>
                      <TableCell>
                        <Typography variant="body2" fontWeight="medium">
                          {request.title}
                        </Typography>
                      </TableCell>
                      <TableCell>{request.department}</TableCell>
                      <TableCell>
                        <Chip
                          label={request.priority}
                          color={getPriorityColor(request.priority) as any}
                          size="small"
                        />
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
                    <TableCell colSpan={5} align="center">
                      <Typography color="textSecondary">
                        No job requests found
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
