package queries

// Employee queries
const (
	// GetAllEmployees retrieves all employees with their basic information
	GetAllEmployees = `
		SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.created_at, e.updated_at
		FROM employees e
		ORDER BY e.name
	`

	// GetEmployeeByID retrieves a specific employee by ID
	GetEmployeeByID = `
		SELECT e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.created_at, e.updated_at
		FROM employees e
		WHERE e.id = $1
	`

	// CreateEmployee inserts a new employee and returns the created record
	CreateEmployee = `
		INSERT INTO employees (name, email, department, level, location, bio)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	// UpdateEmployee updates an existing employee
	UpdateEmployee = `
		UPDATE employees 
		SET name = $1, email = $2, department = $3, level = $4, location = $5, bio = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	// DeleteEmployee removes an employee from the database
	DeleteEmployee = `
		DELETE FROM employees WHERE id = $1
	`

	// GetEmployeeSkills retrieves all skills for a specific employee
	GetEmployeeSkills = `
		SELECT s.id, s.name, s.category, es.proficiency_level, es.years_experience
		FROM skills s
		JOIN employee_skills es ON s.id = es.skill_id
		WHERE es.employee_id = $1
		ORDER BY s.name
	`

	// AddEmployeeSkill creates a relationship between employee and skill
	AddEmployeeSkill = `
		INSERT INTO employee_skills (employee_id, skill_id, proficiency_level, years_experience)
		VALUES ($1, $2, $3, $4)
	`

	// RemoveEmployeeSkills removes all skill relationships for an employee
	RemoveEmployeeSkills = `
		DELETE FROM employee_skills WHERE employee_id = $1
	`
)

// Job Request queries
const (
	// GetAllJobRequests retrieves all job requests ordered by creation date
	GetAllJobRequests = `
		SELECT id, title, description, department, required_skills, preferred_skills,
		       experience_level, location, priority, status, created_by, created_at, updated_at
		FROM job_requests
		ORDER BY created_at DESC
	`

	// GetJobRequestByID retrieves a specific job request by ID
	GetJobRequestByID = `
		SELECT id, title, description, department, required_skills, preferred_skills,
		       experience_level, location, priority, status, created_by, created_at, updated_at
		FROM job_requests
		WHERE id = $1
	`

	// CreateJobRequest inserts a new job request
	CreateJobRequest = `
		INSERT INTO job_requests (title, description, department, required_skills, preferred_skills,
		                         experience_level, location, priority, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	// UpdateJobRequest updates an existing job request
	UpdateJobRequest = `
		UPDATE job_requests 
		SET title = $1, description = $2, department = $3, required_skills = $4, preferred_skills = $5,
		    experience_level = $6, location = $7, priority = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
	`

	// DeleteJobRequest removes a job request
	DeleteJobRequest = `
		DELETE FROM job_requests WHERE id = $1
	`
)

// Skill queries
const (
	// GetAllSkills retrieves all available skills
	GetAllSkills = `
		SELECT id, name, category FROM skills ORDER BY name
	`

	// GetSkillByID retrieves a skill by ID
	GetSkillByID = `
		SELECT id, name, category FROM skills WHERE id = $1
	`

	// GetSkillByName retrieves a skill by name
	GetSkillByName = `
		SELECT id, name, category FROM skills WHERE name = $1
	`

	// CreateSkill inserts a new skill
	CreateSkill = `
		INSERT INTO skills (name, category)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	// UpdateSkill updates an existing skill
	UpdateSkill = `
		UPDATE skills 
		SET name = $1, category = $2
		WHERE id = $3
		RETURNING id, name, category
	`

	// DeleteSkill removes a skill
	DeleteSkill = `
		DELETE FROM skills WHERE id = $1
	`
)

// Match queries
const (
	// GetMatchesByJobRequestID retrieves all matches for a job request
	GetMatchesByJobRequestID = `
		SELECT m.id, m.job_request_id, m.employee_id, m.match_score, m.matching_skills, m.notes, m.created_at,
		       e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.created_at, e.updated_at
		FROM matches m
		JOIN employees e ON m.employee_id = e.id
		WHERE m.job_request_id = $1
		ORDER BY m.match_score DESC
	`

	// CreateMatch inserts a new match record
	CreateMatch = `
		INSERT INTO matches (job_request_id, employee_id, match_score, matching_skills, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	// DeleteMatchesByJobRequestID removes all matches for a job request
	DeleteMatchesByJobRequestID = `
		DELETE FROM matches WHERE job_request_id = $1
	`

	// GetEmployeeSkillsForMatch retrieves skills for match display
	GetEmployeeSkillsForMatch = `
		SELECT s.id, s.name, s.category
		FROM skills s
		JOIN employee_skills es ON s.id = es.skill_id
		WHERE es.employee_id = $1
		ORDER BY s.name
	`
)
