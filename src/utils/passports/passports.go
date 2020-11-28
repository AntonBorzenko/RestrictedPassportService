package passports

import (
	"errors"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func RemovePreviousFiles() error {
	files, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if len(name) < 14 {
			continue
		}
		prefix := name[:10]
		postfix := name[len(name) - 4:]

		if prefix == "passports_" && postfix == ".bz2" {
			log.Printf("Removing file '%v'\n", name)
			err := os.Remove(path.Join(os.TempDir(), name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetPassportsGenerator(tempFileName string, batchSize int) chan uint64 {
	result := make(chan uint64, batchSize)

	stringChan, err := utils.StringChannelFromBzip(tempFileName)
	e.Check(err)

	go func () {
		for passportString := range stringChan {
			passportNumber, err := ConvertPassportToUint64(passportString)
			if err == nil {
				result <- passportNumber
			}
		}
		close(result)
	}()

	return result
}