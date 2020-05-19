/*

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
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mux"
)

// Command-line flags.
var httpAddr = flag.String("http", ":6060", "Listen address")

const baseChangeURL = "https://go.googlesource.com/go/+/"

type Action struct {
	id          int
	name        string
	slug        string
	description string
}

var actions = []Action{
	Action{1, "Facebook", "facebook", "tu te co à facebook"},
	Action{2, "Gmail", "gmail", "tu te co à gmail"},
	Action{3, "Outlook", "outlook", "tu te co à tes mails"},
}

// NotImplemented will simply return the message "Not Implemented"
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

// BasicPage implemente la page de base du serveur
var BasicPage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)

	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Print(err)
	}
})

// BasicPing implemente la réponse pong quand on appel /ping
var BasicPing = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)
	w.Write([]byte("Pong"))
})

// List renvoie les actions du serveur
var List = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(actions)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

// FeedBack gère les requêtes post vers /actions/{slug}/feedback
var FeedBack = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Fprintf(w, "Post from website! Param = ")

	var currentAct Action
	toto := mux.Vars(r)
	slug := toto["slug"]

	for _, p := range actions {
		if p.slug == slug {
			currentAct = p
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if currentAct.slug != "" {
		payload, _ := json.Marshal(currentAct)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}
})

func main() {

	flag.Parse()

	//mux := http.NewServeMux()

	r := mux.NewRouter()

	r.Handle("/status", NotImplemented).Methods("GET")

	//interface de base
	r.Handle("/", BasicPage)

	//envoie les actions du serveur
	r.Handle("/actions", List).Methods("GET")

	//modifie les actions du serveur
	r.Handle("/actions/{slug}/feedback", FeedBack).Methods("POST")

	//répond par un simple pong à un call à "ping"
	//r.Handle("/ping", BasicPing)

	//répond aux requet post et get sur /ping
	//r.Handle("/ping", http.HandlerFunc(RequestHandler))

	log.Fatal(http.ListenAndServe(*httpAddr, r))
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
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

// RequestHandler gère toute les requêtes GET et Post pour le moment /ping
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(&w)

	switch r.Method {
	case "GET":
		payload, _ := json.Marshal(actions)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(payload))
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		r.ParseForm()
		fmt.Fprintf(w, "Post from website! Param = %s\n", r.Form.Get("user"))
		fmt.Fprintf(w, "Post from website! Param = %s\n", r.Form.Get("passwd"))
	case "OPTIONS":
		fmt.Fprintf(w, "Options command")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
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
