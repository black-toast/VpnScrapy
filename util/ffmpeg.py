import subprocess

ffmpeg_command_placeholder = "ffmpeg -loop 1 -i %s -i %s -c:v libx264 -c:a aac -b:a 192k -shortest %s"


def makeImageVideo(videoImage, ttsMp3Output, ttsMp4Output):
    ffmpeg_command = ffmpeg_command_placeholder % (
        videoImage, ttsMp3Output, ttsMp4Output)
    result = subprocess.run(ffmpeg_command, shell=True, capture_output=True)
    if result.returncode != 0:
        print("FFmpeg command failed!")
        print("Error Output:", result.stderr.decode('utf-8'))
