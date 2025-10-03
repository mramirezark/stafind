// Test Spanish candidate extraction
const testSpanishResume = `
María González
Desarrolladora Senior de Software
maria.gonzalez@email.com
+34 612 345 678
Madrid, España

EXPERIENCIA PROFESIONAL
Desarrolladora Senior - TechCorp (2020-2024)
• 4 años de experiencia desarrollando aplicaciones web con React y Node.js
• Liderazgo de equipos de desarrollo ágiles
• Experiencia con AWS, Docker y Kubernetes
• Comunicación efectiva con stakeholders internacionales

Desarrolladora Full Stack - StartupXYZ (2018-2020)
• 2 años desarrollando aplicaciones con Python, Django y PostgreSQL
• Trabajo en equipo en metodologías Scrum
• Resolución de problemas complejos de rendimiento

HABILIDADES TÉCNICAS
Lenguajes de Programación: JavaScript, TypeScript, Python, Java, Go
Tecnologías Web: React, Angular, Vue.js, Node.js, Express, Django
Bases de Datos: PostgreSQL, MongoDB, Redis, Elasticsearch
Cloud y DevOps: AWS, Azure, Docker, Kubernetes, Jenkins, Terraform
Herramientas: Git, GitHub, GitLab, Jira, Confluence

HABILIDADES BLANDAS
• Liderazgo y gestión de equipos
• Comunicación efectiva
• Trabajo en equipo
• Resolución de problemas
• Pensamiento analítico
• Creatividad e innovación
• Adaptabilidad
• Gestión del tiempo
• Mentoring de desarrolladores junior

EDUCACIÓN
Ingeniería en Sistemas - Universidad Politécnica de Madrid (2014-2018)
Certificación AWS Solutions Architect (2021)
`;

async function testSpanishExtraction() {
    try {
        console.log('🧪 Testing Spanish candidate extraction...\n');
        
        const response = await fetch('http://localhost:8080/api/v1/extract/process', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': 'dev-api-key-12345'
            },
            body: JSON.stringify({
                message_text: "Necesitamos un desarrollador senior con experiencia en React y Python para un nuevo proyecto",
                text: testSpanishResume,
                file_name: "maria_gonzalez_resume.pdf",
                file_url: "https://example.com/resumes/maria_gonzalez.pdf",
                processing_type: "candidate_extraction",
                extraction_source: "test",
                total_files: 1,
                file_number: 1,
                teams_message_id: "spanish-test-123",
                channel_id: "test-channel",
                user_id: "test-user",
                user_name: "Test User",
                metadata: {
                    language: "spanish",
                    format: "resume",
                    source: "test"
                }
            })
        });

        const result = await response.json();
        
        if (response.ok) {
            console.log('✅ Spanish extraction successful!');
            console.log('📊 Response Status:', response.status);
            console.log('📝 Message:', result.message);
            
            if (result.extraction_result) {
                const extraction = JSON.parse(result.extraction_result.processed_content);
                console.log('\n🔍 Extracted Information:');
                console.log('👤 Name:', extraction.candidate_name);
                console.log('📧 Email:', extraction.contact_info.email);
                console.log('📱 Phone:', extraction.contact_info.phone);
                console.log('📍 Location:', extraction.contact_info.location);
                console.log('💼 Current Position:', extraction.current_position);
                console.log('🎯 Seniority Level:', extraction.seniority_level);
                console.log('⏱️ Years Experience:', extraction.years_experience);
                console.log('🌍 Languages Detected:', extraction.languages);
                console.log('💻 Programming Languages:', extraction.skills.programming_languages);
                console.log('🌐 Web Technologies:', extraction.skills.web_technologies);
                console.log('🗄️ Databases:', extraction.skills.databases);
                console.log('☁️ Cloud & DevOps:', extraction.skills.cloud_devops);
                console.log('🤝 Soft Skills:', extraction.skills.soft_skills);
                console.log('🎓 Education Level:', extraction.education.level);
                console.log('📚 Institutions:', extraction.education.institutions);
                console.log('📈 Total Skills Found:', extraction.total_skills);
                console.log('🎯 Confidence Score:', extraction.confidence);
            }
            
            if (result.candidate_result) {
                console.log('\n👤 Candidate Processing:');
                console.log('🆔 Employee ID:', result.candidate_result.employee_id);
                console.log('⚡ Action:', result.candidate_result.action);
                console.log('🔄 Changes Detected:', result.candidate_result.changes_detected);
                console.log('📝 Status:', result.candidate_result.status);
                console.log('💬 Message:', result.candidate_result.message);
            }
            
            if (result.matching_result) {
                console.log('\n🎯 Matching Results:');
                console.log('📋 Requirements:', result.matching_result.requirements);
                console.log('🔧 Required Skills:', result.matching_result.required_skills);
                console.log('👥 Total Candidates:', result.matching_result.total_candidates);
                console.log('🎯 Matches Found:', result.matching_result.matches.length);
                if (result.matching_result.matches.length > 0) {
                    console.log('🥇 Top Match:', result.matching_result.matches[0].employee_name, 
                              `(Score: ${result.matching_result.matches[0].match_score}%)`);
                }
            }
            
        } else {
            console.log('❌ Spanish extraction failed!');
            console.log('📊 Status:', response.status);
            console.log('📝 Error:', result.error || result.message);
        }
        
    } catch (error) {
        console.error('💥 Test failed:', error.message);
    }
}

// Run the test
testSpanishExtraction();
