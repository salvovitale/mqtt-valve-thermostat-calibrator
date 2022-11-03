package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/calibrator"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/config"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var configPath string

func init() {
	flag.StringVar(&configPath, "path", "input/config.yml", "config file")
}

func main() {
	flag.Parse()

	// read config file
	cfg, err := config.ReadConfigFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	// initialize calibrator manager
	manager := calibrator.New(cfg)
	if err := manager.Start(); err != nil {
		log.Fatal().Err(err).Msg("Error starting calibration manager")
	}

	// initialize signal handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	// stop calibrator manager
	if err := manager.Stop(); err != nil {
		log.Fatal().Err(err).Msg("Error stopping calibration manager")
	}
	log.Info().Msg("Exiting...")
}
