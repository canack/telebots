package mediaProcessing

import (
	"fmt"
	"image"
	"log"
	"os/exec"
)

// text to video using ffmpeg

//type VideoProcessing struct {
//	Command    string
//	Screenshot []image.Image
//	Sound      []io.Reader
//}

type VideoProcessing struct {
	StartAlias   string      // [0:v] for initialize after than should [tmp]
	EndAlias     string      // Always [tmp]
	Filename     string      // Filename of video
	Img          image.Image // nth image
	Overlay      string      // (main_w-overlay_w)/2:(main_h-overlay_h) for center
	BetweenStart string      // Start time of showing image
	BetweenEnd   string      // End time of showing image
}

func ProcessVideo(input []image.Image) []VideoProcessing {
	var VP []VideoProcessing

	initNum := 1
	startAlias := fmt.Sprintf("[0:v][%d:v]", initNum)
	endAlias := "[tmp]"

	for _, v := range input {

		VP = append(VP, VideoProcessing{
			StartAlias:   startAlias,
			EndAlias:     endAlias,
			Img:          v,
			Overlay:      "(main_w-overlay_w)/2:(main_h-overlay_h)",
			BetweenStart: "",
			BetweenEnd:   "",
		})

		initNum++
		startAlias = fmt.Sprintf("[%s][%d:v]", endAlias, initNum)

	}

	return VP
}

func GenerateVideoCode(input []VideoProcessing) string {

	code := ""

	for _, v := range input {
		code += fmt.Sprintf("-i %s ", v.Filename)
	}

	code += `-filter_complex "`

	for _, v := range input {
		code += fmt.Sprintf("%s overlay=%s:enable='between(t,%s,%s)' %s;",
			v.StartAlias, v.Overlay, v.StartAlias, v.EndAlias, v.EndAlias)
	}

	code += `"`

	return code
}

func GenerateImageProcessedVideo(videoInputFilename, videoOutputFilename, command string) error {
	code := "ffmpeg -i " + videoInputFilename + " " + command + " " + videoOutputFilename

	executeCommand := exec.Command(code)

	err := executeCommand.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Image processed video generated")
	return nil
}

//func createCommand() {
//	command := `ffmpeg -i input.mp4 -i 1.png -i 2.png -filter_complex
//	"[0:v][1:v] overlay=10:10:enable='between(t,1,2)' [tmp];
//	[tmp][2:v] overlay=20:20:enable='between(t,2,3)'" output.mp4`
//}

// For centered image
//overlay=(main_w-overlay_w)/2:(main_h-overlay_h)
