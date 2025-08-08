///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

///////////////////////////////////////////////////////////////////////////////////////
//                          part one · constants and globals                         //
///////////////////////////////////////////////////////////////////////////////////////

// 1.1: command-line flags
///////////////////////////

// MainFlags is the default command-line flag parser.
var MainFlags = flag.NewFlagSet("stnts", flag.ExitOnError)

// Defined command-line flags.
var (
	FlagAddr  = MainFlags.String("addr", "localhost:8000", "host:port address")
	FlagCache = MainFlags.String("cache", "./cache", "cache folder")
	FlagConf  = MainFlags.String("conf", "", "configuration file")
	FlagDebug = MainFlags.Bool("debug", false, "enable debug mode")
)

// 1.2: main control variables
///////////////////////////////

// MainLogs is the default system logger.
var MainLogs *log.Logger

// MainMux is the default server handler mux.
var MainMux *http.ServeMux

// MainServer is the default system server.
var MainServer *http.Server

///////////////////////////////////////////////////////////////////////////////////////
//                             part ??? · main functions                             //
///////////////////////////////////////////////////////////////////////////////////////

// init initialises the stnts program before start.
func init() {
	// Parse command-line flags.
	MainFlags.Parse(os.Args[1:])

	// Initialise MainLogs and print configuration.
	MainLogs = log.New(os.Stdout, "", log.LstdFlags)
	MainLogs.Printf("-addr=%q", *FlagAddr)
	MainLogs.Printf("-cache=%q", *FlagCache)
	MainLogs.Printf("-conf=%q", *FlagConf)
	MainLogs.Printf("-debug is %t", *FlagDebug)

	// Initialise MainMux.
	MainMux = http.NewServeMux()

	// Initialise MainServer.
	MainServer = &http.Server{Addr: *FlagAddr, Handler: MainMux}
}

// main runs the main stnts program.
func main() {
	MainLogs.Printf("starting server on %q", *FlagAddr)
	if err := MainServer.ListenAndServe(); err != nil {
		MainLogs.Fatal(err)
	}
}
