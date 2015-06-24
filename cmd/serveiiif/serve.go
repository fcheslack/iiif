package main

import (
	"encoding/json"
	"flag"
	"github.com/disintegration/imaging"
	"github.com/fcheslack/iiif"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var pathPrefix = flag.String("prefix", "", "prefix for the iiif server")

func main() {
	flag.Parse()
	imageDir := flag.Arg(0)
	imageDir, err := filepath.Abs(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	//add a leading slash to the prefix if there isn't one already
	if !strings.HasPrefix(*pathPrefix, "/") {
		*pathPrefix = "/" + *pathPrefix
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request for %s", r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fullPath := r.URL.Path
		if !strings.HasPrefix(fullPath, *pathPrefix) {
			http.Error(w, "Request without appropriate path prefix", 404)
			return
		}
		path := strings.TrimPrefix(fullPath, *pathPrefix)

		iiifparams, err := iiif.ParseURI(path)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//fill out params based on predefined or http parsed values
		scheme := r.URL.Scheme
		if scheme == "" {
			scheme = "http"
		}
		host := r.Host
		iiifparams.Scheme = scheme
		iiifparams.Server = host
		iiifparams.Prefix = *pathPrefix

		filename, err := url.QueryUnescape(iiifparams.Identifier)
		if err != nil {
			http.Error(w, "source image not found", 404)
			return
		}

		if iiifparams.Info {
			info, err := iiif.Info(filename)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			info.Id = iiifparams.ImageID()

			encw := json.NewEncoder(w)
			err = encw.Encode(info)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			return
		}

		//pass the file straight through if no modifications are needed
		if iiifparams.Unmodified() {
			log.Print("unmodified, passing through")
			http.ServeFile(w, r, filepath.Join(imageDir, filename))
			return
		}

		imgPath := filepath.Join(imageDir, filename)
		img, err := iiif.Process(imgPath, iiifparams)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Print("Encoding to responsewriter")
		imaging.Encode(w, img, iiifparams.Format.ToImagingFormat())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
