-- Create employees table
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(100),
    level VARCHAR(50),
    location VARCHAR(100),
    bio TEXT,
    current_project VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Create skills table (without category column)
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Create skills_categories junction table (many-to-many relationship)
CREATE TABLE skills_categories (
    skill_id INTEGER REFERENCES skills(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (skill_id, category_id)
);

-- Create employee_skills junction table
CREATE TABLE employee_skills (
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    skill_id INTEGER REFERENCES skills(id) ON DELETE CASCADE,
    proficiency_level INTEGER CHECK (proficiency_level >= 1 AND proficiency_level <= 5),
    years_experience DECIMAL(3,1),
    PRIMARY KEY (employee_id, skill_id)
);

-- Create matches table (for AI agent matches)
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees(id) ON DELETE CASCADE,
    match_score DECIMAL(5,2),
    matching_skills JSONB,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_employee_skills_employee_id ON employee_skills(employee_id);
CREATE INDEX idx_employee_skills_skill_id ON employee_skills(skill_id);
CREATE INDEX idx_matches_employee_id ON matches(employee_id);
CREATE INDEX idx_matches_score ON matches(match_score DESC);
CREATE INDEX idx_employees_current_project ON employees(current_project);

-- Indexes for categories and skills_categories
CREATE INDEX idx_skills_categories_skill_id ON skills_categories(skill_id);
CREATE INDEX idx_skills_categories_category_id ON skills_categories(category_id);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_skills_name ON skills(name);
