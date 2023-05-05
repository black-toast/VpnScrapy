package main

import (
	"VpnScrapy/util"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-26.txt`, wd+`\output\SpyMageSystem\ch-26.mp3`)
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-27.txt`, wd+`\output\SpyMageSystem\ch-27.mp3`)
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-28.txt`, wd+`\output\SpyMageSystem\ch-28.mp3`)
}
