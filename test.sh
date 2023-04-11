#!/usr/bin/env bash

function cleanup {
  echo "Cleaning up SSM port forwarding..."
  kill %1
  docker kill zigbee2mqtt
  docker rm zigbee2mqtt
}

trap cleanup EXIT

TARGET_CLUSTER_NAME=$(aws ecs list-clusters | jq -r '.clusterArns[0] | split("/")[-1]')
TASK_ARN=$(aws ecs list-tasks --cluster $TARGET_CLUSTER_NAME | jq -r '.taskArns[0]')
TASK_ID="${TASK_ARN##*/}"
RUNTIME_ID=$(aws ecs describe-tasks --cluster $TARGET_CLUSTER_NAME --tasks $TASK_ID --query 'tasks[0].containers[0].runtimeId' --output text)

TARGET="ecs:${TARGET_CLUSTER_NAME}_${TASK_ID}_${RUNTIME_ID}"

echo "$TARGET"

aws ssm start-session --target $TARGET\
  --parameters '{"portNumber":["1883"], "localPortNumber":["1883"]}'\
  --document-name AWS-StartPortForwardingSession &

wait %1

docker kill zigbee2mqtt
docker rm zigbee2mqtt

SECRET_ID="MqttCreds63F27723-1cD4KPLregs4"
SECRET=$(aws secretsmanager get-secret-value --secret-id $SECRET_ID --query 'SecretString' --output text)
ZIGBEE2MQTT_CONFIG_MQTT_USER=$(echo "$SECRET" | jq -r '.user')
ZIGBEE2MQTT_CONFIG_MQTT_PASSWORD=$(echo "$SECRET" | jq -r '.password')

echo "$ZIGBEE2MQTT_CONFIG_MQTT_USER"
echo "$ZIGBEE2MQTT_CONFIG_MQTT_PASSWORD"

./bin/cli --mqtt-password "$ZIGBEE2MQTT_CONFIG_MQTT_PASSWORD" --mqtt-user "$ZIGBEE2MQTT_CONFIG_MQTT_USER" set 0

docker run \
  --name zigbee2mqtt \
  --network="host" \
  -v $(pwd)/configs/zigbee2mqtt-data:/app/data \
  -v /run/udev:/run/udev:ro \
  --device /dev/ttyUSB0:/dev/ttyUSB0 \
  -e TZ=America/Toronto \
  -e ZIGBEE2MQTT_CONFIG_MQTT_USER="$ZIGBEE2MQTT_CONFIG_MQTT_USER" \
  -e ZIGBEE2MQTT_CONFIG_MQTT_PASSWORD="$ZIGBEE2MQTT_CONFIG_MQTT_PASSWORD" \
  --restart unless-stopped \
  koenkk/zigbee2mqtt
