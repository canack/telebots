package handler

import (
	"bytes"
	"github.com/canack/telebots/services/minipolly/speech"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (m *mediaFile) textToSound() error {
	speechReader, err := speech.TextToSpeech(&m.text)
	if err != nil {
		log.Println(err)
		return err
	}

	speechBytes := new(bytes.Buffer)
	io.Copy(speechBytes, speechReader)
	if err := ioutil.WriteFile(m.rawSoundFilename, speechBytes.Bytes(), 0640); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *mediaFile) changeSoundPitch() error {
	//cmd := exec.Command("sox", m.rawSoundFilename, m.pitchedSoundFilename, "speed", "1.06", "pitch", m.pitch)
	cmd := exec.Command("sox", m.rawSoundFilename, m.pitchedSoundFilename, "speed", "1.04", "pitch", "20")
	//cmd := exec.Command("sox", m.rawSoundFilename, m.pitchedSoundFilename) // "speed", "1.06", "pitch", m.pitch)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *mediaFile) registerSoundLength() error {
	cmd := exec.Command("soxi", "-D", m.pitchedSoundFilename)
	lengthBytes, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return err
	}
	rawLength := strings.Trim(string(lengthBytes), "\n")
	f, err := strconv.ParseFloat(rawLength, 64)
	if err != nil {
		log.Println(err)
		return err
	}

	m.soundLength = f

	return nil
}

func (m *mediaFile) registerRandomVideo() error {
	log.Println("Registering random video")

	// get random raw videos in disk
	files, err := ioutil.ReadDir("videos/raw")
	if err != nil {
		log.Println(err)
		return err

	}
	var fileList []string

	for f := range files {
		fileList = append(fileList, files[f].Name())
	}

	rand.Seed(time.Now().UnixNano())
	luckyVideo := rand.Intn(len(fileList))

	m.rawVideoFilename = "videos/raw/" + fileList[luckyVideo]

	return nil
}

// cut video by given time
func (m *mediaFile) cutVideo(endTime string) error {
	m.cuttedVideoFilename = "tmp/" + randomString(28) + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", m.rawVideoFilename,
		"-ss", "00:00:00", "-to", endTime, "-c", "copy",
		"-c:v", "libx264", "-crf", "28", "-preset", "veryfast", m.cuttedVideoFilename)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (m *mediaFile) putSoundOnRandomVideo() error {

	// seconds to minutes and seconds
	minutes := int(m.soundLength / 60)
	seconds := int(m.soundLength - (float64(minutes) * 60))
	seconds += 4 // add 6 seconds to make sure sound is not cutted

	//minutes to string
	minutesString := strconv.Itoa(minutes)
	//zerofill minutes
	if minutes < 10 {
		minutesString = "0" + minutesString
	}
	//seconds to string
	secondsString := strconv.Itoa(seconds)
	//zerofill seconds
	if seconds < 10 {
		secondsString = "0" + secondsString
	}

	//endtime
	endTime := "00" + ":" + minutesString + ":" + secondsString

	if err := m.cutVideo(endTime); err != nil {
		log.Println(err)
		return err
	}
	log.Println("Video cutted")
	// merge sound and video
	cmd := exec.Command("ffmpeg", "-i", m.cuttedVideoFilename, "-i", m.pitchedSoundFilename, "-c:v",
		"libx264", "-c", "copy", "-crf", "28", "-preset", "veryfast", "-map", "0:v", "-map", "1:a", m.processedVideoFilename)

	if err := cmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *mediaFile) GenerateMedia() error {
	if err := m.textToSound(); err != nil {
		return err
	}
	log.Println("Text to speech generated")
	if err := m.changeSoundPitch(); err != nil {
		return err
	}
	log.Println("Sound pitch changed")
	if err := m.registerSoundLength(); err != nil {
		return err
	}
	log.Println("Sound length registered")
	if err := m.registerRandomVideo(); err != nil {
		return err
	}
	log.Println("Random video registered")
	if err := m.putSoundOnRandomVideo(); err != nil {
		return err
	}
	log.Println("Sound put on random video")
	return nil
}

// DeleteTempFiles delete temporary files
func (m *mediaFile) DeleteTempFiles() {
	log.Println("Deleting temporary files")
	os.Remove(m.rawSoundFilename)
	os.Remove(m.pitchedSoundFilename)
	os.Remove(m.cuttedVideoFilename)
	os.Remove(m.processedVideoFilename)
}

// generate random string by given range
func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
