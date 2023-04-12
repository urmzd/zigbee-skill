# Use the official Mosquitto image as the base image
FROM eclipse-mosquitto:latest

# Install required packages
RUN apk add --no-cache openssl

# Copy configuration file
COPY config/mosquitto.conf /mosquitto/config/mosquitto.conf

# Copy entrypoint script
COPY ./entrypoint.sh /app/entrypoint.sh

# Copy the health check script
COPY ./health_check.sh /app/health_check.sh

# Copy websocat binary
ADD https://github.com/vi/websocat/releases/download/v1.11.0/websocat_max.x86_64-unknown-linux-musl	/app/websocat

# Make websocat executable
RUN chmod +x /app/websocat

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
