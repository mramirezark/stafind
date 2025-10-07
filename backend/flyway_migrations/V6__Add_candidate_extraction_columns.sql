-- Add candidate extraction related columns to employees table
ALTER TABLE employees 
ADD COLUMN original_text TEXT,
ADD COLUMN extracted_data JSONB,
ADD COLUMN extraction_timestamp TIMESTAMP,
ADD COLUMN extraction_source VARCHAR(100),
ADD COLUMN extraction_status VARCHAR(50) DEFAULT 'pending',
ADD COLUMN resume_url VARCHAR(500);

-- Add indexes for better performance
CREATE INDEX idx_employees_extraction_status ON employees(extraction_status);
CREATE INDEX idx_employees_extraction_timestamp ON employees(extraction_timestamp);
CREATE INDEX idx_employees_extraction_source ON employees(extraction_source);
CREATE INDEX idx_employees_resume_url ON employees(resume_url);
