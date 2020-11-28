package passports

import (
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	"testing"
)

var (
	goodPassports = []string{
		"1234 123456",
		"1234,123456",
		"1234.123456",
		"1234_123456",
		"1234123456",
		"7890_876543",
	}
	badPassports = []string{
		"1234:123456",
		"1234_12345a",
		"123412345",
		"12341234567",
	}
	parseResults = []uint64{
		1234_123456,
		1234_123456,
		1234_123456,
		1234_123456,
		1234_123456,
		7890_876543,
	}
	testFileExpectedValues = []uint64 {
		2900088627,
		2900088630,
		2900088631,
		2900088633,
		2900088635,
		2900088637,
		2900088643,
		2900088647,
		2900088649,
	}
)

func TestCheckPassportNumber(t *testing.T) {
	for _, passport := range goodPassports {
		if !CheckPassportNumber(passport) {
			t.Errorf("Passport number '%v' is checked as wrong", passport)
		}
	}

	for _, passport := range badPassports {
		if CheckPassportNumber(passport) {
			t.Errorf("Passport number '%v' is checked as right", passport)
		}
	}
}

func TestConvertPassportToUint64(t *testing.T) {
	for i, passportString := range goodPassports {
		parsedNumber := e.CheckUI64(ConvertPassportToUint64(passportString))
		if parseResults[i] != parsedNumber {
			t.Errorf("Parse passport number '%v' to %v number, expected %v",
				passportString, parsedNumber, parseResults[i])
		}
	}
}

func TestGetPassportsGenerator(t *testing.T) {
	count := 0
	for found := range GetPassportsGenerator("passport_generator_test.csv.bz2", 10) {
		count += 1
		if count > len(testFileExpectedValues) {
			t.Errorf("Expected generator with %v of elements", len(testFileExpectedValues))
		}
		if found != testFileExpectedValues[count - 1] {
			t.Errorf("Expected number %v, found %v", testFileExpectedValues[count - 1], found)
		}
	}
	if count != len(testFileExpectedValues) {
		t.Errorf("Expected %v elements, found %v", len(testFileExpectedValues), count)
	}
}