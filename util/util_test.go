package util

import (
	"os"
	"testing"
)

func TestEdgeTts(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	EdgeTts(wd+`\output\SpyMageSystem\ch-1.txt`, wd+`\output\SpyMageSystem\ch-1.mp3`)
}

func TestFfmpeg(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	MakeImageVideo(wd+`\output\SpyMageSystem\cover.jpg`, wd+`\output\SpyMageSystem\ch-1.mp3`, wd+`\output\SpyMageSystem\ch-1.mp4`)
}
