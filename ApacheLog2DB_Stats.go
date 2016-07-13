package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB_Stats/date"
	"github.com/DataDrake/ApacheLog2DB_Stats/stats"
	global2 "kgcoe-git.rit.edu/btmeme/DWSim/global"
	"os"
	"time"
)

func usage() {
	fmt.Println("USAGE: ApacheLog2DB_Stats [OPTION]... CMD DEST START END HTTP_SOURCE [HTTPS_SOURCE]")
	flag.PrintDefaults()
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

	dest_file, err := os.OpenFile(args[1], os.O_RDWR|os.O_CREATE|os.O_SYNC, 00644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open output file, reason: %s\n", err.Error())
	}
	dest := csv.NewWriter(dest_file)

	start, err := time.Parse(global2.SQL_TIME_LAYOUT, args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Start time, should resemble: %s\n", global2.SQL_TIME_LAYOUT)
		goto DEST_CLEANUP
	}

	end, err := time.Parse(global2.SQL_TIME_LAYOUT, args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Start time, should resemble: %s\n", global2.SQL_TIME_LAYOUT)
		goto DEST_CLEANUP
	}

	httpdb, err := global.OpenDatabase(args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open http DB, reason: %s\n", err.Error())
		goto DEST_CLEANUP
	}

	var httpsdb *sql.DB
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

	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
