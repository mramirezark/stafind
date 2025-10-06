// Test Spanish candidate extraction
const testSpanishResume = `
Servicios rest - Desarrollo de APIS / Web services
Arquitectura de base de datos
Business Intelligence y modelado de datos
Gestión y mejora de procesos
Á R E A S D E E S P E C I A L I Z A C I Ó N
GUSTAVO
GARRIDO
Especialista en procesos vinculado al desarrollo de aplicaciones
e inteligencia de negocios. Soy un apasionado por la tecnología
y la mejora. Mi mayor fortaleza es la versatilidad y capacidad
de entender los sistemas desde la visión del negocio. Cada
proyecto lo veo como un desafío para fortalecer mis habilidades
como desarrollador y gestor.
COORDINADOR DE PROCESOS Y LIDER DE
IMPLEMENTACIÓN
Análisis y relevamiento técnico para poder llevar a cabo la
transformación digital de procesos.
Lider de proyectos de implementación de software y mejora de
procesos.
Implementación de soluciones de ETL (Extracción, Transformación y
Carga de Datos)
Desarrollo de software para procesos complementarios
Soporte técnico a equipo BI
Desarrollo y e implementación de CRM para gestión de reclamos. -
backend desarrollador en node JS frontend desarrollado en Angular.
Implementación de software en todos los procesos ambulatorios.
Armado y publicación de SUITE de reportes en Qliksense con acceso a
los datos e información de todos los sistemas implementados en la
organización.
Principales logros:
H O S P I T A L S A N J U A N D E D I O S | J U N I O 2 0 1 6 - P R E S E N T E
T R A Y E C T O R I A P R O F E S I O N A L
Wenceslao de tata
4672, caseros,
Provincia de buenos
aires
ingindustrial.gustavo@gmail.c
om
https://www.linkedin.co
m/in/gustavo-garrido-
96694393/:
@unsitiogenial
+54 11-2403-3763
F U L L S T A C K
D E V E L O P E R
COORDINADOR DE PROYECTOS Y CONSULTOR SR
Lider de equipo para la implementación de proyectos de mejora de
procesos
Certificación norma API Spec Q1 en empresa de fabricación de valvulas
bridadas.
Implementación de proceso justing time en empresa de fabricación de
jugos.
Principales logros:
G E S T I Ó N 3 6 0 | J U L I O 2 0 1 4 - J U N I O 2 0 1 6
D a t o s d e
c o n t a c t o

UNIVERSIDAD TECNOLÓGICA NACIONAL -
ARGENTINA
E S P E C I A L I S T A E N I N G E N I E R Í A E N C A L I D A D , P R O M O C I Ó N
D E 2 0 1 8
F O R M A C I Ó N A C A D É M I C A
C O M P E T E N C I A S T E C N I C A S
UNIVERSIDAD DEL MAGDALENA - COLOMBIA
I N G E N I E R O I N D U S T R I A L , P R O M O C I Ó N D E 2 0 1 2
UNIVERSIDAD TECNOLÓGICA NACIONAL -
ARGENTINA
I O N I C 6 - C R E A R A P L I C A C I O N E S I O S , A N D R O I D Y P W A C O N
A N G U L A R - 2 0 2 0
F U L L S T A C K D E V E L O P E R 2 0 1 9
TUV RHEINLAND
A U D I T O R L I D E R E N S I S T E M A S D E G E S T I Ó N D E L A C A L I D A D
I S O 9 0 0 1
UNIVERSIDAD DEL MAGDALENA - COLOMBIA
D I P L O M A T U R A G E S T I Ó N D E P R O Y E C T O S , P R O M O C I Ó N D E
2 0 1 1
C U R S O S Y O T R A S
F O R M A C I O N E S
Python Avanzado
R avanzado
React Avanzado
AWS avanzado
Testing unitario
Ingles Oral y
escrito Avanzado
con examen toefl
aprobado
Plan de carrera
para los proximos 2
años
Idiomas:
Ingles:
lectura: Avanzado
Escritura: Intermedio
Oral: Intermedio
Node JS
Angular
Git
SLQ en motores PL/SQL, MySQL, SQL SERVER
Typescript
Javascript
HTML
CSS
Procesos ETL
QLiksence
UNIVERSIDAD TECNOLÓGICA NACIONAL -
ARGENTINA
P R O G R A M A D O R W E B F R O N T E N D 2 0 1 9
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
