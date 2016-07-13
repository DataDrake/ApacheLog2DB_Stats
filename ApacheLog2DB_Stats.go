package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB_Stats/date"
	"os"
)

func usage() {
	fmt.Println("USAGE: ApacheLog2DB_Stats [OPTION]... CMD DEST HTTP_SOURCE [HTTPS_SOURCE]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = func() { usage() }
	flag.Parse()

	args := flag.Args()

	if len(args) < 3 || len(args) > 4 {
		usage()
		os.Exit(1)
	}

	var m *date.DateModifier
	switch args[0] {
	case "second":
		m = *date.AddSecond
	case "hour":
		m = *date.AddHour
	case "day":
		m = *date.AddDay
	case "week":
		m = *date.AddWeek
	case "month":
		m = *date.AddMonth
	case "year":
		m = *date.AddYear
	default:
		usage()
		os.Exit(1)
	}

	dest_file, err := os.OpenFile(args[1], os.O_RDWR|os.O_CREATE|os.O_SYNC, 00644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open output file, reason: %s\n", err.Error())
		os.Exit(1)
	}
	dest := csv.NewWriter(dest_file)

	httpdb, err := global.OpenDatabase(args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open http DB, reason: %s\n", err.Error())
		os.Exit(1)
	}
	var httpsdb *sql.DB
	if len(args) == 4 {
		httpsdb, err = global.OpenDatabase(args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open https DB, reason: %s\n", err.Error())
			os.Exit(1)
		}
	}

	dest.Flush()
	dest_file.Close()

	os.Exit(0)
}
