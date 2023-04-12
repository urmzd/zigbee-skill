package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pkg "github.com/urmzd/sunrise-lamp/pkg"
)

type ControlEvent struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type DeviceMapping struct {
	Name       string `json:"Name"`
	DeviceName string `json:"DeviceName"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	lambda.Start(handler)
}

func getDeviceName(ctx context.Context, name string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load default config")
		return "", err
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("DEVICE_MAPPING_TABLE")

	result, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{
				Value: name,
			},
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to get item from DynamoDB")
		return "", err
	}

	if result.Item == nil {
		log.Error().Msg("No item found for name")
		return "", err
	}

	var deviceMapping DeviceMapping
	err = attributevalue.UnmarshalMap(result.Item, &deviceMapping)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal item")
		return "", err
	}

	return deviceMapping.DeviceName, nil
}

func handler(ctx context.Context, event ControlEvent) error {
	log.Info().Str("Name", event.Name).Int("Level", event.Level).Msg("Received control event")

	deviceName, err := getDeviceName(ctx, event.Name)
	if err != nil {
		return err
	}

	pkg.MqttServer = os.Getenv("SERVER")
	pkg.MqttUser = os.Getenv("USER")
	pkg.MqttPassword = os.Getenv("PASSWORD")
	pkg.DeviceFriendlyName = deviceName

	client := pkg.NewClient(pkg.MqttServer, pkg.MqttUser, pkg.MqttPassword)
	currentBrightness := pkg.GetCurrentBrightness(client, pkg.DeviceFriendlyName)

	log.Info().Int("CurrentBrightness", currentBrightness).Msg("Current brightness retrieved")

	pkg.SetBrightness(client, pkg.DeviceFriendlyName, currentBrightness+1)

	log.Info().Msg("Brightness updated successfully")

	return nil
}
