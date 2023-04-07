# Sunlight Lamp

## Prerequisities

You need a dimmable Zigbee light bulb.

## Use Cases

- In the morning, it should go from current brightness to max (as sun rises): it increments in sizes of 255 / civil twilight duration.
- We can send text messages to turn it on or off.

## Infrastructure

- IOT Core
- AWS VPC
- Event Bridge
- AWS Secrets Manager (IOT Core MQTT credentials)
