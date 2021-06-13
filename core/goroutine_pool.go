package core

import (
	"github.com/panjf2000/ants/v2"
	"mqtts/utils"
	"os"
)

func CreatePool(poolSize int) *ants.Pool {
	options := ants.Options{
		ExpiryDuration:   0,
		PreAlloc:         false,
		MaxBlockingTasks: 0,
		Nonblocking:      false,
		PanicHandler:     nil,
		Logger:           nil,
	}
	p, err := ants.NewPool(poolSize, ants.WithOptions(options))
	if err != nil {
		utils.OutputErrorMessageWithoutOption("Create goroutine pool error: " + err.Error())
		os.Exit(0)
	}
	return p
}
