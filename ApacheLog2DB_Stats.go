package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB_Stats/date"
	"github.com/DataDrake/ApacheLog2DB_Stats/stats"
	_ "github.com/mattn/go-sqlite3"
	global2 "kgcoe-git.rit.edu/btmeme/DWSim/global"
	"os"
	"time"
)

func usage() {
	fmt.Println("USAGE: ApacheLog2DB_Stats [OPTION]... STEP DEST START END HTTP_SOURCE [HTTPS_SOURCE]\n")
	fmt.Println("Generate request statistics for an ApacheLog2DB import\n")
	fmt.Println("\tSTEPS:\tSECOND, MINUTE, HOUR, DAY, WEEK, MONTH, YEAR\n")
	fmt.Println("\tDEST:\tpath to a CSV file for output\n")
	fmt.Println("\tSTART:\tstart time in SQL time format\n")
	fmt.Println("\tEND:\tend time in SQL time format\n")
	fmt.Println("\tHTTP_SOURCE:\tlocation of ApacheLog2DB database\n")
	fmt.Println("\tHTTPS_SOURCE:\tseparate location of ApacheLog2DB database for HTTPS\n")
}

func main() {
	flag.Usage = func() { usage() }
	flag.Parse()

	args := flag.Args()

	if len(args) < 5 || len(args) > 6 {
		usage()
		os.Exit(1)
	}

	var m date.DateModifier
	switch args[0] {
	case "second":
		m = date.AddSecond
	case "minute":
		m = date.AddMinute
	case "hour":
		m = date.AddHour
	case "day":
		m = date.AddDay
	case "week":
		m = date.AddWeek
	case "month":
		m = date.AddMonth
	case "year":
		m = date.AddYear
	default:
		usage()
		os.Exit(1)
	}

	var httpdb *sql.DB
	var httpsdb *sql.DB
	var start time.Time
	var end time.Time
	var dest *csv.Writer

	dest_file, err := os.Create(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create output file, reason: %s\n", err.Error())
		goto END
	}
	dest = csv.NewWriter(dest_file)

	start, err = time.Parse(global2.SQL_TIME_LAYOUT, args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Start time, should resemble: %s\n", global2.SQL_TIME_LAYOUT)
		goto DEST_CLEANUP
	}

	end, err = time.Parse(global2.SQL_TIME_LAYOUT, args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Start time, should resemble: %s\n", global2.SQL_TIME_LAYOUT)
		goto DEST_CLEANUP
	}

	httpdb, err = global.OpenDatabase(args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open http DB, reason: %s\n", err.Error())
		goto DEST_CLEANUP
	}

	if len(args) == 6 {
		httpsdb, err = global.OpenDatabase(args[5])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open https DB, reason: %s\n", err.Error())
			goto HTTPDB_CLEANUP
		}
	}

	stats.GenerateStats(httpdb, httpsdb, m, start, end, dest)

	if httpsdb != nil {
		httpsdb.Close()
	}
HTTPDB_CLEANUP:
	httpdb.Close()
DEST_CLEANUP:
	dest.Flush()
	dest_file.Close()
END:
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
