'use client'

import { useState } from 'react'
import { Grid, Paper, Typography } from '@mui/material'
import { EmployeeListProps } from './interfaces'
import { EmployeeCard } from './EmployeeCard'
import { EmployeeTable } from './EmployeeTable'

export function EmployeeList({ employees, viewMode, onEdit, onDelete }: EmployeeListProps) {
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(10)

  const handlePageChange = (event: unknown, newPage: number) => {
    setPage(newPage)
  }

  const handleRowsPerPageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10))
    setPage(0)
  }

  if (employees.length === 0) {
    return (
      <Paper sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          No employees found
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Try adjusting your search criteria or add a new employee.
        </Typography>
      </Paper>
    )
  }

  if (viewMode === 'table') {
    const paginatedEmployees = employees.slice(
      page * rowsPerPage,
      page * rowsPerPage + rowsPerPage
    )

    return (
      <EmployeeTable
        employees={paginatedEmployees}
        onEdit={onEdit}
        onDelete={onDelete}
        page={page}
        rowsPerPage={rowsPerPage}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
      />
    )
  }

  return (
    <Grid container spacing={2}>
      {employees.map((employee) => (
        <Grid item xs={12} sm={6} md={4} lg={3} key={employee.id}>
          <EmployeeCard
            employee={employee}
            viewMode={viewMode}
            onEdit={onEdit}
            onDelete={onDelete}
          />
        </Grid>
      ))}
    </Grid>
  )
}