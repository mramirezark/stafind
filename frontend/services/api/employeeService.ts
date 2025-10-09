/**
 * Employee API Service
 * 
 * Handles all employee-related API calls
 */

import { BaseApiService } from './baseService'
import { Employee, SearchCriteria } from '@/types'

export class EmployeeService extends BaseApiService {
  protected getDomainName(): string {
    return 'employees'
  }

  // ============================================================================
  // EMPLOYEE CRUD OPERATIONS
  // ============================================================================

  /**
   * Get all employees
   */
  async getEmployees(): Promise<Employee[]> {
    return this.request('GET', '/api/v1/employees')
  }

  /**
   * Get all employees (bypassing cache)
   */
  async getEmployeesFresh(): Promise<Employee[]> {
    return this.request('GET', '/api/v1/employees', undefined, false)
  }

  /**
   * Get employee by ID
   */
  async getEmployee(id: number): Promise<Employee> {
    return this.request('GET', `/api/v1/employees/${id}`)
  }

  /**
   * Create new employee
   */
  async createEmployee(employeeData: Partial<Employee>): Promise<Employee> {
    const result = await this.request('POST', '/api/v1/employees', employeeData, false)
    this.clearDomainCache()
    return result as Employee
  }

  /**
   * Update employee
   */
  async updateEmployee(id: number, employeeData: Partial<Employee>): Promise<Employee> {
    const result = await this.request('PUT', `/api/v1/employees/${id}`, employeeData, false)
    this.clearDomainCache()
    return result as Employee
  }

  /**
   * Delete employee
   */
  async deleteEmployee(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/employees/${id}`, undefined, false)
    this.clearDomainCache()
  }

  // ============================================================================
  // EMPLOYEE SEARCH OPERATIONS
  // ============================================================================

  /**
   * Search employees by criteria
   */
  async searchEmployees(searchCriteria: SearchCriteria): Promise<Employee[]> {
    return this.request('POST', '/api/v1/search', searchCriteria, false)
  }

  // ============================================================================
  // EMPLOYEE SKILLS OPERATIONS
  // ============================================================================

  /**
   * Get employee skills
   */
  async getEmployeeSkills(employeeId: number): Promise<any[]> {
    return this.request('GET', `/api/v1/employees/${employeeId}/skills`)
  }

  /**
   * Add skill to employee
   */
  async addEmployeeSkill(employeeId: number, skillData: any): Promise<any> {
    const result = await this.request('POST', `/api/v1/employees/${employeeId}/skills`, skillData, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Update employee skill
   */
  async updateEmployeeSkill(employeeId: number, skillId: number, skillData: any): Promise<any> {
    const result = await this.request('PUT', `/api/v1/employees/${employeeId}/skills/${skillId}`, skillData, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Remove skill from employee
   */
  async removeEmployeeSkill(employeeId: number, skillId: number): Promise<void> {
    await this.request('DELETE', `/api/v1/employees/${employeeId}/skills/${skillId}`, undefined, false)
    this.clearDomainCache()
  }
}

// Export singleton instance
export const employeeService = new EmployeeService()
