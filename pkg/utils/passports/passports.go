package passports

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	"github.com/AntonBorzenko/RestrictedPassportService/utils/net"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

//goland:noinspection GoSnakeCaseUsage
var PASSPORT_REGEX = regexp.MustCompile(`^(\d{4})[ ,._-]?(\d{6})$`)

func CheckPassportNumber(number string) bool {
	return PASSPORT_REGEX.MatchString(number)
}

func ConvertPassportToUint64(number string) (uint64, error) {
	matches := PASSPORT_REGEX.FindStringSubmatch(number)
	if matches == nil {
		return 0, errors.New(`"` + number + `" is not a passport number`)
	}

	result, _ := strconv.ParseUint(matches[1] + matches[2], 10, 64)
	return result, nil
}

func GetPassportsGenerator(fileUrl string, batchSize int) chan uint64 {
	result := make(chan uint64, batchSize)

	tempZipFile := e.CheckFile(ioutil.TempFile("", "passports_*.zip"))
	tempUnzipDirectory := e.CheckString(ioutil.TempDir("", "passports_*"))

	e.Check(net.DownloadFile(tempZipFile, fileUrl))

	fmt.Println(`Unzipping file...`)
	unzippedFiles := e.CheckAString(utils.Unzip(tempZipFile.Name(), tempUnzipDirectory))
	if len(unzippedFiles) != 1 || unzippedFiles[0] != "list_of_expired_passports.csv" {
		panic(fmt.Sprintf("Wrong files structure in zip archive: %v", unzippedFiles))
	}
	e.Check(os.Remove(tempZipFile.Name()))

	csvFileName := filepath.Join(tempUnzipDirectory, "list_of_expired_passports.csv")
	reader := bufio.NewReader(e.CheckFile(os.Open(csvFileName)))
	go func () {
		defer e.Check(os.RemoveAll(tempUnzipDirectory))
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					panic(err)
				} else {
					break
				}
			}

			passportNumber, err := ConvertPassportToUint64(string(line))
			if err == nil {
				result <- passportNumber
			}
		}
		close(result)
	}()

	return result
}