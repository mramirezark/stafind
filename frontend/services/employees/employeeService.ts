import { apiClient } from '../api/client'
import { API_ENDPOINTS } from '../api/endpoints'

// Employee types
export interface Employee {
  id: number
  name: string
  email: string
  department: string
  level: 'junior' | 'mid' | 'senior'
  location: string
  bio: string
  skills: Skill[]
  created_at: string
  updated_at: string
}

export interface Skill {
  id: number
  name: string
  category: string
  created_at?: string
  updated_at?: string
}

export interface CreateEmployeeRequest {
  name: string
  email: string
  department: string
  level: string
  location: string
  bio: string
  skills: EmployeeSkillRequest[]
}

export interface UpdateEmployeeRequest extends Partial<CreateEmployeeRequest> {
  id: number
}

export interface EmployeeSkillRequest {
  skill_name: string
  proficiency_level: number
  years_experience?: number
}

export interface SearchRequest {
  required_skills?: string[]
  preferred_skills?: string[]
  department?: string
  experience_level?: string
  location?: string
  min_match_score?: number
}

export interface Match {
  id: number
  job_request_id?: number
  employee_id: number
  match_score: number
  matching_skills: string[]
  notes?: string
  employee: Employee
  created_at: string
}

/**
 * Employee Service
 * Handles all employee-related API calls
 */
export class EmployeeService {
  /**
   * Get all employees
   */
  static async getAll(): Promise<Employee[]> {
    const response = await apiClient.get(API_ENDPOINTS.EMPLOYEES.LIST)
    return response.data
  }

  /**
   * Get employee by ID
   */
  static async getById(id: number): Promise<Employee> {
    const response = await apiClient.get(API_ENDPOINTS.EMPLOYEES.GET(id))
    return response.data
  }

  /**
   * Create new employee
   */
  static async create(employeeData: CreateEmployeeRequest): Promise<Employee> {
    const response = await apiClient.post(API_ENDPOINTS.EMPLOYEES.CREATE, employeeData)
    return response.data
  }

  /**
   * Update existing employee
   */
  static async update(id: number, employeeData: UpdateEmployeeRequest): Promise<Employee> {
    const response = await apiClient.put(API_ENDPOINTS.EMPLOYEES.UPDATE(id), employeeData)
    return response.data
  }

  /**
   * Delete employee
   */
  static async delete(id: number): Promise<void> {
    await apiClient.delete(API_ENDPOINTS.EMPLOYEES.DELETE(id))
  }

  /**
   * Search employees based on criteria
   */
  static async search(criteria: SearchRequest): Promise<Match[]> {
    const response = await apiClient.post(API_ENDPOINTS.EMPLOYEES.SEARCH, criteria)
    return response.data
  }

  /**
   * Get employees by department
   */
  static async getByDepartment(department: string): Promise<Employee[]> {
    const allEmployees = await this.getAll()
    return allEmployees.filter(emp => emp.department === department)
  }

  /**
   * Get employees by experience level
   */
  static async getByLevel(level: string): Promise<Employee[]> {
    const allEmployees = await this.getAll()
    return allEmployees.filter(emp => emp.level === level)
  }

  /**
   * Get employees by location
   */
  static async getByLocation(location: string): Promise<Employee[]> {
    const allEmployees = await this.getAll()
    return allEmployees.filter(emp => emp.location.toLowerCase().includes(location.toLowerCase()))
  }

  /**
   * Get employees with specific skill
   */
  static async getBySkill(skillName: string): Promise<Employee[]> {
    const allEmployees = await this.getAll()
    return allEmployees.filter(emp => 
      emp.skills.some(skill => skill.name.toLowerCase().includes(skillName.toLowerCase()))
    )
  }

  /**
   * Get employee statistics
   */
  static async getStats(): Promise<{
    total: number
    byDepartment: Record<string, number>
    byLevel: Record<string, number>
    byLocation: Record<string, number>
  }> {
    const employees = await this.getAll()
    
    const stats = {
      total: employees.length,
      byDepartment: {} as Record<string, number>,
      byLevel: {} as Record<string, number>,
      byLocation: {} as Record<string, number>,
    }

    employees.forEach(emp => {
      // Count by department
      stats.byDepartment[emp.department] = (stats.byDepartment[emp.department] || 0) + 1
      
      // Count by level
      stats.byLevel[emp.level] = (stats.byLevel[emp.level] || 0) + 1
      
      // Count by location
      stats.byLocation[emp.location] = (stats.byLocation[emp.location] || 0) + 1
    })

    return stats
  }
}
