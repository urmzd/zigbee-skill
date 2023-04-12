package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pkg "github.com/urmzd/sunrise-lamp/pkg"
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

	deviceMapping := pkg.DeviceMapping{
		Name:       event.Name,
		DeviceName: event.DeviceName,
	}

	s3Client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("CONFIG_BUCKET")
	log.Info().Msg("bucket: " + bucket)

	err = pkg.UpdateDeviceMapping(s3Client, ctx, bucket, event.Name, &deviceMapping)
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
