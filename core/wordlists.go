package core

import (
	"mqtts/utils"
	"strings"
)

var Username []string
var Password []string

func LoadWordLists(usernameFilePath string, passwordFilePath string) {
	var loadUsernameErr error
	var loadPasswordErr error
	if strings.EqualFold(usernameFilePath, "") {
		usernameFilePath = "./username.txt"
	}
	if strings.EqualFold(passwordFilePath, "") {
		passwordFilePath = "./password.txt"
	}
	Username, loadUsernameErr = utils.ReadFileByLine(usernameFilePath)
	if loadUsernameErr != nil {
		utils.OutputErrorMessageWithoutOption("Load username list file failed " + loadUsernameErr.Error())
	}
	Password, loadPasswordErr = utils.ReadFileByLine(passwordFilePath)
	if loadPasswordErr != nil {
		utils.OutputErrorMessageWithoutOption("Load password list file failed " + loadPasswordErr.Error())
	}
}
