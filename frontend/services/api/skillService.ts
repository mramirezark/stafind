/**
 * Skill API Service
 * 
 * Handles all skill-related API calls
 */

import { BaseApiService } from './baseService'
import { Skill } from '@/types'

export class SkillService extends BaseApiService {
  protected getDomainName(): string {
    return 'skills'
  }

  // ============================================================================
  // SKILL CRUD OPERATIONS
  // ============================================================================

  /**
   * Get all skills
   */
  async getSkills(): Promise<Skill[]> {
    return this.request('GET', '/api/v1/skills')
  }

  /**
   * Get skill by ID
   */
  async getSkill(id: number): Promise<Skill> {
    return this.request('GET', `/api/v1/skills/${id}`)
  }

  /**
   * Create new skill
   */
  async createSkill(skillData: Partial<Skill>): Promise<Skill> {
    const result = await this.request('POST', '/api/v1/skills', skillData, false)
    this.clearDomainCache()
    return result as Skill
  }

  /**
   * Update skill
   */
  async updateSkill(id: number, skillData: Partial<Skill>): Promise<Skill> {
    const result = await this.request('PUT', `/api/v1/skills/${id}`, skillData, false)
    this.clearDomainCache()
    return result as Skill
  }

  /**
   * Delete skill
   */
  async deleteSkill(id: number): Promise<void> {
    await this.request('DELETE', `/api/v1/skills/${id}`, undefined, false)
    this.clearDomainCache()
  }

  // ============================================================================
  // SKILL CATEGORY OPERATIONS
  // ============================================================================

  /**
   * Get skills by category
   */
  async getSkillsByCategory(category: string): Promise<Skill[]> {
    return this.request('GET', `/api/v1/skills?category=${category}`)
  }

  /**
   * Get all skill categories
   */
  async getSkillCategories(): Promise<string[]> {
    return this.request('GET', '/api/v1/skills/categories')
  }

  /**
   * Create skill category
   */
  async createSkillCategory(categoryData: { name: string; description?: string }): Promise<any> {
    const result = await this.request('POST', '/api/v1/skills/categories', categoryData, false)
    this.clearDomainCache()
    return result
  }

  // ============================================================================
  // SKILL SEARCH OPERATIONS
  // ============================================================================

  /**
   * Search skills by name
   */
  async searchSkills(query: string): Promise<Skill[]> {
    return this.request('GET', `/api/v1/skills/search?q=${encodeURIComponent(query)}`)
  }

  /**
   * Get popular skills
   */
  async getPopularSkills(limit: number = 10): Promise<Skill[]> {
    return this.request('GET', `/api/v1/skills/popular?limit=${limit}`)
  }

  /**
   * Get trending skills
   */
  async getTrendingSkills(limit: number = 10): Promise<Skill[]> {
    return this.request('GET', `/api/v1/skills/trending?limit=${limit}`)
  }
}

// Export singleton instance
export const skillService = new SkillService()
