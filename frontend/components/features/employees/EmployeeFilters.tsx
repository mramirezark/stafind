'use client'

import { Grid, FormControl, InputLabel, Select, MenuItem, Autocomplete, Chip, IconButton, TextField } from '@mui/material'
import { Clear as ClearIcon } from '@mui/icons-material'
import { FilterControlsProps } from './interfaces'
import { DEPARTMENTS, EXPERIENCE_LEVELS } from '@/types'

const departments = DEPARTMENTS
const experienceLevels = EXPERIENCE_LEVELS

export function EmployeeFilters({
  filters,
  onFiltersChange,
  onClearFilters,
  skills,
  skillsLoading
}: FilterControlsProps) {
  const handleDepartmentChange = (event: any) => {
    onFiltersChange({ departmentFilter: event.target.value })
  }

  const handleLevelChange = (event: any) => {
    onFiltersChange({ levelFilter: event.target.value })
  }

  const handleSkillChange = (event: any, newValue: string[]) => {
    onFiltersChange({ skillFilter: newValue })
  }

  return (
    <Grid container spacing={2} alignItems="center">
      {/* Department Filter */}
      <Grid item xs={12} md={3}>
        <FormControl fullWidth>
          <InputLabel>Department</InputLabel>
          <Select
            value={filters.departmentFilter}
            onChange={handleDepartmentChange}
            label="Department"
          >
            <MenuItem value="">All Departments</MenuItem>
            {departments.map((dept) => (
              <MenuItem key={dept} value={dept}>
                {dept}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

      {/* Level Filter */}
      <Grid item xs={12} md={3}>
        <FormControl fullWidth>
          <InputLabel>Level</InputLabel>
          <Select
            value={filters.levelFilter}
            onChange={handleLevelChange}
            label="Level"
          >
            <MenuItem value="">All Levels</MenuItem>
            {experienceLevels.map((level) => (
              <MenuItem key={level.value} value={level.value}>
                {level.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

      {/* Skills Filter */}
      <Grid item xs={12} md={5}>
        <Autocomplete
          multiple
          options={skills?.map(skill => skill.name) || []}
          value={filters.skillFilter}
          onChange={handleSkillChange}
          loading={skillsLoading}
          renderTags={(value, getTagProps) =>
            value.map((option, index) => (
              <Chip
                variant="outlined"
                label={option}
                size="small"
                {...getTagProps({ index })}
              />
            ))
          }
          renderInput={(params) => (
            <TextField
              {...params}
              label="Skills"
              placeholder="Filter by skills..."
            />
          )}
        />
      </Grid>

      {/* Clear Filters */}
      <Grid item xs={12} md={1}>
        <IconButton onClick={onClearFilters} size="small" title="Clear all filters">
          <ClearIcon />
        </IconButton>
      </Grid>
    </Grid>
  )
}
