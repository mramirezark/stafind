import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Grid,
  Chip,
  Divider,
  Avatar,
  LinearProgress,
  Tabs,
  Tab,
} from '@mui/material';
import {
  Visibility as ViewIcon,
  Error as ErrorIcon,
  CheckCircle as SuccessIcon,
  Schedule as PendingIcon,
} from '@mui/icons-material';
import { CircularProgress } from '@mui/material';
import { AIAgentRequest } from '../../../types';
import { AIAgentResponse } from '../../../services/ai/aiAgentService';

interface RequestDetailsModalProps {
  open: boolean;
  onClose: () => void;
  selectedRequest: AIAgentRequest | null;
  response: AIAgentResponse | null;
}

const RequestDetailsModal: React.FC<RequestDetailsModalProps> = ({
  open,
  onClose,
  selectedRequest,
  response,
}) => {
  const [activeTab, setActiveTab] = React.useState(0);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <SuccessIcon color="success" />;
      case 'failed':
        return <ErrorIcon color="error" />;
      case 'processing':
        return <CircularProgress size={20} />;
      default:
        return <PendingIcon color="action" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'success';
      case 'failed':
        return 'error';
      case 'processing':
        return 'warning';
      default:
        return 'default';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const renderOverviewTab = () => (
    <Grid container spacing={3}>
      {/* Basic Information Card */}
      <Grid item xs={12} md={6}>
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.1 }}
        >
          <Card sx={{
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 3,
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 12px 40px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 3 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                <Avatar sx={{ 
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  color: 'white',
                  width: 40,
                  height: 40
                }}>
                  <ViewIcon />
                </Avatar>
                <Typography variant="h6" sx={{ fontWeight: 700, color: '#1e293b' }}>
                  Basic Information
                </Typography>
              </Box>
              
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', p: 2, background: 'rgba(102, 126, 234, 0.05)', borderRadius: 2 }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, color: '#64748b' }}>Request ID</Typography>
                  <Typography variant="body2" sx={{ fontWeight: 700, color: '#1e293b' }}>#{selectedRequest?.id}</Typography>
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', p: 2, background: 'rgba(102, 126, 234, 0.05)', borderRadius: 2 }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, color: '#64748b' }}>Status</Typography>
                  <Chip
                    icon={getStatusIcon(selectedRequest?.status || '')}
                    label={selectedRequest?.status}
                    color={getStatusColor(selectedRequest?.status || '') as any}
                    size="small"
                    sx={{ fontWeight: 600 }}
                  />
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', p: 2, background: 'rgba(102, 126, 234, 0.05)', borderRadius: 2 }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, color: '#64748b' }}>User</Typography>
                  <Typography variant="body2" sx={{ fontWeight: 700, color: '#1e293b' }}>{selectedRequest?.user_name}</Typography>
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', p: 2, background: 'rgba(102, 126, 234, 0.05)', borderRadius: 2 }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, color: '#64748b' }}>Created</Typography>
                  <Typography variant="body2" sx={{ fontWeight: 700, color: '#1e293b' }}>{formatDate(selectedRequest?.created_at || '')}</Typography>
                </Box>
              </Box>
            </CardContent>
          </Card>
        </motion.div>
      </Grid>

      {/* Message Card */}
      <Grid item xs={12} md={6}>
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
        >
          <Card sx={{
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 3,
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 12px 40px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 3 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                <Avatar sx={{ 
                  background: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)',
                  color: 'white',
                  width: 40,
                  height: 40
                }}>
                  💬
                </Avatar>
                <Typography variant="h6" sx={{ fontWeight: 700, color: '#1e293b' }}>
                  Message Content
                </Typography>
              </Box>
              
              <Typography variant="body2" sx={{ 
                color: '#64748b',
                lineHeight: 1.6,
                whiteSpace: 'pre-wrap',
                p: 2,
                background: 'rgba(245, 158, 11, 0.05)',
                borderRadius: 2,
                border: '1px solid rgba(245, 158, 11, 0.1)'
              }}>
                {selectedRequest?.message_text || 'No message text provided'}
              </Typography>
            </CardContent>
          </Card>
        </motion.div>
      </Grid>
    </Grid>
  );


  const renderAIResponseTab = () => (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <Card sx={{
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 3,
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 12px 40px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 4 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
                <Avatar sx={{ 
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  color: 'white',
                  width: 48,
                  height: 48
                }}>
                  <SuccessIcon />
                </Avatar>
                <Typography variant="h5" sx={{ fontWeight: 700, color: '#1e293b' }}>
                  AI Response
                </Typography>
              </Box>
              
              {response?.summary && (
                <Box sx={{ 
                  p: 3, 
                  background: 'rgba(102, 126, 234, 0.05)', 
                  borderRadius: 2,
                  border: '1px solid rgba(102, 126, 234, 0.1)',
                  mb: 4
                }}>
                  <Typography variant="body1" sx={{ 
                    whiteSpace: 'pre-wrap',
                    lineHeight: 1.6,
                    color: '#1e293b'
                  }}>
                    {response.summary}
                  </Typography>
                </Box>
              )}

              <Box sx={{ 
                p: 3, 
                background: 'rgba(16, 185, 129, 0.05)', 
                borderRadius: 2,
                border: '1px solid rgba(16, 185, 129, 0.1)'
              }}>
                <Typography variant="body2" sx={{ 
                  whiteSpace: 'pre-wrap',
                  lineHeight: 1.6,
                  color: '#1e293b',
                  fontWeight: 600
                }}>
                  <strong>Processing Status:</strong> {response?.status} • <strong>Time:</strong> {response?.processing_time_ms}ms
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </motion.div>
      </Grid>
    </Grid>
  );

  const renderMatchesTab = () => (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <Card sx={{
            background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(255,255,255,0.2)',
            borderRadius: 3,
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
            transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            '&:hover': {
              transform: 'translateY(-4px)',
              boxShadow: '0 12px 40px rgba(0, 0, 0, 0.15)',
            }
          }}>
            <CardContent sx={{ p: 4 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 4 }}>
                <Avatar sx={{ 
                  background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                  color: 'white',
                  width: 48,
                  height: 48,
                  fontSize: '0.875rem'
                }}>
                  {response?.matches?.length || 0}
                </Avatar>
                <Typography variant="h5" sx={{ fontWeight: 700, color: '#1e293b' }}>
                  Top Candidate Matches
                </Typography>
              </Box>
              
              {response?.matches && response.matches.length > 0 && (
                <Grid container spacing={3}>
                  {response.matches.map((match: any, index: number) => (
                    <Grid item xs={12} key={index}>
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.2 + index * 0.1 }}
                      >
                        <Card sx={{
                          background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.05) 0%, rgba(16, 185, 129, 0.02) 100%)',
                          border: '1px solid rgba(16, 185, 129, 0.1)',
                          borderRadius: 3,
                          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                          '&:hover': {
                            transform: 'translateY(-2px)',
                            boxShadow: '0 8px 25px rgba(16, 185, 129, 0.15)',
                          }
                        }}>
                          <CardContent sx={{ p: 3 }}>
                            <Grid container spacing={3} alignItems="flex-start">
                              <Grid item xs={12} md={8}>
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                                  <Avatar sx={{ 
                                    background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                                    color: 'white',
                                    width: 40,
                                    height: 40
                                  }}>
                                    {match.employee_name.charAt(0)}
                                  </Avatar>
                                  <Box>
                                    <Typography variant="h6" sx={{ fontWeight: 700, color: '#059669', mb: 0.5 }}>
                                      {match.employee_name}
                                    </Typography>
                                    <Typography variant="body2" sx={{ color: '#64748b', fontWeight: 500 }}>
                                      {match.position} • {match.seniority} • {match.location}
                                    </Typography>
                                  </Box>
                                </Box>
                                
                                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 2 }}>
                                  <Typography variant="body2" sx={{ mb: 0.5 }}>
                                    <strong>Email:</strong> {match.employee_email}
                                  </Typography>
                                  {match.current_project && (
                                    <Typography variant="body2" sx={{ mb: 0.5 }}>
                                      <strong>Current Project:</strong> {match.current_project}
                                    </Typography>
                                  )}
                                </Box>
                                
                                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                                  <Typography variant="body2" sx={{ fontWeight: 600, color: '#1e293b' }}>
                                    Match Score:
                                  </Typography>
                                  <Box sx={{ flexGrow: 1 }}>
                                    <LinearProgress
                                      variant="determinate"
                                      value={match.match_score}
                                      sx={{
                                        height: 8,
                                        borderRadius: 4,
                                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                                        '& .MuiLinearProgress-bar': {
                                          background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                                          borderRadius: 4,
                                        }
                                      }}
                                    />
                                  </Box>
                                  <Typography variant="body2" sx={{ fontWeight: 700, color: '#059669', minWidth: 50 }}>
                                    {(match.match_score).toFixed(1)}%
                                  </Typography>
                                </Box>
                                
                                {match.ai_summary && (
                                  <Typography variant="body2" sx={{ 
                                    color: '#64748b',
                                    lineHeight: 1.6,
                                    mb: 2,
                                    p: 2,
                                    background: 'rgba(16, 185, 129, 0.05)',
                                    borderRadius: 2,
                                    border: '1px solid rgba(16, 185, 129, 0.1)'
                                  }}>
                                    <strong>AI Summary:</strong> {match.ai_summary}
                                  </Typography>
                                )}
                                
                                {match.matching_skills && match.matching_skills.length > 0 && (
                                  <Box>
                                    <Typography variant="body2" sx={{ mb: 1, fontWeight: 600, color: '#1e293b' }}>
                                      <strong>Matching Skills:</strong>
                                    </Typography>
                                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                                      {match.matching_skills.map((skill: string, skillIndex: number) => (
                                        <motion.div
                                          key={skillIndex}
                                          initial={{ opacity: 0, scale: 0.8 }}
                                          animate={{ opacity: 1, scale: 1 }}
                                          transition={{ delay: 0.3 + index * 0.1 + skillIndex * 0.05 }}
                                        >
                                          <Chip
                                            label={skill}
                                            size="small"
                                            sx={{
                                              background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                                              color: 'white',
                                              fontWeight: 600,
                                              '&:hover': {
                                                transform: 'scale(1.05)',
                                              },
                                              transition: 'all 0.2s ease-in-out'
                                            }}
                                          />
                                        </motion.div>
                                      ))}
                                    </Box>
                                  </Box>
                                )}
                                
                                {match.bio && (
                                  <Box sx={{ 
                                    mt: 2,
                                    p: 2,
                                    background: 'rgba(102, 126, 234, 0.05)',
                                    borderRadius: 2,
                                    border: '1px solid rgba(102, 126, 234, 0.1)'
                                  }}>
                                    <Typography variant="body2" sx={{ 
                                      color: '#64748b',
                                      lineHeight: 1.6,
                                      fontWeight: 600
                                    }}>
                                      <strong>Bio:</strong> {match.bio}
                                    </Typography>
                                  </Box>
                                )}
                              </Grid>
                              
                              <Grid item xs={12} md={4}>
                                {match.resume_link && (
                                  <Button
                                    variant="contained"
                                    href={match.resume_link}
                                    target="_blank"
                                    sx={{
                                      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                                      color: 'white',
                                      fontWeight: 600,
                                      width: '100%',
                                      py: 1.5,
                                      '&:hover': {
                                        background: 'linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%)',
                                        transform: 'translateY(-2px)',
                                        boxShadow: '0 8px 25px rgba(102, 126, 234, 0.3)',
                                      },
                                      transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
                                    }}
                                  >
                                    View Resume
                                  </Button>
                                )}
                              </Grid>
                            </Grid>
                          </CardContent>
                        </Card>
                      </motion.div>
                    </Grid>
                  ))}
                </Grid>
              )}
            </CardContent>
          </Card>
        </motion.div>
      </Grid>
    </Grid>
  );

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      PaperProps={{
        sx: {
          background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
          backdropFilter: 'blur(20px)',
          border: '1px solid rgba(255,255,255,0.2)',
          borderRadius: 4,
          boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
          overflow: 'hidden',
        }
      }}
    >
      <DialogTitle sx={{ 
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        color: 'white',
        fontWeight: 700,
        fontSize: '1.25rem',
        py: 2,
        px: 3,
        position: 'relative',
        '&::before': {
          content: '""',
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
          backdropFilter: 'blur(10px)',
          zIndex: -1
        }
      }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Avatar sx={{ 
            background: 'rgba(255,255,255,0.2)',
            color: 'white',
            width: 40,
            height: 40
          }}>
            <ViewIcon />
          </Avatar>
          AI Agent Request Details
        </Box>
      </DialogTitle>
      
      <DialogContent sx={{ 
        p: 0,
        background: 'linear-gradient(180deg, rgba(248,250,252,0.8) 0%, rgba(255,255,255,0.9) 100%)',
      }}>
        <AnimatePresence>
          {selectedRequest && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs 
                  value={activeTab} 
                  onChange={(e, newValue) => setActiveTab(newValue)}
                  sx={{
                    px: 3,
                    '& .MuiTab-root': {
                      textTransform: 'none',
                      fontWeight: 600,
                      fontSize: '0.95rem',
                      minHeight: 48,
                      '&.Mui-selected': {
                        color: '#667eea',
                      }
                    },
                    '& .MuiTabs-indicator': {
                      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                      height: 3,
                      borderRadius: '2px 2px 0 0'
                    }
                  }}
                >
                  <Tab label="Overview" />
                  {response && <Tab label="AI Response" />}
                  {response?.matches && response.matches.length > 0 && <Tab label="Matches" />}
                </Tabs>
              </Box>
              
              <Box sx={{ p: 3 }}>
                <AnimatePresence mode="wait">
                  {activeTab === 0 && (
                    <motion.div
                      key="overview"
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -20 }}
                      transition={{ duration: 0.3 }}
                    >
                      {renderOverviewTab()}
                    </motion.div>
                  )}
                  
                  {activeTab === 1 && response && (
                    <motion.div
                      key="ai-response"
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -20 }}
                      transition={{ duration: 0.3 }}
                    >
                      {renderAIResponseTab()}
                    </motion.div>
                  )}
                  
                  {activeTab === 2 && response?.matches && response.matches.length > 0 && (
                    <motion.div
                      key="matches"
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -20 }}
                      transition={{ duration: 0.3 }}
                    >
                      {renderMatchesTab()}
                    </motion.div>
                  )}
                </AnimatePresence>
              </Box>
            </motion.div>
          )}
        </AnimatePresence>
      </DialogContent>
      
      <DialogActions sx={{ 
        p: 2, 
        background: 'linear-gradient(180deg, rgba(248,250,252,0.8) 0%, rgba(255,255,255,0.9) 100%)',
        borderTop: '1px solid rgba(226, 232, 240, 0.3)'
      }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.0 }}
        >
          <Button 
            onClick={onClose}
            variant="contained"
            sx={{
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
              fontWeight: 600,
              px: 4,
              py: 1.5,
              borderRadius: 2,
              '&:hover': {
                background: 'linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%)',
                transform: 'translateY(-2px)',
                boxShadow: '0 8px 25px rgba(102, 126, 234, 0.3)',
              },
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
            }}
          >
            Close
          </Button>
        </motion.div>
      </DialogActions>
    </Dialog>
  );
};

export default RequestDetailsModal;