-- Insert sample employees
INSERT INTO employees (name, email, department, level, location, bio, current_project) VALUES
('Alice Johnson', 'alice.johnson@company.com', 'Engineering', 'senior', 'San Francisco, CA', 'Senior full-stack engineer with expertise in React and Node.js', 'E-commerce Platform Redesign'),
('Bob Smith', 'bob.smith@company.com', 'Engineering', 'mid', 'New York, NY', 'Backend developer specializing in Go and PostgreSQL', 'API Gateway Migration'),
('Carol Davis', 'carol.davis@company.com', 'Engineering', 'senior', 'Seattle, WA', 'DevOps engineer with strong AWS and Kubernetes experience', 'Cloud Infrastructure Optimization'),
('David Wilson', 'david.wilson@company.com', 'Engineering', 'junior', 'Austin, TX', 'Frontend developer focused on React and TypeScript', 'Mobile App Development'),
('Eva Brown', 'eva.brown@company.com', 'Data Science', 'senior', 'Boston, MA', 'Data scientist with Python and machine learning expertise', 'Customer Analytics Dashboard');

-- Add some employee skills
INSERT INTO employee_skills (employee_id, skill_id, proficiency_level, years_experience) VALUES
-- Alice Johnson (employee_id: 1)
(1, (SELECT id FROM skills WHERE name = 'JavaScript'), 5, 4.5),
(1, (SELECT id FROM skills WHERE name = 'React'), 5, 4.0),
(1, (SELECT id FROM skills WHERE name = 'Node.js'), 4, 3.5),
(1, (SELECT id FROM skills WHERE name = 'TypeScript'), 4, 3.0),
-- Bob Smith (employee_id: 2)
(2, (SELECT id FROM skills WHERE name = 'Go'), 5, 3.0),
(2, (SELECT id FROM skills WHERE name = 'PostgreSQL'), 4, 2.5),
(2, (SELECT id FROM skills WHERE name = 'Docker'), 4, 2.0),
(2, (SELECT id FROM skills WHERE name = 'Git'), 4, 3.0),
-- Carol Davis (employee_id: 3)
(3, (SELECT id FROM skills WHERE name = 'AWS'), 5, 5.0),
(3, (SELECT id FROM skills WHERE name = 'Kubernetes'), 5, 4.0),
(3, (SELECT id FROM skills WHERE name = 'Docker'), 5, 4.5),
(3, (SELECT id FROM skills WHERE name = 'CI/CD'), 4, 3.5),
-- David Wilson (employee_id: 4)
(4, (SELECT id FROM skills WHERE name = 'JavaScript'), 4, 1.5),
(4, (SELECT id FROM skills WHERE name = 'React'), 4, 1.5),
(4, (SELECT id FROM skills WHERE name = 'TypeScript'), 3, 1.0),
(4, (SELECT id FROM skills WHERE name = 'Git'), 3, 1.5),
-- Eva Brown (employee_id: 5)
(5, (SELECT id FROM skills WHERE name = 'Python'), 5, 6.0),
(5, (SELECT id FROM skills WHERE name = 'PostgreSQL'), 4, 3.0),
(5, (SELECT id FROM skills WHERE name = 'Git'), 4, 4.0);
