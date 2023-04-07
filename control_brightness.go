package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	increaseCmd = &cobra.Command{
		Use:   "increase",
		Short: "Increase brightness",
		Run: func(cmd *cobra.Command, args []string) {
			controlBrightness(mqttServer, mqttUsername, mqttPassword, deviceFriendlyName, "increase")
		},
	}

	decreaseCmd = &cobra.Command{
		Use:   "decrease",
		Short: "Decrease brightness",
		Run: func(cmd *cobra.Command, args []string) {
			controlBrightness(mqttServer, mqttUsername, mqttPassword, deviceFriendlyName, "decrease")
		},
	}

	turnOnCmd = &cobra.Command{
		Use:   "on",
		Short: "Turn on the device",
		Run: func(cmd *cobra.Command, args []string) {
			controlState(mqttServer, mqttUsername, mqttPassword, deviceFriendlyName, "ON")
		},
	}

	turnOffCmd = &cobra.Command{
		Use:   "off",
		Short: "Turn off the device",
		Run: func(cmd *cobra.Command, args []string) {
			controlState(mqttServer, mqttUsername, mqttPassword, deviceFriendlyName, "OFF")
		},
	}
)

func init() {
	rootCmd.AddCommand(increaseCmd)
	rootCmd.AddCommand(decreaseCmd)
	rootCmd.AddCommand(turnOnCmd)
	rootCmd.AddCommand(turnOffCmd)
}

func setBrightness(client mqtt.Client, device string, brightness int) {
	message := make(map[string]int)
	message["brightness"] = brightness
	client.Publish("zigbee2mqtt/"+device+"/set", 0, false, toJSON(message))
}

func setState(client mqtt.Client, device, state string) {
	message := make(map[string]string)
	message["state"] = state
	client.Publish("zigbee2mqtt/"+device+"/set", 0, false, toJSON(message))
}

func toJSON(obj interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func getCurrentBrightness(client mqtt.Client, device string) (int, error) {
	var brightness int
	messageChannel := make(chan mqtt.Message)

	token := client.Subscribe("zigbee2mqtt/"+device, 0, func(client mqtt.Client, msg mqtt.Message) {
		messageChannel <- msg
	})
	token.Wait()

	if token.Error() != nil {
		return 0, token.Error()
	}

	getPayload := map[string]string{
		"brightness": "",
	}
	client.Publish("zigbee2mqtt/"+device+"/get", 0, false, toJSON(getPayload))

	select {
	case msg := <-messageChannel:
		var deviceData map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &deviceData); err != nil {
			return 0, err
		}
		log.Info().Msgf("Device data: %v", deviceData)
		if b, ok := deviceData["brightness"].(float64); ok {
			brightness = int(b)
		}
	case <-time.After(5 * time.Second):
		log.Info().Msg("Timeout waiting for brightness")
		return 0, fmt.Errorf("timeout waiting for brightness")
	}

	client.Unsubscribe("zigbee2mqtt/" + device)

	return brightness, nil
}

func controlBrightness(server, user, pass, device, action string) {
	opts := mqtt.NewClientOptions().AddBroker(server).SetUsername(user).SetPassword(pass)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	currentBrightness, err := getCurrentBrightness(client, device)
	if err != nil {
		log.Error().Err(err).Msg("Error getting current brightness")
		client.Disconnect(250)
		return
	}

	log.Info().Msgf("Current brightness: %d", currentBrightness)

	switch action {
	case "increase":
		newBrightness := currentBrightness + 50
		if newBrightness > 255 {
			newBrightness = 255
		}
		setBrightness(client, device, newBrightness)
	case "decrease":
		newBrightness := currentBrightness - 50
		if newBrightness < 0 {
			newBrightness = 0
		}
		setBrightness(client, device, newBrightness)
	}

	client.Disconnect(250)
}

func controlState(server, user, pass, device, state string) {
	opts := mqtt.NewClientOptions().AddBroker(server).SetUsername(user).SetPassword(pass)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	setState(client, device, state)

	client.Disconnect(250)
}
