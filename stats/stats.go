package stats

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB_Stats/date"
	"os"
	"strconv"
	"time"
)

func GenerateStat(db *sql.DB, last, current time.Time) int {
	temp := 0
	var err error
	row, err := db.Query("SELECT count(*) FROM txns WHERE occured>=? AND occured<?", last, current)
	if err != nil {
		goto ERROR
	}
	if row.Next() {
		err = errors.New("Failed to get count between '" + last.String() + "' and '" + current.String() + "'")
		goto ERROR
	} else {
		err = row.Scan(temp)
		if err != nil {
			goto ERROR
		}
		return temp
	}
ERROR:
	fmt.Fprintf(os.Stderr, "Failed to get http stats between '%t' and '%t', reason: %s", last, current, err.Error())
	return 0
}

func GenerateStats(httpdb, httpsdb *sql.DB, m date.DateModifier, start, End time.Time, dest *csv.Writer) {
	results := make([]string, 2)
	current := start
	last := start
	m(current)
	results[0] = "Time"
	results[1] = "Requests"
	err := dest.Write(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write header, reason: %s\n", err.Error())
	}
	for current.Before(End) {
		count := GenerateStat(httpdb, last, current)
		if httpsdb != nil {
			count += GenerateStat(httpsdb, last, current)
		}
		results[0] = current.String()
		results[1] = strconv.FormatInt(int64(count), 10)
		err := dest.Write(results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write entry, reason: %s\n", err.Error())
		}
		m(current)
		m(last)
	}
}
