/**
 * Resume Parser Service
 * 
 * This service handles parsing of resume files (PDF, DOC, DOCX) to extract
 * employee information for bulk import operations.
 */

export interface ParsedResumeData {
  name: string
  email: string
  phone?: string
  skills: string[]
  experience: string
  location?: string
  bio?: string
  education?: string[]
  workHistory?: WorkExperience[]
  certifications?: string[]
}

export interface WorkExperience {
  company: string
  position: string
  duration: string
  description?: string
}

export class ResumeParserService {
  private static instance: ResumeParserService

  public static getInstance(): ResumeParserService {
    if (!ResumeParserService.instance) {
      ResumeParserService.instance = new ResumeParserService()
    }
    return ResumeParserService.instance
  }

  /**
   * Parse a resume file and extract employee data
   */
  async parseResume(file: File): Promise<ParsedResumeData> {
    const fileType = this.getFileType(file)
    
    switch (fileType) {
      case 'pdf':
        return this.parsePDF(file)
      case 'docx':
        return this.parseDOCX(file)
      case 'doc':
        return this.parseDOC(file)
      default:
        throw new Error(`Unsupported file type: ${fileType}`)
    }
  }

  /**
   * Parse multiple resume files
   */
  async parseMultipleResumes(files: File[]): Promise<ParsedResumeData[]> {
    const results = await Promise.allSettled(
      files.map(file => this.parseResume(file))
    )

    return results
      .filter((result): result is PromiseFulfilledResult<ParsedResumeData> => 
        result.status === 'fulfilled'
      )
      .map(result => result.value)
  }

  /**
   * Get file type from file extension
   */
  private getFileType(file: File): string {
    const extension = file.name.split('.').pop()?.toLowerCase()
    
    switch (extension) {
      case 'pdf':
        return 'pdf'
      case 'docx':
        return 'docx'
      case 'doc':
        return 'doc'
      default:
        throw new Error(`Unsupported file type: ${extension}`)
    }
  }

  /**
   * Parse PDF resume
   */
  private async parsePDF(file: File): Promise<ParsedResumeData> {
    // TODO: Implement PDF parsing using a library like pdf-parse or pdf2pic
    // For now, return mock data
    return this.generateMockData(file.name)
  }

  /**
   * Parse DOCX resume
   */
  private async parseDOCX(file: File): Promise<ParsedResumeData> {
    // TODO: Implement DOCX parsing using a library like mammoth
    // For now, return mock data
    return this.generateMockData(file.name)
  }

  /**
   * Parse DOC resume
   */
  private async parseDOC(file: File): Promise<ParsedResumeData> {
    // TODO: Implement DOC parsing using a library like mammoth
    // For now, return mock data
    return this.generateMockData(file.name)
  }

  /**
   * Generate mock data for development/testing
   */
  private generateMockData(fileName: string): ParsedResumeData {
    const names = [
      'John Smith', 'Jane Doe', 'Michael Johnson', 'Sarah Wilson',
      'David Brown', 'Emily Davis', 'Robert Miller', 'Lisa Garcia',
      'James Martinez', 'Jennifer Anderson'
    ]
    
    const skills = [
      'JavaScript', 'TypeScript', 'React', 'Node.js', 'Python',
      'Java', 'C#', 'SQL', 'MongoDB', 'AWS', 'Docker', 'Kubernetes',
      'Git', 'Agile', 'Scrum', 'Machine Learning', 'Data Analysis'
    ]

    const locations = [
      'San Francisco, CA', 'New York, NY', 'Seattle, WA', 'Austin, TX',
      'Boston, MA', 'Chicago, IL', 'Denver, CO', 'Portland, OR'
    ]

    const experiences = ['Junior', 'Mid-level', 'Senior', 'Lead', 'Principal']

    const randomName = names[Math.floor(Math.random() * names.length)]
    const randomSkills = skills
      .sort(() => 0.5 - Math.random())
      .slice(0, Math.floor(Math.random() * 5) + 3)
    const randomLocation = locations[Math.floor(Math.random() * locations.length)]
    const randomExperience = experiences[Math.floor(Math.random() * experiences.length)]

    return {
      name: randomName,
      email: `${randomName.toLowerCase().replace(' ', '.')}@example.com`,
      phone: `+1-555-${Math.floor(Math.random() * 9000) + 1000}`,
      skills: randomSkills,
      experience: randomExperience,
      location: randomLocation,
      bio: `Experienced ${randomExperience.toLowerCase()} developer with expertise in ${randomSkills.slice(0, 3).join(', ')}.`,
      education: [
        'Bachelor of Science in Computer Science',
        'Master of Science in Software Engineering'
      ],
      workHistory: [
        {
          company: 'Tech Corp',
          position: `${randomExperience} Developer`,
          duration: '2020 - Present',
          description: 'Developed and maintained web applications'
        }
      ],
      certifications: [
        'AWS Certified Developer',
        'Google Cloud Professional'
      ]
    }
  }

  /**
   * Validate parsed resume data
   */
  validateParsedData(data: ParsedResumeData): { isValid: boolean; errors: string[] } {
    const errors: string[] = []

    if (!data.name || data.name.trim().length === 0) {
      errors.push('Name is required')
    }

    if (!data.email || !this.isValidEmail(data.email)) {
      errors.push('Valid email is required')
    }

    if (!data.skills || data.skills.length === 0) {
      errors.push('At least one skill is required')
    }

    if (!data.experience || data.experience.trim().length === 0) {
      errors.push('Experience level is required')
    }

    return {
      isValid: errors.length === 0,
      errors
    }
  }

  /**
   * Validate email format
   */
  private isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }

  /**
   * Extract skills from text using keyword matching
   */
  extractSkillsFromText(text: string): string[] {
    const skillKeywords = [
      'javascript', 'typescript', 'react', 'angular', 'vue', 'node.js', 'python',
      'java', 'c#', 'c++', 'php', 'ruby', 'go', 'rust', 'swift', 'kotlin',
      'sql', 'mongodb', 'postgresql', 'mysql', 'redis', 'elasticsearch',
      'aws', 'azure', 'gcp', 'docker', 'kubernetes', 'jenkins', 'git',
      'agile', 'scrum', 'devops', 'ci/cd', 'microservices', 'api',
      'machine learning', 'ai', 'data science', 'analytics', 'blockchain'
    ]

    const foundSkills: string[] = []
    const lowerText = text.toLowerCase()

    skillKeywords.forEach(skill => {
      if (lowerText.includes(skill)) {
        foundSkills.push(skill.charAt(0).toUpperCase() + skill.slice(1))
      }
    })

    return Array.from(new Set(foundSkills)) // Remove duplicates
  }

  /**
   * Extract experience level from text
   */
  extractExperienceLevel(text: string): string {
    const lowerText = text.toLowerCase()
    
    if (lowerText.includes('senior') || lowerText.includes('lead') || lowerText.includes('principal')) {
      return 'Senior'
    } else if (lowerText.includes('mid') || lowerText.includes('intermediate')) {
      return 'Mid-level'
    } else if (lowerText.includes('junior') || lowerText.includes('entry')) {
      return 'Junior'
    } else if (lowerText.includes('intern') || lowerText.includes('internship')) {
      return 'Intern'
    }
    
    return 'Mid-level' // Default
  }
}

export const resumeParserService = ResumeParserService.getInstance()
