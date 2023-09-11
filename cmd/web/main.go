package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/piyush-mishra-pm/snip-n-share-go/internal/models"
)

type application struct {
	logInfo  *log.Logger
	logError *log.Logger
	snips    *models.SnipModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port address at which server will listen. Avoid reserved ports 0:1023. Example ':8080'")
	dsn := flag.String("dsn", "web:4321@/snipnshare?parseTime=true", "MySQL data source name. Example: ")
	flag.Parse()

	app := &application{
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		logError: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
	}

	db, err := openDB(*dsn)
	if err != nil {
		app.logError.Fatal(err)
	}
	defer db.Close()
	app.snips = &models.SnipModel{DB: db}

	app.logInfo.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logError,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
	app.logError.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
