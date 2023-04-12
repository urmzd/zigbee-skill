package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CreateMappingEvent struct {
	Name       string `json:"name"`
	DeviceName string `json:"deviceName"`
}

func handler(ctx context.Context, event CreateMappingEvent) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DEVICE_MAPPING_TABLE")),
		Item: map[string]types.AttributeValue{
			"Name":       &types.AttributeValueMemberS{Value: event.Name},
			"DeviceName": &types.AttributeValueMemberS{Value: event.DeviceName},
		},
	}

	_, err = svc.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item, %v", err)
	}

	return nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	lambda.Start(handler)
}
