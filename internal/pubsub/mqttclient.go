package pubsub

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

const (
	timeToDisconnect = 250 // milliseconds
)

// the calibrate module should have a calibrate struct in which 2 mqtt client are spinned
// and send message to channels to calibrate if needed
type MqttClient struct {
	client    mqtt.Client
	brokerURL string
	topic     string
	clientID  string
	qos       int
}

func New(brokerURL string, topic string, clientID string, qos int) (*MqttClient, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(brokerURL)
	options.SetClientID(clientID)

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Info().Msgf("Connections established with broker %s from client %s", brokerURL, clientID)
	}

	var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log.Error().Err(err).Msgf("Connection lost from client %s", clientID)
	}
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Error connecting to broker")
		return nil, token.Error()
	}

	return &MqttClient{
		client,
		brokerURL,
		topic,
		clientID,
		qos,
	}, nil
}

func (c *MqttClient) Subscribe(messageHandler mqtt.MessageHandler) error {
	token := c.client.Subscribe(c.topic, byte(c.qos), messageHandler)
	if token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("Error subscribing to topic")
		return token.Error()
	}
	return nil
}

func (c *MqttClient) Publish(payload []byte) error {
	token := c.client.Publish(c.topic, byte(c.qos), false, payload)
	if token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("Error publishing to topic")
		return token.Error()
	}
	return nil
}

func (c *MqttClient) Disconnect() {
	log.Info().Msgf("Disconnecting from client %s", c.clientID)
	c.client.Disconnect(timeToDisconnect)
}
