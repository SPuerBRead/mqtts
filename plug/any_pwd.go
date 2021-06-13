package plug

import (
	"mqtts/core"
	"mqtts/utils"
	"strings"
)

// emqx embox_plugin_template plugin open

func MQTTAnyPwdCheck(opts *core.TargetOptions) bool {
	opts1 := new(core.TargetOptions)
	utils.DeepCopy(opts1, opts)
	opts1.UserName = "fdc6fa8517d34c17"
	opts1.Password = utils.GetRandomString(15, "string")
	if strings.EqualFold(opts1.ClientId, "") {
		opts1.ClientId = core.GenerateClientId(utils.GetRandomString(8, "string"))
	}
	client1 := core.GetMQTTClient(opts1)
	utils.OutputInfoMessage(opts1.Host, opts1.Port, "Try connect server with username 'fdc6fa8517d34c17' and password '957b65edc4c0b1b9'")
	connectError1 := client1.Connect()

	if connectError1 != nil {
		utils.OutputInfoMessage(opts1.Host, opts1.Port, "Authentication failed with username 'fdc6fa8517d34c17' and password '957b65edc4c0b1b9'")
		utils.OutputNotVulnMessage(opts1.Host, opts1.Port, "Any password login vulnerability not exists")
		return false
	} else {
		opts2 := new(core.TargetOptions)
		utils.DeepCopy(opts2, opts)
		opts2.UserName = "3022412031bf49ee"
		opts2.Password = utils.GetRandomString(15, "string")
		if strings.EqualFold(opts2.ClientId, "") {
			opts2.ClientId = core.GenerateClientId(utils.GetRandomString(8, "string"))
		}
		utils.OutputInfoMessage(opts2.Host, opts2.Port, "Try connect server with username '3022412031bf49ee' and password '97a677594d1fcf01'")
		client2 := core.GetMQTTClient(opts2)
		connectError2 := client2.Connect()
		if connectError2 != nil {
			utils.OutputInfoMessage(opts2.Host, opts2.Port, "Authentication failed with username '3022412031bf49ee' and password '97a677594d1fcf01'")
			utils.OutputNotVulnMessage(opts2.Host, opts2.Port, "Any password login vulnerability not exists")
			return false
		} else {
			utils.OutputVulnMessage(opts2.Host, opts2.Port, "Any password login vulnerability exists")
			core.SetResult(opts2, "Any password login vulnerability exists")
			return true
		}
	}
}
