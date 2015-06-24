package iiif

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//{scheme}://{server}{/prefix}/{identifier}/{region}/{size}/{rotation}/{quality}.{format}

//{scheme}://{server}{/prefix}/{identifier}/info.json

type URI struct {
	Info       bool
	Scheme     string
	Server     string
	Prefix     string
	Identifier string
	Region     Region
	Size       Size
	Rotation   Rotation
	Quality    Quality
	Format     Format
}

// build the image ID indicated by the URI
// http://iiif.io/api/image/2.0/#information-request for @id
func (params URI) ImageID() string {
	return fmt.Sprintf("%s://%s%s/%s", params.Scheme, params.Server, params.Prefix, params.Identifier)
}

// true if the URI indicates all defaults such that the source image should be passed straight through
// without modifications
func (params URI) Unmodified() bool {
	if !params.Region.Full {
		return false
	}

	if !params.Size.Full {
		return false
	}

	if params.Rotation.Mirror || (params.Rotation.N != 0.0) {
		return false
	}

	if string(params.Quality) != "default" {
		return false
	}

	//match the format to the filename extension
	ext := strings.TrimPrefix(filepath.Ext(params.Identifier), ".")
	if string(params.Format) != ext {
		return false
	}

	return true
}

type Region struct {
	Full    bool
	Percent bool
	X       string
	Y       string
	W       string
	H       string
}

func (r Region) ParseFloats() (float64, float64, float64, float64, error) {
	x, xerr := strconv.ParseFloat(r.X, 64)
	y, yerr := strconv.ParseFloat(r.Y, 64)
	w, werr := strconv.ParseFloat(r.W, 64)
	h, herr := strconv.ParseFloat(r.H, 64)
	var e error
	switch true {
	case xerr != nil:
		e = xerr
		break
	case yerr != nil:
		e = yerr
		break
	case werr != nil:
		e = werr
		break
	case herr != nil:
		e = herr
		break
	}
	return x, y, w, h, e
}

func (r Region) ParseInts() (int64, int64, int64, int64, error) {
	x, xerr := strconv.ParseInt(r.X, 10, 64)
	y, yerr := strconv.ParseInt(r.Y, 10, 64)
	w, werr := strconv.ParseInt(r.W, 10, 64)
	h, herr := strconv.ParseInt(r.H, 10, 64)
	var e error
	switch true {
	case xerr != nil:
		e = xerr
		break
	case yerr != nil:
		e = yerr
		break
	case werr != nil:
		e = werr
		break
	case herr != nil:
		e = herr
		break
	}
	return x, y, w, h, e
}

type Size struct {
	Form    string
	Full    bool
	Percent bool
	Scale   float64
	W       int64
	H       int64
	BestFit bool //! leading w,h indicates scaling to best fit such that dimensions w and h are not exceeded, while maintaining aspect ratio
}

type Rotation struct {
	N      float64
	Mirror bool
}

type Quality string

type Format string

func (f Format) ToImagingFormat() imaging.Format {
	switch f {
	case "jpg":
		return imaging.JPEG
		break
	case "png":
		return imaging.PNG
		break
	case "tif":
		return imaging.TIFF
	case "gif":
		return imaging.GIF
	default:
		return imaging.JPEG
	}
	return imaging.JPEG
}

func ParseRegion(r string) (Region, error) {
	region := Region{}
	if r == "full" {
		region.Full = true
		return region, nil
	}
	if strings.HasPrefix(r, "pct:") {
		region.Percent = true
		r = strings.TrimPrefix(r, "pct:")
	}
	numbers := strings.Split(r, ",")
	if len(numbers) < 4 {
		return region, errors.New("ParseRegion: not enough numbers")
	}
	region.X = numbers[0]
	region.Y = numbers[1]
	region.W = numbers[2]
	region.H = numbers[3]

	return region, nil
}

func ParseSize(s string) (Size, error) {
	size := Size{}
	var err error
	if s == "full" {
		size.Full = true
		size.Form = "full"
		return size, nil
	}
	if strings.HasPrefix(s, "pct:") {
		size.Percent = true
		size.Form = "pct:n"
		s = strings.TrimPrefix(s, "pct:")
		size.Scale, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return size, err
		}
		return size, nil
	}
	if strings.HasPrefix(s, ",") {
		size.Form = ",h"
		s = strings.TrimPrefix(s, ",")
		size.H, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return size, err
		}
		return size, nil
	}
	if strings.HasSuffix(s, ",") {
		size.Form = "w,"
		s = strings.TrimSuffix(s, ",")
		size.W, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return size, err
		}
		return size, nil
	}
	if strings.HasPrefix(s, "!") {
		size.Form = "!w,h"
		size.BestFit = true
		s = strings.TrimPrefix(s, "!")
		matches := strings.Split(s, ",")
		size.W, err = strconv.ParseInt(matches[0], 10, 64)
		if err != nil {
			return size, err
		}
		size.H, err = strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return size, err
		}
		return size, nil
	}

	size.Form = "w,h"
	matches := strings.Split(s, ",")
	size.W, err = strconv.ParseInt(matches[0], 10, 64)
	if err != nil {
		return size, err
	}
	size.H, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return size, err
	}
	return size, nil
}

func ParseRotation(r string) (Rotation, error) {
	rotation := Rotation{}
	var err error
	if strings.HasPrefix(r, "!") {
		rotation.Mirror = true
		r = strings.TrimPrefix(r, "!")
	}
	rotation.N, err = strconv.ParseFloat(r, 64)
	if err != nil {
		return rotation, err
	}
	return rotation, nil
}

/*
func ParseURIStart(s string, prefix string) {
	startRE := regexp.MustCompile(`([^:]+):\/\/([^/]+)/`)
	matches := startRE.FindStringSubmatch(s, -1)
	if len(matches) < 3 {
		return URI{}, errors.New("Not enough parameters in URI")
		u.Info = true
		u.Identifier = matches[0][1]
		return u, nil
	}
}
*/

//parse the non-server section of an IIIF URI (specifying the image ID and parameters)
func ParseURI(s string) (URI, error) {
	u := URI{}

	//check first for the simpler info.json request
	//(identifier)(region)(size)(rotation)(quality)(format)
	//{scheme}://{server}{/prefix}/{identifier}/info.json
	infoRE := regexp.MustCompile(`/([^/]+)/info\.json$`)
	matches := infoRE.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 && len(matches[0]) == 2 {
		u.Info = true
		u.Identifier = matches[0][1]
		return u, nil
	}

	//capture groups between slashes (and format extension)
	//(identifier)(region)(size)(rotation)(quality)(format)
	paramsRE := regexp.MustCompile(`/([^/]+)/([^/]+)/([^/]+)/([^/]+)/([^/]+)\.([^/]+)`)
	matches = paramsRE.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 || len(matches[0]) < 7 {
		return URI{}, errors.New("Not enough parameters in URI")
	}
	u.Info = false
	u.Identifier = matches[0][1]
	u.Quality = Quality(matches[0][5])
	u.Format = Format(matches[0][6])

	region, err := ParseRegion(matches[0][2])
	if err != nil {
		return u, err
	}
	u.Region = region

	size, err := ParseSize(matches[0][3])
	if err != nil {
		return u, err
	}
	u.Size = size

	rotation, err := ParseRotation(matches[0][4])
	if err != nil {
		return u, err
	}
	u.Rotation = rotation

	return u, nil
}
