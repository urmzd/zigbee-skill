import paho.mqtt.client as mqtt
import json


def set_brightness(client, device, brightness):
    scaled_brightness = brightness * 255 // 100
    message = {"brightness": scaled_brightness}
    topic = "zigbee2mqtt/" + device + "/set"
    payload = json.dumps(message)
    client.publish(topic, payload, qos=0, retain=False)


broker_url = "ws://localhost:9000"
client_id = "my-client-id"

# Create a new MQTT client with WebSocket transport
client = mqtt.Client(client_id=client_id, transport="websockets")


# Define callback functions for MQTT events
def on_connect(client, userdata, flags, rc):
    print("Connected to MQTT broker with result code " + str(rc))


def on_disconnect(client, userdata, rc):
    print("Disconnected from MQTT broker with result code " + str(rc))


client.on_connect = on_connect
client.on_disconnect = on_disconnect

client.username_pw_set("user", "password")

# Connect to the MQTT broker
client.connect(broker_url)

# Set the brightness of the "a19" device to 50%
set_brightness(client, "a19", 50)

# Disconnect from the MQTT broker
client.disconnect()
