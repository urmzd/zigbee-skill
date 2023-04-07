package main

import (
	"github.com/spf13/cobra"
)

const (
	mqttServer          = "mqtt://localhost:1883"
	mqttUsername        = "root"
	mqttPassword        = "pass"
	deviceFriendlyName  = "0xf0d1b80000180b67"
)

var rootCmd = &cobra.Command{
	Use:   "sunlight-lamp",
	Short: "Control Sylvania A19 lightbulb using MQTT",
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
