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
  Paper,
  Avatar,
  Rating,
  Divider,
} from '@mui/material'
import {
  Search as SearchIcon,
  Person as PersonIcon,
  Email as EmailIcon,
  LocationOn as LocationIcon,
  Work as WorkIcon,
  Star as StarIcon,
} from '@mui/icons-material'
import { useSearchEmployees, useSkills } from '@/hooks/useApi'
import { Employee, Skill, SearchCriteria, EXPERIENCE_LEVELS, DEPARTMENTS } from '@/types'

const experienceLevels = EXPERIENCE_LEVELS
const departments = DEPARTMENTS

export function SearchEmployees() {
  const [searchCriteria, setSearchCriteria] = useState<SearchCriteria>({
    required_skills: [],
    preferred_skills: [],
    department: '',
    experience_level: '',
    location: '',
  })

  const { data: skills, loading: skillsLoading, error: skillsError } = useSkills()
  const { searchEmployees, results: searchResults, loading: searching, error: searchError } = useSearchEmployees()

  const error = skillsError || searchError

  const handleInputChange = (field: keyof SearchCriteria) => (event: any) => {
    setSearchCriteria(prev => ({
      ...prev,
      [field]: event.target.value
    }))
  }

  const handleSkillsChange = (field: 'required_skills' | 'preferred_skills') => (event: any, newValue: string[]) => {
    setSearchCriteria(prev => ({
      ...prev,
      [field]: newValue
    }))
  }

  const handleSearch = async () => {
    if (!searchCriteria.required_skills.length && !searchCriteria.preferred_skills.length) {
      return
    }

    try {
      await searchEmployees(searchCriteria)
    } catch (err: any) {
      console.error('Error searching employees:', err)
    }
  }

  const getProficiencyStars = (level?: number) => {
    if (!level) return null
    return (
      <Box display="flex" alignItems="center" gap={0.5}>
        <Rating value={level} max={5} size="small" readOnly />
        <Typography variant="caption" color="text.secondary">
          {level}/5
        </Typography>
      </Box>
    )
  }

  const getExperienceColor = (level: string) => {
    switch (level.toLowerCase()) {
      case 'junior':
        return 'success'
      case 'mid':
        return 'warning'
      case 'senior':
        return 'error'
      default:
        return 'default'
    }
  }

  return (
    <Grid container spacing={3}>
      {/* Search Criteria */}
      <Grid item xs={12} md={4}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Search Criteria
            </Typography>
            
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            <Grid container spacing={2}>
              <Grid item xs={12}>
                <Autocomplete
                  multiple
                  options={skills?.map(skill => skill.name) || []}
                  value={searchCriteria.required_skills}
                  onChange={handleSkillsChange('required_skills')}
                  loading={skillsLoading}
                  renderTags={(value, getTagProps) =>
                    value.map((option, index) => (
                      <Chip
                        variant="filled"
                        label={option}
                        {...getTagProps({ index })}
                        key={option}
                        size="small"
                      />
                    ))
                  }
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Required Skills"
                      placeholder="Select required skills"
                      variant="outlined"
                      size="small"
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12}>
                <Autocomplete
                  multiple
                  options={skills?.map(skill => skill.name) || []}
                  value={searchCriteria.preferred_skills}
                  onChange={handleSkillsChange('preferred_skills')}
                  loading={skillsLoading}
                  renderTags={(value, getTagProps) =>
                    value.map((option, index) => (
                      <Chip
                        variant="outlined"
                        label={option}
                        {...getTagProps({ index })}
                        key={option}
                        size="small"
                      />
                    ))
                  }
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Preferred Skills"
                      placeholder="Select preferred skills"
                      variant="outlined"
                      size="small"
                    />
                  )}
                />
              </Grid>

              <Grid item xs={12}>
                <FormControl fullWidth size="small">
                  <InputLabel>Department</InputLabel>
                  <Select
                    value={searchCriteria.department}
                    onChange={handleInputChange('department')}
                    label="Department"
                  >
                    <MenuItem value="">Any Department</MenuItem>
                    {departments.map((dept) => (
                      <MenuItem key={dept} value={dept}>
                        {dept}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12}>
                <FormControl fullWidth size="small">
                  <InputLabel>Experience Level</InputLabel>
                  <Select
                    value={searchCriteria.experience_level}
                    onChange={handleInputChange('experience_level')}
                    label="Experience Level"
                  >
                    <MenuItem value="">Any Level</MenuItem>
                    {experienceLevels.map((level) => (
                      <MenuItem key={level.value} value={level.value}>
                        {level.label}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Location"
                  value={searchCriteria.location}
                  onChange={handleInputChange('location')}
                  variant="outlined"
                  size="small"
                  placeholder="e.g., San Francisco, Remote"
                />
              </Grid>

              <Grid item xs={12}>
                <Button
                  fullWidth
                  variant="contained"
                  startIcon={searching ? <CircularProgress size={20} /> : <SearchIcon />}
                  onClick={handleSearch}
                  disabled={searching}
                  size="large"
                >
                  {searching ? 'Searching...' : 'Search Employees'}
                </Button>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </Grid>

      {/* Search Results */}
      <Grid item xs={12} md={8}>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Search Results ({searchResults.length})
            </Typography>
            
            {searching ? (
              <Box display="flex" justifyContent="center" py={4}>
                <CircularProgress />
              </Box>
            ) : searchResults.length > 0 ? (
              <Grid container spacing={2}>
                {searchResults.map((employee) => (
                  <Grid item xs={12} key={employee.id}>
                    <Paper elevation={1} sx={{ p: 3 }}>
                      <Box display="flex" gap={2}>
                        <Avatar sx={{ width: 60, height: 60 }}>
                          <PersonIcon sx={{ fontSize: 30 }} />
                        </Avatar>
                        
                        <Box flex={1}>
                          <Box display="flex" justifyContent="space-between" alignItems="start">
                            <Box>
                              <Typography variant="h6" gutterBottom>
                                {employee.name}
                              </Typography>
                              <Box display="flex" gap={1} mb={1}>
                                <Chip
                                  label={employee.level}
                                  color={getExperienceColor(employee.level) as any}
                                  size="small"
                                />
                                <Chip
                                  label={employee.department}
                                  variant="outlined"
                                  size="small"
                                />
                              </Box>
                            </Box>
                            
                            <Box textAlign="right">
                              <Typography variant="body2" color="text.secondary">
                                Match Score
                              </Typography>
                              <Typography variant="h6" color="primary">
                                95%
                              </Typography>
                            </Box>
                          </Box>

                          <Box display="flex" alignItems="center" gap={2} mb={2}>
                            <Box display="flex" alignItems="center" gap={0.5}>
                              <EmailIcon fontSize="small" color="action" />
                              <Typography variant="body2">
                                {employee.email}
                              </Typography>
                            </Box>
                            
                            {employee.location && (
                              <Box display="flex" alignItems="center" gap={0.5}>
                                <LocationIcon fontSize="small" color="action" />
                                <Typography variant="body2">
                                  {employee.location}
                                </Typography>
                              </Box>
                            )}
                          </Box>

                          {employee.bio && (
                            <Typography variant="body2" color="text.secondary" paragraph>
                              {employee.bio}
                            </Typography>
                          )}

                          <Divider sx={{ my: 2 }} />

                          <Typography variant="subtitle2" gutterBottom>
                            Skills & Experience
                          </Typography>
                          <Box display="flex" gap={1} flexWrap="wrap">
                            {employee.skills.map((skill) => (
                              <Box key={skill.id} display="flex" alignItems="center" gap={0.5}>
                                <Chip
                                  label={skill.name}
                                  size="small"
                                  variant="outlined"
                                />
                                {skill.proficiency_level && (
                                  <Box>
                                    {getProficiencyStars(skill.proficiency_level)}
                                  </Box>
                                )}
                                {skill.years_experience && (
                                  <Typography variant="caption" color="text.secondary">
                                    ({skill.years_experience}y)
                                  </Typography>
                                )}
                              </Box>
                            ))}
                          </Box>
                        </Box>
                      </Box>
                    </Paper>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Box textAlign="center" py={4}>
                <SearchIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  No employees found
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Try adjusting your search criteria
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}
