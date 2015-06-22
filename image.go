package iiif

import (
	//"errors"
	"github.com/disintegration/imaging"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

/*
func GetRegion(img image.Image, region Region) (image.Image, error) {

}

func Resize(img image.Image, size Size) (image.Image, error) {

}

func Rotate(img image.Image, r Rotation) (image.Image, error) {

}

func GetQuality(img image.Image, q Quality) (image.Image, error) {

}

func FormatImage(img image.Image, f Format) ([]byte, error) {

}
*/

type Dimensions struct {
	Width  int `json:"width`
	Height int `json:"height`
}

type Profile struct {
	Id        string   `json:""`
	Formats   []string `json:"formats"`
	Qualities []string `json:"qualities"`
	Supports  []string `json:"supports"`
}

type IIIFInfo struct {
	Context  string `json:"@context"`
	Id       string `json:"@id"`
	Protocol string `json:"protocol"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Profile  `json:"profile"`
	//Sizes    []Dimensions `json:"sizes"`
	//Tiles
	//Service
}

func Info(filename string) map[string]string {
	im, err := imaging.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	sourceBounds := im.Bounds()
	sourceWidth := sourceBounds.Dx()
	sourceHeight := sourceBounds.Dy()

	info := make(map[string]string)
	info["width"] = string(sourceWidth)
	info["height"] = string(sourceHeight)
	return info
}

func Process(filename string, params URI) (image.Image, error) {
	im, err := imaging.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	sourceBounds := im.Bounds()
	sourceWidth := sourceBounds.Dx()
	sourceHeight := sourceBounds.Dy()

	//crop to the selected region
	switch true {
	case params.Region.Full:
		break
	case params.Region.Percent:
		x, y, w, h, err := params.Region.ParseFloats()
		if err != nil {
			log.Fatal(err)
		}
		x0 := int(x * float64(sourceWidth))
		y0 := int(y * float64(sourceHeight))
		x1 := x0 + int(w*float64(sourceWidth))
		y1 := y0 + int(h*float64(sourceHeight))
		im = imaging.Crop(im, image.Rect(x0, y0, x1, y1))
		break
	default:
		//x,y,w,h format
		x, y, w, h, err := params.Region.ParseInts()
		if err != nil {
			log.Fatal(err)
		}
		im = imaging.Crop(im, image.Rect(int(x), int(y), int(x+w), int(y+h)))
		break
	}

	//resize region to requested dimensions
	switch params.Size.Form {
	case "full":
		break
	case "w,":
		im = imaging.Resize(im, int(params.Size.W), 0, imaging.Hamming)
		break
	case ",h":
		im = imaging.Resize(im, 0, int(params.Size.H), imaging.Hamming)
		break
	case "pct:n":
		scaledW := int(params.Size.Scale * float64(sourceWidth))
		scaledH := int(params.Size.Scale * float64(sourceHeight))
		im = imaging.Resize(im, scaledW, scaledH, imaging.Hamming)
		break
	case "w,h":
		im = imaging.Resize(im, int(params.Size.W), int(params.Size.H), imaging.Hamming)
		break
	case "!w,h":
		im = imaging.Fit(im, int(params.Size.W), int(params.Size.H), imaging.Hamming)
		break
	default:
		log.Fatal("unknown Size form")
	}

	//rotate image by 90 degree increments
	//iiif specifies rotation clockwise, imaging specifies counterclockwise
	if params.Rotation.Mirror {
		im = imaging.FlipH(im)
	}
	switch params.Rotation.N {
	case 0.0:
		break
	case 90.0:
		im = imaging.Rotate270(im)
		break
	case 180.0:
		im = imaging.Rotate180(im)
		break
	case 270.0:
		im = imaging.Rotate90(im)
		break
	default:
		log.Fatal("unsupported rotation degrees")
	}

	switch params.Quality {
	case "default":
		break
	case "color":
		break
	case "gray":
		log.Fatal("unimplemented")
		break
	case "bitonal":
		log.Fatal("unimplemented")
		break
	default:
		log.Fatal("unknown Quality argument")
	}

	return im, nil
}
