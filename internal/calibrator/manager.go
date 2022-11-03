package calibrator

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/config"
)

type CalibrationManager struct {
	config        *config.Config
	brokerURL     string
	pairedDevices []*PairedSensors
}

func New(cfg *config.Config) *CalibrationManager {
	broker := fmt.Sprintf("tcp://%s:%d", cfg.Mqtt.Host, cfg.Mqtt.Port)
	log.Info().Msgf("Initializing calibration manager for Broker: %s", broker)
	return &CalibrationManager{
		config:    cfg,
		brokerURL: broker,
	}
}

func (c *CalibrationManager) Start() error {
	log.Info().Msg("Starting calibration")
	for _, dev := range c.config.PairedDevices {
		pair := NewPair(dev.Name)
		c.pairedDevices = append(c.pairedDevices, pair)

		go pair.Initialize(c.brokerURL, dev, c.config.Mqtt)
	}
	return nil
}

func (c *CalibrationManager) Stop() error {
	log.Info().Msg("Stopping calibration")
	var wg sync.WaitGroup
	for _, pair := range c.pairedDevices {
		wg.Add(1)
		go pair.Stop(&wg)
	}
	wg.Wait()
	return nil
}
