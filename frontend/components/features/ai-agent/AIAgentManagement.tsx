import React, { useState, useEffect } from 'react';
import {
  Box,
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
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Grid,
  Alert,
  CircularProgress,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  Visibility as ViewIcon,
  PlayArrow as ProcessIcon,
  Error as ErrorIcon,
  CheckCircle as SuccessIcon,
  Schedule as PendingIcon,
} from '@mui/icons-material';
import { useAIAgentRequests, useProcessAIAgentRequest, useAIAgentRequest, useCreateAndProcessAIAgentRequest } from '../../../hooks/useApi';
import { AIAgentRequest, AIAgentMatch, AIAgentResponse, CreateAIAgentRequest, aiAgentService } from '../../../services/ai/aiAgentService';

const AIAgentManagement: React.FC = () => {
  const [selectedRequest, setSelectedRequest] = useState<AIAgentRequest | null>(null);
  const [response, setResponse] = useState<AIAgentResponse | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedRequestId, setSelectedRequestId] = useState<number | null>(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [createForm, setCreateForm] = useState<CreateAIAgentRequest>({
    teams_message_id: '',
    channel_id: 'test-channel',
    user_id: 'test-user',
    user_name: 'Test User',
    message_text: '',
    attachment_url: undefined
  });
  
  // Use hooks for data fetching
  const { data: requests, loading, error: fetchError, refetch } = useAIAgentRequests(50, 0);
  const { processRequest, loading: processing, error: processError } = useProcessAIAgentRequest();
  const { data: requestDetails, loading: detailsLoading } = useAIAgentRequest(selectedRequestId || 0);
  const { processNewRequest, loading: creating, error: createError } = useCreateAndProcessAIAgentRequest();

  // Set combined error
  const combinedError = error || fetchError || processError || createError;

  const handleProcessRequest = async (requestId: number) => {
    try {
      setError(null);
      const data = await processRequest(requestId);
      setResponse(data);
      refetch(); // Refresh the list
    } catch (err) {
      setError('Failed to process request');
    }
  };

  const viewRequest = async (requestId: number) => {
    setSelectedRequestId(requestId);
    setDialogOpen(true);
    
    // Fetch response details if request is completed
    try {
      const responseData = await aiAgentService.getResponse(requestId);
      setResponse(responseData);
    } catch (error) {
      // Response might not exist yet, that's okay
      console.log('No response found for request:', requestId);
    }
  };

  // Update selectedRequest when requestDetails changes
  useEffect(() => {
    if (requestDetails) {
      setSelectedRequest(requestDetails);
    }
  }, [requestDetails]);

  const handleCreateRequest = async () => {
    try {
      setError(null);
      const data = await processNewRequest(createForm);
      setResponse(data);
      setCreateDialogOpen(false);
      refetch(); // Refresh the list
    } catch (err) {
      setError('Failed to create and process AI agent request');
    }
  };

  const resetCreateForm = () => {
    setCreateForm({
      teams_message_id: `test-${Date.now()}`,
      channel_id: 'test-channel',
      user_id: 'test-user',
      user_name: 'Test User',
      message_text: '',
      attachment_url: undefined
    });
  };

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

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          AI Agent Management
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={refetch}
            disabled={loading}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            onClick={() => {
              resetCreateForm();
              setCreateDialogOpen(true);
            }}
            disabled={creating}
          >
            Create Request
          </Button>
        </Box>
      </Box>

      {combinedError && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
          {combinedError}
        </Alert>
      )}

      <Card>
        <CardContent>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>User</TableCell>
                  <TableCell>Message</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Created</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {(requests || []).map((request) => (
                  <TableRow key={request.id}>
                    <TableCell>{request.id}</TableCell>
                    <TableCell>{request.user_name}</TableCell>
                    <TableCell>
                      <Typography variant="body2" noWrap sx={{ maxWidth: 200 }}>
                        {request.message_text || 'No message text'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        icon={getStatusIcon(request.status)}
                        label={request.status}
                        color={getStatusColor(request.status) as any}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>{formatDate(request.created_at)}</TableCell>
                    <TableCell>
                      <Tooltip title="View Details">
                        <IconButton
                          size="small"
                          onClick={() => viewRequest(request.id)}
                        >
                          <ViewIcon />
                        </IconButton>
                      </Tooltip>
                      {request.status === 'pending' && (
                        <Tooltip title="Process Request">
                          <IconButton
                            size="small"
                            onClick={() => handleProcessRequest(request.id)}
                            disabled={processing}
                          >
                            {processing ? (
                              <CircularProgress size={20} />
                            ) : (
                              <ProcessIcon />
                            )}
                          </IconButton>
                        </Tooltip>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      {/* Request Details Dialog */}
      <Dialog
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>AI Agent Request Details</DialogTitle>
        <DialogContent>
          {selectedRequest && (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Typography variant="h6">Basic Information</Typography>
                <Typography variant="body2">
                  <strong>ID:</strong> {selectedRequest.id}
                </Typography>
                <Typography variant="body2">
                  <strong>User:</strong> {selectedRequest.user_name}
                </Typography>
                <Typography variant="body2">
                  <strong>Status:</strong> {selectedRequest.status}
                </Typography>
                <Typography variant="body2">
                  <strong>Created:</strong> {formatDate(selectedRequest.created_at)}
                </Typography>
                {selectedRequest.processed_at && (
                  <Typography variant="body2">
                    <strong>Processed:</strong> {formatDate(selectedRequest.processed_at)}
                  </Typography>
                )}
              </Grid>

              <Grid item xs={12}>
                <Typography variant="h6">Message</Typography>
                <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                  {selectedRequest.message_text || 'No message text'}
                </Typography>
              </Grid>

              {selectedRequest.attachment_url && (
                <Grid item xs={12}>
                  <Typography variant="h6">Attachment</Typography>
                  <Typography variant="body2">
                    <a href={selectedRequest.attachment_url} target="_blank" rel="noopener noreferrer">
                      {selectedRequest.attachment_url}
                    </a>
                  </Typography>
                </Grid>
              )}

              {selectedRequest.extracted_text && (
                <Grid item xs={12}>
                  <Typography variant="h6">Extracted Text</Typography>
                  <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                    {selectedRequest.extracted_text}
                  </Typography>
                </Grid>
              )}

              {selectedRequest.extracted_skills && selectedRequest.extracted_skills.length > 0 && (
                <Grid item xs={12}>
                  <Typography variant="h6">Extracted Skills</Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {selectedRequest.extracted_skills.map((skill, index) => (
                      <Chip key={index} label={skill} size="small" />
                    ))}
                  </Box>
                </Grid>
              )}

              {selectedRequest.error && (
                <Grid item xs={12}>
                  <Typography variant="h6" color="error">Error</Typography>
                  <Typography variant="body2" color="error" sx={{ whiteSpace: 'pre-wrap' }}>
                    {selectedRequest.error}
                  </Typography>
                </Grid>
              )}

              {response && (
                <Grid item xs={12}>
                  <Typography variant="h6">AI Response</Typography>
                  <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap', mb: 2 }}>
                    {response.summary}
                  </Typography>
                  
                  {response.matches && response.matches.length > 0 && (
                    <>
                      <Typography variant="h6" sx={{ mt: 2 }}>Top {response.matches.length} Candidate Matches</Typography>
                      {response.matches.map((match, index) => (
                        <Card key={index} sx={{ mb: 2, p: 2 }}>
                          <Grid container spacing={2}>
                            <Grid item xs={12} md={8}>
                              <Typography variant="h6" color="primary">
                                {match.employee_name}
                              </Typography>
                              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                                {match.position} • {match.seniority} • {match.location}
                              </Typography>
                              <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>Email:</strong> {match.employee_email}
                              </Typography>
                              {match.current_project && (
                                <Typography variant="body2" sx={{ mb: 1 }}>
                                  <strong>Current Project:</strong> {match.current_project}
                                </Typography>
                              )}
                              <Typography variant="body2" sx={{ mb: 1 }}>
                                <strong>Match Score:</strong> {match.match_score.toFixed(2)}
                              </Typography>
                            </Grid>
                            <Grid item xs={12} md={4}>
                              {match.resume_link && (
                                <Button
                                  variant="outlined"
                                  size="small"
                                  href={match.resume_link}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  sx={{ mb: 1, width: '100%' }}
                                >
                                  View Resume
                                </Button>
                              )}
                            </Grid>
                          </Grid>
                          
                          <Typography variant="body2" sx={{ mt: 2, mb: 1, fontStyle: 'italic' }}>
                            <strong>AI Summary:</strong> {match.ai_summary}
                          </Typography>
                          
                          {match.matching_skills.length > 0 && (
                            <Box sx={{ mt: 1 }}>
                              <Typography variant="body2" sx={{ mb: 0.5 }}>
                                <strong>Matching Skills:</strong>
                              </Typography>
                              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                                {match.matching_skills.map((skill, skillIndex) => (
                                  <Chip key={skillIndex} label={skill} size="small" color="primary" />
                                ))}
                              </Box>
                            </Box>
                          )}
                          
                          {match.bio && (
                            <Typography variant="body2" sx={{ mt: 1, fontStyle: 'italic' }}>
                              <strong>Bio:</strong> {match.bio}
                            </Typography>
                          )}
                        </Card>
                      ))}
                    </>
                  )}
                </Grid>
              )}
            </Grid>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Create AI Agent Request Dialog */}
      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Create AI Agent Request</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="User Name"
                value={createForm.user_name}
                onChange={(e) => setCreateForm({ ...createForm, user_name: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Channel ID"
                value={createForm.channel_id}
                onChange={(e) => setCreateForm({ ...createForm, channel_id: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Message Text"
                value={createForm.message_text}
                onChange={(e) => setCreateForm({ ...createForm, message_text: e.target.value })}
                multiline
                rows={4}
                placeholder="Enter a job description or skill request, e.g., 'Looking for React developers with 3+ years experience in TypeScript and Node.js'"
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Attachment URL (Optional)"
                value={createForm.attachment_url || ''}
                onChange={(e) => setCreateForm({ ...createForm, attachment_url: e.target.value || undefined })}
                placeholder="https://example.com/resume.pdf"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleCreateRequest} 
            variant="contained"
            disabled={creating || !createForm.message_text.trim()}
          >
            {creating ? 'Creating...' : 'Create & Process'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default AIAgentManagement;
