package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	logInfo  *log.Logger
	logError *log.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP port address at which server will listen. Avoid reserved ports 0:1023. Example ':8080'")
	flag.Parse()

	app := &application{
		logInfo:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC),
		logError: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snipView)
	mux.HandleFunc("/snippet/create", app.snipCreate)

	app.logInfo.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logError,
		Handler:  mux,
	}
	err := srv.ListenAndServe()
	app.logError.Fatal(err)
}
