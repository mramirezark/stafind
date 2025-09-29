/**
 * Dashboard API Service
 * 
 * Handles all dashboard-related API calls
 */

import { BaseApiService } from './baseService'
import { DashboardStats, Employee, DashboardMetrics, TopSuggestedEmployee, SkillDemandStats } from '@/types'
import { employeeService } from './employeeService'

export class DashboardService extends BaseApiService {
  protected getDomainName(): string {
    return 'dashboard'
  }

  // ============================================================================
  // DASHBOARD STATISTICS
  // ============================================================================

  /**
   * Get dashboard statistics (extracted from metrics)
   */
  async getDashboardStats(): Promise<DashboardStats> {
    try {
      const metrics = await this.getDashboardMetrics()
      return metrics.stats
    } catch (error) {
      console.error('Failed to fetch dashboard stats:', error)
      return {
        totalEmployees: 0,
        totalRequests: 0,
        completedRequests: 0,
        pendingRequests: 0
      }
    }
  }

  /**
   * Get complete dashboard data
   */
  async getDashboardData(): Promise<{
    stats: DashboardStats
    employees: Employee[]
  }> {
    // Fetch all dashboard data in parallel
    const [stats, employees] = await Promise.all([
      this.getDashboardStats(),
      this.getRecentEmployees()
    ])

    return {
      stats,
      employees
    }
  }

  /**
   * Get comprehensive dashboard metrics
   */
  async getDashboardMetrics(): Promise<DashboardMetrics> {
    try {
      const response = await this.request('GET', '/api/v1/dashboard/metrics')
      console.log('Dashboard metrics response:', response)
      
      // Handle different response structures
      if (response && typeof response === 'object') {
        let metricsData = response
        
        // If response has a data property
        if ('data' in response && response.data) {
          metricsData = response.data
        }
        
        // Map snake_case to camelCase if needed
        if (metricsData && typeof metricsData === 'object') {
          console.log('Raw metrics data:', metricsData)
          console.log('stats field:', (metricsData as any).stats)
          
          // Handle stats field mapping
          let statsData = (metricsData as any).stats || metricsData
          console.log('Stats data:', statsData)
          console.log('total_employees in stats:', (statsData as any).total_employees)
          
          const mappedMetrics: DashboardMetrics = {
            stats: {
              totalEmployees: (statsData as any).total_employees || (statsData as any).totalEmployees || 0,
              totalRequests: (statsData as any).total_requests || (statsData as any).totalRequests || 0,
              completedRequests: (statsData as any).completed_requests || (statsData as any).completedRequests || 0,
              pendingRequests: (statsData as any).pending_requests || (statsData as any).pendingRequests || 0
            },
            most_requested_skills: (metricsData as any).most_requested_skills || (metricsData as any).mostRequestedSkills || [],
            top_suggested_employees: (metricsData as any).top_suggested_employees || (metricsData as any).topSuggestedEmployees || [],
            recent_requests: (metricsData as any).recent_requests || (metricsData as any).recentRequests || []
          }
          console.log('Mapped metrics:', mappedMetrics)
          return mappedMetrics
        }
      }
      
      console.warn('Unexpected dashboard metrics response structure:', response)
      return {
        stats: {
          totalEmployees: 0,
          totalRequests: 0,
          completedRequests: 0,
          pendingRequests: 0
        },
        most_requested_skills: [],
        top_suggested_employees: [],
        recent_requests: []
      }
    } catch (error) {
      console.error('Failed to fetch dashboard metrics:', error)
      // Return default metrics on error
      return {
        stats: {
          totalEmployees: 0,
          totalRequests: 0,
          completedRequests: 0,
          pendingRequests: 0
        },
        most_requested_skills: [],
        top_suggested_employees: [],
        recent_requests: []
      }
    }
  }

  /**
   * Get top suggested employees
   */
  async getTopSuggestedEmployees(limit: number = 5): Promise<TopSuggestedEmployee[]> {
    try {
      const response = await this.request('GET', `/api/v1/dashboard/top-employees?limit=${limit}`)
      console.log('Top suggested employees response:', response)
      
      // Handle different response structures
      if (Array.isArray(response)) {
        return response
      }
      
      if (response && typeof response === 'object') {
        // If response has a data property
        if ('data' in response && Array.isArray(response.data)) {
          return response.data
        }
        // If response has employees property
        if ('employees' in response && Array.isArray(response.employees)) {
          return response.employees
        }
      }
      
      console.warn('Unexpected top suggested employees response structure:', response)
      return []
    } catch (error) {
      console.error('Failed to fetch top suggested employees:', error)
      return []
    }
  }

  // ============================================================================
  // RECENT DATA OPERATIONS
  // ============================================================================


  /**
   * Get recent employees
   */
  async getRecentEmployees(limit: number = 5): Promise<Employee[]> {
    try {
      return await this.request('GET', `/api/v1/dashboard/recent-employees?limit=${limit}`)
    } catch (error) {
      console.error('Failed to fetch recent employees:', error)
      return []
    }
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
    try {
      return await this.request('GET', '/api/v1/dashboard/skill-demand')
    } catch (error) {
      console.error('Failed to fetch skill demand stats:', error)
      return []
    }
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
