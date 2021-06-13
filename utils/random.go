package utils

import (
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const intBytes = "1234567890"

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func GetRandomString(length int, stringType string) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if strings.EqualFold(stringType, "string") {
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
		} else if strings.EqualFold(stringType, "int") {
			if idx := int(cache & letterIdxMask); idx < len(intBytes) {
				b[i] = intBytes[idx]
				i--
			}
		} else if strings.EqualFold(stringType, "effectiveNumber") {
			if idx := int(cache & letterIdxMask); idx < len(intBytes) {
				if i == 0 {
					if intBytes[idx] == uint8(48) {
						b[i] = uint8(49)
					} else {
						b[i] = intBytes[idx]
					}
				} else {
					b[i] = intBytes[idx]
				}
				i--
			}
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
