package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rozbeikonazar/gosnippetbox/internal/models"
)

type application struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func goDotEnvVariable(key string) string {
	_ = godotenv.Load(".env")

	return os.Getenv(key)
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	// setting up environment variables
	username := goDotEnvVariable("DB_USERNAME")
	password := goDotEnvVariable("DB_PASSWORD")
	dbname := goDotEnvVariable("DB_NAME")

	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, dbname), "MySQL data source name")

	flag.Parse()

	// Setting up logging
	infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errLog:   errLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
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
