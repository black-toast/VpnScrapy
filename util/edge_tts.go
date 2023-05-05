package util

import (
	"fmt"
	"os"
	"os/exec"
)

func EdgeTts(ttsInput, ttsOutput string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pyExe := "python3"
	edgeTtsPyFile := fmt.Sprintf("%s\\library\\edge-tts.py", wd)

	cmd := exec.Command(pyExe, edgeTtsPyFile, ttsInput, ttsOutput)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
