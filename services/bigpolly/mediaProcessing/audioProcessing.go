package mediaProcessing

import (
	"fmt"
	"github.com/canack/telebots/services/bigpolly/speech"
	"github.com/canack/telebots/services/bigpolly/types"
	"io"
	"log"
	"os/exec"
)

type AudioProcessing struct {
	StartAlias   string    //
	EndAlias     string    //
	Filename     string    //
	Audio        io.Reader // nth audio
	Length       int       // Length of audio
	BetweenStart int       //
	BetweenEnd   int       //
}

// return audio length with using soxi
func getAudioLength(audio io.Reader) int {
	return 0
}

func ProcessAudio(input []types.ImageAndText) []AudioProcessing {
	var AP []AudioProcessing

	for _, v := range input {
		initNum := 1
		startAlias := fmt.Sprintf("[%d:a]", initNum)
		endAlias := fmt.Sprintf("[a%d]", initNum)

		speechReader, err := speech.TextToSpeech(&v.Text)
		if err != nil {
			panic(err)
		}

		audioLength := getAudioLength(speechReader)

		AP = append(AP, AudioProcessing{
			StartAlias:   startAlias,
			EndAlias:     endAlias,
			Audio:        speechReader,
			Length:       audioLength,
			BetweenStart: 0, // Replace here
			BetweenEnd:   0, // Replace here
		})

		initNum++
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
		code += fmt.Sprintf("[%s]atrim=end=%d,asetpts=PTS-STARTPTS[%s];",
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

func GenerateAudioProcessedVideo(videoInputFilename, videoOutputFilename, command string) error {
	code := "ffmpeg -i " + videoInputFilename + " " + command + " " + videoOutputFilename

	executeCommand := exec.Command(code)

	err := executeCommand.Run()

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Audio processed video created")
	return nil
}
