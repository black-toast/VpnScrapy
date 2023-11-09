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

var titles = [12]string{
	"cultivation-online-novel",
	"dual-cultivation-novel",
	"immortal-and-martial-dual-cultivation",
	"invincible-divine-dragons-cultivation-system",
	"MarvelsSuperman",
	"necropolis-immortal",
	"shadow-slave",
	"spy-mage-system",
	"supreme-lord-i-can-extract-everything",
	"the-death-mage-who-doesnt-want-a-fourth-time",
	"the-experimental-log-of-the-crazy-lich",
	"trial-marriage-husband-need-to-work-hard",
}

func TestScrapyNovel(t *testing.T) {
	// novel.Scrapy(1, 300, titles[0])
	// novel.Scrapy(-1, -1, titles[1])
	// novel.Scrapy(501, 501, titles[2])
	// novel.Scrapy(1101, 1101, titles[3])
	// novel.Scrapy(1, 1, titles[4])
	// novel.Scrapy(916, 999, titles[5])
	novel.Scrapy(265, 400, titles[6])
	// novel.Scrapy(332, 433, titles[7])
	// novel.Scrapy(101, 300, titles[8])
	// novel.Scrapy(386, 386, titles[9])
	// novel.Scrapy(1, 1, titles[10])
	// novel.Scrapy(1, 1, titles[11])
}

func TestGenerateChapterList(t *testing.T) {
	outputDir := "./output/novels/"
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
		novelJsonPath := fmt.Sprintf("%s\\output\\novels\\novels.json", wd)
		saveFile := storage.Create(novelJsonPath)
		saveFile.WriteString(novelsJson)
		saveFile.Close()

		// fmt.Println("novelsJson", novelsJson)
	}
}

func TestTransformMp3(t *testing.T) {
	expectedNovelIndex := 6
	expectedChapterStartIndex := 101
	expectedChapterEndIndex := 300
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	outputDir := wd + "\\output\\novels\\"
	chapterPrefix := "ch-"
	chapterSuffix := ".txt"
	novelDirs, err := storage.ReadDir(outputDir)
	if err != nil {
		panic(err)
	}

	for novelIndex, novelDir := range novelDirs {
		if novelDir.Name() == "novels.json" || novelIndex != expectedNovelIndex {
			continue
		}

		novelRootDir := outputDir + "\\" + novelDir.Name()
		novelChapters, err := storage.ReadDir(novelRootDir)
		if err != nil {
			panic(err)
		}

		for novelChapterIndex, novelChapter := range novelChapters {
			if !strings.HasPrefix(novelChapter.Name(), chapterPrefix) ||
				!strings.HasSuffix(novelChapter.Name(), chapterSuffix) {
				continue
			}

			realChapterStr, _ := strings.CutPrefix(novelChapter.Name(), chapterPrefix)
			realChapterStr, _ = strings.CutSuffix(realChapterStr, chapterSuffix)
			realChapterIndex, _ := strconv.Atoi(realChapterStr)
			if realChapterIndex < expectedChapterStartIndex ||
				realChapterIndex > expectedChapterEndIndex {
				continue
			}

			chapterFile := fmt.Sprintf("ch-%d.txt", realChapterIndex)

			if novelChapterIndex == 0 {
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
			novel.TransformMp3(realChapterIndex-1, novelRootDir)

			endCost := time.Now().Unix()
			fmt.Printf("👆======================cost %ds=======================👆\n", (endCost - startCost))
		}
	}
}

func TestTransformMp4(t *testing.T) {
	expectedNovelIndex := 6
	expectedChapterStartIndex := 360
	expectedChapterEndIndex := 384
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	outputDir := wd + "\\output\\novels\\"
	chapterPrefix := "ch-"
	chapterSuffix := ".mp3"
	novelDirs, err := storage.ReadDir(outputDir)
	if err != nil {
		panic(err)
	}

	for novelIndex, novelDir := range novelDirs {
		if novelDir.Name() == "novels.json" || novelIndex != expectedNovelIndex {
			continue
		}

		novelRootDir := outputDir + "\\" + novelDir.Name()
		novelChapters, err := storage.ReadDir(novelRootDir)
		if err != nil {
			panic(err)
		}

		for novelChapterIndex, novelChapter := range novelChapters {
			if !strings.HasPrefix(novelChapter.Name(), chapterPrefix) ||
				!strings.HasSuffix(novelChapter.Name(), chapterSuffix) {
				continue
			}

			realChapterStr, _ := strings.CutPrefix(novelChapter.Name(), chapterPrefix)
			realChapterStr, _ = strings.CutSuffix(realChapterStr, chapterSuffix)
			realChapterIndex, _ := strconv.Atoi(realChapterStr)
			if realChapterIndex < expectedChapterStartIndex ||
				realChapterIndex > expectedChapterEndIndex {
				continue
			}

			chapterFile := fmt.Sprintf("ch-%d.txt", realChapterIndex)

			if novelChapterIndex == 0 {
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
			novel.TransformMp4(realChapterIndex-1, novelRootDir)

			endCost := time.Now().Unix()
			fmt.Printf("👆======================cost %ds=======================👆\n", (endCost - startCost))
		}
	}
}
