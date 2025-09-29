import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Chip,
  Alert,
  CircularProgress,
  Grid,
  Paper,
} from '@mui/material';
import { useExtractSkills } from '../../../hooks/useApi';
import { SkillExtractionResponse } from '../../../services/ai/aiAgentService';

const SkillExtractionTool: React.FC = () => {
  const [text, setText] = useState('');
  const [response, setResponse] = useState<SkillExtractionResponse | null>(null);
  const [localError, setLocalError] = useState<string | null>(null);
  const { extractSkills, loading, error } = useExtractSkills();

  const handleExtractSkills = async () => {
    if (!text.trim()) {
      setLocalError('Please enter some text to extract skills from');
      return;
    }

    try {
      setLocalError(null);
      const data = await extractSkills(text);
      setResponse(data);
    } catch (err) {
      setLocalError('Failed to extract skills. Please try again.');
    }
  };

  const combinedError = localError || error;

  const clearResults = () => {
    setResponse(null);
    setLocalError(null);
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Skill Extraction Tool
      </Typography>
      
      <Typography variant="body1" color="text.secondary" paragraph>
        Use this tool to extract technical skills from job descriptions, resumes, or any text content.
        The AI will analyze the text and identify relevant technical skills.
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Input Text
              </Typography>
              <TextField
                fullWidth
                multiline
                rows={8}
                variant="outlined"
                placeholder="Enter job description, resume text, or any content to extract skills from..."
                value={text}
                onChange={(e) => setText(e.target.value)}
                disabled={loading}
              />
              <Box sx={{ mt: 2, display: 'flex', gap: 2 }}>
                <Button
                  variant="contained"
                  onClick={handleExtractSkills}
                  disabled={loading || !text.trim()}
                  startIcon={loading ? <CircularProgress size={20} /> : undefined}
                >
                  {loading ? 'Extracting...' : 'Extract Skills'}
                </Button>
                <Button
                  variant="outlined"
                  onClick={clearResults}
                  disabled={loading}
                >
                  Clear
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Extracted Skills
              </Typography>
              
              {combinedError && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {combinedError}
                </Alert>
              )}

              {loading && (
                <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
                  <CircularProgress />
                </Box>
              )}

              {response && !loading && (
                <Box>
                  <Typography variant="body2" color="text.secondary" paragraph>
                    Found {response.skills.length} skills:
                  </Typography>
                  
                  {response.skills.length > 0 ? (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                      {response.skills.map((skill, index) => (
                        <Chip
                          key={index}
                          label={skill}
                          color="primary"
                          variant="outlined"
                        />
                      ))}
                    </Box>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      No skills were extracted from the text.
                    </Typography>
                  )}

                  <Paper sx={{ mt: 3, p: 2, bgcolor: 'grey.50' }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Original Text:
                    </Typography>
                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                      {response.text}
                    </Typography>
                  </Paper>
                </Box>
              )}

              {!response && !loading && !combinedError && (
                <Typography variant="body2" color="text.secondary">
                  Enter text above and click "Extract Skills" to see results.
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Box sx={{ mt: 3 }}>
        <Typography variant="h6" gutterBottom>
          How it works
        </Typography>
        <Typography variant="body2" color="text.secondary" paragraph>
          The skill extraction tool uses OpenAI's GPT model to analyze text and identify technical skills.
          It compares the text against our database of known skills and returns the most relevant matches.
        </Typography>
        <Typography variant="body2" color="text.secondary">
          <strong>Tips for better results:</strong>
        </Typography>
        <ul>
          <li>Include specific technical terms and technologies</li>
          <li>Mention programming languages, frameworks, and tools</li>
          <li>Include years of experience when possible</li>
          <li>Be specific about skill levels (beginner, intermediate, expert)</li>
        </ul>
      </Box>
    </Box>
  );
};

export default SkillExtractionTool;
