package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gesedels/stnts/stnts/items/site"
	"github.com/gesedels/stnts/stnts/tools/file"
	"github.com/gesedels/stnts/stnts/tools/resp"
	"github.com/gesedels/stnts/stnts/tools/temp"
	"github.com/gesedels/stnts/stnts/tools/ware"
)

var (
	//go:embed files/**/*
	MainFS   embed.FS
	MainSite *site.Site
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	tobj, err := temp.Parse(MainFS, "files/html/_base.html", "files/html/index.html")
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, "template error")
		log.Printf("error: %s", err)
		return
	}

	resp.Template(w, tobj, http.StatusOK, MainSite)
}

// TODO: Add resp.File for serving FS files.
func GetCSS(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	path := filepath.Join("files/css", name)
	bytes, err := MainFS.ReadFile(path)
	if err != nil {
		resp.Error(w, http.StatusNotFound, "%q not found", name)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func main() {
	fset := flag.NewFlagSet("stnts", flag.ExitOnError)
	addr := fset.String("addr", "localhost:8000", "host:port address")
	conf := fset.String("conf", "", "configuration file")
	fset.Parse(os.Args[1:])

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", ware.Wrap(GetIndex))
	mux.HandleFunc("GET /css/{name...}", ware.Wrap(GetCSS))

	// TODO: Add site.Parse to parse JSON files into Sites.
	MainSite = new(site.Site)
	if err := file.ReadJSON(*conf, MainSite); err != nil {
		log.Fatalf("error: %s", err)
	}

	srv := &http.Server{Addr: *addr, Handler: mux}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
