///////////////////////////////////////////////////////////////////////////////////////
//                       stephen's tiny new tab server · v0.0.0                      //
///////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"flag"
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

// 1.2: global cache variables
///////////////////////////////

// Cache is a global map of downloaded external icons.
var Cache = make(map[string][]byte)

// CacheMutex is the global mutex for writing to the Cache.
var CacheMutex = new(sync.Mutex)
