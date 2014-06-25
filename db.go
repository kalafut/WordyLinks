package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pq "github.com/lib/pq"
)

const UNIQUE_VIOLATION = "23505"

type PostgresDB struct {
	db         *sql.DB
	expiration time.Duration
}

func NewPostgresDB(expiration time.Duration) PostgresDB {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", os.Getenv("WL_DBNAME"), os.Getenv("WL_DBUSER"), os.Getenv("WL_DBPASSWORD"))
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return PostgresDB{db: db, expiration: expiration}
}

func (db PostgresDB) Load(keywords string) (string, error) {
	var url string

	err := db.db.QueryRow("UPDATE urls SET (expiration) = ($2) WHERE words = $1 RETURNING url", keywords, db.expireTime().Unix()).Scan(&url)
	switch {
	case err == sql.ErrNoRows:
		err = KeywordsNotFound

	case err != nil:
		log.Fatal(err)
	}

	return url, err
}

// Save entry if keywords don't exist
// Update entry if they've expired
func (db PostgresDB) Save(keywords string, url string) error {
	_, err := db.db.Query("INSERT INTO urls(words, url, expiration) VALUES($1, $2, $3)", keywords, url, db.expireTime().Unix())

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == UNIQUE_VIOLATION {
		result, _ := db.db.Exec("UPDATE urls SET (url, expiration) = ($3, $2) WHERE words = $1 AND expiration <= $4", keywords, db.expireTime().Unix(), url, time.Now().Unix())

		if ra, _ := result.RowsAffected(); ra == 0 {
			err = KeywordsInUse
		} else {
			err = nil
		}
	}

	return err
}

func (db PostgresDB) expireTime() time.Time {
	return time.Now().Add(db.expiration)
}
