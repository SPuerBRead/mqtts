package core

import (
	"flag"
	"fmt"
	"mqtts/utils"
	"os"
	"strconv"
	"strings"
)

type ScanArgs struct {
	CommonScan        bool
	UnauthScan        bool
	AnyPwdScan        bool
	SystemInfo        bool
	BruteScan         bool
	AutoScan          bool
	TopicsList        bool
	WaitTime          int
	UserPath          string
	PwdPath           string
	TargetFile        string
	GoroutinePoolSize int
}

func CmdArgsParser() *ScanArgs {
	var host string
	var port int
	var username string
	var password string
	var protocol string
	var clientId string
	var unauthScan bool
	var anyPwdScan bool
	var autoScan bool
	var systemInfo bool
	var topicsList bool
	var userPath string
	var pwdPath string
	var waitTime int
	var bruteScan bool
	var targetFile string
	var goroutinePoolSize int
	var showVersion bool
	flag.StringVar(&host, "t", "", "input target ip or host")
	flag.IntVar(&port, "p", 1883, "input target port default is 1883/tcp")
	flag.StringVar(&username, "username", "", "input username")
	flag.StringVar(&password, "password", "", "input password")
	flag.StringVar(&protocol, "protocol", "", "input protocol tcp/ssl/ws/wss")
	flag.StringVar(&clientId, "clientid", "", "input password default is 6 random string")
	flag.BoolVar(&unauthScan, "u", false, "unauth scan (support batch scanning)")
	flag.BoolVar(&anyPwdScan, "a", false, "any username/password login scan for EMQX emqx_plugin_template plugin (support batch scanning)")
	flag.BoolVar(&bruteScan, "b", false, "username and password brute force")
	flag.BoolVar(&autoScan, "au", false, "automatic scanning according to service conditions")
	flag.BoolVar(&systemInfo, "s", false, "mqtt server system topic info scan (batch scanning is not supported)")
	flag.BoolVar(&topicsList, "ts", false, "mqtt server topic list scan (batch scanning is not supported)")
	flag.StringVar(&userPath, "nf", "", "brute force username list file path, default is ./username.txt")
	flag.StringVar(&pwdPath, "pf", "", "brute force password list file path, default is ./password.txt")
	flag.IntVar(&waitTime, "w", 15, "systemInfo scan and topics scan wait time, unit: seconds, default 15s")
	flag.StringVar(&targetFile, "tf", "", "batch scan target file, line format split with \\t host port [protocol | clientId | username | password]")
	flag.IntVar(&goroutinePoolSize, "g", 10, "batch scan goroutine pool size")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()
	if showVersion {
		fmt.Println(utils.Version)
		os.Exit(0)
	}
	if strings.EqualFold(host, "") && strings.EqualFold(targetFile, "") {
		flag.Usage()
		os.Exit(0)
	}
	if !strings.EqualFold(host, "") && !strings.EqualFold(targetFile, "") {
		utils.OutputErrorMessageWithoutOption("Single targets and batch targets cannot be set at the same time")
		os.Exit(0)
	}
	if !strings.EqualFold(host, "") {
		setSingleTarget(protocol, host, port, clientId, username, password)
	}
	if !strings.EqualFold(targetFile, "") {
		setTargets(targetFile)
	}
	if !unauthScan && !anyPwdScan && !systemInfo && !topicsList && !bruteScan && !autoScan {
		utils.OutputErrorMessageWithoutOption("Must specify the type of scan")
		os.Exit(0)
	}
	if !strings.EqualFold(targetFile, "") && (systemInfo || topicsList) {
		utils.OutputErrorMessageWithoutOption("Topic info scanning and topic list scanning do not support batch")
		os.Exit(0)
	}
	return &ScanArgs{
		UnauthScan:        unauthScan,
		AnyPwdScan:        anyPwdScan,
		BruteScan:         bruteScan,
		AutoScan:          autoScan,
		SystemInfo:        systemInfo,
		TopicsList:        topicsList,
		WaitTime:          waitTime,
		UserPath:          userPath,
		PwdPath:           pwdPath,
		GoroutinePoolSize: goroutinePoolSize,
	}
}

func setSingleTarget(protocol string, host string, port int, clientId string, username string, password string) {
	SetTargetInfo(protocol, host, port, clientId, username, password)
}

func setTargets(targetFile string) {
	lines, err := utils.ReadFileByLine(targetFile)
	if err != nil {
		utils.OutputErrorMessageWithoutOption("Load target file failed")
		os.Exit(0)
	}
	for num, line := range lines {
		targetInfo := strings.Split(line, " ")
		if len(targetInfo) < 2 {
			utils.OutputErrorMessageWithoutOption("Target format or data error in line " + strconv.Itoa(num+1) + ": " + line)
		}
		if len(targetInfo) < 6 {
			for range utils.IterRange(6 - len(targetInfo)) {
				targetInfo = append(targetInfo, "")
			}
		}
		host := targetInfo[0]
		port, err := strconv.Atoi(targetInfo[1])
		if err != nil {
			utils.OutputErrorMessageWithoutOption("Target port parse error in line " + strconv.Itoa(num+1) + ": " + line)
			continue
		}
		protocol := targetInfo[2]
		clientId := targetInfo[3]
		username := targetInfo[4]
		password := targetInfo[5]
		SetTargetInfo(protocol, host, port, clientId, username, password)
	}

}
