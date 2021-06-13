package plug

import (
	"mqtts/core"
	"mqtts/utils"
)

func MQTTUnauthCheck(opts *core.TargetOptions) bool {
	client := core.GetMQTTClient(opts)
	connectError := client.Connect()
	if connectError != nil {
		utils.OutputNotVulnMessage(opts.Host, opts.Port, "Unauthorized access vulnerability not Exists")
		return false
	} else {
		core.SetClientToken(opts.Host, opts.Port, *client)
		utils.OutputVulnMessage(opts.Host, opts.Port, "Unauthorized access vulnerability exists")
		core.SetResult(opts, "Unauthorized access vulnerability exists")
		return true
	}
}
