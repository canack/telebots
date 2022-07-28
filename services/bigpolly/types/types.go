package types

import "image"

var VideoPath = "videos/"
var TempPath = "tmp/"

type Webpage struct {
}

type ImageAndText struct {
	Image image.Image
	Text  string
}

type PrepareForProcessing struct {
	VideoFilename, ImageFilename, AudioFilename, Text string
	AudioLength                                       float64
}
