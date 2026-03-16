package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	var b bytes.Buffer
	cmdOut := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmdOut.Stdout = &b
	cmdOut.Stderr = os.Stderr

	err := cmdOut.Run()
	if err != nil {
		log.Printf("Error running ffprobe command: %s", err)
		return "", err
	}


	type AspectsStruct struct {
		Streams []struct {
			Width int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	var data AspectsStruct	
	if err := json.Unmarshal(b.Bytes(), &data); err != nil {
		log.Print("Error unmarshalling to json")
		return "", err
	}

	ratio := float64(data.Streams[0].Width) / float64(data.Streams[0].Height)
	if ratio >= 1.77 && ratio <= 1.78 {
		return "16:9", nil
	} else if ratio >= 0.56 && ratio<= 0.57 {
		return "9:16", nil
	} else {
		return "other", nil
	}
}

func processVideoForFastStart(filepath string) (string, error) {
	outPath := filepath + ".processing"

	cmdOut := exec.Command("ffmpeg", "-i", filepath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outPath)

	err := cmdOut.Run()
	if err != nil {
		log.Printf("Error running ffmpeg command: %s", err)
		return "", err
	}

	return outPath, nil
}

