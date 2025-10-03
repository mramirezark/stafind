// Test the consolidated extract endpoint with candidate storage
async function testCandidateExtraction() {
    const baseURL = 'http://localhost:8080';
    
    // Test data for candidate extraction
    const testData = {
        message_text: "I need a senior Python developer with React experience for a new project",
        text: `John Smith
Senior Python Developer
Email: john.smith@email.com
Phone: +1-555-0123
Location: San Francisco, CA

EXPERIENCE:
- 5 years of experience in Python development
- 3 years with React and JavaScript
- Experience with Django, Flask, FastAPI
- Database: PostgreSQL, MongoDB
- Cloud: AWS, Docker, Kubernetes
- Version Control: Git, GitHub

EDUCATION:
- Bachelor's in Computer Science from Stanford University (2018)

SKILLS:
- Programming Languages: Python, JavaScript, TypeScript, SQL
- Frameworks: Django, Flask, React, Node.js
- Databases: PostgreSQL, MongoDB, Redis
- Cloud & DevOps: AWS, Docker, Kubernetes, CI/CD
- Tools: Git, GitHub, Jira, Slack

CURRENT PROJECT:
- Leading development of microservices architecture for e-commerce platform
- Team size: 8 developers
- Technologies: Python, React, PostgreSQL, AWS

SUMMARY:
Experienced full-stack developer with strong background in Python and modern web technologies. Proven track record of delivering scalable solutions and leading development teams.`,
        file_name: "john_smith_resume.pdf",
        file_url: "https://example.com/resumes/john_smith.pdf",
        processing_type: "candidate_extraction",
        extraction_source: "resume",
        total_files: 1,
        file_number: 1,
        metadata: {
            language: "en",
            format: "resume"
        }
    };

    try {
        console.log('🚀 Testing candidate extraction and storage...');
        console.log('Request data:', JSON.stringify(testData, null, 2));
        
        const response = await fetch(`${baseURL}/api/v1/extract/process`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': 'dev-api-key-12345' // You'll need to replace this with a valid API key
            },
            body: JSON.stringify(testData)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('\n✅ Success! Response:');
        console.log(JSON.stringify(data, null, 2));
        
        // Check if employee was created/updated
        if (data.candidate_result) {
            console.log('\n📊 Candidate Processing Summary:');
            console.log(`- Employee ID: ${data.candidate_result.employee_id}`);
            console.log(`- Action: ${data.candidate_result.action}`);
            console.log(`- Changes Detected: ${data.candidate_result.changes_detected}`);
            console.log(`- Status: ${data.candidate_result.status}`);
            console.log(`- Message: ${data.candidate_result.message}`);
            
            if (data.candidate_result.changes_summary) {
                console.log('- Changes Summary:');
                data.candidate_result.changes_summary.forEach(change => {
                    console.log(`  • ${change}`);
                });
            }
        }
        
    } catch (error) {
        console.error('\n❌ Error testing candidate extraction:');
        console.error('Error:', error.message);
    }
}

// Test updating an existing employee
async function testUpdateExistingEmployee() {
    const baseURL = 'http://localhost:8080';
    
    // Test data with updated information
    const testData = {
        message_text: "Update search for JavaScript developers",
        text: `John Smith
Senior Full-Stack Developer
Email: john.smith@email.com
Phone: +1-555-0123
Location: San Francisco, CA

EXPERIENCE:
- 6 years of experience in Python and JavaScript development
- 4 years with React and TypeScript
- Experience with Django, Flask, FastAPI, Next.js
- Database: PostgreSQL, MongoDB, Redis
- Cloud: AWS, Azure, Docker, Kubernetes
- Version Control: Git, GitHub, GitLab

EDUCATION:
- Bachelor's in Computer Science from Stanford University (2018)
- AWS Certified Solutions Architect (2023)

SKILLS:
- Programming Languages: Python, JavaScript, TypeScript, SQL, Go
- Frameworks: Django, Flask, React, Next.js, Node.js, Express
- Databases: PostgreSQL, MongoDB, Redis, Elasticsearch
- Cloud & DevOps: AWS, Azure, Docker, Kubernetes, CI/CD, Terraform
- Tools: Git, GitHub, GitLab, Jira, Slack, Figma

CURRENT PROJECT:
- Leading development of microservices architecture for e-commerce platform
- Team size: 12 developers
- Technologies: Python, React, Next.js, PostgreSQL, AWS, Docker

SUMMARY:
Experienced full-stack developer with strong background in Python and modern web technologies. Recently expanded expertise to include Go and additional cloud platforms. Proven track record of delivering scalable solutions and leading development teams.`,
        file_name: "john_smith_resume_updated.pdf",
        file_url: "https://example.com/resumes/john_smith_updated.pdf",
        processing_type: "candidate_extraction",
        extraction_source: "resume",
        total_files: 1,
        file_number: 1,
        metadata: {
            language: "en",
            format: "resume"
        }
    };

    try {
        console.log('\n\n🔄 Testing update of existing employee...');
        console.log('Request data:', JSON.stringify(testData, null, 2));
        
        const response = await fetch(`${baseURL}/api/v1/extract/process`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': 'dev-api-key-12345' // You'll need to replace this with a valid API key
            },
            body: JSON.stringify(testData)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log('\n✅ Success! Response:');
        console.log(JSON.stringify(data, null, 2));
        
        // Check if employee was updated
        if (data.candidate_result) {
            console.log('\n📊 Update Processing Summary:');
            console.log(`- Employee ID: ${data.candidate_result.employee_id}`);
            console.log(`- Action: ${data.candidate_result.action}`);
            console.log(`- Changes Detected: ${data.candidate_result.changes_detected}`);
            console.log(`- Status: ${data.candidate_result.status}`);
            console.log(`- Message: ${data.candidate_result.message}`);
            
            if (data.candidate_result.changes_summary) {
                console.log('- Changes Summary:');
                data.candidate_result.changes_summary.forEach(change => {
                    console.log(`  • ${change}`);
                });
            }
        }
        
    } catch (error) {
        console.error('\n❌ Error testing employee update:');
        console.error('Error:', error.message);
    }
}

// Run tests
async function runTests() {
    console.log('🚀 Starting candidate extraction and storage tests...\n');
    
    await testCandidateExtraction();
    await testUpdateExistingEmployee();
    
    console.log('\n✨ Tests completed!');
}

runTests().catch(console.error);
