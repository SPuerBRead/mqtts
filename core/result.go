package core

import (
	"sync"
)

var resultMap map[*TargetOptions][]string
var lock sync.Mutex

func InitResultMap() {
	resultMap = make(map[*TargetOptions][]string)
}

func SetResult(opts *TargetOptions, vulnInfo string) {
	lock.Lock()
	defer lock.Unlock()
	resultMap[opts] = append(resultMap[opts], vulnInfo)
}

func GetResult() map[*TargetOptions][]string {
	return resultMap
}
