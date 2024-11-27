package queue

import (
	"brain_vault/shared"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSSenderInterface interface {
	SendMessage(ctx context.Context, message string) (string, error)
}

var globalSQSSender SQSSenderInterface

func Initialize(appCtx shared.AppContext) {
	ctx := context.Background()
	if !appCtx.IsDev {
		sqsSender, err := NewSQSSender(ctx, SQSConfig{
			Region:   os.Getenv("AWS_REGION"),
			QueueURL: os.Getenv("SQS_QUEUE_URL"),
		})
		if err != nil {
			log.Fatalf("failed to initialize global SQS sender: %v", err)
		}
		globalSQSSender = sqsSender
	} else {
		globalSQSSender = &LocalSQSSender{}
	}
}

// SQSSender handles sending messages to an SQS queue
type SQSSender struct {
	client   *sqs.Client
	queueURL string
}

// SQSConfig provides configuration options for SQS sender
type SQSConfig struct {
	Region   string
	QueueURL string
}

// NewSQSSender creates a new SQS sender with flexible configuration
func NewSQSSender(ctx context.Context, cfg SQSConfig) (*SQSSender, error) {

	// Validate required configuration
	if cfg.Region == "" || cfg.QueueURL == "" {
		return nil, fmt.Errorf("region and queue URL are required: region=%s, queueURL=%s", cfg.Region, cfg.QueueURL)
	}

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config: %w", err)
	}

	return &SQSSender{
		client:   sqs.NewFromConfig(awsCfg),
		queueURL: cfg.QueueURL,
	}, nil
}

// SendMessage sends a message to the configured SQS queue
func (s *SQSSender) SendMessage(ctx context.Context, messageBody string) (string, error) {
	// Validate input
	if messageBody == "" {
		return "", fmt.Errorf("message body cannot be empty")
	}

	// Prepare send message input
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(messageBody),
	}

	// Send the message
	result, err := s.client.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to send message to SQS: %w", err)
	}

	return aws.ToString(result.MessageId), nil
}

// SendMessageWithDefaults is a convenience function for sending a message
func SendMessageWithDefaults(messageBody string) error {
	ctx := context.Background()

	messageID, err := globalSQSSender.SendMessage(ctx, messageBody)
	if err != nil {
		return fmt.Errorf("message send failed: %w", err)
	}

	log.Printf("INFO: Message sent successfully with ID: %s", messageID)
	return nil
}

// LocalSQSSender is a local implementation of the SQSSenderInterface
type LocalSQSSender struct{}

// SendMessage logs the message instead of sending it to SQS
func (l *LocalSQSSender) SendMessage(ctx context.Context, message string) (string, error) {
	// Make a POST request to localhost:5004
	url := "http://localhost:5004"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return "local-message-id", nil

}
