#!/bin/ash

touch /mosquitto/config/mosquitto_passwd && \
    mosquitto_passwd -b /mosquitto/config/mosquitto_passwd ${MQTT_USER} ${MQTT_PASSWORD}

# Start Mosquitto and run health check script in background
mosquitto -c /mosquitto/config/mosquitto.conf & /app/health_check.sh
