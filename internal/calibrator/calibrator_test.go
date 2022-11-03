package calibrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc                     string
		tempSensor               float64
		thermostatTemp           float64
		expectedCalibrationValue float64
	}{
		{
			desc:                     "21.6 - 17.0 = 4.6 -> 4.5",
			tempSensor:               21.6,
			thermostatTemp:           17.0,
			expectedCalibrationValue: 4.5,
		},
		{
			desc:                     "21.6 - 24.0 = -2.4 -> -2.5",
			tempSensor:               21.6,
			thermostatTemp:           24.0,
			expectedCalibrationValue: -2.5,
		},
		{
			desc:                     "21.6 - 21.0 = 0.6 -> 0.5",
			tempSensor:               21.6,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 0.5,
		},
		{
			desc:                     "21.6 - 21.5 = 0.1 -> 0.0",
			tempSensor:               21.6,
			thermostatTemp:           21.5,
			expectedCalibrationValue: 0.0,
		},
		{
			desc:                     "21.6 - 21.5 = 0.1 -> 0.0",
			tempSensor:               21.6,
			thermostatTemp:           21.5,
			expectedCalibrationValue: 0.0,
		},
		{
			desc:                     "21.33 - 21.0 = 0.33 -> 0.0",
			tempSensor:               21.33,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 0.0,
		},
		{
			desc:                     "21.34 - 21.0 = 0.34 -> 0.5",
			tempSensor:               21.34,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 0.5,
		},
		{
			desc:                     "21.66 - 21.0 = 0.66 -> 0.5",
			tempSensor:               21.66,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 0.5,
		},
		{
			desc:                     "21.67 - 21.0 = 0.67 -> 1.0",
			tempSensor:               21.67,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 1.0,
		},
		{
			desc:                     "22.33 - 21.0 = 1.33 -> 1.0",
			tempSensor:               22.33,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 1.0,
		},
		{
			desc:                     "22.34 - 21.0 = 1.34 -> 1.5",
			tempSensor:               22.34,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 1.5,
		},
		{
			desc:                     "22.66 - 21.0 = 1.66 -> 1.5",
			tempSensor:               22.66,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 1.5,
		},
		{
			desc:                     "22.67 - 21.0 = 1.67 -> 2.0",
			tempSensor:               22.67,
			thermostatTemp:           21.0,
			expectedCalibrationValue: 2.0,
		},
		{
			desc:                     "21.00 - 22.33 = -1.33 -> -1.0",
			tempSensor:               21.00,
			thermostatTemp:           22.33,
			expectedCalibrationValue: -1.0,
		},
		{
			desc:                     "21.00 - 22.34 = -1.34 -> -1.5",
			tempSensor:               21.00,
			thermostatTemp:           22.34,
			expectedCalibrationValue: -1.5,
		},
		{
			desc:                     "21.00 - 22.66 = -1.66 -> -1.5",
			tempSensor:               21.00,
			thermostatTemp:           22.66,
			expectedCalibrationValue: -1.5,
		},
		{
			desc:                     "21.00 - 22.67 = -1.67 -> -2.0",
			tempSensor:               21.00,
			thermostatTemp:           22.67,
			expectedCalibrationValue: -2.0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := calibrate(tC.tempSensor, tC.thermostatTemp)
			assert.Equal(t, tC.expectedCalibrationValue, result)
		})
	}
}
