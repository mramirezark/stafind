import { endpoints } from '@/lib/api'

export interface AIAgentRequest {
  id: number
  teams_message_id: string
  channel_id: string
  user_id: string
  user_name: string
  message_text: string
  attachment_url?: string
  extracted_text?: string
  extracted_skills?: string[]
  status: 'pending' | 'processing' | 'completed' | 'failed'
  error?: string
  created_at: string
  processed_at?: string
}

export interface AIAgentMatch {
  employee_id: number
  employee_name: string
  employee_email: string
  position: string
  seniority: string
  location: string
  current_project: string
  resume_link: string
  match_score: number
  matching_skills: string[]
  ai_summary: string
  bio: string
}

export interface AIAgentResponse {
  request_id: number
  matches: AIAgentMatch[]
  summary: string
  processing_time_ms: number
  status: string
  error?: string
}

export interface SkillExtractionResponse {
  skills: string[]
  text: string
}

export interface CreateAIAgentRequest {
  teams_message_id: string
  channel_id: string
  user_id: string
  user_name: string
  message_text: string
  attachment_url?: string
}

class AIAgentService {
  async getRequests(limit: number = 50, offset: number = 0): Promise<AIAgentRequest[]> {
    const response = await endpoints.aiAgent.getRequests({ limit, offset })
    return response.data
  }

  async getRequest(id: number): Promise<AIAgentRequest> {
    const response = await endpoints.aiAgent.getRequest(id)
    return response.data
  }

  async getResponse(id: number): Promise<AIAgentResponse> {
    const response = await endpoints.aiAgent.getResponse(id)
    return response.data
  }

  async processRequest(id: number): Promise<AIAgentResponse> {
    const response = await endpoints.aiAgent.processRequest(id)
    return response.data
  }

  async processNewRequest(data: CreateAIAgentRequest): Promise<AIAgentResponse> {
    const response = await endpoints.aiAgent.process(data)
    return response.data
  }

  async extractSkills(text: string): Promise<SkillExtractionResponse> {
    const response = await endpoints.aiAgent.extractSkills({ text })
    return response.data
  }

  clearCache(pattern?: string) {
    // No-op for now, can implement caching later if needed
  }
}

export const aiAgentService = new AIAgentService()
