'use client'

import { useState, useEffect } from 'react'
import {
  Grid,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Autocomplete,
  Chip,
} from '@mui/material'
import { EmployeeFormProps, EmployeeFormData } from './interfaces'
import { DEPARTMENTS, EXPERIENCE_LEVELS } from '@/types'

const departments = DEPARTMENTS
const experienceLevels = EXPERIENCE_LEVELS

export function EmployeeForm({
  formData,
  onFormDataChange,
  skills,
  skillsLoading,
  error
}: EmployeeFormProps) {
  const [localFormData, setLocalFormData] = useState<EmployeeFormData>(formData)

  useEffect(() => {
    setLocalFormData(formData)
  }, [formData])

  const handleInputChange = (field: keyof EmployeeFormData) => (event: any) => {
    const newData = {
      ...localFormData,
      [field]: event.target.value
    }
    setLocalFormData(newData)
    onFormDataChange(newData)
  }

  const handleSkillsChange = (event: any, newValue: string[]) => {
    const newData = {
      ...localFormData,
      skills: newValue
    }
    setLocalFormData(newData)
    onFormDataChange(newData)
  }

  return (
    <Grid container spacing={2} sx={{ mt: 1 }}>
      <Grid item xs={12} sm={6}>
        <TextField
          fullWidth
          label="Name"
          value={localFormData.name}
          onChange={handleInputChange('name')}
          required
        />
      </Grid>
      <Grid item xs={12} sm={6}>
        <TextField
          fullWidth
          label="Email"
          type="email"
          value={localFormData.email}
          onChange={handleInputChange('email')}
          required
        />
      </Grid>
      <Grid item xs={12} sm={6}>
        <FormControl fullWidth>
          <InputLabel>Department</InputLabel>
          <Select
            value={localFormData.department}
            onChange={handleInputChange('department')}
            label="Department"
            required
          >
            {departments.map((dept) => (
              <MenuItem key={dept} value={dept}>
                {dept}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={12} sm={6}>
        <FormControl fullWidth>
          <InputLabel>Experience Level</InputLabel>
          <Select
            value={localFormData.level}
            onChange={handleInputChange('level')}
            label="Experience Level"
            required
          >
            {experienceLevels.map((level) => (
              <MenuItem key={level.value} value={level.value}>
                {level.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={12} sm={6}>
        <TextField
          fullWidth
          label="Location"
          value={localFormData.location}
          onChange={handleInputChange('location')}
        />
      </Grid>
      <Grid item xs={12} sm={6}>
        <TextField
          fullWidth
          label="Current Project"
          value={localFormData.current_project || ''}
          onChange={handleInputChange('current_project')}
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          fullWidth
          label="Bio"
          multiline
          rows={3}
          value={localFormData.bio}
          onChange={handleInputChange('bio')}
        />
      </Grid>
      <Grid item xs={12}>
        <Autocomplete
          multiple
          options={skills?.map(skill => skill.name) || []}
          value={localFormData.skills}
          onChange={handleSkillsChange}
          loading={skillsLoading}
          renderTags={(value, getTagProps) =>
            value.map((option, index) => (
              <Chip variant="outlined" label={option} {...getTagProps({ index })} />
            ))
          }
          renderInput={(params) => (
            <TextField
              {...params}
              label="Skills"
              placeholder="Select skills..."
            />
          )}
        />
      </Grid>
    </Grid>
  )
}