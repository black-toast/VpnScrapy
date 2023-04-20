package util

import "testing"

func TestEdgeTts(t *testing.T) {
	EdgeTts("D:\\python_workspace\\tutorial\\tutorial\\speech\\edge-tts.py")
}

func TestFfmpeg(t *testing.T) {
	dirPath := "D:\\DesktopData\\gp_ad_project\\youtube\\CULTIVATION ONLINE\\"
	MakeImageVideo(dirPath+"cultivation-online-novel.jpg", dirPath+"ch-1170.mp3", dirPath+"ch-1170.mp4")
}
