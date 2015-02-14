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
import "os"

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
		log.Printf("Creating " + outputName)
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created " + outputName)
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
		log.Printf("Creating " + outputName + " using " + listFilename)
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created " + outputName)
	}
	done <- true
}

func addAudio(silentFilename string, audioFilename string, outputName string, done chan bool) {
	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-i", silentFilename, "-i", audioFilename, "-map", "0:0", "-map", "1:0", "-codec", "copy", "-shortest", outputName)
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Creating " + outputName + " from " + silentFilename + " and " + audioFilename)
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	} else {
		log.Printf("Created " + outputName)
	}
	done <- true
}

func main() {
	timestampsFilenamePtr := flag.String("t", "timestamps.txt", "a file with timestamps and image names")
	audioFilenamePtr := flag.String("a", "audio.mp3", "audio file name")
	outputFilenamePtr := flag.String("o", "finalOut.mp4", "output file name")
	flag.Parse()

	videoListFilename := "videoList.txt"
	silentFilename := "silent.mp4"

	dat, err := ioutil.ReadFile(*timestampsFilenamePtr)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(dat), "\n")
	outLines := make([]string, len(lines))

	nextTimestamp := 0
	timestamp := 0

	for i, line := range lines {
		if line != "" {
			wg.Add(1)

			imgName := strings.Split(lines[i], " ")[1]
			timestamp = nextTimestamp
			var imgDuration float64
			if (i == len(lines)-1) || (lines[i+1] == "") {
				imgDuration = float64(timestamp) / 1000.0
			} else {
				nextLineSplit := strings.Split(lines[i+1], " ")
				nextTimestamp, err := strconv.Atoi(nextLineSplit[0])
				if err != nil {
					log.Fatal(err)
				}
				imgDuration = float64(nextTimestamp - timestamp) / 1000.0
			}
			imgDurationString := strconv.FormatFloat(imgDuration, 'f', 3, 64)
			outputName := "out" + strconv.Itoa(i+1) + ".mp4"
			outLines[i] = "file '" + outputName + "'"

			go img2video(imgName, imgDurationString, outputName)
		}
	}

	outData := []byte(strings.Join(outLines, "\n"))
	err = ioutil.WriteFile(videoListFilename, outData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	done := make(chan bool, 1)

	go concatVideos(len(lines), videoListFilename, silentFilename, done)

	<-done

	go addAudio(silentFilename, *audioFilenamePtr, *outputFilenamePtr, done)

	for i, line := range lines {
		if line != "" {
			outputFilename := "out" + strconv.Itoa(i+1) + ".mp4"
			err = os.Remove(outputFilename)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Deleted " + outputFilename)
			}
		}
	}

	<-done
}
