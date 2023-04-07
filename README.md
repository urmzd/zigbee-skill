# Sunlight Lamp

## Prerequisities

Tested On:
- Zigbee 3.0 USB Dongle Plus
- Syvania A19 70052

## Usage

```bash
# Run Dependencies
docker compose up

# CLI
./sunlight-lamp increase
```

## Notes

- In the morning, it should go from current brightness to max (as sun rises): it increments in sizes of 255 / civil twilight duration.
- While the sun sets, turn the lights on.
- We can send text messages to turn it on or off.

### Infrastructure: 

Note: We can create a simple cloudfront application (over api gateway) 
with a textbox to send commands such as

DESIRED COMMANDS:
- ON
- OFF
- INCREASE
- DECREASE
- SET SUNRISE 
- DELETE SUNRISE
- SET SUNSET
- DELETE SUNSET

Compute:
- Step Functions
- Lambda

Storage: 
- S3: Stores the configuration file.

Network:
- API Gateway: Protects the API + Static Content

General:
- Secrets Manager
- IOT Core
- CloudWatch Events (schedule events)

