package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pions/webrtc/pkg/media"
)

func savegob(rtspchan chan media.RTCSample, savectrlchan chan string, devicename string) {
	var path = "static/clips/" + devicename
	var file *os.File
	var enc *gob.Encoder
	var state = "stopped"
	defer func() {
		file.Close()
	}()

	for {
		select {
		case gob := <-rtspchan:
			if state == "started" {
				enc.Encode(gob)
			}
		case ctrl := <-savectrlchan:
			switch ctrl {

			case "stop":
				state = "stopped"

			case "start":
				file, err := os.Create(path + "/" + devicename + time.Now().Format(time.RFC3339) + ".gob")
				if err != nil {
					fmt.Printf("Could not create file: %s", err.Error())
					return
				}
				//			fmt.Println("Received start and created file for ")
				enc = gob.NewEncoder(file)
				state = "started"
			case "destroy":
				return

			}
		}
	}

}

func playgob(playctrlchan chan string) {
	var state = "stopped"
	var file *os.File
	var path = "static/clips/"
	var decoder *gob.Decoder
	defer func() {
		file.Close()
	}()
	var err error
	ticker := time.NewTicker(time.Millisecond * 25)

	//	var cnt = 0
	for {
		select {
		case filename := <-namechan:
			if state == "playing" {
				fmt.Println("Canceled Play")
				file.Close()
				state = "stopped"
			}
			file, err = os.Open(path + filename)
			if err != nil {
				fmt.Printf("Could not open file: %s", err.Error())
				continue
			}
			decoder = gob.NewDecoder(file)

			fmt.Println("Changed state to playing for  file: " + filename)
			state = "playing"
		case ctrl := <-playctrlchan:
			if ctrl == "stop" {
				if state == "playing" {
					fmt.Println("Canceled Play")
					file.Close()
					state = "stopped"

				}
			}
		case <-ticker.C:
			if state == "playing" {
				goblet := new(media.RTCSample)

				err = decoder.Decode(goblet)
				if err != nil {
					if err.Error() == "EOF" {
						file.Seek(0, 0)
						fmt.Println("Looping")
						decoder = gob.NewDecoder(file)

					} else {
						fmt.Printf("Decode GOB Error: %s", err.Error())
						continue
					}

				}

				vp8Track.Samples <- *goblet
			}
		}
	}

}

func IOReadDir(root string) ([]string, error) {
	var files []string
	var path = "static/clips/"
	fileInfo, err := ioutil.ReadDir(path + root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
