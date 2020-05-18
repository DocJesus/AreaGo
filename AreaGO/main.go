/*
Copyright 2014 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// outyet is a web server that announces whether or not a particular Go version
// has been tagged.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Command-line flags.
var (
	httpAddr = flag.String("http", ":6060", "Listen address")
)

const baseChangeURL = "https://go.googlesource.com/go/+/"

type timeHandler struct {
	format string
}

func main() {

	flag.Parse()

	mux := http.NewServeMux()
	//th := &timeHandler{format: time.RFC1123}

	//mux.Handle("/time", th)
	mux.Handle("/", http.HandlerFunc(ServeHTTP))
	//mux.Handle("/ping", http.HandlerFunc(ServeHTTP))

	mux.Handle("/ping", http.HandlerFunc(RequestHandler))

	//http.Handle("/", NewServer(*version, changeURL, *pollPeriod))
	//http.Handle("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(*httpAddr, mux))
}

// isTagged makes an HTTP HEAD request to the given URL and reports whether it
// returned a 200 OK response.
func isTagged(url string) bool {
	r, err := http.Head(url)
	if err != nil {
		log.Print(err)
		return false
	}
	return r.StatusCode == http.StatusOK
}

// EnableCORS autorise le cors sinon pas call extérieur
// /!\ trouver une manière plus propore de le faire
func EnableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// ServeHTTP implements the HTTP user interface.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Print(err)
	}
}

// PingPongHandler répond Pong à une deande /ping
func PingPongHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)
	w.Write([]byte("Pong"))
}

// RequestHandler gère toute les requêtes GET et Post pour le moment /ping
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)

	if r.URL.Path != "/ping" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		//http.ServeFile(w, r, "form.html")
		w.Write([]byte("Stonks"))
		fmt.Fprintf(w, "stronks")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	case "OPTIONS":
		fmt.Fprintf(w, "Options command")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//ServeHTTP requête /time
func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)

	tm := time.Now().Format(th.format)
	data := struct {
		TIME string
	}{
		tm,
	}
	err := tmpl2.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

// tmpl is the HTML template that drives the user interface.
var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Welcome to my humble server</h2>
	<h1>
		I will try to send "Pong" if someone send me "Ping"
	</h1>
</center></body></html>
`))

var tmpl2 = template.Must(template.New("tmpl2").Parse(`
<!DOCTYPE html><html><body><center>
	<h2>Welcome to my humble server</h2>
	<h1>
		Time Today is {{.TIME}}
	</h1>
</center></body></html>
`))
