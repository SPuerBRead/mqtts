package utils

import (
	"bufio"
	"os"
	"strings"
)

func ReadFileByLine(path string) ([]string, error) {
	var data []string
	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.EqualFold(line, "") {
			data = append(data, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return data, err
	}
	return data, nil
}
