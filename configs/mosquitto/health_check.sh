#!/bin/ash

# Function to check MQTT status using websocat
function check_mqtt_status {
  result=$(/app/websocat "ws://localhost:9001")
  
  if [ $? -ne 0 ]; then
    echo "HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\nError checking MQTT status" 
  else
    echo -e "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n$result"
  fi
}

while true; do
  echo -e "$(check_mqtt_status)" | nc -l -p 8081
done
