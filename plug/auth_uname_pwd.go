package plug

import (
	"mqtts/core"
	"mqtts/utils"
)

func MQTTAuthUnamePwdCheck(opts core.TargetOptions) bool {
	client := core.GetMQTTClient(&opts)
	connectError := client.Connect()
	if connectError != nil {
		utils.OutputErrorMessage(opts.Host, opts.Port, "Username or password is wrong")
		return false
	} else {
		utils.OutputInfoMessage(opts.Host, opts.Port, "Correct username and password")
		core.SetClientToken(opts.Host, opts.Port, *client)
		return true
	}
}
