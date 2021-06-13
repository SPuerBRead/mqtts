package plug

import (
	"crypto/x509"
	"fmt"
	certinfo "github.com/grantae/certinfo"
	"mqtts/core"
	"mqtts/utils"
	"strings"
	"sync"
)

var SystemTopic = []string{
	"$SYS/broker/version",
	"$SYS/broker/timestamp",
	"$SYS/broker/uptime",
	"$SYS/broker/subscriptions/count",
	"$SYS/broker/clients/connected",
	"$SYS/broker/clients/expired",
	"$SYS/broker/clients/disconnected",
	"$SYS/broker/clients/maximum",
	"$SYS/broker/clients/total",
}

var serverInfo = make(map[string]string)
var certificate []x509.Certificate

func GetMQTTServerInfo(opts *core.TargetOptions, waitTime int) {
	utils.OutputInfoMessage(opts.Host, opts.Port, "Start getting server info")
	client := core.GetClientToken(opts, false)
	certificate = append(certificate, client.Certificate)
	if client.CurrentClient.IsConnected() {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go utils.RefreshInfo("Waiting for %s seconds...", waitTime, &wg)
		err := client.Subscribe(func(c *core.Client, msg *core.Message) {
			serverInfo[msg.Topic] = msg.Msg
		}, 0, SystemTopic...)
		if err != nil {
			utils.OutputErrorMessage(opts.Host, opts.Port, "Subscribe server info topic failed")
			utils.OutputErrorMessage(opts.Host, opts.Port, err.Error())
		}
		wg.Wait()
	}
	outputServerInfo(opts)
}

func outputServerInfo(opts *core.TargetOptions) {
	if strings.EqualFold(opts.Protocol, "ssl") || strings.EqualFold(opts.Protocol, "wss") {
		if len(certificate) > 0 {
			utils.OutputSuccessMessage(opts.Host, opts.Port, "Server certificate Info:")
			certificateText, parseCertificateErr := certinfo.CertificateText(&certificate[0])
			if parseCertificateErr != nil {
				utils.OutputErrorMessage(opts.Host, opts.Port, "Parser certificate failed "+parseCertificateErr.Error())
			}
			fmt.Print(certificateText)
		}
	} else {
		utils.OutputInfoMessage(opts.Host, opts.Port, opts.Protocol+" service no certificate info")
	}
	if len(serverInfo) > 0 {
		utils.OutputSuccessMessage(opts.Host, opts.Port, "Server system info:")
		utils.TableLogger([]string{"TOPIC", "VALUE"}, serverInfo)
	} else {
		utils.OutputInfoMessage(opts.Host, opts.Port, "No server system topics message receive")
	}
}
