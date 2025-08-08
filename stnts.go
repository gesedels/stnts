///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

// IconCache is the global icon file cache map.
var IconCache map[string][]byte

// IconCacheLock is the global mutex for IconCache.
var IconCacheLock *sync.Mutex

// IconTypes is a list of valid icon filenames.
var IconTypes = []string{
	"apple-touch-icon.png",
	"apple-touch-icon-180x180.png",
	"favicon.png",
	"favicon.svg",
	"favicon.ico",
}

///////////////////////////////////////////////////////////////////////////////////////
//                       part two · file and download functions                      //
///////////////////////////////////////////////////////////////////////////////////////

// DownloadURL returns a URL's contents as a byteslice.
func DownloadURL(addr string) ([]byte, error) {
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get(addr)
	switch {
	case err != nil:
		return nil, fmt.Errorf("cannot download %q - %w", addr, err)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("cannot download %q - status %d", addr, resp.StatusCode)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// DownloadURLs returns the first valid URL's contents as a byteslice.
func DownloadURLs(host string, names ...string) ([]byte, error) {
	for _, name := range names {
		addr := fmt.Sprintf("https://%s/%s", host, name)
		bytes, err := DownloadURL(addr)
		if err == nil {
			return bytes, nil
		}
	}

	return nil, fmt.Errorf("cannot download from %q - no accessible URLs", host)
}

// ReadJSON unmarshals a JSON file into an object.
func ReadJSON(orig string, data any) error {
	bytes, err := os.ReadFile(orig)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
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

///////////////////////////////////////////////////////////////////////////////////////
//                            part ??? · handler functions                           //
///////////////////////////////////////////////////////////////////////////////////////

// GetIcon returns a new or cached remote icon file.
func GetIcon(w http.ResponseWriter, r *http.Request) {
	host := r.PathValue("host")
	if _, ok := IconCache[host]; !ok {
		bytes, err := DownloadURLs(host, IconTypes...)
		if err != nil {
			WriteError(w, http.StatusNotFound, "%q icon not found", host)
			return
		}

		IconCacheLock.Lock()
		IconCache[host] = bytes
		IconCacheLock.Unlock()
	}

	mime := http.DetectContentType(IconCache[host])
	w.Header().Set("Content-Type", mime)
	w.WriteHeader(http.StatusOK)
	w.Write(IconCache[host])
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
	Server = &http.Server{Addr: *FlagAddr, Handler: Mux}

	// Initialise content variables.
	IconCache = make(map[string][]byte)
	IconCacheLock = new(sync.Mutex)

	// Register handler routes.
	Mux.HandleFunc("GET /icon/{host...}", LoggingWare(GetIcon))
}

// main runs the main stnts program.
func main() {
	Log.Printf("starting server on %q", *FlagAddr)
	if err := Server.ListenAndServe(); err != nil {
		Log.Fatal(err)
	}
}
