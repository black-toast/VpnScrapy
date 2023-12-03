
from net.request import request
from util.edge_tts import edgeTts
import os
import re
import time
import datetime
import json
import sys
from scrapy import Selector
from util.ffmpeg import makeImageVideo


class Crawl():

    _fileTitle = "title.txt"
    _fileAuthor = "author.txt"
    _fileDesc = "desc.txt"
    _fileChapterList = "chapter_list.txt"
    _fileChapter = "ch-%d.txt"

    def __init__(self):
        rootPath = os.getenv('novels_root')
        if rootPath is None:
            rootPath = os.getcwd()
        self.novelsOutputPath = f'{rootPath}{os.path.sep}output{os.path.sep}novels'
        self.novelDomainUrl = "https://novelhulk.com"
        self.novelPathUrl = self.novelDomainUrl + "/nb/"
        self.chapterListPath = "/ajax/chapter-archive?novelId="
        self.xpathTitle = "//*[@id=\"novel\"]/div[1]/div[1]/div[3]/h3/text()"
        self.xpathFirstLineTitle = "//*[@id=\"chr-content\"]/h4/text()"
        self.xpathNovelId = "//div[@id=\"rating\"]/@data-novel-id"
        self.xpathAuthor = "//div[contains(@class, \"desc\")]/ul/li[2]/a/text()"
        self.xpathDesc = "//*[@id=\"tab-description\"]/div/text()"
        self.xpathReadUrl = "//a[contains(@class, \"btn-read-now\")]/@href"

        self.xpathChapterListUrl = "//a/@href"
        self.xpathChapterListTitle = "//a/@title"
        self.xpathChapterTitle = "//a[@class=\"chr-title\"]/span/text()"
        self.xpathChapter = "//*[@id=\"chr-content\"]/p[not(@style)]/text()|//*[@id=\"chr-content\"]/text()|//div[@id=\"novelArticle2\"]/p[not(@style)]/text()"

    def request(self, method, url):
        response = request(method, url)
        response.encoding = 'utf-8'
        if response.status_code == 200:
            return response
        print("illegal status code", response.status_code, "\n")
        return None

    def save(self, path, content):
        with open(path, 'w', encoding='utf-8') as f:
            f.write(content)

    def reFindAll(self, searchStr, pattern):
        pattern = re.compile(pattern=pattern)
        searchResult = pattern.findall(string=searchStr)
        return searchResult

    def parseTitle(self, selector):
        return selector.xpath(self.xpathTitle).extract_first().strip("\n ")

    def parseNovelId(self, selector):
        return selector.xpath(self.xpathNovelId).extract_first()

    def parseAuthor(self, selector):
        return selector.xpath(self.xpathAuthor).extract_first().strip("\n ")

    def parseDesc(self, selector):
        return selector.xpath(self.xpathDesc).extract_first()

    def parseChapterList(self, selector):
        listUrlElems = selector.xpath(self.xpathChapterListUrl)
        listTitleElems = selector.xpath(self.xpathChapterListTitle)

        chapterList = []
        if len(listUrlElems) <= 0:
            return chapterList

        for title, url in zip(listTitleElems, listUrlElems):
            chapterList.append(
                {'title': title.extract(), 'url': url.extract()})
        return chapterList

    def parseChapterTitle(self, selector):
        chapterTitle = selector.xpath(
            self.xpathChapterTitle).extract_first().strip("\n ")
        findResults = self.reFindAll(chapterTitle, r'(\d{1,}).?-?:?')
        findLen = len(findResults)
        if findLen <= 0:
            return chapterTitle
        findResult = findResults[findLen - 1]
        findPos = chapterTitle.rfind(findResult) + len(findResult)
        return chapterTitle[findPos:].strip(": ").strip("- ")

    # 解析文章内容第一行中可能存在的章节标题
    def parseFirstLineChapterTitle(self, selector):
        content = selector.xpath(self.xpathFirstLineTitle).extract_first()
        if content == None:
            return ""
        content = content.strip("\n ")
        # extract chapter number
        # example: https://novelhulk.com/nb/spy-mage-system-book/cchapter-1
        findResults = self.reFindAll(content, r'(\d{1,}).?-?:?')
        findLen = len(findResults)
        if findLen <= 0:
            return content
        findResult = findResults[findLen - 1]
        findPos = content.rfind(findResult) + len(findResult)
        return content[findPos:].strip(": ").strip("- ")

    def removeSpecialChars(self, content):
        content = content.replace('\\\"', '')
        content = content.replace(
            'this content of novelfullbook.com, if you reading this content please go to website novelfullbook.com to continue reading, fastest update hourly', '')
        # remove ..... format
        # example: https://novelhulk.com/nb/spy-mage-system/cchapter-1
        findResults = self.reFindAll(content, r'^"?\.{2,}"?$')
        findLen = len(findResults)
        if findLen > 0:
            return ""

        # example: https://novelhulk.com/nb/spy-mage-system/cchapter-27
        if content.startswith("The source of this content is"):
            return ""
        if "I've created a discord server" in content:
            return ""

        if content.lower().startswith("translator:"):
            return ""

        findResults = self.reFindAll(content, r'(…{2,})')
        findLen = len(findResults)
        if findLen > 0:
            return ""

        findResults = self.reFindAll(content, r'^(-{3,})$')
        findLen = len(findResults)
        if findLen > 0:
            return ""

        findResults = self.reFindAll(content, r'(\*{3,})')
        findLen = len(findResults)
        if findLen > 0:
            return ""

        if content == "[…]" or content == "-" or content == "“…”"\
                or content == "__" or content == "–" or content == "—"\
                or content == "…..." or content == "…" or content == "[ ... ]":
            return ""
        
        lowerContent = content.lower()
        if "note:" in lowerContent or "notes:" in lowerContent or "author’s note:" in lowerContent:
            return ""
        if "endnote:" in lowerContent or "endnote" in lowerContent or "endnotes:" in lowerContent:
            return ""
        if "preview:" in lowerContent:
            return ""
        if "62d67767f92eb560e77c9100" in content:
            return ""

        webContent = content.replace("[", "").replace("]", "")\
            .replace("/", "").strip(" ")
        findResults = self.reFindAll(
            webContent, r'(\.[c|n|𝑪|𝓬|𝐂|𝑐|𝐜|𝕔|𝒸|𝗰|𝔠|𝚌|𝓒|𝒞][o|e|𝞸|𝚘|𝑶|𝓞|𝔬|𝐎|𝗈|𝒐|𝑂|𝑜|𝒪][m|t|𝓜|𝑚|𝚖|𝔪|𝓶|𝑀|𝗆|𝕞|𝓂|𝐦|𝐌])')
        findLen = len(findResults)
        if findLen > 0:
            return ""
        return content

    def isEndLine(self, content):
        # lowerContent = content.lower()
        # if "note:" in lowerContent or "notes:" in lowerContent or "author’s note:" in lowerContent:
        #     return True
        # if "endnote:" in lowerContent or "endnote" in lowerContent or "endnotes:" in lowerContent:
        #     return True
        # if "preview:" in lowerContent:
        #     return True
        return False

    def parseChapter(self, selector):
        chapterElems = selector.xpath(self.xpathChapter)
        if len(chapterElems) <= 0:
            return ""
        chapter = ""
        for index, elem in enumerate(chapterElems):
            lineContent = elem.extract().strip("\n ")
            line = self.removeSpecialChars(lineContent).strip("\n ")
            if self.isEndLine(line):
                break
            if line == "":
                continue
            chapter += line
            if index+1 != len(chapterElems):
                chapter += "\n"
        return chapter

    def saveChapter(self, novelDir, realIndex, chapterTitle, content):
        content = 'Chapter %d: %s\n%s' % (realIndex, chapterTitle, content)
        saveChapterFileName = self._fileChapter % realIndex
        self.save(novelDir + saveChapterFileName, content)
        print("save chapter complete")

    def generateChapterMp3(self, novelDir, realIndex):
        ttsInput = '%sch-%d.txt' % (novelDir, realIndex)
        ttsMp3Output = '%sch-%d.mp3' % (novelDir, realIndex)
        # 保证执行edge-tts命令不会提示该文件已存在
        if os.path.exists(ttsMp3Output):
            os.remove(ttsMp3Output)
        edgeTts(ttsInput, ttsMp3Output)
        print("transform mp3 complete")

    def generateChapterMp4(self, novelDir, realIndex):
        videoImage = '%scover.jpg' % novelDir
        ttsMp3Output = '%sch-%d.mp3' % (novelDir, realIndex)
        ttsMp4Output = '%sch-%d.mp4' % (novelDir, realIndex)
        # 保证执行ffmpeg命令不会提示该文件已存在
        if os.path.exists(ttsMp4Output):
            os.remove(ttsMp4Output)
        makeImageVideo(videoImage, ttsMp3Output, ttsMp4Output)
        print("transform mp4 complete")

    def requestNovelChapterList(self, startChapter, endChapter, novelDir,
                                novelChapterUrl, makeMp3, makeMp4):
        # request novel chapter list
        print("wait 5s, then request novel chapter list")
        time.sleep(5)
        response = self.request("GET", novelChapterUrl)
        if response == None:
            return

        selector = Selector(response)
        chapterList = self.parseChapterList(selector)
        if len(chapterList) <= 0:
            print('chapter list is empty.')
            return
        print("chapter list number is %d" % len(chapterList))
        self.save(novelDir + self._fileChapterList, json.dumps(chapterList))
        if startChapter <= -1:
            startChapter = 1

        if endChapter <= -1:
            endChapter = len(chapterList)

        # request novel chapter
        for index, chapter in enumerate(chapterList):
            if index+1 < startChapter:
                continue
            if index+1 > endChapter:
                break
            if index == 0:
                print("wait 2s, then request novel chapter(%s)" %
                      chapter['title'])
            else:
                print("wait 2s, then request next novel chapter(%s)" %
                      chapter['title'])
            time.sleep(2)

            currentTime = datetime.datetime.now()
            print("current time: %s" % currentTime)

            response = self.request("GET", chapter['url'])
            if response == None:
                break
            selector = Selector(response)
            chapterTitle = self.parseChapterTitle(selector)
            # 解析第一行标题
            if chapterTitle == "":
                chapterTitle = self.parseFirstLineChapterTitle(selector)
            chapter = self.parseChapter(selector)

            self.saveChapter(novelDir, index + 1, chapterTitle, chapter)
            if makeMp3:
                self.generateChapterMp3(novelDir, index + 1)
            if makeMp4:
                self.generateChapterMp4(novelDir, index + 1)

            if makeMp4:
                ttsMp3Output = f'{novelDir}ch-{index + 1}.mp3'
                if os.path.exists(ttsMp3Output):
                    os.remove(ttsMp3Output)

            endCost = datetime.datetime.now().timestamp()
            print("👆======================cost %.2fs=======================👆" %
                  (endCost - currentTime.timestamp()))

    def scrapy(self, path, startChapter, endChapter, makeMp3, makeMp4):
        if startChapter == 0 or endChapter == 0:
            print(f"skip novel path is {path}")
            return
        if startChapter < -1 or endChapter < -1:
            print(
                f"start or end chapter is illegal chapter and path is {path}")
            return

        # request novel introduction
        print("request novel introduction")
        novelUrl = self.novelPathUrl + path
        response = self.request("GET", novelUrl)
        if response == None:
            return

        # parse html
        selector = Selector(response)
        title = self.parseTitle(selector)
        novelId = self.parseNovelId(selector)

        titlePath = title.replace(" ", "").replace(":", "").replace("!", "")
        novelDir = '%s%s%s%s' % (
            self.novelsOutputPath, os.path.sep, titlePath, os.path.sep)
        if not os.path.exists(novelDir):
            os.makedirs(novelDir)

        # save novel introduction
        saveNovelIntroductSwitch = False
        if saveNovelIntroductSwitch:
            author = self.parseAuthor(selector)
            desc = self.parseDesc(selector)
            self.save('%s%s' % (novelDir, self._fileTitle), title)
            self.save('%s%s' % (novelDir, self._fileAuthor), author)
            self.save('%s%s' % (novelDir, self._fileDesc), desc)

        novelChapterUrl = '%s%s%s' % (
            self.novelDomainUrl, self.chapterListPath, novelId)
        self.requestNovelChapterList(startChapter, endChapter,  novelDir,
                                     novelChapterUrl, makeMp3, makeMp4)

    def offlineMake(self, path, startChapter, endChapter, makeMp3, makeMp4):
        if startChapter == 0 or endChapter == 0:
            print(f"skip novel path is {path}")
            return
        if startChapter < -1 or endChapter < -1:
            print(
                f"start or end chapter is illegal chapter and path is {path}")
            return
        if makeMp3 == False and makeMp4 == False:
            print(f"novel({path}) don't offline make mp3 or mp4")
            return

        # request novel introduction
        print("request novel introduction")
        novelUrl = self.novelPathUrl + path
        response = self.request("GET", novelUrl)
        if response == None:
            return

        # parse html
        selector = Selector(response)
        title = self.parseTitle(selector)

        titlePath = title.replace(" ", "").replace(":", "").replace("!", "")
        novelDir = '%s%s%s%s' % (
            self.novelsOutputPath, os.path.sep, titlePath, os.path.sep)
        print(f"novel({title}) offline make mp3 or mp4")

        for file in os.listdir(novelDir):
            if makeMp3 == True and file.startswith("ch-") and file.endswith(".txt"):
                chapterNum = int(file.replace("ch-", "").replace(".txt", ""))
                if chapterNum >= startChapter and chapterNum <= endChapter:
                    currentTime = datetime.datetime.now()
                    print(f"current time: {currentTime} and make chapter {chapterNum} mp3")
                    self.generateChapterMp3(novelDir, chapterNum)
                    endCost = datetime.datetime.now().timestamp()
                    print("👆======================cost %.2fs=======================👆" %
                        (endCost - currentTime.timestamp()))

        for file in os.listdir(novelDir):
            if makeMp4 == True and file.startswith("ch-") and file.endswith(".mp3"):
                chapterNum = int(file.replace("ch-", "").replace(".mp3", ""))
                if chapterNum >= startChapter and chapterNum <= endChapter:
                    currentTime = datetime.datetime.now()
                    print(f"current time: {currentTime} and make chapter {chapterNum} mp4")
                    self.generateChapterMp4(novelDir, chapterNum)

                    ttsMp3Output = f'{novelDir}ch-{chapterNum}.mp3'
                    if os.path.exists(ttsMp3Output):
                        os.remove(ttsMp3Output)
                        
                    endCost = datetime.datetime.now().timestamp()
                    print("👆======================cost %.2fs=======================👆" %
                        (endCost - currentTime.timestamp()))

    def generateNovelsJson(self):
        print("generate novels json")

        novelsJson = []
        novelPathList = os.listdir(self.novelsOutputPath)
        if "novels.json" in novelPathList:
            novelPathList.remove("novels.json")

        for index, novelFile in enumerate(novelPathList):
            novelDir = f"{self.novelsOutputPath}{os.path.sep}{novelFile}{os.path.sep}"
            with open(f'{novelDir}title.txt', 'r', encoding='utf-8') as f:
                novelName = f.read()

            with open(f'{novelDir}desc.txt', 'r', encoding='utf-8') as f:
                novelDesc = f.read()
            
            novel = {"index": index + 1, "name": novelName, "chapters": []}

            for chapterFile in os.listdir(novelDir):
                if not chapterFile.startswith("ch-") or not chapterFile.endswith(".txt"):
                    continue

                f = open(f'{novelDir}{chapterFile}', 'r', encoding='utf-8')
                firstLineTitle = f.readline()
                firstLineTitleSplit = firstLineTitle.split(": ")
                chapterIndex = int(firstLineTitleSplit[0].replace("Chapter ", ""))
                chapterTitle = firstLineTitleSplit[1]

                title = f'{novelName} CH-{chapterIndex}'
                desc = f'NOVEL INFO: {novelName}\nChapter {chapterIndex}: {chapterTitle}Description: {novelDesc}'
                novel["chapters"].append({"index":chapterIndex, "title": title, "desc": desc})

            novelsJson.append(novel)

        self.save(f'{self.novelsOutputPath}{os.path.sep}novels.json', json.dumps(novelsJson))