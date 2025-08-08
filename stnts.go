///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	FlagLogs = FlagSet.String("logs", "", "file to write logs to")
	FlagWarm = FlagSet.Bool("warm", false, "warm icon cache before start")
)

// 1.2: global configuration variables
///////////////////////////////////////

// MainSite is the global Site configuration object.
var MainSite *Site

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

// DownloadIcon returns a URL's favicon as a byteslice.
func DownloadIcon(addr string) ([]byte, error) {
	uobj, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	icon := fmt.Sprintf("%s://%s/favicon.ico", uobj.Scheme, uobj.Hostname())
	return DownloadURL(icon)
}

///////////////////////////////////////////////////////////////////////////////////////
//                        part four · configuration data types                       //
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

// Link.Host returns the Link's host URL.
func (l *Link) Host() (string, error) {
	uobj, err := url.Parse(l.Addr)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s", uobj.Scheme, uobj.Hostname()), nil
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
