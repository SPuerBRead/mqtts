package core

import (
	"strconv"
	"sync"
)

type clientTokenMap struct {
	ClientInfo map[string]Client
}

var ctMap *clientTokenMap
var clientOnce sync.Once

func getClientMap() *clientTokenMap {
	clientOnce.Do(func() {
		ctMap = &clientTokenMap{}
		ctMap.ClientInfo = make(map[string]Client)
	})
	return ctMap
}

func SetClientToken(host string, port int, client Client) {
	if ctMap != nil && ctMap.ClientInfo != nil {
		ctMap.ClientInfo[host+":"+strconv.Itoa(port)] = client
	} else {
		getClientMap()
		ctMap.ClientInfo[host+":"+strconv.Itoa(port)] = client
	}
}

func GetClientToken(opts *TargetOptions, existing bool) Client {
	if ctMap != nil && ctMap.ClientInfo != nil {
		if _, existKey := ctMap.ClientInfo[opts.Host+":"+strconv.Itoa(opts.Port)]; existKey {
			if ctMap.ClientInfo[opts.Host+":"+strconv.Itoa(opts.Port)].CurrentClient.IsConnected() {
				return ctMap.ClientInfo[opts.Host+":"+strconv.Itoa(opts.Port)]
			}
		}
	}
	if !existing {
		return ConnectWithOpts(opts)
	} else {
		return Client{}
	}
}
