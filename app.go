package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	urllib "net/url"
	"regexp"
	"strings"
	"time"
)

const (
	SAVE        = iota
	saveSuccess = iota
	LOAD        = iota
)
const COOKIE_NAME = "post"
const expirationDuration = 30 * 24 * time.Hour

var db PostgresDB
var numRequests int
var templates = template.Must(template.ParseGlob("templates/*.html"))
var indexCache []byte
var debug bool

type TmplData struct {
	SaveSuccess                                      bool
	LinkValue, LinkErrorClass, LinkErrorMsg          string
	KeywordValue, KeywordErrorClass, KeywordErrorMsg string
	SaveClass, LoadClass                             string
}

type CookieData struct {
	Operation    int
	LinkError    string
	KeywordError string
	Link         string
	Keywords     string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	http.HandleFunc("/load", load)
	http.HandleFunc("/save", save)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	panic(http.ListenAndServe(":3000", nil))
}

func init() {
	db = NewPostgresDB(expirationDuration)
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
}

func about(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "about.html", &TmplData{})
}

func load(w http.ResponseWriter, r *http.Request) {
	words := r.FormValue("words")
	cleanedWords, err := cleanKeywords(words)

	if err != nil {
		cookie := CookieData{Operation: LOAD, KeywordError: err.Error(), Keywords: words}
		cookie.setCookie(w)
		http.Redirect(w, r, "/", 302)
	} else {
		url, err := db.Load(cleanedWords)

		switch {
		case err == KeywordsNotFound:
			cookie := CookieData{Operation: LOAD, KeywordError: err.Error(), Keywords: words}
			cookie.setCookie(w)
			http.Redirect(w, r, "/", 302)

		case err != nil:
			log.Print(err)

		default:
			http.Redirect(w, r, url, 302)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := TmplData{
		SaveClass: "btn-sel",
		LoadClass: "btn-nonsel",
	}

	if debug {
		templates = template.Must(template.ParseGlob("templates/*.html"))
	}

	cookie, err := getAndClearCookie(w, r)

	if err == nil {
		if cookie.KeywordError != "" {
			tmpl.KeywordErrorClass = "has-error"
			tmpl.KeywordErrorMsg = cookie.KeywordError
		}

		if cookie.LinkError != "" {
			tmpl.LinkErrorClass = "has-error"
			tmpl.LinkErrorMsg = cookie.LinkError
		}

		if cookie.Operation == saveSuccess {
			tmpl.SaveSuccess = true
		}

		if cookie.Operation == LOAD {
			tmpl.SaveClass = "btn-nonsel"
			tmpl.LoadClass = "btn-sel"
			tmpl.KeywordValue = cookie.Keywords
		} else if cookie.Operation == SAVE {
			tmpl.LinkValue = cookie.Link
			tmpl.KeywordValue = cookie.Keywords
		}
		templates.ExecuteTemplate(w, "index.html", &tmpl)
	} else {
		if indexCache == nil || debug {
			var s bytes.Buffer

			templates.ExecuteTemplate(&s, "index.html", &tmpl)
			indexCache = s.Bytes()
		}
		w.Write(indexCache)
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	valid := true
	url, URLerr := formatURL(r.PostFormValue("url"))
	words, keywordErr := cleanKeywords(r.PostFormValue("words"))
	cookie := CookieData{Operation: SAVE, Link: r.PostFormValue("url")}

	if keywordErr != nil {
		cookie.KeywordError = keywordErr.Error()
		valid = false
	}

	if URLerr != nil {
		cookie.LinkError = URLerr.Error()
		valid = false
	}

	if valid {
		err := db.Save(words, url)

		switch err {
		case nil:
			cookie.Operation = saveSuccess

		case KeywordsInUse:
			cookie.KeywordError = err.Error()

		default:
			log.Print(err)
		}
	}

	cookie.setCookie(w)

	http.Redirect(w, r, "/", 303)
}

func cleanKeywords(keywords string) (string, error) {
	var err error
	cleanString := strings.Trim(strings.ToLower(keywords), " ")

	re := regexp.MustCompile(" {2,}")
	cleanString = re.ReplaceAllString(cleanString, " ")

	if len(cleanString) == 0 {
		err = missingKeywordsError
	}

	return cleanString, err
}

func formatURL(rawURL string) (string, error) {
	var err error
	url := strings.Trim(rawURL, " ")

	_, err = urllib.Parse(url)

	if err != nil {
		err = invalidURLError
	} else if len(url) == 0 {
		err = emptyURLError
	} else {
		re := regexp.MustCompile("\\w+://")
		if re.FindStringIndex(url) == nil {
			url = "http://" + url
		}
	}

	return url, err
}

func (cookie CookieData) setCookie(w http.ResponseWriter) {
	jsonCookie, _ := json.Marshal(cookie)
	http.SetCookie(w, &http.Cookie{Name: COOKIE_NAME, Value: base64.URLEncoding.EncodeToString(jsonCookie)})
}

func getAndClearCookie(w http.ResponseWriter, r *http.Request) (cookie *CookieData, err error) {
	cookie = new(CookieData)
	rawCookie, err := r.Cookie(COOKIE_NAME)

	if err == nil {
		// clear cookie
		http.SetCookie(w, &http.Cookie{Name: COOKIE_NAME, MaxAge: -1})

		decodedCookie, _ := base64.URLEncoding.DecodeString(rawCookie.Value)
		err = json.Unmarshal(decodedCookie, cookie)
	}

	return
}
