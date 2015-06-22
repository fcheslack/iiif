package main

import (
	"flag"
	"github.com/disintegration/imaging"
	"github.com/fcheslack/iiif"
	"log"
)

func main() {
	flag.Parse()
	imgPath := flag.Arg(0)
	iiifpath := flag.Arg(1)
	outputImgPath := flag.Arg(2)

	uriparams, err := iiif.ParseURI(iiifpath)
	if err != nil {
		log.Fatal(err)
	}

	img, err := iiif.Process(imgPath, uriparams)
	err = imaging.Save(img, outputImgPath)
	if err != nil {
		log.Fatal(err)
	}
}
