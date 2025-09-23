'use client'

import { useState, useEffect } from 'react'
import {
  Grid,
  Card,
  CardContent,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Box,
  Typography,
  Alert,
  CircularProgress,
  Autocomplete,
} from '@mui/material'
import { Add as AddIcon, Save as SaveIcon } from '@mui/icons-material'
import { api, endpoints } from '@/lib/api'
import { Skill, JobRequestFormData, EXPERIENCE_LEVELS, PRIORITIES, DEPARTMENTS } from '@/types'

const experienceLevels = EXPERIENCE_LEVELS
const priorities = PRIORITIES
const departments = DEPARTMENTS

export function JobRequestForm() {
  const [formData, setFormData] = useState<JobRequestFormData>({
    title: '',
    description: '',
    department: '',
    required_skills: [],
    preferred_skills: [],
    experience_level: '',
    location: '',
    priority: 'medium',
    created_by: '',
  })

  const [skills, setSkills] = useState<Skill[]>([])
  const [loading, setLoading] = useState(false)
  const [skillsLoading, setSkillsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  useEffect(() => {
    const fetchSkills = async () => {
      try {
        setSkillsLoading(true)
        const response = await endpoints.skills.list()
        setSkills(response.data)
      } catch (err) {
        console.error('Error fetching skills:', err)
        setError('Failed to load skills')
      } finally {
        setSkillsLoading(false)
      }
    }

    fetchSkills()
  }, [])

  const handleInputChange = (field: keyof JobRequestFormData) => (event: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: event.target.value
    }))
  }

  const handleSkillsChange = (field: 'required_skills' | 'preferred_skills') => (event: any, newValue: string[]) => {
    setFormData(prev => ({
      ...prev,
      [field]: newValue
    }))
  }

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    
    if (!formData.title || !formData.department) {
      setError('Title and department are required')
      return
    }

    try {
      setLoading(true)
      setError(null)
      
      await endpoints.jobRequests.create(formData)
      
      setSuccess('Job request created successfully!')
      
      // Reset form
      setFormData({
        title: '',
        description: '',
        department: '',
        required_skills: [],
        preferred_skills: [],
        experience_level: '',
        location: '',
        priority: 'medium',
        created_by: '',
      })
    } catch (err: any) {
      console.error('Error creating job request:', err)
      setError(err.response?.data?.error || 'Failed to create job request')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Create New Job Request
            </Typography>
            
            {error && (
              <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
                {error}
              </Alert>
            )}
            
            {success && (
              <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess(null)}>
                {success}
              </Alert>
            )}

            <form onSubmit={handleSubmit}>
              <Grid container spacing={3}>
                <Grid item xs={12} md={6}>
                  <TextField
                    fullWidth
                    label="Job Title"
                    value={formData.title}
                    onChange={handleInputChange('title')}
                    required
                    variant="outlined"
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <FormControl fullWidth required>
                    <InputLabel>Department</InputLabel>
                    <Select
                      value={formData.department}
                      onChange={handleInputChange('department')}
                      label="Department"
                    >
                      {departments.map((dept) => (
                        <MenuItem key={dept} value={dept}>
                          {dept}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>

                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    label="Job Description"
                    value={formData.description}
                    onChange={handleInputChange('description')}
                    multiline
                    rows={4}
                    variant="outlined"
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <FormControl fullWidth>
                    <InputLabel>Experience Level</InputLabel>
                    <Select
                      value={formData.experience_level}
                      onChange={handleInputChange('experience_level')}
                      label="Experience Level"
                    >
                      {experienceLevels.map((level) => (
                        <MenuItem key={level.value} value={level.value}>
                          {level.label}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>

                <Grid item xs={12} md={6}>
                  <FormControl fullWidth>
                    <InputLabel>Priority</InputLabel>
                    <Select
                      value={formData.priority}
                      onChange={handleInputChange('priority')}
                      label="Priority"
                    >
                      {priorities.map((priority) => (
                        <MenuItem key={priority.value} value={priority.value}>
                          {priority.label}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>

                <Grid item xs={12} md={6}>
                  <TextField
                    fullWidth
                    label="Location"
                    value={formData.location}
                    onChange={handleInputChange('location')}
                    variant="outlined"
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <TextField
                    fullWidth
                    label="Created By"
                    value={formData.created_by}
                    onChange={handleInputChange('created_by')}
                    variant="outlined"
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <Autocomplete
                    multiple
                    options={skills.map(skill => skill.name)}
                    value={formData.required_skills}
                    onChange={handleSkillsChange('required_skills')}
                    loading={skillsLoading}
                    renderTags={(value, getTagProps) =>
                      value.map((option, index) => (
                        <Chip
                          variant="filled"
                          label={option}
                          {...getTagProps({ index })}
                          key={option}
                        />
                      ))
                    }
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        label="Required Skills"
                        placeholder="Select required skills"
                        variant="outlined"
                      />
                    )}
                  />
                </Grid>

                <Grid item xs={12} md={6}>
                  <Autocomplete
                    multiple
                    options={skills.map(skill => skill.name)}
                    value={formData.preferred_skills}
                    onChange={handleSkillsChange('preferred_skills')}
                    loading={skillsLoading}
                    renderTags={(value, getTagProps) =>
                      value.map((option, index) => (
                        <Chip
                          variant="outlined"
                          label={option}
                          {...getTagProps({ index })}
                          key={option}
                        />
                      ))
                    }
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        label="Preferred Skills"
                        placeholder="Select preferred skills"
                        variant="outlined"
                      />
                    )}
                  />
                </Grid>

                <Grid item xs={12}>
                  <Box display="flex" gap={2}>
                    <Button
                      type="submit"
                      variant="contained"
                      startIcon={loading ? <CircularProgress size={20} /> : <SaveIcon />}
                      disabled={loading}
                    >
                      {loading ? 'Creating...' : 'Create Job Request'}
                    </Button>
                    
                    <Button
                      variant="outlined"
                      onClick={() => window.location.reload()}
                    >
                      Reset
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </form>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}
