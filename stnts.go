///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

///////////////////////////////////////////////////////////////////////////////////////
//                          part one · constants and globals                         //
///////////////////////////////////////////////////////////////////////////////////////

// 1.1: command-line flags
///////////////////////////

// FlagSet is the system command-line flag parser.
var FlagSet = flag.NewFlagSet("stnts", flag.ExitOnError)

var (
	FlagAddr = FlagSet.String("addr", "localhost:8000", "host:port to listen on")
	FlagConf = FlagSet.String("conf", "", "configuration file to use")
	FlagLogs = FlagSet.String("logs", "", "file to write logs to")
	FlagWarm = FlagSet.Bool("warm", false, "warm icon cache before start")
)

// 1.2: global content variables
/////////////////////////////////

// MainSite is the global Site configuration object.
var MainSite = new(Site)

// MainTemp is the global Template object.
var MainTemp = template.Must(template.ParseFiles("template.html"))

// 1.3: global cache variables
///////////////////////////////

// Cache is a global map of downloaded external icons.
var Cache = make(map[string][]byte)

// CacheMutex is the global mutex for writing to the Cache.
var CacheMutex = new(sync.Mutex)

///////////////////////////////////////////////////////////////////////////////////////
//                           part two · json file functions                          //
///////////////////////////////////////////////////////////////////////////////////////

// ReadJSON unmarshals a JSON file into an object.
func ReadJSON(orig string, data any) error {
	bytes, err := os.ReadFile(orig)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

// WriteJSON marshals an object into a JSON file.
func WriteJSON(dest string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dest, bytes, 0666)
}

///////////////////////////////////////////////////////////////////////////////////////
//                     part three · cache and download functions                     //
///////////////////////////////////////////////////////////////////////////////////////

// DownloadURL returns the contents of a URL as a byteslice.
func DownloadURL(addr string) ([]byte, error) {
	resp, err := http.Get(addr)
	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("download %s failed with %d", addr, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

///////////////////////////////////////////////////////////////////////////////////////
//                        part four · http response functions                        //
///////////////////////////////////////////////////////////////////////////////////////

// WriteFromCache writes a Cache value to a ResponseWriter.
func WriteFromCache(w http.ResponseWriter, code int, name, ctyp string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", ctyp)
	w.Write(Cache[name])
}

// WriteError writes a text/plain error to a ResponseWriter.
func WriteError(w http.ResponseWriter, code int, text string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "error %d: %s", code, text)
}

// WriteHTML writes a rendered Template to a ResponseWriter.
func WriteHTML(w http.ResponseWriter, code int, temp *template.Template, pipe any) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := temp.Execute(w, pipe); err != nil {
		log.Printf("template error - %s", err)
	}
}

///////////////////////////////////////////////////////////////////////////////////////
//                        part five · configuration data types                       //
///////////////////////////////////////////////////////////////////////////////////////

// 4.1: the Conf type
//////////////////////

// Conf is a single configuration map.
type Conf struct {
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Footer   string `json:"footer"`
	TimeZone string `json:"timezone"`
}

// 4.2: the Link type
//////////////////////

// Link is a single named website link.
type Link struct {
	Name string `json:"name"`
	From string `json:"from"`
	Addr string `json:"addr"`
}

// Link.Host returns the Link's hostname.
func (l *Link) Host() string {
	uobj, _ := url.Parse(l.Addr)
	return uobj.Hostname()
}

// Link.Root returns the Link's base URL.
func (l *Link) Root() string {
	uobj, _ := url.Parse(l.Addr)
	return fmt.Sprintf("%s://%s", uobj.Scheme, uobj.Hostname())
}

// 4.3: the List type
//////////////////////

// List is a single ordered list of Links.
type List struct {
	Name  string  `json:"name"`
	Links []*Link `json:"links"`
}

// 4.3: the Site type
//////////////////////

// Site is a complete container of configuration and content data.
type Site struct {
	Conf  *Conf   `json:"conf"`
	Icons []*Link `json:"icons"`
	Lists []*List `json:"lists"`
}

///////////////////////////////////////////////////////////////////////////////////////
//                            part six · handler functions                           //
///////////////////////////////////////////////////////////////////////////////////////

// GetIndex returns the index page.
func GetIndex(w http.ResponseWriter, r *http.Request) {
	WriteHTML(w, http.StatusOK, MainTemp, MainSite)
}

///////////////////////////////////////////////////////////////////////////////////////
//                         part seven · middleware functions                         //
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

// main runs the main stnts program.
func main() {
	// Parse command-line flags.
	FlagSet.Parse(os.Args[1:])

	// Configure system logging.
	if *FlagLogs != "" {
		file, err := os.Open(*FlagLogs)
		if err != nil {
			log.Fatalf("cannot open logfile %s - %s", *FlagLogs, err)
		}

		defer file.Close()
		log.Printf("setting logfile to %s", *FlagLogs)
		log.SetOutput(file)
	} else {
		log.Printf("setting logfile to STDOUT")
		log.SetOutput(os.Stdout)
	}

	// Parse configuration file.
	log.Printf("parsing configuration file %s", *FlagConf)
	if err := ReadJSON(*FlagConf, MainSite); err != nil {
		log.Fatalf("cannot read configuration file %s - %s", *FlagConf, err)
	}

	// Warm cache if -warm is true.
	if *FlagWarm {
		log.Printf("-warm is true, warming icon cache")
		for _, link := range MainSite.Icons {
			log.Printf("downloading icon for %s", link.Host())
			addr := fmt.Sprintf("%s/favicon.ico", link.Root())

			bytes, err := DownloadURL(addr)
			if err != nil {
				log.Printf("cannot download icon for %s - %s", addr, err)
			}

			CacheMutex.Lock()
			Cache[link.Host()] = bytes
			CacheMutex.Unlock()
		}
	}

	// Configure muxer.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", LoggingWare(GetIndex))

	// Configure and run server.
	srv := &http.Server{Addr: *FlagAddr, Handler: mux}
	log.Printf("now serving stnts on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed - %s", err)
	}
}
