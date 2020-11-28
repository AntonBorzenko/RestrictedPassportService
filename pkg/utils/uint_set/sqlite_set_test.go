package uint_set

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	_ "github.com/mattn/go-sqlite3"
)

func TestSqliteSet_Has_Insert(t *testing.T) {
	tempFile := e.CheckFile(ioutil.TempFile("", "db_*.sqlite"))
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	set := NewSqliteSet(tempFile.Name(), true, true)
	defer set.Close()

	var (
		testPassportNumber uint64 = 1234_123456
		foundNumber uint64
	)

	e.Check(set.Insert(testPassportNumber))
	query := fmt.Sprintf("SELECT * FROM %v WHERE NUMBER=`%v`", set.tableName, testPassportNumber)
	row := set.db.QueryRow(query)
	e.Check(row.Scan(&foundNumber))
	if foundNumber != testPassportNumber {
		t.Errorf("got %v expected %v", foundNumber, testPassportNumber)
	}

	var (
		testPassportNumber2 uint64 = 4567_890123
		testPassportNumber3 uint64 = 5678_345678
	)

	result, err := set.Has(testPassportNumber3)
	e.Check(err)
	if result {
		t.Errorf("Unexpected number %v in set", testPassportNumber3)
	}

	e.Check(
		set.Insert(testPassportNumber2),
		set.Insert(testPassportNumber3))
	for _, value := range []uint64{testPassportNumber, testPassportNumber2, testPassportNumber3} {
		result, err = set.Has(value)
		e.Check(err)
		if !result {
			t.Errorf("Set does not have number %v", value)
		}
	}
}

func TestSqliteSet_CreateDB_CreateIndex(t *testing.T) {
	tempFile, err := ioutil.TempFile("dir", "db_*.sqlite")
	e.Check(err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	set := NewSqliteSet(tempFile.Name(), false, false)
	defer set.Close()

	_, err = set.db.Exec(fmt.Sprintf("SELECT * FROM %v LIMIT 2", set.tableName))
	if err == nil {
		t.Errorf("Table %v is exists", set.tableName)
	}

	set.CreateDB()
	_, err = set.db.Exec(fmt.Sprintf("SELECT * FROM %v LIMIT 2", set.tableName))
	e.Check(err)

	countFoundIndexes := func() int {
		var foundRows int
		row := set.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='set_index'")
		e.Check(row.Scan(&foundRows))
		return foundRows
	}
	if countFoundIndexes() != 0 {
		t.Errorf("Index with name 'set_index' is found")
	}

	set.CreateIndex()
	if countFoundIndexes() != 1 {
		t.Errorf("Index with name 'set_index' is not found")
	}
}

func sliceToChan(array []uint64) chan uint64 {
	result := make(chan uint64, len(array))
	for _, value := range array {
		result <- value
	}
	close(result)
	return result
}

func TestSqliteSet_InsertMultiple(t *testing.T) {
	tempFile, err := ioutil.TempFile("dir", "db_*.sqlite")
	e.Check(err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	set := NewSqliteSet(tempFile.Name(), true, true)
	defer set.Close()

	insertableValues := []uint64{4, 5, 6, 1234_567890, 3456_234567, 9876_123456}
	nonInsertableValues := []uint64{1, 2, 3, 5678_901234, 1234_123456}
	e.Check(set.InsertMultiple(sliceToChan(insertableValues), false))

	for _, value := range insertableValues {
		if !e.CheckBool(set.Has(value)) {
			t.Errorf("Value %v is not found in set", value)
		}
	}

	for _, value := range nonInsertableValues {
		if e.CheckBool(set.Has(value)) {
			t.Errorf("Value %v is found in set", value)
		}
	}
}
