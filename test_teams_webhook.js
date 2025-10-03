/**
 * Test script for Teams Message Receiver n8n workflow
 * This script simulates Teams messages being sent to the webhook
 */

const fetch = require('node-fetch');

// Configuration
const WEBHOOK_URL = 'http://localhost:5678/webhook/teams-message-receiver'; // Update with your n8n instance URL
const API_KEY = 'dev-api-key-12345';

// Test messages
const testMessages = [
  {
    name: "Resume Extraction Request",
    data: {
      messageId: `msg_${Date.now()}_1`,
      channelId: "channel_engineering",
      userId: "user_john_doe",
      userName: "John Doe",
      messageText: "Please extract this resume: Maria Gonzalez is a senior developer with 5 years of experience in JavaScript, React, and Node.js. She has worked at TechCorp and has a degree in Computer Science.",
      timestamp: new Date().toISOString(),
      channelName: "Engineering",
      teamId: "team_engineering",
      teamName: "Engineering Team",
      messageType: "text",
      attachments: [],
      mentions: []
    }
  },
  {
    name: "Candidate Search Request",
    data: {
      messageId: `msg_${Date.now()}_2`,
      channelId: "channel_hr",
      userId: "user_jane_smith",
      userName: "Jane Smith",
      messageText: "We need to find candidates for our Python developer position. Looking for someone with Django and PostgreSQL experience.",
      timestamp: new Date().toISOString(),
      channelName: "Human Resources",
      teamId: "team_hr",
      teamName: "HR Team",
      messageType: "text",
      attachments: [],
      mentions: []
    }
  },
  {
    name: "CV Processing Request",
    data: {
      messageId: `msg_${Date.now()}_3`,
      channelId: "channel_recruiting",
      userId: "user_mike_wilson",
      userName: "Mike Wilson",
      messageText: "Can you extract this CV? Carlos Rodriguez - Senior Software Engineer with 8 years of experience in Java, Spring Boot, and microservices architecture.",
      timestamp: new Date().toISOString(),
      channelName: "Recruiting",
      teamId: "team_recruiting",
      teamName: "Recruiting Team",
      messageType: "text",
      attachments: [],
      mentions: []
    }
  },
  {
    name: "Non-Extraction Message",
    data: {
      messageId: `msg_${Date.now()}_4`,
      channelId: "channel_general",
      userId: "user_sarah_jones",
      userName: "Sarah Jones",
      messageText: "Good morning everyone! How's the project going?",
      timestamp: new Date().toISOString(),
      channelName: "General",
      teamId: "team_general",
      teamName: "General Team",
      messageType: "text",
      attachments: [],
      mentions: []
    }
  }
];

async function testTeamsWebhook() {
  console.log('🧪 Testing Teams Message Receiver Webhook...\n');
  
  for (const testCase of testMessages) {
    try {
      console.log(`📤 Testing: ${testCase.name}`);
      console.log(`📝 Message: ${testCase.data.messageText.substring(0, 100)}...`);
      
      const response = await fetch(WEBHOOK_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': API_KEY
        },
        body: JSON.stringify(testCase.data)
      });
      
      const result = await response.json();
      
      if (response.ok) {
        console.log('✅ Success!');
        console.log(`📊 Status: ${result.status}`);
        console.log(`💬 Message: ${result.message}`);
        
        if (result.summary) {
          console.log(`👤 Employee ID: ${result.summary.employee_id || 'N/A'}`);
          console.log(`⚡ Action: ${result.summary.action || 'N/A'}`);
          console.log(`🔧 Skills Found: ${result.summary.total_skills || 0}`);
          console.log(`🎯 Matches Found: ${result.summary.matches_found || 0}`);
        }
      } else {
        console.log('❌ Failed!');
        console.log(`📊 Status: ${response.status}`);
        console.log(`💬 Error: ${result.error || 'Unknown error'}`);
      }
      
      console.log('─'.repeat(80));
      
      // Wait 2 seconds between requests
      await new Promise(resolve => setTimeout(resolve, 2000));
      
    } catch (error) {
      console.error(`🚨 Error testing ${testCase.name}:`, error.message);
      console.log('─'.repeat(80));
    }
  }
  
  console.log('\n🎉 Testing completed!');
}

// Run the test
testTeamsWebhook().catch(console.error);
