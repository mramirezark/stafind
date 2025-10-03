'use client'

import React, { useState, useMemo } from 'react'
import {
  Box,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TextField,
  InputAdornment,
  Chip,
  IconButton,
  Tooltip,
  Menu,
  MenuItem,
  FormControl,
  InputLabel,
  Select,
  SelectChangeEvent,
  Typography,
  Alert,
  CircularProgress,
  Stack,
  Button,
  Collapse,
} from '@mui/material'
import {
  Search as SearchIcon,
  MoreVert as MoreVertIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
  FilterList as FilterIcon,
  Sort as SortIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  UnfoldMore as ExpandAllIcon,
  UnfoldLess as CollapseAllIcon,
} from '@mui/icons-material'

import { Skill, Category, SkillFilters } from '@/types'

interface SkillListProps {
  skills: Skill[]
  categories: Category[]
  onEdit: (skill: Skill) => void
  onDelete: (id: number) => void
  onRefresh: () => void
  loading?: boolean
}

const SkillList: React.FC<SkillListProps> = ({
  skills,
  categories,
  onEdit,
  onDelete,
  onRefresh,
  loading = false,
}) => {
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(10)
  const [filters, setFilters] = useState<SkillFilters>({
    search: '',
    category: '',
    sortBy: 'name',
    sortOrder: 'asc',
  })
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [selectedSkill, setSelectedSkill] = useState<Skill | null>(null)
  const [deleteConfirm, setDeleteConfirm] = useState(false)
  const [groupByCategory, setGroupByCategory] = useState(true)
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set())

  const filteredSkills = useMemo(() => {
    let filtered = [...skills]

    // Search filter
    if (filters.search) {
      filtered = filtered.filter(skill =>
        skill.name.toLowerCase().includes(filters.search.toLowerCase())
      )
    }

    // Category filter
    if (filters.category) {
      filtered = filtered.filter(skill =>
        skill.categories?.some(cat => cat.name === filters.category)
      )
    }

    // Sort
    filtered.sort((a, b) => {
      let aValue: any, bValue: any

      switch (filters.sortBy) {
        case 'name':
          aValue = a.name.toLowerCase()
          bValue = b.name.toLowerCase()
          break
        case 'employee_count':
          aValue = a.employee_count || 0
          bValue = b.employee_count || 0
          break
        case 'created_at':
          aValue = new Date(a.created_at || 0).getTime()
          bValue = new Date(b.created_at || 0).getTime()
          break
        default:
          aValue = a.name.toLowerCase()
          bValue = b.name.toLowerCase()
      }

      if (filters.sortOrder === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
      }
    })

    return filtered
  }, [skills, filters])

  // Group skills by category
  const groupedSkills = useMemo(() => {
    const groups: { [key: string]: Skill[] } = {}
    
    filteredSkills.forEach(skill => {
      // Use the first category as the primary category for grouping
      const primaryCategory = skill.categories?.[0]?.name || 'Uncategorized'
      
      if (!groups[primaryCategory]) {
        groups[primaryCategory] = []
      }
      groups[primaryCategory].push(skill)
    })

    // Sort categories alphabetically
    const sortedGroups = Object.keys(groups).sort().reduce((acc, key) => {
      acc[key] = groups[key]
      return acc
    }, {} as { [key: string]: Skill[] })

    return sortedGroups
  }, [filteredSkills])

  // Initialize expanded categories when groupedSkills changes
  React.useEffect(() => {
    if (groupByCategory && Object.keys(groupedSkills).length > 0) {
      const allCategories = Object.keys(groupedSkills)
      setExpandedCategories(new Set(allCategories))
    }
  }, [groupedSkills, groupByCategory])

  // Toggle category expansion
  const toggleCategory = (categoryName: string) => {
    setExpandedCategories(prev => {
      const newSet = new Set(prev)
      if (newSet.has(categoryName)) {
        newSet.delete(categoryName)
      } else {
        newSet.add(categoryName)
      }
      return newSet
    })
  }

  // Expand all categories
  const expandAllCategories = () => {
    const allCategories = Object.keys(groupedSkills)
    setExpandedCategories(new Set(allCategories))
  }

  // Collapse all categories
  const collapseAllCategories = () => {
    setExpandedCategories(new Set())
  }

  // Show pagination only in ungrouped view or when there are many categories
  const shouldShowPagination = !groupByCategory || Object.keys(groupedSkills).length > 10

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFilters(prev => ({ ...prev, search: event.target.value }))
    setPage(0)
  }

  const handleCategoryChange = (event: SelectChangeEvent) => {
    setFilters(prev => ({ ...prev, category: event.target.value }))
    setPage(0)
  }

  const handleSortChange = (sortBy: SkillFilters['sortBy']) => {
    setFilters(prev => ({
      ...prev,
      sortBy,
      sortOrder: prev.sortBy === sortBy && prev.sortOrder === 'asc' ? 'desc' : 'asc',
    }))
  }

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, skill: Skill) => {
    setAnchorEl(event.currentTarget)
    setSelectedSkill(skill)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
    setSelectedSkill(null)
  }

  const handleEdit = () => {
    if (selectedSkill) {
      onEdit(selectedSkill)
    }
    handleMenuClose()
  }

  const handleDeleteClick = () => {
    setDeleteConfirm(true)
    handleMenuClose()
  }

  const handleDeleteConfirm = async () => {
    if (selectedSkill) {
      try {
        await onDelete(selectedSkill.id)
        setDeleteConfirm(false)
      } catch (error) {
        console.error('Error deleting skill:', error)
      }
    }
  }

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage)
  }

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10))
    setPage(0)
  }

  const clearFilters = () => {
    setFilters({
      search: '',
      category: '',
      sortBy: 'name',
      sortOrder: 'asc',
    })
    setPage(0)
  }

  return (
    <Box>
      {/* Filters and Search */}
      <Paper sx={{ p: 2, mb: 2 }}>
        <Stack direction="row" spacing={2} alignItems="center" flexWrap="wrap">
          <TextField
            placeholder="Search skills..."
            value={filters.search}
            onChange={handleSearchChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
            sx={{ minWidth: 200 }}
          />

          <FormControl sx={{ minWidth: 150 }}>
            <InputLabel>Category</InputLabel>
            <Select
              value={filters.category}
              onChange={handleCategoryChange}
              label="Category"
            >
              <MenuItem value="">All Categories</MenuItem>
              {categories.map(category => (
                <MenuItem key={category.id} value={category.name}>
                  {category.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Button
            startIcon={<SortIcon />}
            onClick={() => handleSortChange('name')}
            variant={filters.sortBy === 'name' ? 'contained' : 'outlined'}
            size="small"
          >
            Name {filters.sortBy === 'name' && (filters.sortOrder === 'asc' ? '↑' : '↓')}
          </Button>

          <Button
            startIcon={<FilterIcon />}
            onClick={() => handleSortChange('employee_count')}
            variant={filters.sortBy === 'employee_count' ? 'contained' : 'outlined'}
            size="small"
          >
            Usage {filters.sortBy === 'employee_count' && (filters.sortOrder === 'asc' ? '↑' : '↓')}
          </Button>

          <Button
            onClick={clearFilters}
            variant="outlined"
            size="small"
          >
            Clear Filters
          </Button>

          <Button
            onClick={() => setGroupByCategory(!groupByCategory)}
            variant={groupByCategory ? 'contained' : 'outlined'}
            size="small"
            startIcon={<FilterIcon />}
          >
            {groupByCategory ? 'Grouped' : 'List'}
          </Button>

          {groupByCategory && (
            <>
              <Button
                onClick={expandAllCategories}
                variant="outlined"
                size="small"
                startIcon={<ExpandAllIcon />}
              >
                Expand All
              </Button>
              <Button
                onClick={collapseAllCategories}
                variant="outlined"
                size="small"
                startIcon={<CollapseAllIcon />}
              >
                Collapse All
              </Button>
            </>
          )}

          <Box sx={{ flexGrow: 1 }} />

          <Tooltip title="Refresh">
            <IconButton onClick={onRefresh} disabled={loading}>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
        </Stack>
      </Paper>

      {/* Skills Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Categories</TableCell>
                <TableCell align="right">Employee Count</TableCell>
                <TableCell align="right">Created</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : (groupByCategory ? Object.keys(groupedSkills).length === 0 : filteredSkills.length === 0) ? (
                <TableRow>
                  <TableCell colSpan={5} align="center" sx={{ py: 4 }}>
                    <Typography variant="body2" color="text.secondary">
                      {filters.search || filters.category
                        ? 'No skills match your filters'
                        : 'No skills found'}
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : groupByCategory ? (
                // Grouped view
                Object.entries(groupedSkills).map(([categoryName, categorySkills]) => {
                  const isExpanded = expandedCategories.has(categoryName)
                  
                  return (
                    <React.Fragment key={categoryName}>
                      {/* Category Header Row */}
                      <TableRow 
                        sx={{ 
                          backgroundColor: 'grey.50',
                          cursor: 'pointer',
                          '&:hover': { backgroundColor: 'grey.100' }
                        }}
                        onClick={() => toggleCategory(categoryName)}
                      >
                        <TableCell colSpan={5} sx={{ py: 1 }}>
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <IconButton
                              size="small"
                              onClick={(e) => {
                                e.stopPropagation()
                                toggleCategory(categoryName)
                              }}
                              sx={{ p: 0.5 }}
                            >
                              {isExpanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                            </IconButton>
                            <Typography variant="subtitle2" fontWeight="bold" color="primary">
                              {categoryName}
                            </Typography>
                            <Chip
                              label={`${categorySkills.length} skill${categorySkills.length !== 1 ? 's' : ''}`}
                              size="small"
                              variant="outlined"
                              color="primary"
                            />
                          </Box>
                        </TableCell>
                      </TableRow>
                      {/* Skills in this category - with collapse animation */}
                      <TableRow>
                        <TableCell colSpan={5} sx={{ p: 0, border: 0 }}>
                          <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                            <Table size="small">
                              <TableBody>
                                {categorySkills.map(skill => (
                                  <TableRow key={skill.id} hover>
                                    <TableCell sx={{ pl: 4, borderLeft: '4px solid', borderLeftColor: 'primary.main' }}>
                                      <Typography variant="body2" fontWeight="medium">
                                        {skill.name}
                                      </Typography>
                                    </TableCell>
                                    <TableCell>
                                      <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                                        {skill.categories?.map(category => (
                                          <Chip
                                            key={category.id}
                                            label={category.name}
                                            size="small"
                                            variant="outlined"
                                            color={category.name === categoryName ? 'primary' : 'default'}
                                          />
                                        )) || (
                                          <Typography variant="caption" color="text.secondary">
                                            No categories
                                          </Typography>
                                        )}
                                      </Box>
                                    </TableCell>
                                    <TableCell align="right">
                                      <Typography variant="body2">
                                        {skill.employee_count || 0}
                                      </Typography>
                                    </TableCell>
                                    <TableCell align="right">
                                      <Typography variant="caption" color="text.secondary">
                                        {skill.created_at
                                          ? new Date(skill.created_at).toLocaleDateString()
                                          : '-'}
                                      </Typography>
                                    </TableCell>
                                    <TableCell align="center">
                                      <IconButton
                                        onClick={(e) => handleMenuOpen(e, skill)}
                                        size="small"
                                      >
                                        <MoreVertIcon />
                                      </IconButton>
                                    </TableCell>
                                  </TableRow>
                                ))}
                              </TableBody>
                            </Table>
                          </Collapse>
                        </TableCell>
                      </TableRow>
                    </React.Fragment>
                  )
                })
              ) : (
                // Ungrouped view (original list view)
                filteredSkills.map(skill => (
                  <TableRow key={skill.id} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {skill.name}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap' }}>
                        {skill.categories?.map(category => (
                          <Chip
                            key={category.id}
                            label={category.name}
                            size="small"
                            variant="outlined"
                          />
                        )) || (
                          <Typography variant="caption" color="text.secondary">
                            No categories
                          </Typography>
                        )}
                      </Box>
                    </TableCell>
                    <TableCell align="right">
                      <Typography variant="body2">
                        {skill.employee_count || 0}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Typography variant="caption" color="text.secondary">
                        {skill.created_at
                          ? new Date(skill.created_at).toLocaleDateString()
                          : '-'}
                      </Typography>
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        onClick={(e) => handleMenuOpen(e, skill)}
                        size="small"
                      >
                        <MoreVertIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        {shouldShowPagination && (
          <TablePagination
            rowsPerPageOptions={[5, 10, 25, 50]}
            component="div"
            count={filteredSkills.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
        )}
      </Paper>

      {/* Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleEdit}>
          <EditIcon sx={{ mr: 1 }} />
          Edit
        </MenuItem>
        <MenuItem onClick={handleDeleteClick} sx={{ color: 'error.main' }}>
          <DeleteIcon sx={{ mr: 1 }} />
          Delete
        </MenuItem>
      </Menu>

      {/* Delete Confirmation Dialog */}
      {deleteConfirm && selectedSkill && (
        <Alert
          severity="warning"
          action={
            <Stack direction="row" spacing={1}>
              <Button
                size="small"
                onClick={() => setDeleteConfirm(false)}
              >
                Cancel
              </Button>
              <Button
                size="small"
                color="error"
                onClick={handleDeleteConfirm}
              >
                Delete
              </Button>
            </Stack>
          }
          sx={{ mt: 2 }}
        >
          Are you sure you want to delete the skill "{selectedSkill.name}"? This action cannot be undone.
        </Alert>
      )}
    </Box>
  )
}

export default SkillList
