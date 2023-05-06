import asyncio
import time
import sys
import edge_tts

tts_path = ""
voice_path = ""
VOICE = "en-US-SteffanNeural"

def readTranText():
    file = open(tts_path, 'r', encoding='utf-8')  # 打开文件
    file_data = file.readlines()  # 读取所有行

    file_content = ""
    for line in file_data:
        file_content = file_content + line
    return file_content

async def tran_voice(tran_text) -> None:
    # communicate = edge_tts.Communicate(tran_text, VOICE)
    communicate = edge_tts.Communicate(tran_text, VOICE, proxy="http://127.0.0.1:10809")
    await communicate.save(voice_path)

if __name__ == '__main__':
    argv_len = len(sys.argv)
    if argv_len != 3:
        print("Expected argv length to be equal to 3")
    else:
        tts_path = sys.argv[1]
        voice_path = sys.argv[2]

        # edge-tts
        tts_content = str(readTranText())
        asyncio.run(tran_voice(str(tts_content)))
