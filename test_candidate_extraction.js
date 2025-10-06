// Test the consolidated extract endpoint with candidate storage
async function testCandidateExtraction() {
    const baseURL = 'http://localhost:8080';
    
    // Test data for candidate extraction
    const testData = {
        message_text: "I need a senior Python developer with React experience for a new project",
        text: `Jesus Manuel Ramirez Mendez
Senior Software Engineer with over 15 years of experience in developing scalable web applications and
backend services working with clients from both the US and Mexico. Strong problem-solving skills and
innovative developer with a passion for building scalable, efficient, user-centric software solutions using
modern technology stacks. Collaborative team player with experience in Agile methodologies.
SKILLS
● Programming languages
○ 10 years: Java, JavaScript, SQL
○ 5 years: C#, VB, TypeScript
● Frameworks
○ 7 years: Spring Boot
○ 
5 years: Angular, jQuery, ExtJs, ASP.NET Core
○ 
1 year: React, Node.js
● Database
○ 7 years: PostgreSQL, SQL Server, MySQL
○ 5 years: Oracle
○ 1 year: Apache Cassandra, Kinetica
● Development tools
○ 10 years: Eclipse, IntelliJ
○ 5 years: Visual Studio, VS Code
● Version Control
○ 7 years: Git, Github, Bitbucket, SVN
● Others
○ 10 years: Layered Architecture, RESTful APIs, HTML, CSS, Maven
○ 7 years: Postman, Microservices, Docker, Scrum, Kanban, JIRA
EXPERIENCE
ArkusNexus
Senior Software Engineer | Nov 2018 - Present
● 
Developed and optimized a full-stack application using React and TypeScript, increasing efficiency
through code refactoring and performance enhancements using some AI tools.
● 
Developed and maintained robust web applications, microservices and RESTful APIs with Spring Boot.
● 
Optimized backend processes with Spring Boot, resulting in a several reduction in server response
times.
● 
Collaborated with cross-functional teams in Agile (Scrum/Kanban) sprints to plan and analyze the
implementation of new features and functionalities within deadlines.
● 
Identified and resolved complex software issues and enhanced overall reliability.
● 
Performed unit, integration, performance, and E2E testing reducing post-deployment defects.
● 
Analyze, prototyping and implementing new features on schedule.
Prianti Consulting
Senior Software Engineer | Jan 2023 - Nov 2018
● 
Developed, maintained, and deployed high-performance web applications using Java (Spring MVC,
JEE), ExtJS, and JavaScript, improving user experience through responsive UI optimizations.
● 
Partnered with product to gather requirements, prototype solutions, and launch new features.
● 
Designed and optimized RESTful APIs with Spring MVC enabling seamless integration across multiple
internal applications.
● 
Conducted extensive testing and quality assurance.
● 
Analyze, prototyping and implementing new modules for internal applications.

Hildebrando
Senior Software Engineer | Apr 2012 - Dec 2012
● 
Developed and maintained internal web applications using .NET Technologies.
● 
Performed user acceptance testing (UAT) with cross-functional teams to validate business
requirements.
Interfactura
Senior Software Engineer | Dic 2008 - Mar 2012
● 
Designed, developed, and deployed full-stack web applications using .NET Technologies and
JavaScript.
● 
Partnered with product managers, business analysts to gather requirements, analyze feasibility, and
implement new features, improving user efficiency.
● 
Engineered high-performance Windows Services to automate massive invoice generation cutting
processing time.
● 
Conducted unit, integration, and E2E testing reducing post-deployment defects.
EDUCATION
● 
Instituto Tecnologico de Zacatecas / Ingeniero en Sistemas Computacionales 2001-2005`,
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
        text: `Jesus Manuel Ramirez Mendez
Senior Software Engineer with over 15 years of experience in developing scalable web applications and
backend services working with clients from both the US and Mexico. Strong problem-solving skills and
innovative developer with a passion for building scalable, efficient, user-centric software solutions using
modern technology stacks. Collaborative team player with experience in Agile methodologies.
SKILLS
● Programming languages
○ 10 years: Java, JavaScript, SQL
○ 5 years: C#, VB, TypeScript
● Frameworks
○ 7 years: Spring Boot
○ 
5 years: Angular, jQuery, ExtJs, ASP.NET Core
○ 
1 year: React, Node.js
● Database
○ 7 years: PostgreSQL, SQL Server, MySQL
○ 5 years: Oracle
○ 1 year: Apache Cassandra, Kinetica
● Development tools
○ 10 years: Eclipse, IntelliJ
○ 5 years: Visual Studio, VS Code
● Version Control
○ 7 years: Git, Github, Bitbucket, SVN
● Others
○ 10 years: Layered Architecture, RESTful APIs, HTML, CSS, Maven
○ 7 years: Postman, Microservices, Docker, Scrum, Kanban, JIRA
EXPERIENCE
ArkusNexus
Senior Software Engineer | Nov 2018 - Present
● 
Developed and optimized a full-stack application using React and TypeScript, increasing efficiency
through code refactoring and performance enhancements using some AI tools.
● 
Developed and maintained robust web applications, microservices and RESTful APIs with Spring Boot.
● 
Optimized backend processes with Spring Boot, resulting in a several reduction in server response
times.
● 
Collaborated with cross-functional teams in Agile (Scrum/Kanban) sprints to plan and analyze the
implementation of new features and functionalities within deadlines.
● 
Identified and resolved complex software issues and enhanced overall reliability.
● 
Performed unit, integration, performance, and E2E testing reducing post-deployment defects.
● 
Analyze, prototyping and implementing new features on schedule.
Prianti Consulting
Senior Software Engineer | Jan 2023 - Nov 2018
● 
Developed, maintained, and deployed high-performance web applications using Java (Spring MVC,
JEE), ExtJS, and JavaScript, improving user experience through responsive UI optimizations.
● 
Partnered with product to gather requirements, prototype solutions, and launch new features.
● 
Designed and optimized RESTful APIs with Spring MVC enabling seamless integration across multiple
internal applications.
● 
Conducted extensive testing and quality assurance.
● 
Analyze, prototyping and implementing new modules for internal applications.

Hildebrando
Senior Software Engineer | Apr 2012 - Dec 2012
● 
Developed and maintained internal web applications using .NET Technologies.
● 
Performed user acceptance testing (UAT) with cross-functional teams to validate business
requirements.
Interfactura
Senior Software Engineer | Dic 2008 - Mar 2012
● 
Designed, developed, and deployed full-stack web applications using .NET Technologies and
JavaScript.
● 
Partnered with product managers, business analysts to gather requirements, analyze feasibility, and
implement new features, improving user efficiency.
● 
Engineered high-performance Windows Services to automate massive invoice generation cutting
processing time.
● 
Conducted unit, integration, and E2E testing reducing post-deployment defects.
EDUCATION
● 
Instituto Tecnologico de Zacatecas / Ingeniero en Sistemas Computacionales 2001-2005`,
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
