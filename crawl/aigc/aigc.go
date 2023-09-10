package aigc

import (
	"VpnScrapy/bean"
	"VpnScrapy/http"
	"VpnScrapy/storage"
	"VpnScrapy/util"
	"fmt"
	"os"
	"strings"
	"time"
)

func Scrapy() {
	aigcPath := getAigcPath()
	fmt.Println("aigcPath:", aigcPath)
	api := "https://api.esheep.com/gateway/post/detail?id="
	for id := 1; id <= 23683; id++ {
		fmt.Printf("request aigc(id=%d)\n", id)

		url := fmt.Sprintf("%s%d", api, id)
		content, error := Request(url, "GET")
		if error != nil {
			panic(error)
		}

		// 解析详情json
		result, err := util.Parse(content, bean.EsheepResp{})
		if err != nil {
			fmt.Println("parse aigc detail data failure, err=", err)
			return
		}

		if result.Code != 0 {
			fmt.Println("detail json code isn't 0")
			continue
		}
		if len(result.Data.Images) <= 0 {
			fmt.Println("detail json images is empty")
			continue
		}

		// 创建id目录
		aigcIdPath := fmt.Sprintf("%s%s%d", aigcPath, string(os.PathSeparator), id)
		storage.Mkdir(aigcIdPath)

		// 保存详情
		aigcIdDetailPath := fmt.Sprintf("%s%sdetail.json", aigcIdPath, string(os.PathSeparator))
		saveFile := storage.Create(aigcIdDetailPath)
		saveFile.WriteString(string(content))
		saveFile.Close()

		for imgIndex, image := range result.Data.Images {
			imageByte, error := Request(image.Url, "GET")
			if error != nil {
				panic(error)
			}

			imgFormat := "png"
			if strings.Contains(image.Url, ".jpg") {
				imgFormat = "jpg"
			}

			aigcIdImgPath := fmt.Sprintf("%s%s%d.%s", aigcIdPath,
				string(os.PathSeparator), imgIndex, imgFormat)
			storage.WriteFile(aigcIdImgPath, imageByte)
		}

		fmt.Println("wait 5s, then request next aigc")
		time.Sleep(5 * time.Second)

		fmt.Println("===============================")
	}
}

func getAigcPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	wd = strings.ReplaceAll(wd, "\\crawl\\aigc", "")
	return fmt.Sprintf("%s%s%s%s%s", wd, string(os.PathSeparator),
		"output", string(os.PathSeparator), "aigc")
}

func generateBuvid() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s%sinfoc", "A09D5C17", "B7D8",
		"1FBD", "9A46", "9826D2515C01", "30312")
}

func Request(url string, method string) ([]byte, error) {
	novelContent, err := http.Request(&http.RequestConfig{
		Method:    method,
		Url:       url,
		Transport: http.V2rayProxy(),
		Headers:   map[string]string{"Cookie": "buvid=" + generateBuvid()},
	})
	if err != nil {
		fmt.Printf("request %s fail.\n %v\n", url, err)
		return nil, err
	}

	return novelContent, nil
}
