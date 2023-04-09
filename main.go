package main

import (
	"github.com/spf13/cobra"
)

var (
	mqttServer          string
	mqttUsername        string
	mqttPassword        string
	deviceFriendlyName  string
)

var rootCmd = &cobra.Command{
	Use:   "sunlight-lamp",
	Short: "Control Sylvania A19 lightbulb using MQTT",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		mqttServer, _ = cmd.Flags().GetString("mqtt-server")
		mqttUsername, _ = cmd.Flags().GetString("mqtt-username")
		mqttPassword, _ = cmd.Flags().GetString("mqtt-password")
		deviceFriendlyName, _ = cmd.Flags().GetString("device-friendly-name")
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&mqttServer, "mqtt-server", "mqtt://localhost:1883", "MQTT server address")
	rootCmd.PersistentFlags().StringVar(&mqttUsername, "mqtt-username", "root", "MQTT username")
	rootCmd.PersistentFlags().StringVar(&mqttPassword, "mqtt-password", "pass", "MQTT password")
	rootCmd.PersistentFlags().StringVar(&deviceFriendlyName, "device-friendly-name", "a19", "Device friendly name")

	cobra.CheckErr(rootCmd.Execute())
}
