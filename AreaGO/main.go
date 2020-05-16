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
	"expvar"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// Command-line flags.
var (
	httpAddr   = flag.String("http", ":6060", "Listen address")
	pollPeriod = flag.Duration("poll", 5*time.Second, "Poll period")
	version    = flag.String("version", "1.8", "Go version")
)

const baseChangeURL = "https://go.googlesource.com/go/+/"

type timeHandler struct {
	format string
}

func main() {

	flag.Parse()
	changeURL := fmt.Sprintf("%sgo%s", baseChangeURL, *version)

	mux := http.NewServeMux()
	th := &timeHandler{format: time.RFC1123}
	mux.Handle("/time", th)
	mux.Handle("/", NewServer(*version, changeURL, *pollPeriod))
	mux.Handle("/ping", http.HandlerFunc(ServeHTTP))

	//http.Handle("/", NewServer(*version, changeURL, *pollPeriod))
	//http.Handle("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(*httpAddr, mux))
}

// Exported variables for monitoring the server.
// These are exported via HTTP as a JSON object at /debug/vars.
var (
	hitCount       = expvar.NewInt("hitCount")
	pollCount      = expvar.NewInt("pollCount")
	pollError      = expvar.NewString("pollError")
	pollErrorCount = expvar.NewInt("pollErrorCount")
)

// Server implements the outyet server.
// It serves the user interface (it's an http.Handler)
// and polls the remote repository for changes.
type Server struct {
	version string
	url     string
	period  time.Duration

	mu  sync.RWMutex // protects the yes variable
	yes bool
}

// NewServer returns an initialized outyet server.
func NewServer(version, url string, period time.Duration) *Server {
	s := &Server{version: version, url: url, period: period}

	//goroutine
	//go s.poll()
	return s
}

// Hooks that may be overridden for integration tests.
var (
	pollSleep = time.Sleep
	pollDone  = func() {}
)

// poll polls the change URL for the specified period until the tag exists.
// Then it sets the Server's yes field true and exits.
func (s *Server) poll() {
	for !isTagged(s.url) {
		pollSleep(s.period)
	}
	s.mu.Lock()
	s.yes = true
	s.mu.Unlock()
	pollDone()
}

// isTagged makes an HTTP HEAD request to the given URL and reports whether it
// returned a 200 OK response.
func isTagged(url string) bool {
	pollCount.Add(1)
	r, err := http.Head(url)
	if err != nil {
		log.Print(err)
		pollError.Set(err.Error())
		pollErrorCount.Add(1)
		return false
	}
	return r.StatusCode == http.StatusOK
}

// EnableCORS autorise le cors sinon pas call extérieur
// /!\ trouver une manière plus propore de le faire
func EnableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hitCount.Add(1)
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Print(err)
	}
}

// ServeHTTP répond Pong à une deande /ping
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/* 	hitCount.Add(1)
	   	err := tmpl.Execute(w, nil)
	   	if err != nil {
	   		log.Print(err)
		   } */
	EnableCORS(&w)
	w.Write([]byte("Pong"))

}

//ServeHTTP requête /time
func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
