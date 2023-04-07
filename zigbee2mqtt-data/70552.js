const fz = require('zigbee-herdsman-converters/converters/fromZigbee');
const tz = require('zigbee-herdsman-converters/converters/toZigbee');
const exposes = require('zigbee-herdsman-converters/lib/exposes');
const ota = require('zigbee-herdsman-converters/lib/ota');
const reporting = require('zigbee-herdsman-converters/lib/reporting');
const extend = require('zigbee-herdsman-converters/lib/extend');

const e = exposes.presets;
const ea = exposes.access;

const definition = {
    zigbeeModel: ['A19 W non CEC'], // Update this line
    model: '70552',
    vendor: 'Sylvania',
    description: 'Smart A19 LED bulb',
    extend: extend.ledvance.light_onoff_brightness(),
    ota: ota.ledvance,
};

module.exports = definition;
