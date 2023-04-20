package util

import (
	"fmt"
	"os/exec"
)

func MakeImageVideo(imagePath, audioPath, outputVideoPath string) {
	// 设置 ffmpeg 命令行参数
	args := []string{
		"-loop",
		"1",
		"-i",
		imagePath,
		"-i",
		audioPath,
		"-c:v",
		"libx264",
		"-c:a",
		"aac",
		"-b:a",
		"192k",
		"-shortest",
		outputVideoPath,
	}

	fmt.Println(args)

	// 创建 *exec.Cmd
	cmd := exec.Command("ffmpeg", args...)

	// // 运行 ffmpeg 命令
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("转码成功")
}
