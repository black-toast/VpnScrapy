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

class WuxiaCrawl():
    
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
        self.novelPathUrl = "https://wuxia.click/chapter/%s-%d"

        self.xpathChapterTitle = "//h1[contains(@class,\"mantine-Title-root\")]/text()"
        self.xpathChapter = "//div[@id=\"chapterText\"]/text()"

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

        if content.lower().startswith("translator:") or content.lower().startswith("editor:"):
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
    
    def parseChapter(self, selector):
        chapterElems = selector.xpath(self.xpathChapter)
        if len(chapterElems) <= 0:
            return ""
        chapter = ""
        for index, elem in enumerate(chapterElems):
            lineContent = elem.extract().strip("\n ")
            line = self.removeSpecialChars(lineContent).strip("\n ")
            if index == 0:
                continue
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
        ttsInput = '%sch-%s.txt' % (novelDir, realIndex)
        ttsMp3Output = '%sch-%s.mp3' % (novelDir, realIndex)
        # 保证执行edge-tts命令不会提示该文件已存在
        if os.path.exists(ttsMp3Output):
            os.remove(ttsMp3Output)
        edgeTts(ttsInput, ttsMp3Output)
        print("transform mp3 complete")

    def generateChapterMp4(self, novelDir, realIndex):
        videoImage = '%scover.jpg' % novelDir
        ttsMp3Output = '%sch-%s.mp3' % (novelDir, realIndex)
        ttsMp4Output = '%sch-%s.mp4' % (novelDir, realIndex)
        # 保证执行ffmpeg命令不会提示该文件已存在
        if os.path.exists(ttsMp4Output):
            os.remove(ttsMp4Output)
        makeImageVideo(videoImage, ttsMp3Output, ttsMp4Output)
        print("transform mp4 complete")

    def requestNovelChapterList(self, startChapter, endChapter, novelDir, path,
                                novelChapterUrl, makeMp3, makeMp4):
        # request novel chapter list
        print("wait 5s, then request novel chapter list")
        time.sleep(5)

        # request novel chapter
        for index in range(startChapter, endChapter + 1):
            if index == startChapter:
                print("wait 2s, then request novel chapter(ch-%d)" % index)
            else:
                print("wait 2s, then request next novel chapter(ch-%s)" % index)
            time.sleep(2)

            currentTime = datetime.datetime.now()
            print("current time: %s" % currentTime)

            response = self.request("GET", novelChapterUrl % (path, index))
            if response == None:
                break
            selector = Selector(response)
            chapterTitle = self.parseChapterTitle(selector)
            chapter = self.parseChapter(selector)
            # print("chapterTitle:", chapterTitle)
            # print("chapter:", chapter)

            self.saveChapter(novelDir, index, chapterTitle, chapter)
            if makeMp3:
                self.generateChapterMp3(novelDir, f'{index}')
            if makeMp4:
                self.generateChapterMp4(novelDir, f'{index}')

            if makeMp4:
                ttsMp3Output = f'{novelDir}ch-{index}.mp3'
                if os.path.exists(ttsMp3Output):
                    os.remove(ttsMp3Output)

            endCost = datetime.datetime.now().timestamp()
            print("👆======================cost %.2fs=======================👆" %
                  (endCost - currentTime.timestamp()))
            
    def scrapy(self, title, path, startChapter, endChapter, makeMp3, makeMp4):
        if startChapter == 0 or endChapter == 0:
            print(f"skip novel path is {path}")
            return
        if startChapter < -1 or endChapter < -1:
            print(f"start or end chapter is illegal chapter and path is {path}")
            return

        titlePath = title.replace(" ", "").replace(":", "").replace("!", "")
        novelDir = '%s%s%s%s' % (
            self.novelsOutputPath, os.path.sep, titlePath, os.path.sep)
        if not os.path.exists(novelDir):
            os.makedirs(novelDir)

        self.requestNovelChapterList(startChapter, endChapter,  novelDir, path,
                                     self.novelPathUrl, makeMp3, makeMp4)
        
    def offlineMake(self, title, path, startChapter, endChapter, makeMp3, makeMp4):
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

        titlePath = title.replace(" ", "").replace(":", "").replace("!", "")
        novelDir = '%s%s%s%s' % (
            self.novelsOutputPath, os.path.sep, titlePath, os.path.sep)
        print(f"novel({title}) offline make mp3 or mp4")

        for file in os.listdir(novelDir):
            if makeMp3 == True and file.startswith("ch-") and file.endswith(".txt"):
                chapterNum = file.replace("ch-", "").replace(".txt", "")
                chapterRealNum = chapterNum
                chapterRealNumSplit = chapterRealNum.split("-")
                if len(chapterRealNumSplit) > 1:
                    chapterRealNum = int(chapterRealNumSplit[0])
                else:
                    chapterRealNum = int(chapterNum)
                if chapterRealNum >= startChapter and chapterRealNum <= endChapter:
                    currentTime = datetime.datetime.now()
                    print(f"current time: {currentTime} and make chapter {chapterNum} mp3")
                    self.generateChapterMp3(novelDir, f'{chapterNum}')
                    endCost = datetime.datetime.now().timestamp()
                    print("👆======================cost %.2fs=======================👆" %
                        (endCost - currentTime.timestamp()))

        for file in os.listdir(novelDir):
            if makeMp4 == True and file.startswith("ch-") and file.endswith(".mp3"):
                chapterNum = file.replace("ch-", "").replace(".mp3", "")
                chapterRealNum = chapterNum
                chapterRealNumSplit = chapterRealNum.split("-")
                if len(chapterRealNumSplit) > 1:
                    chapterRealNum = int(chapterRealNumSplit[0])
                else:
                    chapterRealNum = int(chapterNum)
                if chapterRealNum >= startChapter and chapterRealNum <= endChapter:
                    currentTime = datetime.datetime.now()
                    print(f"current time: {currentTime} and make chapter {chapterNum} mp4")
                    self.generateChapterMp4(novelDir, f'{chapterNum}')

                    ttsMp3Output = f'{novelDir}ch-{chapterNum}.mp3'
                    if os.path.exists(ttsMp3Output):
                        os.remove(ttsMp3Output)
                        
                    endCost = datetime.datetime.now().timestamp()
                    print("👆======================cost %.2fs=======================👆" %
                        (endCost - currentTime.timestamp()))