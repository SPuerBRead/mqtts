package utils

import (
	"fmt"
	"github.com/pterm/pterm"
	"strconv"
	"sync"
	"time"
)

func OutputInfoMessage(host string, port int, msg string) {
	pterm.Info.Println("[" + host + ":" + strconv.Itoa(port) + "] " + msg)
}

func OutputErrorMessage(host string, port int, msg string) {
	prefixPrinter := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DescriptionMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.DescriptionPrefixStyle,
			Text:  " ERROR ",
		},
	}
	prefixPrinter.Println("[" + host + ":" + strconv.Itoa(port) + "] " + msg)
}

func OutputErrorMessageWithoutOption(msg string) {
	prefixPrinter := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DebugMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.DebugPrefixStyle,
			Text:  " ERROR ",
		},
	}
	prefixPrinter.Println(msg)
}

func OutputSuccessMessage(host string, port int, msg string) {
	pterm.Success.Println("[" + host + ":" + strconv.Itoa(port) + "] " + msg)
}

func OutputVulnMessage(host string, port int, msg string) {
	prefixPrinter := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.ErrorMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.ErrorPrefixStyle,
			Text:  " VULN ",
		},
	}
	prefixPrinter.Println("[" + host + ":" + strconv.Itoa(port) + "] " + msg)
}

func OutputNotVulnMessage(host string, port int, msg string) {
	prefixPrinter := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.WarningMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.WarningPrefixStyle,
			Text:  "SAFE",
		},
	}
	prefixPrinter.Println("[" + host + ":" + strconv.Itoa(port) + "] " + msg)
}

func RefreshInfo(stencil string, count int, wg *sync.WaitGroup) {
	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf(stencil, strconv.Itoa(count)))
	for c := count - 1; c >= 0; c-- {
		time.Sleep(time.Second)
		introSpinner.UpdateText(fmt.Sprintf(stencil, strconv.Itoa(c)))
	}
	introSpinner.Stop()
	wg.Done()
}

func TableLogger(header []string, data map[string]string) {
	var tableData [][]string
	tableData = append(tableData, header)
	for topic, value := range data {
		tableData = append(tableData, []string{topic, value})
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func ScanResultLogger(data [][]string) {
	var tableData [][]string
	header := []string{"protocol", "host", "port", "clientId", "username", "password", "vuln"}
	tableData = append(tableData, header)
	for _, item := range data {
		tableData = append(tableData, item)
	}
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}
