package iiif

import (
	"encoding/json"
	"errors"
	"github.com/disintegration/imaging"
	"image"
	_ "image/jpeg"
	_ "image/png"
	//	"log"
)

//http://iiif.io/api/image/2.0/#image-information sizes (currently unused by this implementation)
type Dimensions struct {
	Width  int `json:"width`
	Height int `json:"height`
}

//http://iiif.io/api/image/2.0/#image-information profiles subsection
type Profile struct {
	Id        string   `json:""`
	Formats   []string `json:"formats"`
	Qualities []string `json:"qualities"`
	Supports  []string `json:"supports"`
}

func (p Profile) MarshalJSON() ([]byte, error) {
	if p.Id != "" {
		return json.Marshal(p.Id)
	}

	m := make(map[string][]string)
	m["formats"] = p.Formats
	m["qualities"] = p.Qualities
	m["supports"] = p.Supports
	return json.Marshal(m)
}

// http://iiif.io/api/image/2.0/#information-request
type IIIFInfo struct {
	Context  string    `json:"@context"`
	Id       string    `json:"@id"`
	Protocol string    `json:"protocol"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Profile  []Profile `json:"profile"`
	//Sizes    []Dimensions `json:"sizes"`
	//Tiles
	//Service
}

//Generate a default info struct for this server implementation
func DefaultInfo() IIIFInfo {
	profiles := make([]Profile, 2)
	profiles[0] = Profile{Id: "http://iiif.io/api/image/2/level0.json"}
	profiles[1] = Profile{
		Formats: []string{
			"jpg",
			"png",
			"tif",
			"gif",
		},
		Qualities: []string{"default"},
		Supports:  []string{"cors", "mirroring", "regionByPx", "regionByPct", "rotationBy90s", "sizeAboveFull", "sizeByWhListed", "sizeByH", "sizeByPct", "sizeByW", "sizeByWh"},
	}

	return IIIFInfo{
		Context:  "http://iiif.io/api/image/2/context.json",
		Protocol: "http://iiif.io/api/image",
		Profile:  profiles,
	}
}

//Get a IIIFInfo struct about the file passed in
func Info(filename string) (IIIFInfo, error) {
	im, err := imaging.Open(filename)
	if err != nil {
		return IIIFInfo{}, err
	}

	sourceBounds := im.Bounds()
	sourceWidth := sourceBounds.Dx()
	sourceHeight := sourceBounds.Dy()

	iinfo := DefaultInfo()
	iinfo.Width = sourceWidth
	iinfo.Height = sourceHeight

	return iinfo, nil
}

//perform the specified transformations from a IIIF URI on an image
func Process(filename string, params URI) (image.Image, error) {
	im, err := imaging.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, err
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
			return nil, err
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
		return nil, errors.New("Unknown Size form")
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
		return nil, errors.New("unsupported rotation degrees")
	}

	switch params.Quality {
	case "default":
		break
	case "color":
		break
	case "gray":
		return nil, errors.New("unimplemented")
		break
	case "bitonal":
		return nil, errors.New("unimplemented")
		break
	default:
		return nil, errors.New("unknown Quality argument")
	}

	return im, nil
}
