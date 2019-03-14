package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"
)

var pipeFile = "/tmp/wspipe.tmp"

func fifo() {

	os.Remove(pipeFile)
	err := syscall.Mkfifo(pipeFile, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	go writeMsg()
	defer func() {
		fmt.Printf("closing fifo")
	}()
	fmt.Println("open a named pipe file for read.")
	file, err := os.OpenFile(pipeFile, os.O_CREATE|os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error Named Pipe")

			continue
		}

		fmt.Printf("Received from pipe: %s", line)
		raw := new(rawmsg)
		err = json.Unmarshal(line, raw)
		if err != nil {
			fmt.Printf("failed to decode motion\n")

		}
		if raw.Value.(string) == "on" {
			fmt.Println("starting recorder")
			//		fmt.Printf("starting record: %s\n", time.Now().Format(time.RFC3339Nano))
			go record()
		}
		if raw.Value.(string) == "off" {
			fmt.Println("sending record cancel")
			cancelRecordChan <- true
		}

		msg := new(message)
		msg.Type = "motion"
		msg.Msg = line
		//	fmt.Println("sending msg to writemsg")
		messageChan <- *msg
		//	fmt.Println("sent msg to writemsg")

	}
}

func writeMsg() {

	defer func() {
		fmt.Printf("closing writeMsg\n")
	}()
	for {
		select {

		// drain the channel
		case msg := <-messageChan:
			if c != nil {
				if online {
					fmt.Println("sending msg to server")

					err := c.WriteJSON(msg)
					fmt.Println("sent msg to server")

					if err != nil {
						fmt.Printf("Error while sending: %s ", err.Error())
					}
				}
			}
		case <-cancelwriteChan:
			return
		}
	}
}
