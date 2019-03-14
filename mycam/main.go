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
	Type    string
	Targets []string
	Msg     json.RawMessage
}

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

var messageChan = make(chan message)
var cancelwriteChan = make(chan bool)
var cancelholdChan = make(chan bool)
var cancelChan = make(chan bool)
var cancelRecordChan = make(chan bool)
var c *websocket.Conn
var err error
var devicedesc, devicename *string
var online bool

func main() {

	devicename = flag.String("name", "testcam", " -name The name of this camera")
	devicedesc = flag.String("desc", "test camera", " -desc The description of this camera")
	wsServer := flag.String("server", "192.168.50.13:8080", "-server The webserver to contact")

	flag.Parse()

	//	sig := make(chan os.Signal, 1)
	//	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	//	go StartHTTPServer()
	go rtsp2webrtc()
	go fifo()
	u := url.URL{Scheme: "ws", Host: *wsServer, Path: "/ws"}

Start:
	online = false
	log.Printf("connecting to %s", u.String())

	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		online = false
		time.Sleep(10 * time.Second)
		goto Start
	}
	dv := device{Type: "camera", Name: *devicename, Desc: *devicedesc, State: "online"}
	c.SetReadLimit(10000)
	dev, _ := json.Marshal(dv)
	status := message{Type: "status", Msg: dev}
	err = c.WriteJSON(&status)

	if err != nil {
		fmt.Println("I am sorry but I cannot send to the server.")
		online = false
		time.Sleep(10 * time.Second)
		goto Start

	}
	fmt.Println("Sent online status")
	online = true
	msg := new(message)

	for {
		err = c.ReadJSON(&msg)
		if err != nil {
			online = false
			fmt.Println("I am sorry but I cannot receive from the server.")
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
			case "nightmode":
				cmd := exec.Command("/usr/scripts/nightmode.sh", ctrl.Value.(string))
				cmd.Run()

			case "offer":
				fmt.Println("received offer")
				session := ctrl.Value.(map[string]interface{})
				sdp := session["data"].(string)
				from := session["from"].(string)
				go Rtc(sdp, from)
			case "oscommand":

				oscmd := ctrl.Value.(map[string]interface{})
				command := oscmd["cmd"].(string)
				args := oscmd["args"].([]interface{})
				argstr := make([]string, len(args))
				for k, v := range args {
					argstr[k] = v.(string)
				}

				cmd := exec.Command(command, argstr...)
				fmt.Printf("Command: %v", cmd)
				err := cmd.Start()
				log.Printf("Command finished with error: %v", err)
			}
		}
	}
}
