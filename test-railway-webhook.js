// Test script for Railway webhook endpoint
const RAILWAY_WEBHOOK_URL = 'https://primary-production-674a.up.railway.app/webhook-test/teams-optimized-extraction';

async function testRailwayWebhook() {
  console.log('🚀 Testing Railway Webhook Endpoint\n');
  console.log('='.repeat(60));
  console.log('Webhook URL:', RAILWAY_WEBHOOK_URL);
  
  // Test 1: Job description with skills
  const jobDescriptionMessage = {
    message_id: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    channel_id: 'teams_channel_123',
    user_id: 'teams_user_456',
    user_name: 'John Doe',
    message_text: 'Looking for a Senior React Developer with 5+ years experience in JavaScript, TypeScript, Node.js, AWS, Docker, and Kubernetes. Must have MySQL and Redis experience. Experience with microservices architecture preferred.',
    message_type: 'message',
    timestamp: new Date().toISOString(),
    attachments: []
  };
  
  console.log('\n📤 Test 1: Job Description Message');
  console.log('Message ID:', jobDescriptionMessage.message_id);
  console.log('Channel ID:', jobDescriptionMessage.channel_id);
  console.log('User:', jobDescriptionMessage.user_name);
  console.log('Message:', jobDescriptionMessage.message_text.substring(0, 100) + '...');
  
  try {
    const response = await fetch(RAILWAY_WEBHOOK_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(jobDescriptionMessage)
    });
    
    console.log('\n✅ Railway Webhook Response:');
    console.log('Status:', response.status);
    console.log('Status Text:', response.statusText);
    
    const responseText = await response.text();
    console.log('Response Body:', responseText);
    
    if (response.ok) {
      try {
        const responseData = JSON.parse(responseText);
        console.log('Parsed Response:', JSON.stringify(responseData, null, 2));
        
        if (responseData.success && responseData.data) {
          const data = responseData.data;
          console.log('\n🚀 Railway Workflow Results:');
          console.log(`Total skills found: ${data.extracted_skills?.total_skills_found || 0}`);
          console.log(`Total employees found: ${data.total_matches || 0}`);
          console.log(`Processing time: ${data.processing_time || 'N/A'}`);
          console.log(`Extraction method: ${data.extracted_skills?.skill_extraction_method || 'N/A'}`);
          console.log(`AI Confidence: ${data.extracted_skills?.ai_confidence || 'N/A'}`);
          
          if (data.matching_employees && data.matching_employees.length > 0) {
            console.log('\n👥 Top Matching Employees:');
            data.matching_employees.slice(0, 3).forEach((match, index) => {
              const employee = match.employee;
              const skills = employee.skills?.map(s => s.name).join(', ') || 'N/A';
              console.log(`${index + 1}. ${employee.name} (${employee.email})`);
              console.log(`   Match Score: ${match.match_score.toFixed(2)}`);
              console.log(`   Skills: ${skills}`);
              console.log(`   Level: ${employee.level || 'N/A'}`);
              console.log(`   Location: ${employee.location || 'N/A'}`);
              console.log('');
            });
          }
          
          return true;
        } else {
          console.log('❌ No data received from Railway response');
          return false;
        }
      } catch (parseError) {
        console.log('❌ Error parsing JSON response:');
        console.log('Parse Error:', parseError.message);
        console.log('Raw Response:', responseText);
        return false;
      }
    } else {
      console.log('❌ Railway webhook returned error status');
      return false;
    }
    
  } catch (error) {
    console.log('❌ Error occurred:');
    console.log('Error:', error.message);
    return false;
  }
}

async function testRailwayWebhookWithPDF() {
  console.log('\n' + '='.repeat(60));
  console.log('📄 Test 2: PDF Attachment Message\n');
  
  const pdfMessage = {
    message_id: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    channel_id: 'teams_channel_456',
    user_id: 'teams_user_789',
    user_name: 'Jane Smith',
    message_text: 'Looking for Python developers. Please find the detailed job requirements in the attached PDF.',
    message_type: 'message',
    timestamp: new Date().toISOString(),
    attachments: [
      {
        name: 'job_description.pdf',
        contentType: 'application/pdf',
        contentUrl: 'https://example.com/job_description.pdf',
        size: 1024000
      }
    ]
  };
  
  console.log('📤 Test 2: PDF Attachment Message');
  console.log('Message ID:', pdfMessage.message_id);
  console.log('User:', pdfMessage.user_name);
  console.log('Message:', pdfMessage.message_text);
  console.log('Attachments:', pdfMessage.attachments.length);
  
  try {
    const response = await fetch(RAILWAY_WEBHOOK_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(pdfMessage)
    });
    
    console.log('\n✅ PDF Railway Webhook Response:');
    console.log('Status:', response.status);
    console.log('Status Text:', response.statusText);
    
    const responseText = await response.text();
    console.log('Response Body:', responseText);
    
    if (response.ok) {
      try {
        const responseData = JSON.parse(responseText);
        console.log('Parsed Response:', JSON.stringify(responseData, null, 2));
        
        if (responseData.success && responseData.data) {
          const data = responseData.data;
          console.log('\n📄 PDF Railway Workflow Results:');
          console.log(`Total skills found: ${data.extracted_skills?.total_skills_found || 0}`);
          console.log(`Total employees found: ${data.total_matches || 0}`);
          console.log(`Processing time: ${data.processing_time || 'N/A'}`);
          
          if (data.matching_employees && data.matching_employees.length > 0) {
            console.log('\n👥 Matching Employees:');
            data.matching_employees.forEach((match, index) => {
              const employee = match.employee;
              const skills = employee.skills?.map(s => s.name).join(', ') || 'N/A';
              console.log(`${index + 1}. ${employee.name} (${employee.email})`);
              console.log(`   Match Score: ${match.match_score.toFixed(2)}`);
              console.log(`   Skills: ${skills}`);
            });
          }
          
          return true;
        } else {
          console.log('❌ No data received from PDF Railway response');
          return false;
        }
      } catch (parseError) {
        console.log('❌ Error parsing JSON response:');
        console.log('Parse Error:', parseError.message);
        console.log('Raw Response:', responseText);
        return false;
      }
    } else {
      console.log('❌ Railway PDF webhook returned error status');
      return false;
    }
    
  } catch (error) {
    console.log('❌ Error occurred:');
    console.log('Error:', error.message);
    return false;
  }
}

async function testRailwayWebhookInvalidMessage() {
  console.log('\n' + '='.repeat(60));
  console.log('❌ Test 3: Invalid Message (Should be ignored)\n');
  
  const invalidMessage = {
    message_id: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    channel_id: 'teams_channel_789',
    user_id: 'teams_user_999',
    user_name: 'Invalid User',
    message_text: 'Hello everyone! How is your day going?',
    message_type: 'message',
    timestamp: new Date().toISOString(),
    attachments: []
  };
  
  console.log('📤 Test 3: Invalid Message (Regular chat)');
  console.log('Message ID:', invalidMessage.message_id);
  console.log('User:', invalidMessage.user_name);
  console.log('Message:', invalidMessage.message_text);
  
  try {
    const response = await fetch(RAILWAY_WEBHOOK_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(invalidMessage)
    });
    
    console.log('\n✅ Invalid Message Railway Response:');
    console.log('Status:', response.status);
    console.log('Status Text:', response.statusText);
    
    const responseText = await response.text();
    console.log('Response Body:', responseText);
    
    if (response.ok) {
      try {
        const responseData = JSON.parse(responseText);
        console.log('Parsed Response:', JSON.stringify(responseData, null, 2));
        
        if (responseData.status === 'ignored' || responseData.status === 'filtered') {
          console.log('✅ Correctly ignored non-job message');
          return true;
        } else if (responseData.success && responseData.data && responseData.data.total_matches === 0) {
          console.log('✅ Correctly processed but found no matches (as expected for non-job message)');
          return true;
        } else {
          console.log('❌ Should have been ignored or found no matches');
          return false;
        }
      } catch (parseError) {
        console.log('❌ Error parsing JSON response:');
        console.log('Parse Error:', parseError.message);
        console.log('Raw Response:', responseText);
        return false;
      }
    } else {
      console.log('❌ Railway invalid message webhook returned error status');
      return false;
    }
    
  } catch (error) {
    console.log('❌ Error occurred:');
    console.log('Error:', error.message);
    return false;
  }
}

async function runRailwayWebhookTests() {
  console.log('🧪 Railway Webhook Tests\n');
  
  const test1 = await testRailwayWebhook();
  const test2 = await testRailwayWebhookWithPDF();
  const test3 = await testRailwayWebhookInvalidMessage();
  
  console.log('\n' + '='.repeat(60));
  console.log('🎯 Railway Webhook Test Results:');
  console.log('================================');
  console.log(`${test1 ? '✅' : '❌'} Job Description: ${test1 ? 'Working' : 'Failed'}`);
  console.log(`${test2 ? '✅' : '❌'} PDF Attachment: ${test2 ? 'Working' : 'Failed'}`);
  console.log(`${test3 ? '✅' : '❌'} Invalid Message: ${test3 ? 'Working' : 'Failed'}`);
  
  console.log('\n🚀 Railway Webhook Configuration:');
  console.log('URL:', RAILWAY_WEBHOOK_URL);
  console.log('Method: POST');
  console.log('Content-Type: application/json');
  
  console.log('\n💡 Expected Workflow Behavior:');
  console.log('1. Receives Teams message via webhook');
  console.log('2. Filters for job description messages');
  console.log('3. Extracts text from message or PDF');
  console.log('4. Calls optimized NER endpoint');
  console.log('5. Returns matching employees to Teams');
  
  console.log('\n🔧 n8n Workflow Setup:');
  console.log('1. Import teams-optimized-extraction.json into n8n');
  console.log('2. Set webhook URL to:', RAILWAY_WEBHOOK_URL);
  console.log('3. Configure Teams connector to send messages to webhook');
  console.log('4. Test with real Teams messages');
  
  if (test1 || test2 || test3) {
    console.log('\n🎉 Railway Webhook is Working!');
    console.log('The webhook endpoint is responding and processing messages.');
  } else {
    console.log('\n❌ Railway Webhook Issues Detected');
    console.log('Please check the webhook configuration and n8n workflow setup.');
  }
}

// Run the Railway webhook tests
runRailwayWebhookTests().catch(console.error);
