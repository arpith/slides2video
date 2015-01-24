# slides2video
Go script to create a video from slides and audio using ffmpeg. 

## Usage
Create a file called slideDurations.txt with each line looking like "imagePath imageDuration" (no quotes). 

This script creates a video for each slide[1] then concatenates them[2] and finally adds the audio[3] to create finalVideo.mp4

[1] http://trac.ffmpeg.org/wiki/Create%20a%20video%20slideshow%20from%20images

[2] http://trac.ffmpeg.org/wiki/Concatenate

[3] http://stackoverflow.com/questions/11779490/ffmpeg-how-to-add-new-audio-not-mixing-in-video
