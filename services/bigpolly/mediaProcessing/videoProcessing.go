package mediaProcessing

import (
	"fmt"
	"github.com/canack/telebots/services/bigpolly/types"
	"log"
	"os/exec"
)

type VideoProcessing struct {
	StartAlias   string  // [0:v] for initialize after than should [tmp]
	EndAlias     string  // Always [tmp]
	Filename     string  // Filename of video
	Overlay      string  // (main_w-overlay_w)/2:(main_h-overlay_h)/2 for center
	BetweenStart float64 // Start time of showing image
	BetweenEnd   float64 // End time of showing image
}

func ProcessVideo(input []types.PrepareForProcessing) []VideoProcessing {
	var VP []VideoProcessing

	initNum := 1
	startAlias := fmt.Sprintf("[0:v][%d:v]", initNum)
	endAlias := "[tmp]"

	var start float64
	var end float64
	for _, v := range input {

		end = start + v.AudioLength

		VP = append(VP, VideoProcessing{
			StartAlias:   startAlias,
			EndAlias:     endAlias,
			Filename:     v.ImageFilename,
			Overlay:      "(main_w-overlay_w)/2:(main_h-overlay_h)/2",
			BetweenStart: start,
			BetweenEnd:   end,
		})

		start += v.AudioLength

		initNum++
		startAlias = fmt.Sprintf("%s[%d:v]", endAlias, initNum)

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
		code += fmt.Sprintf("%s overlay=%s:enable='between(t,%f,%f)' %s;",
			v.StartAlias, v.Overlay, v.BetweenStart, v.BetweenEnd, v.EndAlias)
	}

	code = code[:len(code)-6]
	code += `"`

	return code
}

func ExecuteVideoCode(videoInputFilename, videoOutputFilename, command string) error {
	code := "ffmpeg -i " + videoInputFilename + " " + command + " " + videoOutputFilename
	executeCommand := exec.Command("bash", "-c", code)

	err := executeCommand.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Image processed video generated")
	return nil
}
