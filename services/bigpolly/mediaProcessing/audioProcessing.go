package mediaProcessing

import (
	"fmt"
	"github.com/canack/telebots/services/bigpolly/types"
	"log"
	"os/exec"
)

type AudioProcessing struct {
	StartAlias   string  //
	EndAlias     string  //
	Filename     string  // Filename of video
	Length       float64 // Length of audio
	BetweenStart float64 // Start time of showing image
	BetweenEnd   float64 // End time of showing image
}

func ProcessAudio(input []types.PrepareForProcessing) []AudioProcessing {
	var AP []AudioProcessing

	initNum := 1
	startAlias := fmt.Sprintf("[%d:a]", initNum)
	endAlias := fmt.Sprintf("[a%d]", initNum)

	var start float64
	var end float64

	for _, v := range input {

		end = start + v.AudioLength

		AP = append(AP, AudioProcessing{
			StartAlias:   startAlias,
			EndAlias:     endAlias,
			Filename:     v.AudioFilename,
			Length:       v.AudioLength,
			BetweenStart: start,
			BetweenEnd:   end,
		})

		start += v.AudioLength

		initNum++
		startAlias = fmt.Sprintf("[%d:a]", initNum)
		endAlias = fmt.Sprintf("[a%d]", initNum)
	}

	return AP

}

func GenerateAudioCode(input []AudioProcessing) string {

	code := ""

	for _, v := range input {
		code += fmt.Sprintf("-i %s ", v.Filename)
	}

	code += `-filter_complex "`

	for _, v := range input {
		code += fmt.Sprintf("%satrim=end=%f,asetpts=PTS-STARTPTS%s;",
			v.StartAlias, v.Length, v.EndAlias)
	}

	audioLen := len(input)

	for i := 1; i <= audioLen; i++ {
		code += fmt.Sprintf("[a%d]", i)
	}

	code += "concat=n=" + fmt.Sprintf("%d", audioLen) + ":v=0:a=1[a]"
	code += `"`

	code += " -map 0:v -map \"[a]\" -codec:v copy -codec:a libmp3lame -shortest"
	return code

}

func ExecuteAudioCode(videoInputFilename, videoOutputFilename, command string) error {
	code := "ffmpeg -i " + videoInputFilename + " " + command + " " + videoOutputFilename
	executeCommand := exec.Command("bash", "-c", code)

	err := executeCommand.Run()

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Audio processed video created")
	return nil
}
