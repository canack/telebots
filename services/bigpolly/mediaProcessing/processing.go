package mediaProcessing

import (
	"github.com/canack/telebots/services/bigpolly/speech"
	"github.com/canack/telebots/services/bigpolly/types"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func registerRandomVideo() (string, error) {
	log.Println("Registering random video")

	// get random raw videos in disk
	files, err := ioutil.ReadDir(types.VideoPath)
	if err != nil {
		log.Println(err)
		return "", err

	}
	var fileList []string

	for f := range files {
		fileList = append(fileList, files[f].Name())
	}

	luckyVideo := rand.Intn(len(fileList))

	video := types.VideoPath + fileList[luckyVideo]
	return video, nil
}

// write image to file
func saveImage(image image.Image) string {
	filename := types.TempPath + RandomString(32) + ".jpg"
	f, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer f.Close()
	err = jpeg.Encode(f, image, nil)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return filename
}

func saveAudio(audio io.ReadCloser) (string, error) {
	s := types.TempPath + RandomString(32) + ".mp3"
	f, err := os.Create(s)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()

	io.Copy(f, audio)
	return s, nil
}

// return audio length with using soxi
func getAudioLength(audioFilename string) (float64, error) {
	code := "soxi -D " + audioFilename
	executeCommand := exec.Command("bash", "-c", code)
	out, err := executeCommand.Output()
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	// parse newline out
	out = out[:len(out)-1]

	audioLength, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return audioLength, nil

}
func PreProcess(media *[]types.ImageAndText) []types.PrepareForProcessing {
	var prepare []types.PrepareForProcessing
	video, err := registerRandomVideo()

	if err != nil {
		log.Println(err)
		panic(err)
	}

	for _, v := range *media {
		speechReader, err := speech.TextToSpeech(&v.Text)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		audioFilename, err := saveAudio(speechReader)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		audioLength, err := getAudioLength(audioFilename)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		imgFilename := saveImage(v.Image)
		prepare = append(prepare, types.PrepareForProcessing{
			VideoFilename: video,
			AudioFilename: audioFilename,
			AudioLength:   audioLength,
			ImageFilename: imgFilename,
			Text:          v.Text,
		})
	}
	return prepare
}

func Process(input *[]types.ImageAndText) (string, error) {
	prep := PreProcess(input)

	ap := ProcessAudio(prep)
	acode := GenerateAudioCode(ap)

	vp := ProcessVideo(prep)
	vcode := GenerateVideoCode(vp)

	inputVideo, err := registerRandomVideo()
	if err != nil {
		return "", err
	}

	out := types.TempPath + RandomString(32) + ".mp4"
	final := types.TempPath + RandomString(32) + ".mp4"

	if err := ExecuteAudioCode(inputVideo, out, acode); err != nil {
		return "", err
	}

	if err := ExecuteVideoCode(out, final, vcode); err != nil {
		return "", err
	}

	return final, nil
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
