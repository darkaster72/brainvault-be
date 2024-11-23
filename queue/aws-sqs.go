package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSSender is the service struct that will handle sending messages to SQS.
type SQSSender struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSSender creates a new instance of SQSSender.
func NewSQSSender(region, queueURL string) (*SQSSender, error) {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)

	return &SQSSender{
		client:   client,
		queueURL: queueURL,
	}, nil
}

// SendMessage sends a message to the configured SQS queue.
func (s *SQSSender) SendMessage(messageBody string) (string, error) {
	// Prepare the input
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(messageBody),
	}

	// Send the message
	result, err := s.client.SendMessage(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to send message, %v", err)
	}

	// Return the message ID on successful send
	return *result.MessageId, nil
}

func SendMessage(messageBody string) {
	// Initialize the SQS sender service
	region := "ap-south-1"
	queueURL := "https://sqs.ap-south-1.amazonaws.com/202533503295/brainvault-article-created"
	sqsSender, err := NewSQSSender(region, queueURL)
	if err != nil {
		log.Fatalf("Error initializing SQS sender: %v", err)
	}

	// Send the message
	messageID, err := sqsSender.SendMessage(messageBody)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	// Print the message ID
	fmt.Printf("Message sent successfully with ID: %s\n", messageID)
}
