#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { SunriseLampStack } from "../lib/infrastructure-stack";

const app = new cdk.App();
const suffix="-1681280422669"
new SunriseLampStack(app, `SunriseLamp${suffix}`)
