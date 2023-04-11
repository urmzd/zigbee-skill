# Use the official Mosquitto image as the base image
FROM eclipse-mosquitto:latest

# Set the working directory
WORKDIR /app

ARG MQTT_USERNAME
ARG MQTT_PASSWORD

# Copy configuration file
COPY configs/mosquitto/config/mosquitto.conf /mosquitto/config/mosquitto.conf

COPY ./entrypoint.sh /app/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]
