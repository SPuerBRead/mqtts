package plug

import (
	"encoding/json"
	"mqtts/core"
	"mqtts/utils"
	"strconv"
)

func BruteUsernamePwd(opts *core.TargetOptions) {
	utils.OutputInfoMessage(opts.Host, opts.Port, "Start brute force cracking user name and password")
	utils.OutputInfoMessage(opts.Host, opts.Port, "Number of brute force cracks: "+strconv.Itoa(len(core.Username)*len(core.Password))+" ...")
	connectResult := make(chan map[string]interface{})
	for _, username := range core.Username {
		for _, password := range core.Password {
			tmpOpts := new(core.TargetOptions)
			utils.DeepCopy(tmpOpts, opts)
			tmpOpts.UserName = username
			tmpOpts.Password = password
			tmp := make(map[string]interface{})
			tmp["username"] = username
			tmp["password"] = password
			go func() {
				client := core.GetMQTTClient(tmpOpts)
				err := client.Connect()
				if err == nil {
					tmp["result"] = true
				} else {
					tmp["result"] = false
				}
				connectResult <- tmp
			}()
		}
	}
	realAuthInfo := make(map[string]interface{})
	for range utils.Iter(len(core.Username) * len(core.Password)) {
		scanResult := <-connectResult
		if scanResult["result"].(bool) {
			realAuthInfo = scanResult
			break
		}
	}
	outputBruteResult(opts, realAuthInfo)
	if len(realAuthInfo) > 0 {
		result, err := json.Marshal(realAuthInfo)
		if err == nil {
			core.SetResult(opts, string(result))
		} else {
			utils.OutputErrorMessageWithoutOption("json Marshal realAuthInfo error")
		}
	}
}

func outputBruteResult(opts *core.TargetOptions, result map[string]interface{}) {
	if len(result) == 0 {
		utils.OutputNotVulnMessage(opts.Host, opts.Port, "Brute force cracking of the username and password failed, and no correct username and password were found")
	} else {
		utils.OutputVulnMessage(opts.Host, opts.Port, "Found the correct username and password")
		auth := make(map[string]string)
		auth[result["username"].(string)] = result["password"].(string)
		utils.TableLogger([]string{"username", "password"}, auth)
	}
}
