package main

import "errors"

var emptyURLError = errors.New("URL missing")
var invalidURLError = errors.New("Invalid URL")
var missingKeywordsError = errors.New("Keywords missing")
var KeywordsNotFound = errors.New("Keywords not found")
var KeywordsInUse = errors.New("Keywords already used")
