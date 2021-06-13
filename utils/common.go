package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func Iter(n int) []int {
	return make([]int, n)
}

func IterRange(args ...int) chan int {
	var start, stop int
	var step int = 1
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, stop, step = args[0], args[1], args[2]
	}

	ch := make(chan int)
	go func() {
		if step > 0 {
			for start < stop {
				ch <- start
				start = start + step
			}
		} else {
			for start > stop {
				ch <- start
				start = start + step
			}
		}
		close(ch)
	}()

	return ch
}

func OutputStruct(s interface{}) string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	structStr, err := json.Marshal(data)
	if err == nil {
		return string(structStr)
	} else {
		OutputErrorMessageWithoutOption("struct to string error:" + err.Error())
		return ""
	}
}
