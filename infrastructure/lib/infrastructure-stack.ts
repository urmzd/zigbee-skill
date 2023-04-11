import * as cdk from "aws-cdk-lib";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { Construct } from "constructs";
import * as path from "node:path"

export class InfrastructureStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Define VPC with a private and a public subnet
    const vpc = new ec2.Vpc(this, "VPC", {
      maxAzs: 2,
      subnetConfiguration: [
        {
          cidrMask: 24,
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
          name: "Private",
        },
      ],
    });

    // Define S3 bucket for the configuration file
    const configBucket = new s3.Bucket(this, "ConfigBucket");

    // Define ECS cluster and service in the private subnet
    const cluster = new ecs.Cluster(this, "MQTTBrokerCluster", {
      vpc: vpc,
    });

    const mqttCreds = new secretsmanager.Secret(this, "MqttCreds", {
      generateSecretString: {
        secretStringTemplate: JSON.stringify({
          username: "root",
        }),
        generateStringKey: "password",
      },
    });

    const taskDefinition = new ecs.FargateTaskDefinition(
      this,
      "MQTTBrokerTask"
    );

    const mqttPassword = mqttCreds.secretValueFromJson("password").toString()
    const mqttUsername = mqttCreds.secretValueFromJson("username").toString()

    taskDefinition
      .addContainer("MQTTBrokerContainer", {
        image: ecs.ContainerImage.fromAsset(path.resolve(process.cwd(), "../")),
        environment: {
          MQTT_USERNAME: mqttUsername,
          MQTT_PASSWORD: mqttPassword
        },
        logging: new ecs.AwsLogDriver({
          streamPrefix: "mqtt-broker",
        }),
      })
      .addPortMappings({
        containerPort: 1883,
      });

    const mqttBroker = new ecs.FargateService(this, "MQTTBrokerService", {
      cluster,
      taskDefinition,
      assignPublicIp: false,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
      },
    });

    // Create an Application Load Balancer
    const mqttLoadBalancer = new elbv2.NetworkLoadBalancer(this, "MQTTLoadBalancer", {
      vpc,
      internetFacing: false,
    });

    // Add a listener to the Application Load Balancer
    const mqttListener = mqttLoadBalancer.addListener("MQTTListener", {
      port: 1883,
      protocol: elbv2.Protocol.TCP,
    });

    // Add the MQTT broker as a target for the listener
    mqttListener.addTargets("MQTTBrokerTarget", {
      port: 1883,
      targets: [mqttBroker],
    });

    //const coreEnv = {
      //BUCKET: configBucket.bucketName,
      //SERVER: `mqtt://${mqttLoadBalancer.loadBalancerDnsName}:1883`,
      //CREDS: mqttCreds.secretArn,
      //DEVICE: "a19",
    //}

    //// We call this function several times as scheduled by the schedule lambda.
    //const increaseFunction = new lambda.Function(
      //this,
      //"LambdaIncreaseFunction",
      //{
        //runtime: lambda.Runtime.GO_1_X,
        //handler: "increase",
        //code: lambda.Code.fromAsset("../../bin"),
        //environment: coreEnv,
        //vpc: vpc,
        //vpcSubnets: {
          //subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
        //},
      //}
    //);

    //// This is scheduled every day at midnight, and schedules multiple increases for sunrise and sunset.
    //// We number the invocations to prevent redundant logic from running.
    //const scheduleFunction = new lambda.Function(
      //this,
      //"LambdaScheduleFunction",
      //{
        //runtime: lambda.Runtime.GO_1_X,
        //handler: "schedule",
        //code: lambda.Code.fromAsset("../../bin"),
        //environment: {
          //...coreEnv,
          //INCREASE_FUNCTION_ARN: increaseFunction.functionArn,
        //},
        //vpc: vpc,
        //vpcSubnets: {
          //subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
        //},
      //}
    //);

    //// Define EventBridge rule for midnight schedule
    //const midnightRule = new events.Rule(this, "MidnightSchedule", {
      //schedule: events.Schedule.rate(cdk.Duration.days(1)),
    //});

    //// Add Lambda function as a target for the EventBridge rule
    //midnightRule.addTarget(new targets.LambdaFunction(scheduleFunction));

    //// Grant required permissions to Lambda functions
    //configBucket.grantReadWrite(increaseFunction);

    //const mqttCredentials = secretsmanager.Secret.fromSecretNameV2(
      //this,
      //"MQTTCredentials",
      //"mqtt-credentials"
    //);

    //mqttCredentials.grantRead(increaseFunction);
  }
}
