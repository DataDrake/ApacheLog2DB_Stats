package stats

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB_Stats/date"
	"kgcoe-git.rit.edu/btmeme/DWSim/global"
	"os"
	"strconv"
	"time"
)

func GenerateStat(db *sql.DB, last, current time.Time) int {
	temp := int(0)
	var err error
	row := db.QueryRow("SELECT count(*) as count FROM txns WHERE occured >= ? AND occured < ?", last, current)
	if row == nil {
		err = errors.New("Query failed")
	}
	err = row.Scan(&temp)
	if err != nil {
		goto ERROR
	}
	return temp
ERROR:
	fmt.Fprintf(os.Stderr, "Failed to get http stats between '%t' and '%t', reason: %s\n", last.Format(global.SQL_TIME_LAYOUT), current.Format(global.SQL_TIME_LAYOUT), err.Error())
	return 0
}

func GenerateStats(httpdb, httpsdb *sql.DB, m date.DateModifier, start, End time.Time, dest *csv.Writer) {
	results := make([]string, 2)
	current := start
	last := start
	current = m(current)
	results[0] = "Time"
	results[1] = "Requests"
	err := dest.Write(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write header, reason: %s\n", err.Error())
	}
	for last.Before(End) {
		count := GenerateStat(httpdb, last, current)
		if httpsdb != nil {
			count += GenerateStat(httpsdb, last, current)
		}
		results[0] = last.Format(global.SQL_TIME_LAYOUT)
		results[1] = strconv.FormatInt(int64(count), 10)
		err := dest.Write(results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write entry, reason: %s\n", err.Error())
		}
		dest.Flush()
		current = m(current)
		last = m(last)
	}
}
