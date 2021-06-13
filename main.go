package main

import (
	"mqtts/core"
	"mqtts/plug"
	"mqtts/utils"
	"strconv"
	"strings"
	"sync"
)

func main() {
	scanArgs := core.CmdArgsParser()
	core.ShowBanner()
	targetsInfo := core.GetTargetInfo().TargetInfo
	core.InitResultMap()
	if scanArgs.AutoScan || scanArgs.BruteScan {
		core.LoadWordLists(scanArgs.UserPath, scanArgs.PwdPath)
	}
	pool := core.CreatePool(scanArgs.GoroutinePoolSize)
	var wg sync.WaitGroup
	defer pool.Release()
	for target := targetsInfo.Front(); target != nil; target = target.Next() {
		opts := (target.Value).(*core.TargetOptions)
		scanHandler := func() {
			startScan(&wg, opts, scanArgs)
		}
		wg.Add(1)
		submitTaskErr := pool.Submit(scanHandler)
		if submitTaskErr != nil {
			utils.OutputErrorMessageWithoutOption("Submit task failed, task info: " + utils.OutputStruct(*((target.Value).(*core.TargetOptions))))
		}
	}
	wg.Wait()
	outputResult()
}

func startScan(wg *sync.WaitGroup, opts *core.TargetOptions, scanArgs *core.ScanArgs) {
	defer func() {
		if err := recover(); err != nil {
			utils.OutputErrorMessageWithoutOption("Unknown error has occurred: " + err.(error).Error())
			wg.Done()
		}
	}()
	if strings.EqualFold(opts.Protocol, "") {
		plug.ServiceCheck(opts)
	}
	if strings.EqualFold(opts.Protocol, "") {
		wg.Done()
		return
	}
	if !plug.ClientIdCheck(opts) {
		availableClientId := plug.FuzzAvailableClientId(opts)
		if len(availableClientId) == 0 {
			wg.Done()
			return
		}
	}
	if scanArgs.UnauthScan {
		plug.MQTTUnauthCheck(opts)
	}
	if scanArgs.AnyPwdScan {
		plug.MQTTAnyPwdCheck(opts)
	}
	if scanArgs.SystemInfo {
		plug.GetMQTTServerInfo(opts, scanArgs.WaitTime)
	}
	if scanArgs.TopicsList {
		plug.GetMQTTTopicInfo(opts, scanArgs.WaitTime)
	}
	if scanArgs.BruteScan {
		plug.BruteUsernamePwd(opts)
	}
	if scanArgs.AutoScan {
		isAuth := plug.MQTTUnauthCheck(opts)
		if !isAuth {
			isAnyPwd := plug.MQTTAnyPwdCheck(opts)
			if !isAnyPwd {
				plug.BruteUsernamePwd(opts)
			}
		}
	}
	client := core.GetClientToken(opts, true)
	if client.CurrentClient != nil && client.CurrentClient.IsConnected() {
		client.CurrentClient.Disconnect(200)
	}
	wg.Done()
}

func outputResult() {
	var result [][]string
	for opts, vuln := range core.GetResult() {
		for _, item := range vuln {
			result = append(result, []string{opts.Protocol, opts.Host, strconv.Itoa(opts.Port), opts.ClientId, opts.UserName, opts.Password, item})
		}
	}
	if len(result) > 0 {
		utils.ScanResultLogger(result)
	}
}
