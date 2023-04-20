package util

import (
	"fmt"
	"os/exec"
)

func EdgeTts(ttsTextPath string) {
	pyExe := "python3"
	pyFile := "D:\\python_workspace\\tutorial\\tutorial\\speech\\edge-tts.py"
	out, err := exec.Command(pyExe, pyFile).Output()
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}
	println(fmt.Sprintf("out: %s", out))
}
