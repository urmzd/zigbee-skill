# Infrastructure

## Overview
DESIRED COMMANDS:
- SET 
- GET 
- SET/DELETE SUNRISE
- SET/DELETE SUNSET

ECS Holding MQTT Broker
CloudWatch Events (for SUNRISE/SUNSET)
Secrets Manager
Lambda 
S3
Step Functions
CloudFormation (via CDK)


```go
type Config struct {
    Brightness: uint8 // between 0 - 10, 0 will set ON: false, and 1-10 will set ON: true
    Sunrise: {
        Set: bool,
        Start: time
        End: time
        Brightness: uint8 // Goes from current to this value.
    },
    Sunset: {
        Set: bool
        Start: time
        End: time
        Brightness: uint8 // Goes from current to this value.
    }
}
```

