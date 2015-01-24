package main

import (
	"ioutil",
	"log",
	"os/exec"
)

func main() {

	dat, err := ioutil.ReadFile("slideDurations.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(dat), "\n")
	outLines := make([]string,len(lines))

	for i,line := range lines {
		imgName := strings.Split(line," ")[0]
		imgDuration := strings.Split(line," ")[1]
		outLines[i] = "file 'out"+i".mp4'"

		cmd := exec.Command("bash", "-c", "ffmpeg -loop 1 -i "+imgName+" -c:v libx264 -t "+imgDuration+" -pix_fmt yuv420p out"+i+"+.mp4")
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Waiting to create out"+i+".mp4 from "+imgName)
		err = cmd.Wait()
		log.Printf("Created out"+i+".mp4 from "+imgName)
	}
	
	outData := []byte(strings.Join(outLines,"\n"))
	err := ioutil.WriteFile("videoList.txt", outData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("bash", "-c", "ffmpeg -f concat -i videoList.txt -c copy withoutAudio.mp4")

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Waiting to create withoutAudio.mp4 using concat")
	err = cmd.Wait()
	log.Printf("Created withoutAudio.mp4")

	cmd := exec.Command("bash", "-c", "ffmpeg -i withoutAudio.mp4 -i audio.mp3 -map 0 -map 1 -codec copy -shortest finalVideo.mp4")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Waiting to create finalVideo.mp4 from withoutAudio.mp4 and audio.mp3")
	err = cmd.Wait()
	log.Printf("Created finalVideo.mp4")
	
}


