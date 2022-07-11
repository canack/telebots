package speech

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"io"
)

var pollyClient *polly.Client
var sampleRate = "24000"

func SetupAWS() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	pollyClient = polly.NewFromConfig(cfg)
	return nil
}

func TextToSpeech(text *string) (io.ReadCloser, error) {
	out, err := pollyClient.SynthesizeSpeech(context.TODO(), &polly.SynthesizeSpeechInput{
		LanguageCode: types.LanguageCodeTrTr,
		Engine:       types.EngineStandard,
		Text:         text,
		OutputFormat: types.OutputFormatMp3,
		TextType:     types.TextTypeText,
		VoiceId:      types.VoiceIdFiliz,
		SampleRate:   &sampleRate,
	})
	if err != nil {
		return nil, err
	}
	return out.AudioStream, nil
}
