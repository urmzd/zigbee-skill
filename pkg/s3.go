package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type DeviceMapping struct {
	Name       string `json:"Name"`
	DeviceName string `json:"DeviceName"`
}

func LoadDeviceMapping(client *s3.Client, ctx context.Context, bucket, key string) (*DeviceMapping, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	resp, err := client.GetObject(ctx, input)
	if err != nil {
		var notFound *types.NoSuchKey
		if errors.As(err, &notFound) {
			return nil, err
		}
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	var deviceMapping DeviceMapping
	err = json.Unmarshal(buf.Bytes(), &deviceMapping)
	if err != nil {
		return nil, err
	}

	return &deviceMapping, nil
}

func UpdateDeviceMapping(client *s3.Client, ctx context.Context, bucket, key string, deviceMapping *DeviceMapping) error {
	deviceMappingBytes, err := json.Marshal(deviceMapping)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader(deviceMappingBytes),
	}

	_, err = client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
