// License: AGPL-3.0-only
// (c) 2024 Dakota Walsh <kota@nilsu.org>
package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"git.sr.ht/~kota/lists/models"
	"git.sr.ht/~kota/lists/ui"
)

type application struct {
	infoLog   *log.Logger
	errLog    *log.Logger
	templates map[string]*template.Template

	lists *models.ListModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "lists.db", "Sqlite data source string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err, *dsn)
	}
	defer db.Close()

	templates, err := ui.Templates()
	if err != nil {
		errLog.Fatal(err)
	}

	app := &application{
		infoLog:   infoLog,
		errLog:    errLog,
		templates: templates,
		lists:     &models.ListModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Println("starting server on", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite3", dsn)
}
