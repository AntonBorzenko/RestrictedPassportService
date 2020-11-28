package uint_set

import (
	"database/sql"
	"fmt"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteSet struct {
	tableName string
	db *sql.DB
}

func NewSqliteSet(filename string, createDb bool, createIndex bool) *SqliteSet {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	set := SqliteSet{"NumbersSet", db}
	if createDb {
		set.CreateDB()
	}
	if createIndex {
		set.CreateIndex()
	}
	return &set
}

func (set *SqliteSet) InsertMultiple(numbers chan uint64, ignoreErrors bool) error {
	// batchGenerator := utils.GetBatchGenerator(numbers, config.Cfg.DBBatchSize)
	stmt, err := set.db.Prepare(fmt.Sprintf("INSERT OR IGNORE INTO `%v` VALUES (?)", set.tableName))
	e.Check(err)

	for number := range numbers {
		_, err := stmt.Exec(number)
		if !ignoreErrors {
			return err
		}
	}

	return nil
}

func (set *SqliteSet) mustExec(query string) sql.Result {
	result, err := set.db.Exec(query)
	if err != nil {
		panic(err)
	}
	return result
}

func (set *SqliteSet) CreateDB() {
	set.mustExec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%v` (NUMBER INTEGER NOT NULL)", set.tableName))
}

func (set *SqliteSet) Close() error {
	return set.db.Close()
}

func (set *SqliteSet) CreateIndex() {
	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS set_index ON `%v`(NUMBER)", set.tableName)
	set.mustExec(query)
}

func (set *SqliteSet) Has(number uint64) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM `%v` WHERE NUMBER=%v", set.tableName, number)
	queryRow := set.db.QueryRow(query)

	var amount int
	err := queryRow.Scan(&amount)
	return amount != 0, err
}

func (set *SqliteSet) Insert(number uint64) error {
	query := fmt.Sprintf("INSERT OR IGNORE INTO `%v` VALUES(%v)", set.tableName, number)
	_, err := set.db.Exec(query)
	return err
}
