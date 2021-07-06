// Copyright (C) 2021 Pau Sanchez
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"syncdbdocs/lib"
)

func main() {
	var dbhost string
	var dbport uint
	var dbuser string
	var dbpass string
	var dbname string
	var dbtype string

	var inputFile string
	var outputFile string
	var format string
	var lineLength int
	// TODO: var syncToDb bool

	flag.StringVar(&dbhost, "h", "127.0.0.1", "Host you want to connect to")
	flag.UintVar(&dbport, "p", 0, "Port on given host you want to connect to")
	flag.StringVar(&dbuser, "u", "", "Username credentials (password should be set via DBPASSWORD env var)")
	flag.StringVar(&dbname, "d", "", "Database name you want to connect to")
	flag.StringVar(&dbtype, "t", "auto", "Database type: auto | pg | mysql | mariadb | mssql | sqlite")
	flag.StringVar(&inputFile, "i", "", "Use given input file to extend on")
	flag.StringVar(&outputFile, "o", "", "Output file to generate")
	flag.StringVar(&format, "format", "", "Output format (text | markdown | dbml)")
	flag.IntVar(&lineLength, "line-length", 80, "Set line length for the text/markdown representation")
	// TODO: flag.BoolVar(&syncToDb, "sync-to-db", false, "Update database comments from markdown")

	// dbhostEnv := os.Getenv("DB_HOST")
	// dbportEnv := os.Getenv("DB_PORT")
	dbuserEnv := os.Getenv("DB_USER")
	dbpass = os.Getenv("DB_PASSWORD")

	flag.Parse()

	if dbname == "" {
		fmt.Println("You should provide database name with -d flag")
		os.Exit(-1)
	}

	// TODO: guess type from port, or just try
	if dbport == 0 {
		dbport = 5432
	}

	if len(dbuser) == 0 && len(dbuserEnv) == 0 {
		dbuser = "root"
	} else if len(dbuser) == 0 {
		dbuser = dbuserEnv
	}

	conn, err := lib.DbConnect(dbtype, dbhost, dbport, dbuser, dbpass, dbname)
	if err != nil {
		connString := fmt.Sprintf("%s://%s:*****@%s:%d/%s", dbtype, dbuser, dbhost, dbport, dbname)
		if conn != nil {
			connString = strings.ReplaceAll(conn.GetConnectionString(), dbpass, "*****")
		}

		fmt.Println(err)
		fmt.Println("ERROR: Cannot connect to the database: ", connString)
		os.Exit(-2)
	}
	defer conn.Close()

	dbLayout, err := conn.GetLayout()
	if err != nil {
		fmt.Println("ERROR: cannot create layout. ", err)
		os.Exit(-3)
	}

	// ensure all new items are always appended in order
	dbLayout.Sort()

	// if output file exists and no input is specified, let's set input as
	// the output so it will be rewritten but keeping the same order
	if inputFile == "" && outputFile != "" {
		if f, _ := os.Open(outputFile); f != nil {
			inputFile = outputFile
			f.Close()
		}
	}

	if inputFile != "" {
		fileLayout, err := lib.NewDbLayoutFromParsedFile(inputFile)
		if err != nil {
			fmt.Printf("ERROR: cannot read input file %s: %s\n", inputFile, err)
			os.Exit(-4)
		}

		fileLayout.MergeFrom(dbLayout)
		dbLayout = fileLayout
	}

	var outStream io.Writer = os.Stdout
	if outputFile != "" {
		ofile, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("ERROR: cannot open output file %s: %s\n", outputFile, err)
			os.Exit(-5)
		}
		defer ofile.Close()

		outStream = ofile
	}

	switch strings.ToLower(format) {
	case "md", "markdown":
		dbLayout.PrintMarkdown(outStream, lineLength)
	case "txt", "plain", "text":
		dbLayout.PrintText(outStream, lineLength)
	default:
		dbLayout.PrintText(outStream, lineLength)
	}
}
