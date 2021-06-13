package plug

import (
	"mqtts/core"
	"mqtts/utils"
	"strings"
)

// token errors: paho.mqtt.golang@v1.3.4/packets/packets.go
// error types

// unacceptable protocol version
// identifier rejected
// server Unavailable
// bad user name or password
// not Authorized
// network Error
// protocol Violation

func ClientIdCheck(opts *core.TargetOptions) bool {
	utils.OutputInfoMessage(opts.Host, opts.Port, "Check if current clientId available...")
	client := core.GetMQTTClient(opts)
	err := client.Connect()
	if err != nil && strings.EqualFold(err.Error(), "identifier rejected") {
		utils.OutputInfoMessage(opts.Host, opts.Port, "ClientId "+client.ClientOptions.ClientID+" unavailable")
		return false
	} else {
		utils.OutputInfoMessage(opts.Host, opts.Port, "ClientId "+client.ClientOptions.ClientID+" available")
		return true
	}
}
