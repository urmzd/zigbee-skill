package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pkg "github.com/urmzd/sunrise-lamp/pkg"
)

type ControlEvent struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting Lambda function")
	lambda.Start(handler)
}

func getDeviceName(ctx context.Context, name string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load default config")
		return "", err
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("CONFIG_BUCKET")

	deviceMapping, err := pkg.LoadDeviceMapping(s3Client, ctx, bucket, name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load device mapping")
		return "", err
	}

	log.Info().Str("DeviceName", deviceMapping.DeviceName).Msg("Device name retrieved")
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

	log.Info().Str("MqttServer", pkg.MqttServer).Str("MqttUser", pkg.MqttUser).Str("DeviceFriendlyName", pkg.DeviceFriendlyName).Msg("Environment variables set")

	client := pkg.NewClient(pkg.MqttServer, pkg.MqttUser, pkg.MqttPassword)
	currentBrightness := pkg.GetCurrentBrightness(client, pkg.DeviceFriendlyName)

	log.Info().Int("CurrentBrightness", currentBrightness).Msg("Current brightness retrieved")

	pkg.SetBrightness(client, pkg.DeviceFriendlyName, currentBrightness+1)

	log.Info().Msg("Brightness updated successfully")

	return nil
}
