package main

import (
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

		//scheme := r.URL.Scheme
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

		filename, err := url.QueryUnescape(iiifparams.Identifier)
		if err != nil {
			http.Error(w, "source image not found", 404)
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
