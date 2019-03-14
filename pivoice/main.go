/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

type device struct {
	Type  string
	Name  string
	Desc  string
	State string
}

type rawmsg struct {
	Type  string
	Value interface{}
}

type message struct {
	Type string
	Msg  json.RawMessage
}

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func main() {
	angelos := flag.String("name", "angelos10", " -name The name of this Angelos")
	devicedesc := flag.String("desc", "angelos Ten", " -desc The description of this Angelos")

	wsServer := flag.String("server", "192.168.50.13:8080", "-server The webserver to contact")
	audioPath := flag.String("apath", "/home/pi/audio", "-apath pathtoaudio")

	flag.Parse()
	speech := Speech{Folder: *audioPath, Language: "en"}
	//	origin := "http://" + *angelos + ".local/"

	//	sig := make(chan os.Signal, 1)
	//	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
Start:

	log.SetFlags(0)

	//	interrupt := make(chan os.Signal, 1)
	//	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *wsServer, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {

		speech.Speak("I am sorry but I cannot connect to the server.")

		//	log. Fatal(err)
		log.Printf("dial: %s", err.Error())

		time.Sleep(10 * time.Second)
		goto Start

	}
	defer c.Close()

	//done := make(chan struct{})

	dv := device{Type: "hmi", Name: *angelos, Desc: *devicedesc, State: "online"}
	dev, _ := json.Marshal(dv)
	status := message{Type: "status", Msg: dev}
	err = c.WriteJSON(&status)

	//	err = websocket.JSON.Send(ws, &status)
	if err != nil {

		speech.Speak("I am sorry but I cannot send to the server.")

		time.Sleep(10 * time.Second)

		goto Start

	}

	msg := new(message)
	for {
		log.Println("waiting for message from server")
		err = c.ReadJSON(&msg)
		//		fmt.Printf("%v", msg)
		//	err = websocket.JSON.Receive(ws, &rcvbuf)
		if err != nil {

			speech.Speak("I am sorry but I cannot receive from the server.")

			time.Sleep(10 * time.Second)

			goto Start
		}
		switch msg.Type {
		case "control":
			var ctrl = new(rawmsg)

			err := json.Unmarshal(msg.Msg, &ctrl)
			if err != nil {
				fmt.Printf("failed to decode control: %s\n", err.Error())

			}

			switch ctrl.Type {

			case "sound":
				if ctrl.Value.(string) == "off" {
					speech.Speak("Sound will be set to " + ctrl.Value.(string))
				}
				cmd := exec.Command("amixer", "set", "'Master'", ctrl.Value.(string))
				cmd.Run()
				if ctrl.Value == "on" {
					speech.Speak("Sound will be set to " + ctrl.Value.(string))
				}

			case "volume":
				cmd := exec.Command("amixer", "set", "'Master'", ctrl.Value.(string))
				cmd.Run()
				speech.Speak("Volume has been set to " + ctrl.Value.(string))

			case "notification":

				err = speech.Speak(ctrl.Value.(string))

				if err != nil {
					fmt.Printf("speak returned error: %s", err.Error())
					continue
				}

			}
		}

	}
}
