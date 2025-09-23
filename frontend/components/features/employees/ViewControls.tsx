'use client'

import { ToggleButton, ToggleButtonGroup, IconButton, Tooltip } from '@mui/material'
import { ViewList as ListIcon, ViewModule as GridIcon, TableChart as TableIcon, Clear as ClearIcon } from '@mui/icons-material'
import { ViewControlsProps } from './interfaces'

export function ViewControls({ viewMode, onViewModeChange, onClearFilters }: ViewControlsProps) {
  return (
    <>
      <ToggleButtonGroup
        value={viewMode}
        exclusive
        onChange={(_, newMode) => newMode && onViewModeChange(newMode)}
        size="small"
      >
        <ToggleButton value="grid">
          <Tooltip title="Grid View">
            <GridIcon />
          </Tooltip>
        </ToggleButton>
        <ToggleButton value="list">
          <Tooltip title="List View">
            <ListIcon />
          </Tooltip>
        </ToggleButton>
        <ToggleButton value="table">
          <Tooltip title="Table View">
            <TableIcon />
          </Tooltip>
        </ToggleButton>
      </ToggleButtonGroup>
      
      <Tooltip title="Clear all filters">
        <IconButton onClick={onClearFilters} size="small">
          <ClearIcon />
        </IconButton>
      </Tooltip>
    </>
  )
}
