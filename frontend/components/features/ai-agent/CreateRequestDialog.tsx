import React from 'react';
import { motion } from 'framer-motion';
import {
  Box,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Grid,
  Avatar,
  Typography,
} from '@mui/material';
import {
  Add as AddIcon,
} from '@mui/icons-material';
import { CreateAIAgentRequest } from '../../../services/ai/aiAgentService';

interface CreateRequestDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateAIAgentRequest) => void;
  loading: boolean;
  formData: CreateAIAgentRequest;
  onFormChange: (data: CreateAIAgentRequest) => void;
}

const CreateRequestDialog: React.FC<CreateRequestDialogProps> = ({
  open,
  onClose,
  onSubmit,
  loading,
  formData,
  onFormChange,
}) => {
  const handleSubmit = () => {
    onSubmit(formData);
  };

  const handleInputChange = (field: keyof CreateAIAgentRequest, value: string | undefined) => {
    onFormChange({
      ...formData,
      [field]: value,
    });
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
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
        background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
        color: 'white',
        fontWeight: 700,
        fontSize: '1.5rem',
        py: 3,
        px: 4,
        position: 'relative',
        '&::before': {
          content: '""',
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          background: 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
        }
      }}>
        <Box sx={{ position: 'relative', zIndex: 1, display: 'flex', alignItems: 'center', gap: 2 }}>
          <Avatar sx={{ 
            background: 'rgba(255,255,255,0.2)', 
            color: 'white',
            width: 40,
            height: 40
          }}>
            <AddIcon />
          </Avatar>
          Create AI Agent Request
        </Box>
      </DialogTitle>
      
      <DialogContent sx={{ 
        p: 4,
        background: 'linear-gradient(180deg, rgba(248,250,252,0.8) 0%, rgba(255,255,255,0.9) 100%)',
      }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <Grid container spacing={3} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="User Name"
                value={formData.user_name}
                onChange={(e) => handleInputChange('user_name', e.target.value)}
                required
                sx={{
                  '& .MuiOutlinedInput-root': {
                    background: 'rgba(255, 255, 255, 0.7)',
                    borderRadius: 2,
                    '&:hover': {
                      background: 'rgba(255, 255, 255, 0.8)',
                    },
                    '&.Mui-focused': {
                      background: 'rgba(255, 255, 255, 0.9)',
                    }
                  }
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Channel ID"
                value={formData.channel_id}
                onChange={(e) => handleInputChange('channel_id', e.target.value)}
                required
                sx={{
                  '& .MuiOutlinedInput-root': {
                    background: 'rgba(255, 255, 255, 0.7)',
                    borderRadius: 2,
                    '&:hover': {
                      background: 'rgba(255, 255, 255, 0.8)',
                    },
                    '&.Mui-focused': {
                      background: 'rgba(255, 255, 255, 0.9)',
                    }
                  }
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Message Text"
                value={formData.message_text}
                onChange={(e) => handleInputChange('message_text', e.target.value)}
                multiline
                rows={4}
                placeholder="Enter a job description or skill request, e.g., 'Looking for React developers with 3+ years experience in TypeScript and Node.js'"
                required
                sx={{
                  '& .MuiOutlinedInput-root': {
                    background: 'rgba(255, 255, 255, 0.7)',
                    borderRadius: 2,
                    '&:hover': {
                      background: 'rgba(255, 255, 255, 0.8)',
                    },
                    '&.Mui-focused': {
                      background: 'rgba(255, 255, 255, 0.9)',
                    }
                  }
                }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Attachment URL (Optional)"
                value={formData.attachment_url || ''}
                onChange={(e) => handleInputChange('attachment_url', e.target.value || undefined)}
                placeholder="https://example.com/resume.pdf"
                sx={{
                  '& .MuiOutlinedInput-root': {
                    background: 'rgba(255, 255, 255, 0.7)',
                    borderRadius: 2,
                    '&:hover': {
                      background: 'rgba(255, 255, 255, 0.8)',
                    },
                    '&.Mui-focused': {
                      background: 'rgba(255, 255, 255, 0.9)',
                    }
                  }
                }}
              />
            </Grid>
          </Grid>
        </motion.div>
      </DialogContent>
      
      <DialogActions sx={{ 
        p: 4, 
        background: 'linear-gradient(180deg, rgba(248,250,252,0.8) 0%, rgba(255,255,255,0.9) 100%)',
        borderTop: '1px solid rgba(226, 232, 240, 0.3)'
      }}>
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          style={{ display: 'flex', gap: '16px' }}
        >
          <Button 
            onClick={onClose}
            variant="outlined"
            sx={{
              borderColor: 'rgba(102, 126, 234, 0.3)',
              color: '#667eea',
              fontWeight: 600,
              px: 4,
              py: 1.5,
              borderRadius: 2,
              '&:hover': {
                borderColor: 'rgba(102, 126, 234, 0.5)',
                background: 'rgba(102, 126, 234, 0.05)',
                transform: 'translateY(-2px)',
              },
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
            }}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleSubmit}
            variant="contained"
            disabled={loading || !formData.message_text.trim()}
            sx={{
              background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
              color: 'white',
              fontWeight: 600,
              px: 4,
              py: 1.5,
              borderRadius: 2,
              '&:hover': {
                background: 'linear-gradient(135deg, #059669 0%, #047857 100%)',
                transform: 'translateY(-2px)',
                boxShadow: '0 8px 25px rgba(16, 185, 129, 0.3)',
              },
              '&:disabled': {
                background: 'rgba(156, 163, 175, 0.5)',
                color: 'rgba(255, 255, 255, 0.7)',
              },
              transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)'
            }}
          >
            {loading ? 'Creating...' : 'Create & Process'}
          </Button>
        </motion.div>
      </DialogActions>
    </Dialog>
  );
};

export default CreateRequestDialog;
