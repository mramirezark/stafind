-- Insert sample job requests
INSERT INTO job_requests (title, description, department, required_skills, preferred_skills, experience_level, location, priority, created_by) VALUES
('Senior React Developer', 'We are looking for a senior React developer to join our frontend team. Must have experience with modern React patterns and state management.', 'Engineering', 
 '["React", "JavaScript", "TypeScript"]', 
 '["Redux", "Jest", "Webpack"]', 
 'senior', 'San Francisco, CA', 'high', 'HR Team'),

('Go Backend Engineer', 'Join our backend team to build scalable microservices using Go. Experience with PostgreSQL and Docker required.', 'Engineering',
 '["Go", "PostgreSQL", "Docker"]',
 '["Kubernetes", "gRPC", "Redis"]',
 'mid', 'Remote', 'medium', 'Tech Lead'),

('DevOps Engineer', 'Looking for a DevOps engineer to manage our AWS infrastructure and implement CI/CD pipelines.', 'Engineering',
 '["AWS", "Kubernetes", "Docker", "CI/CD"]',
 '["Terraform", "Jenkins", "GitLab CI"]',
 'senior', 'Seattle, WA', 'high', 'DevOps Manager'),

('Junior Frontend Developer', 'Great opportunity for a junior developer to work with our experienced team on modern web applications.', 'Engineering',
 '["JavaScript", "React"]',
 '["TypeScript", "CSS", "HTML"]',
 'junior', 'Austin, TX', 'low', 'Frontend Lead'),

('Data Scientist', 'Join our data team to build machine learning models and analyze large datasets.', 'Data Science',
 '["Python", "PostgreSQL"]',
 '["TensorFlow", "Pandas", "NumPy"]',
 'senior', 'Boston, MA', 'medium', 'Data Team Lead');
