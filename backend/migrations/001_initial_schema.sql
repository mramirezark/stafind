-- Create employees table
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(100),
    level VARCHAR(50),
    location VARCHAR(100),
    bio TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create skills table
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create employee_skills junction table
CREATE TABLE employee_skills (
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    skill_id INTEGER REFERENCES skills(id) ON DELETE CASCADE,
    proficiency_level INTEGER CHECK (proficiency_level >= 1 AND proficiency_level <= 5),
    years_experience DECIMAL(3,1),
    PRIMARY KEY (employee_id, skill_id)
);

-- Create job_requests table
CREATE TABLE job_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    department VARCHAR(100),
    required_skills JSONB,
    preferred_skills JSONB,
    experience_level VARCHAR(50),
    location VARCHAR(100),
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'open',
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create matches table
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    job_request_id INTEGER REFERENCES job_requests(id) ON DELETE CASCADE,
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    match_score DECIMAL(5,2),
    matching_skills JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(job_request_id, employee_id)
);

-- Create indexes for better performance
CREATE INDEX idx_employee_skills_employee_id ON employee_skills(employee_id);
CREATE INDEX idx_employee_skills_skill_id ON employee_skills(skill_id);
CREATE INDEX idx_job_requests_status ON job_requests(status);
CREATE INDEX idx_job_requests_department ON job_requests(department);
CREATE INDEX idx_matches_job_request_id ON matches(job_request_id);
CREATE INDEX idx_matches_employee_id ON matches(employee_id);
CREATE INDEX idx_matches_score ON matches(match_score DESC);

-- Insert some sample skills
INSERT INTO skills (name, category) VALUES
('JavaScript', 'Programming Language'),
('Python', 'Programming Language'),
('Go', 'Programming Language'),
('React', 'Frontend Framework'),
('Node.js', 'Backend Framework'),
('PostgreSQL', 'Database'),
('Docker', 'DevOps'),
('AWS', 'Cloud Platform'),
('Kubernetes', 'DevOps'),
('GraphQL', 'API'),
('TypeScript', 'Programming Language'),
('Java', 'Programming Language'),
('Spring Boot', 'Backend Framework'),
('Vue.js', 'Frontend Framework'),
('Angular', 'Frontend Framework'),
('MongoDB', 'Database'),
('Redis', 'Database'),
('Git', 'Version Control'),
('CI/CD', 'DevOps'),
('Microservices', 'Architecture');
