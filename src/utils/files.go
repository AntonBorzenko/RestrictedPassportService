package utils

import (
	"bufio"
	"compress/bzip2"
	"io/ioutil"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateTempFile(pattern string)  (result string, err error) {
	tempFile, err := ioutil.TempFile("", pattern)
	if err != nil {
		return
	}
	err = tempFile.Close()
	if err != nil {
		return
	}

	result = tempFile.Name()
	return
}

func StringChannelFromBzip(fileName string) (chan string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bzip2.NewReader(bufio.NewReader(file)))
	result := make(chan string)

	go func() {
		for scanner.Scan() {
			result <- scanner.Text()
		}
		close(result)
	}()

	return result, nil
}