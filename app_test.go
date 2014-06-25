package main

import (
	"database/sql"
	"testing"
)

type testCase struct {
	input  string
	output string
	err    error
}

func TestFormatURL(t *testing.T) {
	v := []testCase{
		{"cnn.com", "http://cnn.com", nil},
		{"http://cnn.com", "http://cnn.com", nil},
		{"   http://cnn.com    ", "http://cnn.com", nil},
		{"   ftp://cnn.com    ", "ftp://cnn.com", nil},
		{"http://c%nn.com", "http://c%nn.com", invalidURLError},
		{"", "", emptyURLError},
		{"  ", "", emptyURLError},
	}

	checkTable(v, formatURL, t)
}

func TestCleanKeywords(t *testing.T) {
	v := []testCase{
		{"a b c", "a b c", nil},
		{"alpha BravO chArlie", "alpha bravo charlie", nil},
		{"  a    b       c ", "a b c", nil},
		{"  汉字  日本語  ", "汉字 日本語", nil},
		{"", "", missingKeywordsError},
		{"    ", "", missingKeywordsError},
	}

	checkTable(v, cleanKeywords, t)
}

func checkTable(table []testCase, testFn func(string) (string, error), t *testing.T) {
	for _, test := range table {
		if output, err := testFn(test.input); output != test.output || err != test.err {
			t.Error("Expected:", test.output, test.err, "Actual:", output, err)
		}
	}
}

func prepareTestDB() {
	db, _ := sql.Open("postgres", "user=kalafut dbname=wordylinks password=homebound sslmode=disable")
	db.Query("DROP DATABASE IF EXISTS testdb")
	db.Query("CREATE DATABASE testdb")
	db.Query(
		`CREATE TABLE urls
(
  id serial NOT NULL,
  words character varying(100) NOT NULL,
  url character varying(100) NOT NULL,
  expiration integer NOT NULL,
  CONSTRAINT urls_pkey PRIMARY KEY (id)
)`)
}
