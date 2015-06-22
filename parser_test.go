package iiif

import (
	"log"
	"testing"
)

var _ = log.Print

var testPathResults = map[string]URI{
	"/ident4321/full/full/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/4321/10,20,100,200/full/0/default.jpg": URI{
		Identifier: "4321",
		Region:     Region{X: "10", Y: "20", W: "100", H: "200"},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/pct:0,0,5.5,10.10/full/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Percent: true, X: "0", Y: "0", W: "5.5", H: "10.10"},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/200,/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Form: "w,", W: 200},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/,200/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Form: ",h", H: 200},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/pct:25.0/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Form: "pct:n", Scale: 25.0, Percent: true},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/200,400/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Form: "w,h", W: 200, H: 400},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/!200,400/0/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Form: "!w,h", W: 200, H: 400, BestFit: true},
		Rotation:   Rotation{},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/full/22.5/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{N: 22.5},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/full/!22.5/default.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{Mirror: true, N: 22.5},
		Quality:    Quality("default"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/full/!22.5/color.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{Mirror: true, N: 22.5},
		Quality:    Quality("color"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/full/!22.5/gray.jpg": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{Mirror: true, N: 22.5},
		Quality:    Quality("gray"),
		Format:     Format("jpg"),
	},
	"/ident4321/full/full/!22.5/default.pdf": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{Mirror: true, N: 22.5},
		Quality:    Quality("default"),
		Format:     Format("pdf"),
	},
	"/ident4321/full/full/!22.5/bimodal.png": URI{
		Identifier: "ident4321",
		Region:     Region{Full: true},
		Size:       Size{Full: true, Form: "full"},
		Rotation:   Rotation{Mirror: true, N: 22.5},
		Quality:    Quality("bimodal"),
		Format:     Format("png"),
	},
}

func TestValidPaths(t *testing.T) {
	for testPath, expected := range testPathResults {
		uriparams, err := ParseURI(testPath)
		if err != nil {
			t.Error("Error when parsing valid path")
		}
		if uriparams != expected {
			log.Printf("Input: %s", testPath)
			log.Printf("\nResult  : %+v\nExpected: %+v\n", uriparams, expected)
			t.Error("Unexpected result parsing valid test path")
		}
	}
}
