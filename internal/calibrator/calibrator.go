package calibrator

import "math"

func calibrate(sensorTemperature float64, thermostatTemperature float64) float64 {
	newCalibrationValue := sensorTemperature - thermostatTemperature
	intPart, fracPart := math.Modf(newCalibrationValue)
	return intPart + round(fracPart)
}

func round(x float64) float64 {
	if x < 0 {
		return -round(-x)
	}
	if x <= 1.0/3.0 {
		return 0.0
	}
	if x > 1.0/3.0 && x <= 2.0/3.0 {
		return 1.0 / 2.0
	}
	return 1
}
