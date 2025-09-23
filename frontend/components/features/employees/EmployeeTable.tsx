'use client'

import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Avatar,
  Chip,
  IconButton,
  Tooltip,
  Box,
  Typography,
  TablePagination,
} from '@mui/material'
import {
  Person as PersonIcon,
  Email as EmailIcon,
  LocationOn as LocationIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material'
import { Employee } from '@/types'
import { EmployeeTableProps } from './interfaces'

export function EmployeeTable({ 
  employees, 
  onEdit, 
  onDelete, 
  page, 
  rowsPerPage, 
  onPageChange, 
  onRowsPerPageChange 
}: EmployeeTableProps) {
  return (
    <Paper>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Employee</TableCell>
              <TableCell>Department</TableCell>
              <TableCell>Level</TableCell>
              <TableCell>Location</TableCell>
              <TableCell>Skills</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {employees.map((employee) => (
              <TableRow key={employee.id} hover>
                <TableCell>
                  <Box display="flex" alignItems="center" gap={2}>
                    <Avatar sx={{ width: 40, height: 40 }}>
                      <PersonIcon />
                    </Avatar>
                    <Box>
                      <Typography variant="subtitle2" fontWeight="medium">
                        {employee.name}
                      </Typography>
                      <Box display="flex" alignItems="center" gap={0.5}>
                        <EmailIcon fontSize="small" color="action" />
                        <Typography variant="body2" color="text.secondary">
                          {employee.email}
                        </Typography>
                      </Box>
                    </Box>
                  </Box>
                </TableCell>
                <TableCell>
                  <Chip 
                    label={employee.department} 
                    size="small" 
                    color="primary" 
                    variant="outlined"
                  />
                </TableCell>
                <TableCell>
                  <Chip 
                    label={employee.level} 
                    size="small" 
                    variant="outlined"
                  />
                </TableCell>
                <TableCell>
                  {employee.location ? (
                    <Box display="flex" alignItems="center" gap={0.5}>
                      <LocationIcon fontSize="small" color="action" />
                      <Typography variant="body2">
                        {employee.location}
                      </Typography>
                    </Box>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      -
                    </Typography>
                  )}
                </TableCell>
                <TableCell>
                  {employee.skills && employee.skills.length > 0 ? (
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
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      No skills
                    </Typography>
                  )}
                </TableCell>
                <TableCell align="right">
                  <Box display="flex" gap={0.5}>
                    {onEdit && (
                      <Tooltip title="Edit Employee">
                        <IconButton 
                          size="small" 
                          onClick={() => onEdit(employee)}
                        >
                          <EditIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                    {onDelete && (
                      <Tooltip title="Delete Employee">
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => onDelete(employee.id)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                  </Box>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[5, 10, 25, 50]}
        component="div"
        count={employees.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onPageChange={onPageChange}
        onRowsPerPageChange={onRowsPerPageChange}
      />
    </Paper>
  )
}
