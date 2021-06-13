package plug

import (
	"mqtts/core"
	"mqtts/utils"
	"strings"
)

const iterStartDigits = 5
const iterEndDigits = 10

var clientIdTypes = []string{"string", "int", "effectiveNumber"}

func FuzzAvailableClientId(opts *core.TargetOptions) []string {
	utils.OutputInfoMessage(opts.Host, opts.Port, "Start detecting available clientId...")
	var availableClientIds []string
	availableClientIdChan := make(chan string)
	for k := range utils.IterRange(iterStartDigits, iterEndDigits, 1) {
		for _, clientIdType := range clientIdTypes {
			tmpOpts := new(core.TargetOptions)
			utils.DeepCopy(tmpOpts, opts)
			tmpOpts.ClientId = utils.GetRandomString(k, clientIdType)
			go func() {
				client := core.GetMQTTClient(tmpOpts)
				err := client.Connect()
				if err == nil || !strings.EqualFold(err.Error(), "identifier rejected") {
					availableClientIdChan <- tmpOpts.ClientId
				} else {
					availableClientIdChan <- ""
				}
			}()
		}
	}
	for range utils.Iter((iterEndDigits - iterStartDigits) * len(clientIdTypes)) {
		clientId := <-availableClientIdChan
		if clientId != "" {
			availableClientIds = append(availableClientIds, clientId)
		}
	}
	if len(availableClientIds) == 0 {
		utils.OutputErrorMessage(opts.Host, opts.Port, "No available clientId found, you can try use -clientid **** in command to set clientid")
	} else {
		utils.OutputInfoMessage(opts.Host, opts.Port, "Available clientId: "+strings.Join(availableClientIds, ","))
		opts.ClientId = availableClientIds[0]
	}
	return availableClientIds
}
