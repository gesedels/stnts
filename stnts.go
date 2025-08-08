package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gesedels/stnts/stnts/items/site"
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

func GetFile(w http.ResponseWriter, r *http.Request) {
	name := "files/" + r.PathValue("name")
	switch filepath.Ext(name) {
	case ".css":
		resp.File(w, MainFS, http.StatusOK, name, "text/css; charset=utf-8")
	case ".html":
		resp.Error(w, http.StatusNotFound, "%q not found", name)
	default:
		resp.File(w, MainFS, http.StatusOK, name, "")
	}
}

func main() {
	fset := flag.NewFlagSet("stnts", flag.ExitOnError)
	addr := fset.String("addr", "localhost:8000", "host:port address")
	conf := fset.String("conf", "", "configuration file")
	fset.Parse(os.Args[1:])

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", ware.Wrap(GetIndex))
	mux.HandleFunc("GET /file/{name...}", ware.Wrap(GetFile))

	site, err := site.Parse(*conf)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	MainSite = site
	srv := &http.Server{Addr: *addr, Handler: mux}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
