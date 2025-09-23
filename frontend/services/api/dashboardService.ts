/**
 * Dashboard API Service
 * 
 * Handles all dashboard-related API calls
 */

import { BaseApiService } from './baseService'
import { DashboardStats, Employee, JobRequest } from '@/types'
import { employeeService } from './employeeService'
import { jobRequestService } from './jobRequestService'

export class DashboardService extends BaseApiService {
  protected getDomainName(): string {
    return 'dashboard'
  }

  // ============================================================================
  // DASHBOARD STATISTICS
  // ============================================================================

  /**
   * Get dashboard statistics
   */
  async getDashboardStats(): Promise<DashboardStats> {
    return this.request('GET', '/api/v1/dashboard/stats')
  }

  /**
   * Get complete dashboard data
   */
  async getDashboardData(): Promise<{
    stats: DashboardStats
    recentRequests: JobRequest[]
    employees: Employee[]
  }> {
    // Fetch all dashboard data in parallel
    const [stats, recentRequests, employees] = await Promise.all([
      this.getDashboardStats(),
      this.getRecentJobRequests(),
      this.getRecentEmployees()
    ])

    return {
      stats,
      recentRequests,
      employees
    }
  }

  // ============================================================================
  // RECENT DATA OPERATIONS
  // ============================================================================

  /**
   * Get recent job requests
   */
  async getRecentJobRequests(limit: number = 5): Promise<JobRequest[]> {
    return this.request('GET', `/api/v1/dashboard/recent-requests?limit=${limit}`)
  }

  /**
   * Get recent employees
   */
  async getRecentEmployees(limit: number = 5): Promise<Employee[]> {
    return this.request('GET', `/api/v1/dashboard/recent-employees?limit=${limit}`)
  }

  /**
   * Get recent matches
   */
  async getRecentMatches(limit: number = 5): Promise<any[]> {
    return this.request('GET', `/api/v1/dashboard/recent-matches?limit=${limit}`)
  }

  // ============================================================================
  // ANALYTICS OPERATIONS
  // ============================================================================

  /**
   * Get department statistics
   */
  async getDepartmentStats(): Promise<any[]> {
    return this.request('GET', '/api/v1/dashboard/department-stats')
  }

  /**
   * Get skill demand statistics
   */
  async getSkillDemandStats(): Promise<any[]> {
    return this.request('GET', '/api/v1/dashboard/skill-demand')
  }

  /**
   * Get matching success rate
   */
  async getMatchingSuccessRate(): Promise<{ success_rate: number; total_matches: number }> {
    return this.request('GET', '/api/v1/dashboard/matching-success-rate')
  }

  /**
   * Get monthly trends
   */
  async getMonthlyTrends(months: number = 6): Promise<any[]> {
    return this.request('GET', `/api/v1/dashboard/monthly-trends?months=${months}`)
  }

  // ============================================================================
  // ACTIVITY FEED OPERATIONS
  // ============================================================================

  /**
   * Get activity feed
   */
  async getActivityFeed(limit: number = 20): Promise<any[]> {
    return this.request('GET', `/api/v1/dashboard/activity-feed?limit=${limit}`)
  }

  /**
   * Get user activity
   */
  async getUserActivity(userId: number, limit: number = 10): Promise<any[]> {
    return this.request('GET', `/api/v1/dashboard/user-activity/${userId}?limit=${limit}`)
  }
}

// Export singleton instance
export const dashboardService = new DashboardService()
