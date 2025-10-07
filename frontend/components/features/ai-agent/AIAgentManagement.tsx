import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { useAIAgentRequests, useProcessAIAgentRequest, useAIAgentRequest, useCreateAndProcessAIAgentRequest } from '../../../hooks/useApi';
import { AIAgentRequest, AIAgentResponse, CreateAIAgentRequest, aiAgentService } from '../../../services/ai/aiAgentService';
import RequestTable from './RequestTable';
import RequestDetailsModal from './RequestDetailsModal';
import CreateRequestDialog from './CreateRequestDialog';

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

      <RequestTable
        requests={requests || []}
        loading={loading}
        processing={processing}
        onRefresh={refetch}
        onViewRequest={viewRequest}
        onProcessRequest={handleProcessRequest}
      />

      <RequestDetailsModal
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        selectedRequest={selectedRequest}
        response={response}
      />

      <CreateRequestDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        onSubmit={handleCreateRequest}
        loading={creating}
        formData={createForm}
        onFormChange={setCreateForm}
      />
    </Box>
  );
};

export default AIAgentManagement;
