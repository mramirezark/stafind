/**
 * Job Request API Service
 * 
 * Handles all job request-related API calls
 */

import { BaseApiService } from './baseService'
import { JobRequest } from '@/types'

export class JobRequestService extends BaseApiService {
  protected getDomainName(): string {
    return 'job-requests'
  }

  // ============================================================================
  // JOB REQUEST CRUD OPERATIONS
  // ============================================================================

  /**
   * Get all job requests
   */
  async getJobRequests(): Promise<JobRequest[]> {
    return this.request('GET', '/api/v1/job-requests')
  }

  /**
   * Get job request by ID
   */
  async getJobRequest(id: number): Promise<JobRequest> {
    return this.request('GET', `/api/v1/job-requests/${id}`)
  }

  /**
   * Create new job request
   */
  async createJobRequest(jobRequestData: Partial<JobRequest>): Promise<JobRequest> {
    const result = await this.request('POST', '/api/v1/job-requests', jobRequestData, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Update job request
   */
  async updateJobRequest(id: number, jobRequestData: Partial<JobRequest>): Promise<JobRequest> {
    const result = await this.request('PUT', `/api/v1/job-requests/${id}`, jobRequestData, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Delete job request
   */
  async deleteJobRequest(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/job-requests/${id}`, undefined, false)
    this.clearDomainCache()
  }

  // ============================================================================
  // JOB REQUEST MATCHING OPERATIONS
  // ============================================================================

  /**
   * Get matches for a job request
   */
  async getJobRequestMatches(id: number): Promise<any[]> {
    return this.request('GET', `/api/v1/job-requests/${id}/matches`)
  }

  /**
   * Create match for job request
   */
  async createMatch(jobRequestId: number, employeeId: number, matchData: any): Promise<any> {
    const result = await this.request('POST', `/api/v1/job-requests/${jobRequestId}/matches`, {
      employee_id: employeeId,
      ...matchData
    }, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Delete match for job request
   */
  async deleteMatch(jobRequestId: number, matchId: number): Promise<void> {
    await this.request('DELETE', `/api/v1/job-requests/${jobRequestId}/matches/${matchId}`, undefined, false)
    this.clearDomainCache()
  }

  // ============================================================================
  // JOB REQUEST STATUS OPERATIONS
  // ============================================================================

  /**
   * Update job request status
   */
  async updateJobRequestStatus(id: number, status: string): Promise<JobRequest> {
    const result = await this.request('PUT', `/api/v1/job-requests/${id}/status`, { status }, false)
    this.clearDomainCache()
    return result
  }

  /**
   * Get job requests by status
   */
  async getJobRequestsByStatus(status: string): Promise<JobRequest[]> {
    return this.request('GET', `/api/v1/job-requests?status=${status}`)
  }

  /**
   * Get job requests by department
   */
  async getJobRequestsByDepartment(department: string): Promise<JobRequest[]> {
    return this.request('GET', `/api/v1/job-requests?department=${department}`)
  }
}

// Export singleton instance
export const jobRequestService = new JobRequestService()
