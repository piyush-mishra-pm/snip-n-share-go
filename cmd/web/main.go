package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/piyush-mishra-pm/snip-n-share-go/internal/models"
)

type application struct {
	logInfo       *log.Logger
	logError      *log.Logger
	snips         *models.SnipModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port address at which server will listen. Avoid reserved ports 0:1023. Example ':8080'")
	dsn := flag.String("dsn", "web:4321@/snipnshare?parseTime=true", "MySQL data source name. Example: ")
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	logError := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		logError.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logError.Fatal(err)
	}

	app := &application{
		logInfo:       logInfo,
		logError:      logError,
		snips:         &models.SnipModel{DB: db},
		templateCache: templateCache,
	}

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
