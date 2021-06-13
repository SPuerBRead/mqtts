package plug

import (
	"mqtts/core"
	"mqtts/utils"
	"sync"
)

var topics = []string{
	"#",
	"$SYS/#",
}

var topicInfo = make(map[string]string)

func GetMQTTTopicInfo(opts *core.TargetOptions, waitTime int) {
	utils.OutputInfoMessage(opts.Host, opts.Port, "Start getting topic info...")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go utils.RefreshInfo("Waiting for %s seconds...", waitTime, &wg)
	client := core.
		GetClientToken(opts, false)
	if client.CurrentClient.IsConnected() {
		err := client.Subscribe(func(c *core.Client, msg *core.Message) {
			if len(msg.Msg) > 30 {
				topicInfo[msg.Topic] = msg.Msg[1:90] + "..."
			} else {
				topicInfo[msg.Topic] = msg.Msg
			}
		}, 0, topics...)
		if err != nil {
			utils.OutputErrorMessage(opts.Host, opts.Port, "Subscribe topic failed")
			utils.OutputErrorMessage(opts.Host, opts.Port, err.Error())
		}
		wg.Wait()
	}
	outputTopicInfo(opts)
}

func outputTopicInfo(opts *core.TargetOptions) {
	if len(topicInfo) > 0 {
		utils.OutputSuccessMessage(opts.Host, opts.Port, "Server topic info:")
		utils.TableLogger([]string{"TOPIC", "VALUE"}, topicInfo)
	} else {
		utils.OutputErrorMessage(opts.Host, opts.Port, "Get topic info failed")
	}
}
