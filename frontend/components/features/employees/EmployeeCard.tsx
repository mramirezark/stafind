'use client'

import {
  Card,
  CardContent,
  Box,
  Typography,
  Avatar,
  Chip,
  IconButton,
  Tooltip,
} from '@mui/material'
import {
  Person as PersonIcon,
  Email as EmailIcon,
  LocationOn as LocationIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material'
import { EmployeeCardProps } from './interfaces'

export function EmployeeCard({ employee, viewMode, onEdit, onDelete }: EmployeeCardProps) {
  if (viewMode === 'list') {
    return (
      <Card sx={{ height: '100%' }}>
        <CardContent>
          <Box display="flex" alignItems="center" gap={2}>
            <Avatar sx={{ width: 56, height: 56 }}>
              <PersonIcon />
            </Avatar>
            <Box flex={1}>
              <Typography variant="h6" gutterBottom>
                {employee.name}
              </Typography>
              <Box display="flex" alignItems="center" gap={2} mb={1}>
                <Box display="flex" alignItems="center" gap={0.5}>
                  <EmailIcon fontSize="small" color="action" />
                  <Typography variant="body2" color="text.secondary">
                    {employee.email}
                  </Typography>
                </Box>
                {employee.location && (
                  <Box display="flex" alignItems="center" gap={0.5}>
                    <LocationIcon fontSize="small" color="action" />
                    <Typography variant="body2" color="text.secondary">
                      {employee.location}
                    </Typography>
                  </Box>
                )}
              </Box>
              <Box display="flex" alignItems="center" gap={2} mb={1}>
                <Chip label={employee.department} size="small" color="primary" />
                <Chip label={employee.level} size="small" variant="outlined" />
              </Box>
              {employee.skills && employee.skills.length > 0 && (
                <Box display="flex" flexWrap="wrap" gap={0.5}>
                  {employee.skills.slice(0, 3).map((skill, index) => (
                    <Chip
                      key={index}
                      label={skill.name}
                      size="small"
                      variant="outlined"
                    />
                  ))}
                  {employee.skills.length > 3 && (
                    <Chip
                      label={`+${employee.skills.length - 3}`}
                      size="small"
                      variant="outlined"
                    />
                  )}
                </Box>
              )}
            </Box>
            <Box>
              {onEdit && (
                <Tooltip title="Edit Employee">
                  <IconButton onClick={() => onEdit(employee)} size="small">
                    <EditIcon />
                  </IconButton>
                </Tooltip>
              )}
              {onDelete && (
                <Tooltip title="Delete Employee">
                  <IconButton onClick={() => onDelete(employee.id)} size="small" color="error">
                    <DeleteIcon />
                  </IconButton>
                </Tooltip>
              )}
            </Box>
          </Box>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardContent sx={{ flex: 1 }}>
        <Box display="flex" alignItems="center" gap={2} mb={2}>
          <Avatar sx={{ width: 48, height: 48 }}>
            <PersonIcon />
          </Avatar>
          <Box flex={1}>
            <Typography variant="h6" gutterBottom>
              {employee.name}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {employee.department}
            </Typography>
          </Box>
          <Box>
            {onEdit && (
              <Tooltip title="Edit Employee">
                <IconButton onClick={() => onEdit(employee)} size="small">
                  <EditIcon />
                </IconButton>
              </Tooltip>
            )}
            {onDelete && (
              <Tooltip title="Delete Employee">
                <IconButton onClick={() => onDelete(employee.id)} size="small" color="error">
                  <DeleteIcon />
                </IconButton>
              </Tooltip>
            )}
          </Box>
        </Box>

        <Box display="flex" alignItems="center" gap={1} mb={1}>
          <EmailIcon fontSize="small" color="action" />
          <Typography variant="body2" color="text.secondary" noWrap>
            {employee.email}
          </Typography>
        </Box>

        {employee.location && (
          <Box display="flex" alignItems="center" gap={1} mb={2}>
            <LocationIcon fontSize="small" color="action" />
            <Typography variant="body2" color="text.secondary">
              {employee.location}
            </Typography>
          </Box>
        )}

        <Chip label={employee.level} size="small" color="primary" sx={{ mb: 2 }} />

        {employee.skills && employee.skills.length > 0 && (
          <Box>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              Skills:
            </Typography>
            <Box display="flex" flexWrap="wrap" gap={0.5}>
              {employee.skills.slice(0, 4).map((skill, index) => (
                <Chip
                  key={index}
                  label={skill.name}
                  size="small"
                  variant="outlined"
                />
              ))}
              {employee.skills.length > 4 && (
                <Chip
                  label={`+${employee.skills.length - 4}`}
                  size="small"
                  variant="outlined"
                />
              )}
            </Box>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}