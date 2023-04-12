#!/usr/bin/env python3

import argparse
import json
import sys
import paho.mqtt.client as mqtt

def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("Connected to the WebSocket.")
    else:
        print("Connection failed with code:", rc)
        sys.exit(1)

def on_publish(client, userdata, mid):
    print("Message published.")
    client.disconnect()

def on_message(client, userdata, message):
    print("Received message:", message.payload.decode())

def main():
    parser = argparse.ArgumentParser(description="A simple MQTT CLI for connecting to a WebSocket.")
    parser.add_argument("subcommand", choices=["publish"], help="The subcommand to execute.")
    parser.add_argument("-b", "--brightness", type=int, required=True, help="Brightness value to send.")
    args = parser.parse_args()

    if args.subcommand == "publish":
        host = "localhost"
        port = 9001
        username = "root"
        password = "pass"
        topic = "zigbee2mqtt/a19/set"

        client = mqtt.Client(transport="websockets")
        client.username_pw_set(username, password)
        client.on_connect = on_connect
        client.on_publish = on_publish
        client.on_message = on_message

        client.connect(host, port)

        message = json.dumps({"brightness": args.brightness})
        res = client.publish(topic, message, qos=0)

        client.loop_start()
        try:
            while True:
                pass
        except KeyboardInterrupt:
            print("Interrupted by user, stopping...")
        finally:
            client.loop_stop()
            client.disconnect()

if __name__ == "__main__":
    main()
