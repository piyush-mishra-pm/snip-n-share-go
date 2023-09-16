package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/piyush-mishra-pm/snip-n-share-go/internal/models"
)

type application struct {
	logInfo        *log.Logger
	logError       *log.Logger
	snips          *models.SnipModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logInfo:        logInfo,
		logError:       logError,
		snips:          &models.SnipModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	app.logInfo.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logError,
		Handler:  app.routes(),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
