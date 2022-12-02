package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Mqtt          MqttConfig            `yaml:"mqtt"`
	PairedDevices []PairedSensorsConfig `yaml:"paired_devices"`
}

type MqttConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	BaseTopic string `yaml:"base_topic"`
	QoS       int    `yaml:"qos"`
	Delay     int    `yaml:"delay"`
}

type PairedSensorsConfig struct {
	Name                string `yaml:"name"`
	SensorTopic         string `yaml:"sensor_topic"`
	ThermostatTopic     string `yaml:"thermostat_topic"`
	CalibrationSubTopic string `yaml:"calibration_sub_topic"`
}

func ReadConfigFile(configFile string) (*Config, error) {
	cfg := new(Config)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
