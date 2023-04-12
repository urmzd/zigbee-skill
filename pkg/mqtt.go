package pkg

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
	//"github.com/rs/zerolog/log"
)

var (
	MqttServer         string
	MqttUser           string
	MqttPassword       string
	DeviceFriendlyName string
)

func NewClient(server string, user string, pass string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(server).SetUsername(user).SetPassword(pass).SetCleanSession(true).SetConnectTimeout(15 * time.Second)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	return client
}
