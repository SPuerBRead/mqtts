package core

import (
	"container/list"
	"sync"
)

type TargetOptions struct {
	Protocol string
	Host     string
	Port     int
	ClientId string
	UserName string
	Password string
}

type optionList struct {
	TargetInfo *list.List
}

var optsList *optionList
var optsOnce sync.Once

func getTargetInfo() *optionList {
	optsOnce.Do(func() {
		optsList = &optionList{}
		optsList.TargetInfo = list.New()
	})
	return optsList
}

func SetTargetInfo(protocol string, host string, port int, clientId string, username string, password string) {
	if optsList != nil && optsList.TargetInfo != nil {
		opts := &TargetOptions{
			Protocol: protocol,
			Host:     host,
			Port:     port,
			ClientId: clientId,
			UserName: username,
			Password: password,
		}
		optsList.TargetInfo.PushBack(opts)
	} else {
		getTargetInfo()
		opts := &TargetOptions{
			Protocol: protocol,
			Host:     host,
			Port:     port,
			ClientId: clientId,
			UserName: username,
			Password: password,
		}
		optsList.TargetInfo.PushBack(opts)
	}
}

func GetTargetInfo() *optionList {
	return optsList
}
