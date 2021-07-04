// Copyright (C) 2021 Pau Sanchez
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syncdbdocs/lib"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	var dbhost string
	var dbport uint
	var dbuser string
	var dbpass string
	var dbname string
	var dbtype string

	flag.StringVar(&dbhost, "h", "127.0.0.1", "Host you want to connect to")
	flag.UintVar(&dbport, "p", 0, "Port on given host you want to connect to")
	flag.StringVar(&dbuser, "u", "root", "Username credentials (password should be set via DBPASSWORD env var)")
	flag.StringVar(&dbname, "d", "", "Database name you want to connect to")
	flag.StringVar(&dbtype, "t", "auto", "Database type: auto | pg | mysql | mssql")

	dbpass = os.Getenv("DBPASSWORD")

	flag.Parse()

	if dbname == "" {
		fmt.Println("You should provide database name with -d flag")
		os.Exit(-1)
	}

	// TODO: guess type from port, or just try
	if dbport == 0 {
		dbport = 5432
	}

	conn, err := lib.DbConnect(dbtype, dbhost, dbport, dbuser, dbpass, dbname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("ERROR: Cannot connect to the database: ", strings.ReplaceAll(conn.GetConnectionString(), dbpass, "*****"))
		os.Exit(-2)
	}
	defer conn.Close()

	dbLayout, err := conn.GetLayout()
	if err != nil {
		fmt.Println("ERROR: cannot create layout. ", err)
		os.Exit(-3)
	}

	dbLayout.PrintMarkdown(os.Stdout, 80)
}
