package utils

import (
	"fmt"
	"io"

	"github.com/dhowden/tag"
	"github.com/hajimehoshi/go-mp3"
	"github.com/pkg/errors"
)

// HandleMP3File ...
func HandleMP3File(file io.ReadSeeker) (fileImageBytes []byte, duration, artist string, err error) {
	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "HandleMP3File error: can't get metadata")
	}

	fileImage := metadata.Picture()
	if fileImage == nil {
		fileImageBytes = CreateWhiteQuestionImage(200, 200)
	}
	if fileImage != nil {
		fileImageBytes = metadata.Picture().Data
	}

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "HandleMP3File error: can't get decoder")
	}

	audioLength := decoder.Length() / int64(decoder.SampleRate()) / 4
	duration = fmt.Sprintf("%02d:%02d", audioLength/60, audioLength%60)

	artist = metadata.Artist()

	if artist == "" {
		artist = "Unknown artist"
	}

	return fileImageBytes, duration, artist, nil
}
