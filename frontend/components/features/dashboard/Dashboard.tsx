'use client'

import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
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
  LinearProgress,
  Avatar,
  IconButton,
  Tooltip,
} from '@mui/material'
import {
  People as PeopleIcon,
  Work as WorkIcon,
  TrendingUp as TrendingUpIcon,
  Schedule as ScheduleIcon,
  MoreVert as MoreVertIcon,
  Refresh as RefreshIcon,
  TrendingDown as TrendingDownIcon,
  Star as StarIcon,
  CheckCircle as CheckCircleIcon,
  Pending as PendingIcon,
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

  const statsCards = [
    {
      title: 'Total Employees',
      value: stats?.totalEmployees || 0,
      icon: <PeopleIcon />,
      color: '#1e3a8a',
      gradient: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 100%)',
      change: '+12%',
      changeType: 'positive' as const,
      subtitle: 'Active team members'
    },
    {
      title: 'Total Requests',
      value: stats?.totalRequests || 0,
      icon: <WorkIcon />,
      color: '#10b981',
      gradient: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
      change: '+8%',
      changeType: 'positive' as const,
      subtitle: 'This month'
    },
    {
      title: 'Completed',
      value: stats?.completedRequests || 0,
      icon: <CheckCircleIcon />,
      color: '#f59e0b',
      gradient: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
      change: '+15%',
      changeType: 'positive' as const,
      subtitle: 'Success rate: 94%'
    },
    {
      title: 'Pending',
      value: stats?.pendingRequests || 0,
      icon: <PendingIcon />,
      color: '#ef4444',
      gradient: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
      change: '-5%',
      changeType: 'negative' as const,
      subtitle: 'Awaiting review'
    }
  ]

  return (
    <Grid container spacing={4}>
      {/* Modern Stats Cards */}
      {statsCards.map((stat, index) => (
        <Grid item xs={12} sm={6} md={3} key={index}>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1, duration: 0.5 }}
          >
            <Card 
              sx={{ 
                height: '100%',
                position: 'relative',
                overflow: 'hidden',
                background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
                backdropFilter: 'blur(10px)',
                border: '1px solid rgba(255,255,255,0.2)',
                borderRadius: 4,
                '&:hover': {
                  transform: 'translateY(-8px) scale(1.02)',
                  boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
                },
                transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
              }}
            >
              <CardContent sx={{ p: 3 }}>
                <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
                  <Avatar sx={{ 
                    width: 64, 
                    height: 64, 
                    background: stat.gradient,
                    color: 'white',
                    boxShadow: '0 8px 32px rgba(0, 0, 0, 0.12)',
                    border: '3px solid rgba(255, 255, 255, 0.2)',
                    transition: 'all 0.3s ease-in-out',
                    '&:hover': {
                      transform: 'scale(1.1) rotate(5deg)',
                      boxShadow: '0 12px 40px rgba(0, 0, 0, 0.2)',
                    }
                  }}>
                    {stat.icon}
                  </Avatar>
                  <Box textAlign="right">
                    <Box display="flex" alignItems="center" gap={0.5}>
                      {stat.changeType === 'positive' ? (
                        <TrendingUpIcon sx={{ color: '#10b981', fontSize: 20 }} />
                      ) : (
                        <TrendingDownIcon sx={{ color: '#ef4444', fontSize: 20 }} />
                      )}
                      <Typography 
                        variant="body2" 
                        sx={{ 
                          color: stat.changeType === 'positive' ? '#10b981' : '#ef4444',
                          fontWeight: 600
                        }}
                      >
                        {stat.change}
                      </Typography>
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                      vs last month
                    </Typography>
                  </Box>
                </Box>
                
                <Typography variant="h3" sx={{ 
                  fontWeight: 700, 
                  mb: 1,
                  background: stat.gradient,
                  backgroundClip: 'text',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}>
                  {stat.value}
                </Typography>
                
                <Typography variant="h6" color="text.primary" sx={{ fontWeight: 600, mb: 0.5 }}>
                  {stat.title}
                </Typography>
                
                <Typography variant="body2" color="text.secondary">
                  {stat.subtitle}
                </Typography>
              </CardContent>
              
              {/* Decorative gradient line */}
              <Box
                sx={{
                  position: 'absolute',
                  bottom: 0,
                  left: 0,
                  right: 0,
                  height: 4,
                  background: stat.gradient,
                }}
              />
            </Card>
          </motion.div>
        </Grid>
      ))}

      {/* Most Requested Skills */}
      <Grid item xs={12} md={6}>
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.4, duration: 0.5 }}
        >
          <Card sx={{ 
            height: '100%',
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 4,
            boxShadow: '0 10px 40px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 4 }}>
              <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Most Requested Skills
                </Typography>
                <IconButton size="small">
                  <MoreVertIcon />
                </IconButton>
              </Box>
              
              {skillDemandLoading ? (
                <Box display="flex" justifyContent="center" p={4}>
                  <CircularProgress />
                </Box>
              ) : mostRequestedSkills.length > 0 ? (
                <Box>
                  {mostRequestedSkills.slice(0, 5).map((skill, index) => (
                    <motion.div
                      key={skill.skill_name}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: 0.5 + index * 0.1, duration: 0.3 }}
                    >
                      <Box 
                        sx={{ 
                          p: 2, 
                          mb: 2, 
                          borderRadius: 2,
                          border: '1px solid rgba(226, 232, 240, 0.8)',
                          transition: 'all 0.2s ease-in-out',
                            '&:hover': {
                              borderColor: '#1e3a8a',
                              boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                            }
                        }}
                      >
                        <Box display="flex" justifyContent="space-between" alignItems="center">
                          <Box display="flex" alignItems="center" gap={2}>
                            <Avatar sx={{ 
                              width: 32, 
                              height: 32, 
                              bgcolor: 'rgba(30, 58, 138, 0.1)',
                              color: '#1e3a8a',
                              fontWeight: 700
                            }}>
                              {index + 1}
                            </Avatar>
                            <Box>
                              <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                                {skill.skill_name}
                              </Typography>
                              <Chip 
                                label={skill.category} 
                                size="small" 
                                variant="outlined"
                                sx={{ 
                                  mt: 0.5,
                                  fontSize: '0.75rem',
                                  height: 20
                                }}
                              />
                            </Box>
                          </Box>
                          <Box textAlign="right">
                            <Typography variant="h6" sx={{ fontWeight: 700, color: '#1e3a8a' }}>
                              {skill.count}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              requests
                            </Typography>
                          </Box>
                        </Box>
                        <LinearProgress 
                          variant="determinate" 
                          value={(skill.count / (mostRequestedSkills[0]?.count || 1)) * 100}
                          sx={{ 
                            mt: 1,
                            height: 4,
                            borderRadius: 2,
                            bgcolor: 'rgba(30, 58, 138, 0.1)',
                            '& .MuiLinearProgress-bar': {
                              borderRadius: 2,
                              background: 'linear-gradient(135deg, #1e3a8a 0%, #3730a3 100%)',
                            }
                          }}
                        />
                      </Box>
                    </motion.div>
                  ))}
                </Box>
              ) : (
                <Box textAlign="center" py={4}>
                  <Typography color="textSecondary">
                    No skill data available
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </motion.div>
      </Grid>

      {/* Top Suggested Employees */}
      <Grid item xs={12} md={6}>
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.4, duration: 0.5 }}
        >
          <Card sx={{ 
            height: '100%',
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 4,
            boxShadow: '0 10px 40px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 4 }}>
              <Box display="flex" alignItems="center" justifyContent="space-between" mb={3}>
                <Typography variant="h6" sx={{ fontWeight: 700 }}>
                  Top Suggested Employees
                </Typography>
                <IconButton size="small">
                  <MoreVertIcon />
                </IconButton>
              </Box>
              
              {topEmployeesLoading ? (
                <Box display="flex" justifyContent="center" p={4}>
                  <CircularProgress />
                </Box>
              ) : topSuggestedEmployees.length > 0 ? (
                <Box>
                  {topSuggestedEmployees.slice(0, 5).map((employee, index) => (
                    <motion.div
                      key={employee.employee_id}
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: 0.5 + index * 0.1, duration: 0.3 }}
                    >
                      <Box 
                        sx={{ 
                          p: 2, 
                          mb: 2, 
                          borderRadius: 2,
                          border: '1px solid rgba(226, 232, 240, 0.8)',
                          transition: 'all 0.2s ease-in-out',
                          '&:hover': {
                            borderColor: '#10b981',
                            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                          }
                        }}
                      >
                        <Box display="flex" justifyContent="space-between" alignItems="center">
                          <Box display="flex" alignItems="center" gap={2}>
                            <Avatar sx={{ 
                              width: 40, 
                              height: 40, 
                              bgcolor: 'rgba(16, 185, 129, 0.1)',
                              color: '#10b981',
                              fontWeight: 700
                            }}>
                              {employee.employee_name.charAt(0)}
                            </Avatar>
                            <Box>
                              <Typography variant="subtitle2" sx={{ fontWeight: 600 }}>
                                {employee.employee_name}
                              </Typography>
                              <Typography variant="caption" color="text.secondary">
                                {employee.department} • {employee.level}
                              </Typography>
                            </Box>
                          </Box>
                          <Box textAlign="right">
                            <Box display="flex" alignItems="center" gap={0.5}>
                              <StarIcon sx={{ color: '#f59e0b', fontSize: 16 }} />
                              <Typography variant="subtitle2" sx={{ fontWeight: 600, color: '#10b981' }}>
                                {employee.avg_match_score.toFixed(1)}
                              </Typography>
                            </Box>
                            <Typography variant="caption" color="text.secondary">
                              {employee.match_count} matches
                            </Typography>
                          </Box>
                        </Box>
                      </Box>
                    </motion.div>
                  ))}
                </Box>
              ) : (
                <Box textAlign="center" py={4}>
                  <Typography color="textSecondary">
                    No employee data available
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </motion.div>
      </Grid>

      {/* Recent AI Agent Requests */}
      <Grid item xs={12}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.5 }}
        >
          <Card sx={{
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 4,
            boxShadow: '0 10px 40px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-2px)',
              boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 0 }}>
              <Box sx={{ 
                p: 4, 
                borderBottom: '1px solid rgba(226, 232, 240, 0.3)',
                background: 'linear-gradient(135deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0.02) 100%)'
              }}>
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Typography variant="h6" sx={{ fontWeight: 700 }}>
                    Recent AI Agent Requests
                  </Typography>
                  <Box display="flex" gap={1}>
                    <Tooltip title="Refresh">
                      <IconButton size="small" onClick={refetch}>
                        <RefreshIcon />
                      </IconButton>
                    </Tooltip>
                    <IconButton size="small">
                      <MoreVertIcon />
                    </IconButton>
                  </Box>
                </Box>
              </Box>
              
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow sx={{ bgcolor: 'rgba(30, 58, 138, 0.02)' }}>
                      <TableCell sx={{ fontWeight: 600, color: 'text.primary' }}>User</TableCell>
                      <TableCell sx={{ fontWeight: 600, color: 'text.primary' }}>Message</TableCell>
                      <TableCell sx={{ fontWeight: 600, color: 'text.primary' }}>Status</TableCell>
                      <TableCell sx={{ fontWeight: 600, color: 'text.primary' }}>Created</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {recentRequests.length > 0 ? (
                      recentRequests.map((request, index) => (
                        <TableRow 
                          key={request.id}
                          component={motion.tr}
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          transition={{ delay: 0.7 + index * 0.05, duration: 0.3 }}
                          sx={{ 
                            '&:hover': {
                              bgcolor: 'rgba(30, 58, 138, 0.02)',
                            },
                            transition: 'all 0.2s ease-in-out'
                          }}
                        >
                            <TableCell>
                              <Box display="flex" alignItems="center" gap={2}>
                                <Avatar sx={{ 
                                  width: 32, 
                                  height: 32, 
                                  bgcolor: 'rgba(30, 58, 138, 0.1)',
                                  color: '#1e3a8a',
                                  fontSize: '0.875rem',
                                  fontWeight: 600
                                }}>
                                  {request.user_name.charAt(0)}
                                </Avatar>
                                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                  {request.user_name}
                                </Typography>
                              </Box>
                            </TableCell>
                            <TableCell>
                              <Typography 
                                variant="body2" 
                                noWrap 
                                sx={{ 
                                  maxWidth: 200,
                                  color: 'text.secondary'
                                }}
                              >
                                {request.message_text || 'No message text'}
                              </Typography>
                            </TableCell>
                            <TableCell>
                              <Chip
                                label={request.status}
                                color={getStatusColor(request.status) as any}
                                size="small"
                                sx={{ 
                                  fontWeight: 600,
                                  textTransform: 'capitalize'
                                }}
                              />
                            </TableCell>
                            <TableCell>
                              <Typography variant="body2" color="text.secondary">
                                {new Date(request.created_at).toLocaleDateString()}
                              </Typography>
                            </TableCell>
                        </TableRow>
                      ))
                    ) : (
                      <TableRow>
                        <TableCell colSpan={4} align="center" sx={{ py: 4 }}>
                          <Box textAlign="center">
                            <Typography color="textSecondary" variant="body1">
                              No recent requests found
                            </Typography>
                            <Typography color="textSecondary" variant="body2" sx={{ mt: 1 }}>
                              AI agent requests will appear here
                            </Typography>
                          </Box>
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </motion.div>
      </Grid>
    </Grid>
  )
}
