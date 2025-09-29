package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"stafind-backend/internal/constants"
	"stafind-backend/internal/repositories"
	"time"
)

type notificationService struct {
	aiAgentRepo repositories.AIAgentRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(aiAgentRepo repositories.AIAgentRepository) NotificationService {
	return &notificationService{
		aiAgentRepo: aiAgentRepo,
	}
}

func (s *notificationService) SendTeamsMessage(channelID string, message string) error {
	webhookURL := os.Getenv(constants.EnvTeamsWebhookURL)
	if webhookURL == "" {
		return fmt.Errorf("TEAMS_WEBHOOK_URL not set")
	}

	// Create Teams message payload
	payload := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"summary":    "AI Agent Response",
		"themeColor": "0078D4",
		"sections": []map[string]interface{}{
			{
				"activityTitle":    "AI Agent Response",
				"activitySubtitle": "Employee Matching Results",
				"text":             message,
				"markdown":         true,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("teams webhook failed: %s", string(body))
	}

	return nil
}

func (s *notificationService) SendAdminEmail(subject string, body string) error {
	// Email sending disabled for now - just log the notification
	fmt.Printf("EMAIL NOTIFICATION (DISABLED):\nSubject: %s\nBody: %s\n\n", subject, body)
	return nil
}

func (s *notificationService) LogError(requestID int, error string) error {
	// Update the AI agent request with error status
	request, err := s.aiAgentRepo.GetByID(requestID)
	if err != nil {
		return err
	}

	request.Status = "failed"
	request.Error = &error
	now := time.Now()
	request.ProcessedAt = &now

	if err := s.aiAgentRepo.Update(requestID, request); err != nil {
		return err
	}

	// Send admin notification
	subject := fmt.Sprintf("AI Agent Error - Request ID: %d", requestID)
	body := fmt.Sprintf("An error occurred while processing AI agent request %d:\n\n%s", requestID, error)

	if err := s.SendAdminEmail(subject, body); err != nil {
		// Log the email error but don't fail the main operation
		fmt.Printf("Failed to send admin email: %v\n", err)
	}

	return nil
}
