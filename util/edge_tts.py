import asyncio
import os
import re
import edge_tts
import sys
import time


def readTranText(ttsPath):
    file = open(ttsPath, 'r', encoding='utf-8')  # 打开文件
    file_data = file.readlines()  # 读取所有行

    file_content = ""
    for line in file_data:
        file_content = file_content + line
    file.close()
    return file_content


async def tran_voice(voicePath, tran_text) -> None:
    proxy = None
    if sys.platform.startswith('win'):
        proxy = "http://127.0.0.1:10809"
    communicate = edge_tts.Communicate(
        tran_text, "en-US-SteffanNeural", proxy=proxy)
    await communicate.save(voicePath)


def edgeTts(ttsPath, voicePath):
    lastSepIndex = voicePath.rfind(os.path.sep)
    tmpPath = voicePath[0: lastSepIndex + 1] + "tmp.mp3"
    tranTexts = readTranText(ttsPath).split("Chapter ")

    with open(voicePath, "ab+") as audio:
        for index, tranText in enumerate(tranTexts):
            if index == 0:
                continue
            tts_content = 'Chapter %s' % tranText
            # loop = asyncio.get_event_loop()
            # loop.run_until_complete(tran_voice(tmpPath, str(tts_content)))
            asyncio.run(tran_voice(tmpPath, str(tts_content)))

            with open(tmpPath,'rb') as tmpFile:
                audio.write(tmpFile.read())

            time.sleep(60)
    
    # delte tmpPath file
    os.remove(tmpPath)

