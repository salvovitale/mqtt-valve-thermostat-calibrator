package calibrator

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/config"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/model"
	"github.com/salvovitale/mqtt-valve-thermostat-calibrator/internal/pubsub"
)

type PairedSensors struct {
	name           string
	thermostatData model.Thermostat
	senorData      model.Sensor
	doneThermo     chan struct{}
	doneSensor     chan struct{}
	doneProcessing chan struct{}
	donePublishing chan struct{}
	publishingCh   chan float64
	thermostatCh   chan model.Thermostat
	sensorCh       chan model.Sensor
}

func NewPair(name string) *PairedSensors {
	return &PairedSensors{
		name:           name,
		doneThermo:     make(chan struct{}),
		doneSensor:     make(chan struct{}),
		doneProcessing: make(chan struct{}),
		donePublishing: make(chan struct{}),
		publishingCh:   make(chan float64),
		thermostatCh:   make(chan model.Thermostat),
		sensorCh:       make(chan model.Sensor),
	}
}

func (p *PairedSensors) Initialize(brokerUrl string, dev config.PairedSensorsConfig, mqttConfig config.MqttConfig) {

	// starting handling message routine
	go handleMessage(p.thermostatCh, p.sensorCh, p.publishingCh, p.doneProcessing)

	var wg sync.WaitGroup

	// init thermostat mqtt client
	thermostatFullTopic := fmt.Sprintf("%s/%s", mqttConfig.BaseTopic, dev.ThermostatTopic)
	thermostatClientId := fmt.Sprintf("%s-%s", dev.ThermostatTopic, "subscriber")
	wg.Add(1)
	go subscribeToTopic(brokerUrl, thermostatFullTopic, thermostatClientId, mqttConfig.QoS, p.messageHandlerThermostat, p.doneThermo, &wg)

	// init sensor mqtt client
	sensorFullTopic := fmt.Sprintf("%s/%s", mqttConfig.BaseTopic, dev.SensorTopic)
	sensorClientId := fmt.Sprintf("%s-%s", dev.SensorTopic, "subscriber")
	wg.Add(1)
	go subscribeToTopic(brokerUrl, sensorFullTopic, sensorClientId, mqttConfig.QoS, p.messageHandlerSensor, p.doneSensor, &wg)

	// init publish client for calibration
	publishFullTopic := fmt.Sprintf("%s/%s/%s", mqttConfig.BaseTopic, dev.ThermostatTopic, dev.CalibrationSubTopic)
	clientId := fmt.Sprintf("%s-%s", dev.ThermostatTopic, "publisher")
	wg.Add(1)
	go handlePublicationToTopic(brokerUrl, publishFullTopic, clientId, mqttConfig.QoS, p.messageHandlerThermostat, p.donePublishing, p.publishingCh, &wg)
	wg.Wait()
}

func (p *PairedSensors) Stop(wg *sync.WaitGroup) error {
	close(p.doneThermo)
	close(p.doneSensor)
	close(p.doneProcessing)
	close(p.donePublishing)
	time.Sleep(1 * time.Second)
	wg.Done()
	return nil
}

func (p *PairedSensors) messageHandlerThermostat(client mqtt.Client, msg mqtt.Message) {
	go processThermostatMessageFromHandler(msg, p.thermostatCh)
}

func (p *PairedSensors) messageHandlerSensor(client mqtt.Client, msg mqtt.Message) {
	go processSensorMessageFromHandler(msg, p.sensorCh)

}

func subscribeToTopic(brokerUrl string, topic string, clientID string, qos int, messageHandler mqtt.MessageHandler, done chan struct{}, wg *sync.WaitGroup) {
	mqttClient, err := pubsub.New(brokerUrl, topic, clientID, qos)
	if err != nil {
		log.Error().Err(err).Msg("Error creating mqtt client")

	}
	if err := mqttClient.Subscribe(messageHandler); err != nil {
		log.Error().Err(err).Msg("Error subscribing to topic")
	}
	<-done
	mqttClient.Disconnect()
	wg.Done()
}

func handlePublicationToTopic(brokerUrl string, topic string, clientID string, qos int, messageHandler mqtt.MessageHandler, done chan struct{}, publishingCh chan float64, wg *sync.WaitGroup) {
	mqttClient, err := pubsub.New(brokerUrl, topic, clientID, qos)
	if err != nil {
		log.Error().Err(err).Msg("Error creating mqtt client")
	}
	for {
		select {
		case c := <-publishingCh:
			calibrationPayload := fmt.Sprintf("%.1f", c)
			log.Info().Msgf("Publishing calibration payload %s to topic %s", calibrationPayload, topic)
			payload, err := json.Marshal(calibrationPayload)
			if err != nil {
				log.Error().Err(err).Msg("Error marshalling payload")
			}
			if err := mqttClient.Publish(payload); err != nil {
				log.Error().Err(err).Msgf("Error publishing calibration to topic %s", topic)
			}
		case <-done:
			log.Info().Msg("exiting handling publication message routine")
			mqttClient.Disconnect()
			wg.Done()
			return
		}
	}
}

func handleMessage(thermostatCh chan model.Thermostat, sensorCh chan model.Sensor, publishingCh chan float64, done chan struct{}) {
	// if temperature is negative has never been initialized so it cannot be used to calculate calibration
	thermostatMeasuredTemperatureLastValue := -1.0
	sensorTemperatureLastValue := -1.0

	for {
		select {
		case t := <-thermostatCh:
			newMeasuredTemperature := t.Temperature - t.Calibration
			if newMeasuredTemperature != thermostatMeasuredTemperatureLastValue {
				log.Info().Msgf("Thermostat measured temperature: %f", newMeasuredTemperature)
				thermostatMeasuredTemperatureLastValue = newMeasuredTemperature
				if sensorTemperatureLastValue > 0.0 {
					newCalibration := calibrate(sensorTemperatureLastValue, thermostatMeasuredTemperatureLastValue)
					log.Info().Msgf("New calibration: %f", newCalibration)
					publishingCh <- newCalibration
				} else {
					log.Info().Msg("No sensor temperature value available yet")
				}
			}
		case s := <-sensorCh:
			if s.Temperature != sensorTemperatureLastValue {
				log.Info().Msgf("handling sensor message: %v", s)
				sensorTemperatureLastValue = s.Temperature
				if thermostatMeasuredTemperatureLastValue > 0.0 {
					newCalibration := calibrate(sensorTemperatureLastValue, thermostatMeasuredTemperatureLastValue)
					log.Info().Msgf("New calibration: %f", newCalibration)
					publishingCh <- newCalibration
				} else {
					log.Info().Msg("No thermostat temperature value available yet")
				}
			}
		case <-done:
			log.Info().Msg("exiting handling message routine")
			return
		}
	}
}

func processThermostatMessageFromHandler(msg mqtt.Message, thermostatCh chan model.Thermostat) {
	payload := model.ThermostatPayload{}
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Error().Err(err).Msgf("Error unmarshalling message from topic %s", msg.Topic())
		return
	}
	log.Debug().Msgf("handling thermostat message: %v from topic %s", payload, msg.Topic())
	thermostatCh <- model.Thermostat{
		Temperature: payload.LocalTemperature,
		Calibration: payload.LocalTemperatureCalibration,
	}
}

func processSensorMessageFromHandler(msg mqtt.Message, sensorCh chan model.Sensor) {
	payload := model.SensorPayload{}
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Error().Err(err).Msgf("Error unmarshalling message from topic %s", msg.Topic())
		return
	}
	log.Debug().Msgf("handling sensor message: %v from topic %s", payload, msg.Topic())

	sensorCh <- model.Sensor{Temperature: payload.Temperature}
}
