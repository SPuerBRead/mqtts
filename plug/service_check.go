package plug

import (
	"encoding/base64"
	"mqtts/core"
	"mqtts/utils"
	"rsc.io/binaryregexp"
	"strconv"
	"strings"
)

/*
	check if it is a tcp/ssl service by sending a connection packet

	v3.1.1 	http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html
	v5 	  	http://docs.oasis-open.org/mqtt/mqtt/v5.0/cos01/mqtt-v5.0-cos01.html


	// Fixed Header
	0 0 0 1 0 0 0 0		MessageType[4-7] 1(MQTT_CONNECT)/Dup[3] 0/Qos[2] 0/Retain[1] 0
	0 0 0 1 0 0 0 1		remaining length 16 bytes
	// Variable Header
	0 0 0 0 0 0 0 0  	protocol name 4 bytes
	0 0 0 0 0 1 0 0		...
	1 0 0 1 1 1 0 1		M
	0 1 0 1 0 0 0 1 	Q
	0 1 0 1 0 1 0 0		T
	0 1 0 1 0 1 0 0		T
	0 0 0 0 0 1 0 0		protocol version 4 v3.1.1
	0 0 0 0 0 0 1 0		flags clear session 1 other 0
	0 0 0 0 0 0 0 0		keep alive 2 bytes 30s
	// Payload
	0 0 0 1 1 1 1 0		...
	0 0 0 0 0 0 0 0		payload length 5 bytes
	0 0 0 0 0 1 0 1		...
	1 0 0 1 1 1 0 1		M
	0 1 0 1 0 0 0 1 	Q
	0 1 0 1 0 1 0 0		T
	0 1 0 1 0 1 0 0		T
	0 0 1 1 0 1 0 1		S
*/

func ServiceCheck(opts *core.TargetOptions) {
	probePack := "\x10\x11\x00\x04MQTT\x04\x02\x00\x1e\x00\x05MQTTS"
	pattern := "^\x20\x02\x00.$"
	utils.OutputInfoMessage(opts.Host, opts.Port, "Check port protocol type...")

	utils.OutputInfoMessage(opts.Host, opts.Port, "Check if it is TCP protocol")
	tcpErr, tcpBuf := core.ConnectWithSingleProbePack(opts.Host, opts.Port, probePack)
	if tcpErr == nil {
		isTCPMatch, matchErr := binaryregexp.Match(pattern, tcpBuf)
		if matchErr == nil && isTCPMatch {
			opts.Protocol = "tcp"
			utils.OutputSuccessMessage(opts.Host, opts.Port, "Port protocol type is MQTT/TCP")
			return
		}
	}

	utils.OutputInfoMessage(opts.Host, opts.Port, "Check if it is SSL protocol")
	sslErr, sslBuf := core.ConnectWithSingleProbePackTCPTLS(opts.Host, opts.Port, probePack)
	if sslErr == nil {
		isSSLMatch, matchErr := binaryregexp.Match(pattern, sslBuf)
		if matchErr == nil && isSSLMatch {
			opts.Protocol = "ssl"
			utils.OutputSuccessMessage(opts.Host, opts.Port, "Port protocol type is MQTT/SSL")
			return
		}
	}

	utils.OutputInfoMessage(opts.Host, opts.Port, "Check if it is WS protocol")
	wsErr, wsBuf := core.ConnectWithSingleProbePack(opts.Host, opts.Port, constructWebsocketPacket(opts.Host, opts.Port))
	if wsErr == nil {
		wsLowerBuf := strings.ToLower(string(wsBuf))
		if strings.Contains(wsLowerBuf, "http/1.1 101 switching protocols") && strings.Contains(wsLowerBuf, "sec-websocket-protocol: mqtt") {
			opts.Protocol = "ws"
			utils.OutputSuccessMessage(opts.Host, opts.Port, "Port protocol type is MQTT/WS")
			return
		}
	}

	utils.OutputInfoMessage(opts.Host, opts.Port, "Check if it is WSS protocol")
	wssErr, wssBuf := core.ConnectWithSingleProbePackTCPTLS(opts.Host, opts.Port, constructWebsocketPacket(opts.Host, opts.Port))
	if wssErr == nil {
		wssLowerBuf := strings.ToLower(string(wssBuf))
		if strings.Contains(wssLowerBuf, "http/1.1 101 switching protocols") && strings.Contains(wssLowerBuf, "sec-websocket-protocol: mqtt") {
			opts.Protocol = "wss"
			utils.OutputSuccessMessage(opts.Host, opts.Port, "Port protocol type is MQTT/WSS")
			return
		}
	}
	utils.OutputErrorMessage(opts.Host, opts.Port, "Get MQTT protocol type failed, you can try use -protocol wss/ws/tcp/ssl in command to set protocol")
}

func constructWebsocketPacket(ip string, port int) string {
	rawRequest := "GET /mqtt HTTP/1.1"
	header := make(map[string]string)
	header["Sec-WebSocket-Version"] = "13"
	header["Sec-WebSocket-Key"] = base64.StdEncoding.EncodeToString([]byte("mqtts_" + utils.
		GetRandomString(5, "string")))
	header["Connection"] = "Upgrade"
	header["Upgrade"] = "websocket"
	header["Sec-WebSocket-Extensions"] = "permessage-deflate; client_max_window_bits"
	header["Sec-WebSocket-Protocol"] = "mqtt"
	header["Host"] = ip + ":" + strconv.Itoa(port)
	strHeader := ""
	for key, value := range header {
		strHeader += key + ": " + value + "\r\n"
	}
	rawRequest += "\r\n" + strHeader + "\r\n"
	return rawRequest
}
