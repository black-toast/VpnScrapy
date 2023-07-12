package main

import (
	"VpnScrapy/bean"
	"VpnScrapy/crawl/novel"
	"VpnScrapy/storage"
	"VpnScrapy/util"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var titles = [6]string{
	"necropolis-immortal-book",
	"spy-mage-system-book",
	"cultivation-online-novel-book",
	"the-experimental-log-of-the-crazy-lich-book",
	"invincible-divine-dragons-cultivation-system-book",
	"trial-marriage-husband-need-to-work-hard-book",
}

func TestMain(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-26.txt`, wd+`\output\SpyMageSystem\ch-26.mp3`)
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-27.txt`, wd+`\output\SpyMageSystem\ch-27.mp3`)
	util.EdgeTts(wd+`\output\SpyMageSystem\ch-28.txt`, wd+`\output\SpyMageSystem\ch-28.mp3`)
}

func TestScrapyNovel(t *testing.T) {
	// novel.Scrapy(57, 60, titles[0])
	// novel.Scrapy(57, 82, titles[1])
	// novel.Scrapy(329, 329, titles[2])
	// novel.Scrapy(487, 487, titles[3])
	novel.Scrapy(379, 384, titles[4])
	// novel.Scrapy(187, 187, titles[5])
}

func TestGenerateChapterList(t *testing.T) {
	outputDir := "./output/"
	chapterPrefix := "ch-"
	chapterSuffix := ".txt"
	novelDirs, err := storage.ReadDir(outputDir)
	if err != nil {
		panic(err)
	}

	novels := make([]*bean.Novel, 0)

	for novelIndex, novelDir := range novelDirs {
		if novelDir.Name() == "novels.json" {
			continue
		}

		novelDir := fmt.Sprintf("%s%s", outputDir, novelDir.Name())
		novelNamePath := fmt.Sprintf("%s/title.txt", novelDir)
		novelNameByte, err := storage.Read(novelNamePath)
		if err != nil {
			panic(err)
		}
		novelName := string(novelNameByte)

		novelDescPath := fmt.Sprintf("%s/desc.txt", novelDir)
		novelDescByte, err := storage.Read(novelDescPath)
		if err != nil {
			panic(err)
		}
		novelDesc := string(novelDescByte)

		novel := &bean.Novel{
			Index:       novelIndex + 1,
			Name:        novelName,
			ChapterList: make([]*bean.NovelChapter, 0),
		}

		chapterList, err := storage.ReadDir(novelDir)
		if err != nil {
			panic(err)
		}
		for _, chapter := range chapterList {
			if strings.HasPrefix(chapter.Name(), chapterPrefix) &&
				strings.HasSuffix(chapter.Name(), chapterSuffix) {
				chapterPath := fmt.Sprintf("%s/%s", novelDir, chapter.Name())
				chapterByte, err := storage.Read(chapterPath)
				if err != nil {
					panic(err)
				}

				chapterContent := string(chapterByte)
				if chapterContent == "" {
					continue
				}
				chapterContentSplit := strings.Split(chapterContent, "\n")
				chapterTitleSplit := strings.Split(chapterContentSplit[0], ": ")
				chapterIndexStr := strings.ReplaceAll(chapterTitleSplit[0], "Chapter ", "")
				chapterIndex, err := strconv.Atoi(chapterIndexStr)
				if err != nil {
					panic(err)
				}
				chapterTitle := chapterTitleSplit[1]

				title := fmt.Sprintf("%s CH-%d", novel.Name, chapterIndex)
				desc := fmt.Sprintf("NOVEL INFO: %s\nChapter %d: %s\nDescription: %s",
					novel.Name, chapterIndex, chapterTitle, novelDesc)
				novel.ChapterList = append(novel.ChapterList, &bean.NovelChapter{
					Index: chapterIndex,
					Title: title,
					DESC:  desc,
				})
			}
		}
		novels = append(novels, novel)
	}

	saveNovelsJson(novels)
}

func saveNovelsJson(novels []*bean.Novel) {
	if novelsJson, err := util.Format(novels); err != nil {
		panic(err)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		novelJsonPath := fmt.Sprintf("%s\\output\\novels.json", wd)
		saveFile := storage.Create(novelJsonPath)
		saveFile.WriteString(novelsJson)
		saveFile.Close()

		// fmt.Println("novelsJson", novelsJson)
	}
}

func TestMkdir(t *testing.T) {
	dir := "D:\\go_workspace\\src\\VpnScrapy\\output\\MarvelsSuperman\\ch-%d.txt"
	for i := 403; i < 538; i++ {
		path := fmt.Sprintf(dir, i)
		storage.Create(path)
	}
}

func TestTransform(t *testing.T) {
	dir := "D:\\go_workspace\\src\\VpnScrapy\\output\\MarvelsSuperman"
	for index := 412; index < 415; index++ {
		chapterFile := fmt.Sprintf("ch-%d.txt", index+1)

		if index == 0 {
			fmt.Println("wait 2s, then request novel chapter(", chapterFile, ")")
		} else {
			fmt.Println("wait 2s, then request next novel chapter(", chapterFile, ")")
		}
		time.Sleep(2 * time.Second)

		t := time.Now()
		//参数必须是这个时间,格式任意
		// s := t.Format("2006-01-02 15:04:05")
		currentTime := t.Format("2006-01-02 15:04:05")
		fmt.Printf("current time: %s\n", currentTime)

		startCost := t.Unix()
		chapterByte, err := storage.Read(dir + "\\" + chapterFile)
		if err != nil {
			panic(err)
		}
		novel.TransformFormat(index, dir, chapterFile, string(chapterByte))

		endCost := time.Now().Unix()
		fmt.Printf("👆======================cost %ds=======================👆\n", (endCost - startCost))
	}
}
