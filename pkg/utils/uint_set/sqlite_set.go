package uint_set

import (
	"database/sql"
	"fmt"
	"github.com/AntonBorzenko/RestrictedPassportService/config"
	"github.com/AntonBorzenko/RestrictedPassportService/utils"
	e "github.com/AntonBorzenko/RestrictedPassportService/utils/errors"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"time"
)

type SqliteSet struct {
	tableName string
	db *sql.DB
}

func NewSqliteSet(filename string, createDb bool, createIndex bool) *SqliteSet {
	db, err := sql.Open("sqlite3", filename)
	e.Check(err)
	e.CheckDbResult(db.Exec("PRAGMA journal_mode=off"))
	e.CheckDbResult(db.Exec("PRAGMA page_size=65536"))

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
	count := 0
	step := 0
	start := time.Now()

	for batch := range utils.GetBatchGenerator(numbers, config.Cfg.DBBatchSize) {
		batchStrings := make([]string, len(batch))
		for i, value := range batch {
			batchStrings[i] = strconv.FormatUint(value, 10)
		}
		query :=
			"INSERT OR IGNORE INTO `" + set.tableName +
			"` VALUES (" + strings.Join(batchStrings, "),(") + ")"
		_, err := set.db.Exec(query)
		if !ignoreErrors && err != nil {
			return err
		}

		if config.Cfg.DBUpdateVerbose && step % 50 == 0 {
			tpr := float64(time.Now().Sub(start).Seconds()) / float64(count) * 1_000_000
			fmt.Printf("\r%s", strings.Repeat(" ", 50))
			fmt.Printf("\rInserted %v values, %.1fs per 1M insertions", count, tpr)
		}

		count += len(batch)
		step += 1
	}
	if config.Cfg.DBUpdateVerbose {
		fmt.Println()
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
