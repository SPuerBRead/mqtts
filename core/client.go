package core

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"mqtts/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tmpClient = &Client{}
var tmpClientLock sync.Mutex

type Client struct {
	CurrentClient  MQTT.Client
	ClientOptions  *MQTT.ClientOptions
	ClientLocker   *sync.Mutex
	MessageHandler func(client *Client, msg *Message)
	Certificate    x509.Certificate
}

type Message struct {
	Topic string
	Msg   string
}

func ConnectWithOpts(opts *TargetOptions) Client {
	client := GetMQTTClient(opts)
	connectError := client.Connect()
	if connectError == nil {
		SetClientToken(opts.Host, opts.Port, *client)
		return *client
	} else {
		SetClientToken(opts.Host, opts.Port, *client)
		utils.OutputInfoMessage(opts.Host, opts.Port, "Connect mqtt server failed err:"+connectError.Error())
		return *client
	}
}

func GenerateClientId(id string) string {
	return "mqttSecurityCheck" + id
}

func connectHandler(client MQTT.Client) {
	//utils.OutputInfoMessage("Connect MQTT Service Success")
}

func connectLostHandler(client MQTT.Client, err error) {
	reader := client.OptionsReader()
	if len(reader.Servers()) > 0 {
		utils.OutputErrorMessageWithoutOption("[" + reader.Servers()[0].Host + "] Loss MQTT connection")
	}
}

func verifyPeerCertificateHandler(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	cert, _ := x509.ParseCertificate(rawCerts[0])
	tmpClient.Certificate = *cert
	setCertificate(tmpClient)
	return nil
}

func verifyConnectionHandler(state tls.ConnectionState) error {
	return nil
}

func GetMQTTClient(opts *TargetOptions) *Client {
	clientOptions := MQTT.NewClientOptions()
	broker := opts.Protocol + "://" + opts.Host + ":" + strconv.Itoa(opts.Port)
	if strings.EqualFold(opts.Protocol, "tcp") || strings.EqualFold(opts.Protocol, "ssl") {
		clientOptions.AddBroker(broker)
	}
	if strings.EqualFold(opts.Protocol, "ws") || strings.EqualFold(opts.Protocol, "wss") {
		clientOptions.AddBroker(broker + "/mqtt")
	}
	clientOptions.SetUsername(opts.UserName)
	clientOptions.SetPassword(opts.Password)
	if !strings.EqualFold(opts.ClientId, "") {
		clientOptions.SetClientID(opts.ClientId)
	} else {
		clientId := utils.GetRandomString(6, "string")
		clientOptions.SetClientID(utils.GetRandomString(6, "string"))
		opts.ClientId = clientId
	}
	clientOptions.SetConnectTimeout(time.Millisecond * 5000)
	clientOptions.SetOnConnectHandler(connectHandler)
	clientOptions.SetConnectionLostHandler(connectLostHandler)
	tlsConfig := tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyPeerCertificateHandler,
		VerifyConnection:      verifyConnectionHandler,
	}
	clientOptions.SetTLSConfig(&tlsConfig)
	client := MQTT.NewClient(clientOptions)
	return &Client{
		CurrentClient: client,
		ClientOptions: clientOptions,
		ClientLocker:  &sync.Mutex{},
	}
}

func (client *Client) Connect() error {
	if !client.CurrentClient.IsConnected() {
		client.ClientLocker.Lock()
		defer client.ClientLocker.Unlock()
		if !client.CurrentClient.IsConnected() {
			if token := client.CurrentClient.Connect(); token.Wait() && token.Error() != nil {
				client.saveCertificate(currentCertificate().Certificate)
				return token.Error()
			} else {
				client.saveCertificate(currentCertificate().Certificate)
			}
		}
	}
	return nil
}

func defaultMessageHandler(c *Client, message *Message) {
	fmt.Println(message.Msg)
}

func (client *Client) messageHandler(c MQTT.Client, message MQTT.Message) {
	if client.MessageHandler == nil {
		//utils.OutputErrorMessage("Not subscribe message")
		return
	}
	msg := &Message{
		Topic: message.Topic(),
		Msg:   string(message.Payload()),
	}
	client.MessageHandler(client, msg)
}

func (client *Client) Subscribe(handler func(c *Client, message *Message), qos byte, topics ...string) error {
	if client.MessageHandler != nil {
		return errors.New("messageHandler has been bound, have to clear messageHandler first")
	}

	if len(topics) == 0 {
		return errors.New("subscribe method must set topic")
	}

	if handler != nil {
		client.MessageHandler = handler
	} else {
		client.MessageHandler = defaultMessageHandler
	}

	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = qos
	}
	client.CurrentClient.SubscribeMultiple(filters, client.messageHandler)
	return nil
}

func (client *Client) Unsubscribe(topics ...string) {
	client.MessageHandler = nil
	client.CurrentClient.Unsubscribe(topics...)
}

func (client *Client) saveCertificate(certificate x509.Certificate) {
	client.Certificate = certificate
}

func currentCertificate() *Client {
	tmpClientLock.Lock()
	defer tmpClientLock.Unlock()
	return tmpClient
}

func setCertificate(c *Client) {
	tmpClientLock.Lock()
	tmpClient = c
	tmpClientLock.Unlock()
}
