package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"context"
)

func LoadConfig(client *s3.Client, ctx context.Context, bucket, key string) (*Config, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	resp, err := client.GetObject(ctx, input)
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, notFound) {
			config := DefaultConfig();
			err := UpdateConfig(client, ctx, bucket, key, config)
			if err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(buf.Bytes(), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig() *Config {
	sunriseConfig := &BrightnessConfig {
			TargetBrightness: 10,
			Set: true,
		}

	sunsetConfig := &BrightnessConfig{
		TargetBrightness: 10,
		Set: true,
	}

	locationConfig := &LocationConfig{
		// i.e., Halifax Nova Scotia
		Lat: 44.6476,
		Long: -63.5728,
	}

	return &Config{
		Sunrise: *sunriseConfig,
		Sunset: *sunsetConfig,
		Location: *locationConfig,
	}
}

func UpdateConfig(client *s3.Client, ctx context.Context, bucket, key string, config *Config) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(configBytes),
	}

	_, err = client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
