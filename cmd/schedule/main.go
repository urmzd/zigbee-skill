package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/urmzd/sunrise-lamp/pkg"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ScheduleEvent struct {
	TargetBrightness int `json:"target_brightness"`
}

type SunAPIResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
		CivilTwilightBegin string `json:"civil_twilight_begin"`
		CivilTwilightEnd   string `json:"civil_twilight_end"`
	} `json:"results"`
	Status string `json:"status"`
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	var scheduleEvent ScheduleEvent
	err := json.Unmarshal(event.Detail, &scheduleEvent)
	if err != nil {
		return err
	}

	pkg.MqttServer = os.Getenv("MQTT_SERVER")
	pkg.MqttUsername = os.Getenv("MQTT_USERNAME")
	pkg.MqttPassword = os.Getenv("MQTT_PASSWORD")
	pkg.DeviceFriendlyName = os.Getenv("DEVICE_FRIENDLY_NAME")

	client := pkg.NewClient(pkg.MqttServer, pkg.MqttUsername, pkg.MqttPassword)
	currentBrightness := pkg.GetCurrentBrightness(client, pkg.DeviceFriendlyName)

	// Get sunrise and sunset times
	lat := os.Getenv("LAT")
	lng := os.Getenv("LNG")
	sunAPIURL := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%s&lng=%s&timezone=UTC&date=today", lat, lng)
	sunAPIResponse := &SunAPIResponse{}

	resp, err := http.Get(sunAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(sunAPIResponse); err != nil {
		return err
	}

	// Schedule events
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	bucket := os.Getenv("BUCKET")


	s3Client := s3.NewFromConfig(cfg)
	schedulerClient := scheduler.NewFromConfig(cfg)

	// Get the config file

	// Calculate the duration and interval between events
	sunrise, _ := time.Parse(time.RFC3339, sunAPIResponse.Results.Sunrise)
	dawn, _ := time.Parse(time.RFC3339, sunAPIResponse.Results.CivilTwilightBegin)
	sunset, _ := time.Parse(time.RFC3339, sunAPIResponse.Results.Sunset)
	dusk, _ := time.Parse(time.RFC3339, sunAPIResponse.Results.CivilTwilightEnd)

	sunriseDuration := sunrise.Sub(dawn)
	sunsetDuration := dusk.Sub(sunset)

	sunriseInterval := sunriseDuration / 

	return nil

}


func main() {
	lambda.Start(handler)
}
