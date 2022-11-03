package model

type ThermostatPayload struct {
	LocalTemperature            float64 `json:"local_temperature"`
	LocalTemperatureCalibration float64 `json:"local_temperature_calibration"`
}

type SensorPayload struct {
	Temperature float64 `json:"temperature"`
}
