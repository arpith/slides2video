package main

import "io/ioutil"
import "log"
import "os/exec"
import "flag"
import "strings"
import "sync"
import "strconv"
import "fmt"
import "bytes"

var wg sync.WaitGroup

func img2video(imgName string, imgDuration string, outputName string) {
	defer wg.Done()
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", imgName, "-t", imgDuration, "-pix_fmt", "yuv420p", outputName)
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating "+outputName)
	}
	err = cmd.Wait()
	if err != nil {
    		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created "+outputName)
	}
}

func concatVideos(numberOfVideos int, listFilename string, outputName string, done chan bool) {
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-f", "concat", "-i", listFilename, "-c", "copy", outputName)
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating "+outputName+" using "+listFilename)
	}
	err = cmd.Wait()
	if err != nil {
    		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created "+outputName)
	}
	done <- true
}

func addAudio(silentFilename string, audioFilename string, outputName string, done chan bool) {
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-i", silentFilename, "-i", audioFilename, "-map", "0", "-map", "1", "-codec" ,"copy", "-shortest", outputName)
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating "+outputName+" from "+silentFilename+" and "+audioFilename)
	}
	err = cmd.Wait()
	if err != nil {
    		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created "+outputName)
	}
	done <- true
}

func main() {
	dPtr := flag.String("d", "slideDurations.txt", "a file with image names and durations")
	audioFilenamePtr := flag.String("a", "audio.mp3", "audio file name")
	outputFilenamePtr := flag.String("o", "finalOut.mp4", "output file name")
	flag.Parse()

	videoListFilename := "videoList.txt"
	silentFilename := "silent.mp4"

	dat, err := ioutil.ReadFile(*dPtr)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(dat), "\n")
	outLines := make([]string,len(lines))

	for i,line := range lines {
		if line != "" {
			wg.Add(1)
	
			splitLine := strings.Split(line," ")
			imgName := splitLine[0]
			imgDuration := splitLine[1]
			outputName := "out"+strconv.Itoa(i+1)+".mp4"
			outLines[i] = "file '"+outputName+"'"

			go img2video(imgName, imgDuration, outputName)
		}
	}
	
	outData := []byte(strings.Join(outLines,"\n"))
	err = ioutil.WriteFile(videoListFilename, outData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	done := make(chan bool, 1)
	
	go concatVideos(len(lines), videoListFilename, silentFilename, done)
	
	<-done
	
	go addAudio(silentFilename, *audioFilenamePtr, *outputFilenamePtr, done)

	<-done
}


