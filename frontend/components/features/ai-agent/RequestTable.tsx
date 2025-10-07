import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Tooltip,
  CircularProgress,
  Box,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  Visibility as ViewIcon,
  PlayArrow as ProcessIcon,
  Error as ErrorIcon,
  CheckCircle as SuccessIcon,
  Schedule as PendingIcon,
} from '@mui/icons-material';
import { AIAgentRequest } from '../../../types';

interface RequestTableProps {
  requests: AIAgentRequest[];
  loading: boolean;
  processing: boolean;
  onRefresh: () => void;
  onViewRequest: (requestId: number) => void;
  onProcessRequest: (requestId: number) => void;
}

const RequestTable: React.FC<RequestTableProps> = ({
  requests,
  loading,
  processing,
  onRefresh,
  onViewRequest,
  onProcessRequest,
}) => {
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

  return (
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
        <TableContainer 
          component={Paper} 
          sx={{ 
            background: 'transparent',
            boxShadow: 'none',
            '& .MuiTableHead-root': {
              background: 'linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(102, 126, 234, 0.02) 100%)',
            }
          }}
        >
          <Table>
            <TableHead>
              <TableRow>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  ID
                </TableCell>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  User
                </TableCell>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  Message
                </TableCell>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  Status
                </TableCell>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  Created
                </TableCell>
                <TableCell sx={{ 
                  fontWeight: 700, 
                  color: '#1e293b',
                  borderBottom: '1px solid rgba(102, 126, 234, 0.1)'
                }}>
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {requests.map((request) => (
                <TableRow 
                  key={request.id}
                  sx={{
                    '&:hover': {
                      background: 'rgba(102, 126, 234, 0.05)',
                    },
                    '& .MuiTableCell-root': {
                      borderBottom: '1px solid rgba(226, 232, 240, 0.3)',
                    }
                  }}
                >
                  <TableCell sx={{ fontWeight: 600, color: '#1e293b' }}>
                    #{request.id}
                  </TableCell>
                  <TableCell sx={{ fontWeight: 500, color: '#64748b' }}>
                    {request.user_name}
                  </TableCell>
                  <TableCell sx={{ maxWidth: 200 }}>
                    <Typography 
                      variant="body2" 
                      noWrap 
                      sx={{ 
                        color: '#64748b',
                        fontWeight: 500
                      }}
                    >
                      {request.message_text || 'No message text'}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Chip
                      icon={getStatusIcon(request.status)}
                      label={request.status}
                      color={getStatusColor(request.status) as any}
                      size="small"
                      sx={{ fontWeight: 600 }}
                    />
                  </TableCell>
                  <TableCell sx={{ color: '#64748b', fontWeight: 500 }}>
                    {formatDate(request.created_at)}
                  </TableCell>
                  <TableCell>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Tooltip title="View Details">
                        <IconButton
                          size="small"
                          onClick={() => onViewRequest(request.id)}
                          sx={{
                            color: '#667eea',
                            '&:hover': {
                              background: 'rgba(102, 126, 234, 0.1)',
                              transform: 'scale(1.1)',
                            },
                            transition: 'all 0.2s ease-in-out'
                          }}
                        >
                          <ViewIcon />
                        </IconButton>
                      </Tooltip>
                      {request.status === 'pending' && (
                        <Tooltip title="Process Request">
                          <IconButton
                            size="small"
                            onClick={() => onProcessRequest(request.id)}
                            disabled={processing}
                            sx={{
                              color: '#10b981',
                              '&:hover': {
                                background: 'rgba(16, 185, 129, 0.1)',
                                transform: 'scale(1.1)',
                              },
                              '&:disabled': {
                                color: '#9ca3af',
                              },
                              transition: 'all 0.2s ease-in-out'
                            }}
                          >
                            {processing ? (
                              <CircularProgress size={20} />
                            ) : (
                              <ProcessIcon />
                            )}
                          </IconButton>
                        </Tooltip>
                      )}
                    </Box>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  );
};

export default RequestTable;
