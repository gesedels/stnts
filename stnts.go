///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gesedels/stnts/stnts/items/site"
)

///////////////////////////////////////////////////////////////////////////////////////
//                          part one · constants and globals                         //
///////////////////////////////////////////////////////////////////////////////////////

// 1.1: command-line flags
///////////////////////////

// Flags is the default command-line flag parser.
var Flags = flag.NewFlagSet("stnts", flag.ExitOnError)

// Defined command-line flags.
var (
	FlagAddr  = Flags.String("addr", "localhost:8000", "host:port address")
	FlagConf  = Flags.String("conf", "", "configuration file")
	FlagDebug = Flags.Bool("debug", false, "enable debug mode")
)

// 1.2: main control variables
///////////////////////////////

// Log is the global system logger.
var Log *log.Logger

// Mux is the global server handler mux.
var Mux *http.ServeMux

// Server is the global system server.
var Server *http.Server

// 1.3: main storage variables
///////////////////////////////

// MainFS is the global embedded asset filesystem.
//
//go:embed files/**/*
var MainFS embed.FS

// Site is the global site configuration object.
var MainSite *site.Site

// TemplateCache is a live cache of parsed templates.
var TemplateCache = make(map[string]*template.Template)

// TemplateCacheLock is a write-locking mutex for Cache.
var TemplateCacheLock = new(sync.Mutex)

///////////////////////////////////////////////////////////////////////////////////////
//                       part two · file and download functions                      //
///////////////////////////////////////////////////////////////////////////////////////

// ReadJSON unmarshals a JSON file into an object.
func ReadJSON(orig string, data any) error {
	bytes, err := os.ReadFile(orig)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

///////////////////////////////////////////////////////////////////////////////////////
//                      part ??? · template rendering functions                      //
///////////////////////////////////////////////////////////////////////////////////////

// Parse returns a new or cached Template from FS.
func Parse(names ...string) (*template.Template, error) {
	name := strings.Join(names, "|")

	if _, ok := TemplateCache[name]; !ok {
		tobj, err := template.ParseFS(MainFS, names...)
		if err != nil {
			return nil, fmt.Errorf("cannot parse template - %w", err)
		}

		TemplateCacheLock.Lock()
		TemplateCache[name] = tobj
		TemplateCacheLock.Unlock()
	}

	return TemplateCache[name], nil
}

// Render returns a rendered Template as a byteslice.
func Render(tobj *template.Template, pipe any) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := tobj.Execute(buff, pipe); err != nil {
		return nil, fmt.Errorf("cannot render template - %w", err)
	}

	return buff.Bytes(), nil
}

///////////////////////////////////////////////////////////////////////////////////////
//                         part ??? · http response functions                        //
///////////////////////////////////////////////////////////////////////////////////////

// WriteError writes a formatted text/plain error message to a ResponseWriter.
func WriteError(w http.ResponseWriter, code int, text string, elems ...any) {
	text = fmt.Sprintf("error %d: %s", code, text)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, text, elems...)
}

// WriteTemplate writes a rendered HTML Template to a ResponseWriter.
func WriteTemplate(w http.ResponseWriter, code int, tobj *template.Template, pipe any) {
	buff := new(bytes.Buffer)
	if err := tobj.Execute(buff, pipe); err != nil {
		WriteError(w, http.StatusInternalServerError, "template error")
		Log.Printf("error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	w.Write(buff.Bytes())
}

///////////////////////////////////////////////////////////////////////////////////////
//                            part ??? · handler functions                           //
///////////////////////////////////////////////////////////////////////////////////////

// GetIndex returns the index page template.
func GetIndex(w http.ResponseWriter, r *http.Request) {
	tobj, err := Parse("files/html/_base.html", "files/html/index.html")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "template error")
		Log.Printf("error: %s", err)
		return
	}

	WriteTemplate(w, http.StatusOK, tobj, MainSite)
}

// GetCSS returns an embedded CSS file.
func GetCSS(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	path := filepath.Join("files/css", name)
	bytes, err := MainFS.ReadFile(path)
	if err != nil {
		WriteError(w, http.StatusNotFound, "css file %q not found", name)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

///////////////////////////////////////////////////////////////////////////////////////
//                          part ??? · middleware functions                          //
///////////////////////////////////////////////////////////////////////////////////////

// LoggingWare wraps a HandlerFunc in middleware to log every incoming request.
func LoggingWare(hfun http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		hfun(w, r)
	}
}

///////////////////////////////////////////////////////////////////////////////////////
//                             part ??? · main functions                             //
///////////////////////////////////////////////////////////////////////////////////////

// init initialises the stnts program before start.
func init() {
	// Parse command-line flags.
	Flags.Parse(os.Args[1:])

	// Initialise control variables.
	Log = log.New(os.Stdout, "", log.LstdFlags)
	Mux = http.NewServeMux()
	MainSite = new(site.Site)
	Server = &http.Server{Addr: *FlagAddr, Handler: Mux}

	// Register handler routes.
	Mux.HandleFunc("GET /", LoggingWare(GetIndex))
	Mux.HandleFunc("GET /css/{name...}", LoggingWare(GetCSS))

	// Parse site configuration.
	if err := ReadJSON(*FlagConf, MainSite); err != nil {
		Log.Fatalf("error: %s", err)
	}
}

// main runs the main stnts program.
func main() {
	Log.Printf("starting server on %q", *FlagAddr)
	if err := Server.ListenAndServe(); err != nil {
		Log.Fatal(err)
	}
}
