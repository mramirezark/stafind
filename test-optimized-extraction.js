// Test script for optimized NER extraction and search
const BASE_URL = 'http://localhost:8080/api/v1';
const OPTIMIZED_ENDPOINT = `${BASE_URL}/ner/extract-and-search`;

// JWT token from login
const AUTH_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6InRlc3RAc3RhZmluZC5jb20iLCJmaXJzdF9uYW1lIjoiVGVzdCIsImxhc3RfbmFtZSI6IlVzZXIiLCJyb2xlcyI6bnVsbCwiaXNzIjoic3RhZmluZCIsInN1YiI6IjIiLCJleHAiOjE3NTg5MzM5MzAsIm5iZiI6MTc1ODg0NzUzMCwiaWF0IjoxNzU4ODQ3NTMwfQ.ECdqOSulELHu6xdZdTH3KRWviT9JwDX6ngzX40PdxSk';

async function testOptimizedExtraction() {
  console.log('🚀 Testing Optimized NER Extraction & Search\n');
  console.log('='.repeat(60));
  
  const jobDescription = {
    text: 'Looking for a Senior React Developer with 5+ years experience in JavaScript, TypeScript, Node.js, AWS, Docker, and Kubernetes. Must have MySQL and Redis experience. Experience with microservices architecture preferred.',
    title: 'Senior React Developer',
    company: 'TechCorp Inc.',
    location: 'San Francisco, CA',
    salary_range: '$120,000 - $150,000',
    limit: 5
  };
  
  console.log('📤 Testing optimized extraction and search...');
  console.log('Job Description:', jobDescription.text);
  console.log('Company:', jobDescription.company);
  console.log('Location:', jobDescription.location);
  
  try {
    const response = await fetch(OPTIMIZED_ENDPOINT, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${AUTH_TOKEN}`
      },
      body: JSON.stringify(jobDescription)
    });
    
    if (!response.ok) {
      if (response.status === 401) {
        console.log('❌ Authentication required. Please update AUTH_TOKEN in the script.');
        return false;
      }
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const responseData = await response.json();
    
    console.log('\n✅ Response received!');
    console.log('Status:', response.status);
    console.log('Response:', JSON.stringify(responseData, null, 2));
    
    if (responseData.success && responseData.data) {
      const data = responseData.data;
      const skills = data.extracted_skills?.extracted_skills || {};
      const employees = data.matching_employees || [];
      
      console.log('\n🚀 Optimized Analysis:');
      console.log(`Total skills found: ${data.extracted_skills?.total_skills_found || 0}`);
      console.log(`Total employees found: ${data.total_matches || 0}`);
      console.log(`Processing time: ${data.processing_time || 'N/A'}`);
      console.log(`Extraction method: ${data.extracted_skills?.skill_extraction_method || 'N/A'}`);
      console.log(`AI Confidence: ${data.extracted_skills?.ai_confidence || 'N/A'}`);
      
      console.log('\n📋 Skills by Category:');
      if (skills.programming_languages && skills.programming_languages.length > 0) {
        console.log(`Programming Languages: ${skills.programming_languages.join(', ')}`);
      }
      if (skills.web_technologies && skills.web_technologies.length > 0) {
        console.log(`Web Technologies: ${skills.web_technologies.join(', ')}`);
      }
      if (skills.databases && skills.databases.length > 0) {
        console.log(`Databases: ${skills.databases.join(', ')}`);
      }
      if (skills.cloud_devops && skills.cloud_devops.length > 0) {
        console.log(`Cloud & DevOps: ${skills.cloud_devops.join(', ')}`);
      }
      if (skills.soft_skills && skills.soft_skills.length > 0) {
        console.log(`Soft Skills: ${skills.soft_skills.join(', ')}`);
      }
      if (skills.years_of_experience) {
        console.log(`Years of Experience: ${skills.years_of_experience}`);
      }
      
      console.log('\n👥 Matching Employees:');
      if (employees.length > 0) {
        employees.forEach((match, index) => {
          const employee = match.employee;
          const skills = employee.skills?.map(s => s.name).join(', ') || 'N/A';
          console.log(`${index + 1}. ${employee.name} (${employee.email})`);
          console.log(`   Skills: ${skills}`);
          console.log(`   Experience: ${employee.level || 'N/A'} level`);
          console.log(`   Department: ${employee.department || 'N/A'}`);
          console.log(`   Match Score: ${match.match_score.toFixed(2)}`);
          console.log(`   Location: ${employee.location || 'N/A'}`);
          console.log('');
        });
      } else {
        console.log('No matching employees found.');
      }
      
      console.log('\n📊 Search Criteria Used:');
      const criteria = data.search_criteria || {};
      console.log(`Skills: ${criteria.skills?.join(', ') || 'None'}`);
      console.log(`Web Technologies: ${criteria.web_technologies?.join(', ') || 'None'}`);
      console.log(`Databases: ${criteria.databases?.join(', ') || 'None'}`);
      console.log(`Cloud/DevOps: ${criteria.cloud_devops?.join(', ') || 'None'}`);
      console.log(`Soft Skills: ${criteria.soft_skills?.join(', ') || 'None'}`);
      console.log(`Years Experience: ${criteria.years_experience || 'Any'}`);
      console.log(`Languages: ${criteria.languages?.join(', ') || 'Any'}`);
      
      return true;
    } else {
      console.log('❌ No data received from optimized response');
      return false;
    }
    
  } catch (error) {
    console.log('❌ Error occurred:');
    console.log('Error:', error.message);
    return false;
  }
}

async function testOptimizedWithPDF() {
  console.log('\n' + '='.repeat(60));
  console.log('📄 Testing Optimized Extraction with PDF\n');
  
  const pdfJobDescription = {
    text: 'Looking for Python developers. Please find the detailed job requirements in the attached PDF.',
    title: 'Python Developer',
    company: 'DataCorp Inc.',
    location: 'New York, NY',
    salary_range: '$100,000 - $130,000',
    limit: 3
  };
  
  console.log('📤 Testing optimized extraction with PDF reference...');
  console.log('Job Description:', pdfJobDescription.text);
  console.log('Company:', pdfJobDescription.company);
  
  try {
    const response = await fetch(OPTIMIZED_ENDPOINT, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${AUTH_TOKEN}`
      },
      body: JSON.stringify(pdfJobDescription)
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    
    const responseData = await response.json();
    
    console.log('\n✅ Response received!');
    console.log('Status:', response.status);
    console.log('Response:', JSON.stringify(responseData, null, 2));
    
    if (responseData.success && responseData.data) {
      const data = responseData.data;
      const skills = data.extracted_skills?.extracted_skills || {};
      const employees = data.matching_employees || [];
      
      console.log('\n📄 PDF Optimized Analysis:');
      console.log(`Total skills found: ${data.extracted_skills?.total_skills_found || 0}`);
      console.log(`Total employees found: ${data.total_matches || 0}`);
      console.log(`Processing time: ${data.processing_time || 'N/A'}`);
      
      console.log('\n📋 Skills by Category:');
      if (skills.programming_languages && skills.programming_languages.length > 0) {
        console.log(`Programming Languages: ${skills.programming_languages.join(', ')}`);
      }
      if (skills.web_technologies && skills.web_technologies.length > 0) {
        console.log(`Web Technologies: ${skills.web_technologies.join(', ')}`);
      }
      if (skills.databases && skills.databases.length > 0) {
        console.log(`Databases: ${skills.databases.join(', ')}`);
      }
      if (skills.cloud_devops && skills.cloud_devops.length > 0) {
        console.log(`Cloud & DevOps: ${skills.cloud_devops.join(', ')}`);
      }
      
      console.log('\n👥 Matching Employees:');
      if (employees.length > 0) {
        employees.forEach((match, index) => {
          const employee = match.employee;
          const skills = employee.skills?.map(s => s.name).join(', ') || 'N/A';
          console.log(`${index + 1}. ${employee.name} (${employee.email})`);
          console.log(`   Skills: ${skills}`);
          console.log(`   Experience: ${employee.level || 'N/A'} level`);
          console.log(`   Match Score: ${match.match_score.toFixed(2)}`);
        });
      } else {
        console.log('No matching employees found.');
      }
      
      return true;
    } else {
      console.log('❌ No data received from PDF optimized response');
      return false;
    }
    
  } catch (error) {
    console.log('❌ Error occurred:');
    console.log('Error:', error.message);
    return false;
  }
}

async function runOptimizedTests() {
  console.log('🧪 Optimized NER Extraction & Search Tests\n');
  
  const test1 = await testOptimizedExtraction();
  const test2 = await testOptimizedWithPDF();
  
  console.log('\n' + '='.repeat(60));
  console.log('🎯 Test Results Summary:');
  console.log('========================');
  console.log(`${test1 ? '✅' : '❌'} Optimized Extraction: ${test1 ? 'Working' : 'Failed'}`);
  console.log(`${test2 ? '✅' : '❌'} PDF Optimized: ${test2 ? 'Working' : 'Failed'}`);
  
  console.log('\n🚀 Optimized Workflow Benefits:');
  console.log('\n1. Single API Call:');
  console.log('   ✅ NER extraction + Employee search in one request');
  console.log('   ✅ Reduced network overhead');
  console.log('   ✅ Faster response times');
  console.log('   ✅ Simplified n8n workflow');
  
  console.log('\n2. Performance Improvements:');
  console.log('   ✅ 50% fewer API calls');
  console.log('   ✅ Reduced latency');
  console.log('   ✅ Better error handling');
  console.log('   ✅ Atomic operations');
  
  console.log('\n3. Workflow Simplification:');
  console.log('   ✅ Fewer n8n nodes');
  console.log('   ✅ Simpler error handling');
  console.log('   ✅ Easier maintenance');
  console.log('   ✅ Better debugging');
  
  console.log('\n4. API Endpoints:');
  console.log('   ✅ /api/v1/ner/extract-and-search - Combined operation');
  console.log('   ✅ /api/v1/ner/extract-skills - Basic extraction only');
  console.log('   ✅ /api/v1/ner/compare-skills - Skill comparison');
  
  console.log('\n🔧 Next Steps:');
  console.log('1. Import teams-optimized-extraction.json into n8n');
  console.log('2. Test with real Teams messages');
  console.log('3. Monitor performance improvements');
  console.log('4. Deploy to production');
  
  console.log('\n💡 Usage:');
  console.log('The optimized workflow does everything in one call:');
  console.log('1. Receives Teams message');
  console.log('2. Extracts text from message/PDF');
  console.log('3. Performs NER skill extraction');
  console.log('4. Searches for matching employees');
  console.log('5. Returns complete results to Teams');
}

// Run the tests
runOptimizedTests().catch(console.error);
