package errors

import (
	"database/sql"
	"os"
)

func Check(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

func CheckI64(result int64, err error) int64 {
	Check(err)
	return result
}

func CheckUI64(result uint64, err error) uint64 {
	Check(err)
	return result
}

func CheckBool(result bool, err error) bool {
	Check(err)
	return result
}

func CheckString(result string, err error) string {
	Check(err)
	return result
}

func CheckDbResult(result sql.Result, err error) sql.Result {
	Check(err)
	return result
}

func CheckAString(result []string, err error) []string {
	Check(err)
	return result
}

func CheckFile(result *os.File, err error) *os.File {
	Check(err)
	return result
}

func CheckCast(result interface{}, err error) interface{} {
	Check(err)
	return result
}
