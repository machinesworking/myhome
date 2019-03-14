package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

func processmsg() {
	for {
		select {
		case c := <-processchan:
			if c.msg.Type == "status" {
				// update device info with status msg
				var dev device

				err := json.Unmarshal(c.msg.Msg, &dev)
				if err != nil {
					log.Printf("unmarshal error: %s", err.Error())
				}
				dev.Address = c.device.Address
				c.device = &dev

				ti := time.Now()
				//send notification to browsers of status
				note := rawmsg{Type: "notification", Value: ti.Format(time.Stamp) + ": " + c.device.Desc + " connected"}
				jnote, err := json.Marshal(note)
				if err != nil {
					log.Printf("Marshal error: %s", err.Error())
				}

				// load the saved message names to the browsers
				mesg := message{Type: "control", Targets: []string{"text"}, Msg: jnote}
				c.hub.broadcast <- mesg
				if c.device.Type == "text" {
					msglist, _ := loadmsgnames()
					c.send <- *msglist

					ms := new(message)
					ms.Type = "control"
					var devs []string
					for client := range c.hub.clients {
						if client.device.Type == "camera" {
							devs = append(devs, client.device.Name)
						}
					}
					note := rawmsg{Type: "peerlist", Value: devs}

					ms.Msg, _ = json.Marshal(note)
					c.send <- *ms
					id, _ := uuid.NewV4()
					c.device.Name = id.String()
					note = rawmsg{Type: "id", Value: id}
					ms.Msg, _ = json.Marshal(note)
					c.send <- *ms

				}

				// send new camera control to the recorder
				if c.device.Type == "camera" {
					recnote := rawmsg{Type: "create", Value: c.device}
					jrecnote, err := json.Marshal(recnote)
					if err != nil {
						log.Printf("Marshal error: %s", err.Error())
					}
					recmesg := message{Type: "control", Targets: []string{"recorder"}, Msg: jrecnote}
					c.hub.broadcast <- recmesg
				}
				// if the device is the recorder send currently connected cameras
				if c.device.Type == "recorder" {
					for client := range c.hub.clients {
						log.Printf("found %s\n", client.device.Name)

						if client.device.Type == "camera" {
							recnote := rawmsg{Type: "create", Value: client.device}
							jrecnote, err := json.Marshal(recnote)
							if err != nil {
								log.Printf("Marshal error: %s", err.Error())
							}
							recmesg := message{Type: "control", Targets: []string{"recorder"}, Msg: jrecnote}
							c.send <- recmesg

						}

					}
				}

			}

			if c.msg.Type == "motion" {

				val := new(rawmsg)
				err := json.Unmarshal(c.msg.Msg, val)
				if err != nil {
					fmt.Printf("failed to decode motion\n")

				}
				// send motion start to the device node in the recorder
				if val.Value.(string) == "on" {
					rctrl := rawmsg{Type: "start", Value: c.device}
					jrctrl, _ := json.Marshal(rctrl)

					rmesg := message{Type: "control", Targets: []string{"recorder"}, Msg: jrctrl}
					c.hub.broadcast <- rmesg
					// send motion start to angelos and browsers
					note := rawmsg{Type: "notification", Value: "motion detected " + c.device.Desc}
					jnote, _ := json.Marshal(note)
					mesg := message{Type: "control", Targets: []string{"hmi", "text"}, Msg: jnote}
					c.hub.broadcast <- mesg

				}
				//same as above for off command
				if val.Value.(string) == "off" {
					note := rawmsg{Type: "stop", Value: c.device}
					jnote, _ := json.Marshal(note)

					mesg := message{Type: "control", Targets: []string{"recorder"}, Msg: jnote}
					c.hub.broadcast <- mesg
					/*
						note1 := rawmsg{Type: "notification", Value: "no motion  " + c.device.Desc}
						jnote1, _ := json.Marshal(note1)
						mesg1 := message{Type: "control", Targets: []string{"hmi", "text"}, Msg: jnote1}
						c.hub.broadcast <- mesg1
					*/
				}

			}
			if c.msg.Type == "state" {

				val := new(rawmsg)
				err := json.Unmarshal(c.msg.Msg, val)
				if err != nil {
					fmt.Printf("failed to decode state\n")

				}

				note := rawmsg{Type: "notification", Value: c.device.Desc + "is " + val.Value.(string)}
				jnote, _ := json.Marshal(note)

				mesg := message{Type: "control", Targets: []string{"hmi", "text"}, Msg: jnote}
				c.hub.broadcast <- mesg

			}
			//forward a control message
			if c.msg.Type == "control" {
				c.hub.broadcast <- *c.msg

			}

			if c.msg.Type == "command" {
				//get the command

				var cmd command
				err := json.Unmarshal(c.msg.Msg, &cmd)
				if err != nil {
					fmt.Printf("Could not unmarshal  %v\n", cmd)
				}

				switch cmd.Name {

				case "savemsg":
					// save a new message from browser
					log.Printf("Message Value : %v", cmd)
					var filename = cmd.Value["Name"].(string)

					var data = cmd.Value["Value"]

					saveNodes(filename, data)
					msglist, err := loadmsgnames()
					if err != nil {
						fmt.Printf("Could not load config list %s", err.Error())
					}
					c.send <- *msglist

					// delete a message
				case "deletemsg":
					filename := cmd.Value["filename"].(string)
					//		conn := msg["connection"].(*websocket.Conn)
					info, err := os.Lstat(msgfilepath + "/" + filename)
					if err != nil {
						fmt.Println("Can't stat config file")
					} else {
						if !info.IsDir() {
							os.Remove(msgfilepath + "/" + filename)
							msglist, err := loadmsgnames()
							if err != nil {
								fmt.Printf("Could not load config list %s", err.Error())
							}
							c.send <- *msglist
						}
					}

				case "loadmsg":

					// load a message and send to browser
					msg := new(message)
					var ctrl rawmsg
					msg.Type = "control"
					ctrl.Type = "configuration"
					name := cmd.Value["filename"].(string)
					log.Printf("loading filename %s\n", name)
					body, _ := loadnodes(name)

					err = json.Unmarshal(body, &ctrl.Value)
					if err != nil {
						fmt.Printf("could not unmarshal body: %s\n", err.Error())
					}
					msg.Msg, err = json.Marshal(ctrl)
					if err != nil {
						fmt.Printf("could not marshal control value: %s\n", err.Error())
					}
					// send the message to the browser
					select {
					case c.send <- *msg:
					default:
						close(c.send)
						delete(c.hub.clients, c)
					}
					//send the message from the browser
				case "sendmsg":
					msg := new(message)

					name := cmd.Value["filename"].(string)
					log.Printf("loading filename %s\n", name)
					localmsg, _ := loadnodes(name)

					log.Printf("Sending Message: %s", localmsg)
					json.Unmarshal(localmsg, msg)

					select {
					case c.hub.broadcast <- *msg:
					default:
						close(c.send)
						delete(c.hub.clients, c)
					}
				case "listdevices":
					// list all currently connected devices

					ms := new(message)
					ms.Type = "control"
					var devs []device
					//		ti := time.Now()
					for client := range c.hub.clients {
						if "all" == cmd.Value["Type"].(string) || cmd.Value["Type"].(string) == client.device.Type {
							devs = append(devs, *client.device)
						}
					}
					note := rawmsg{Type: "devicelist", Value: devs}

					ms.Msg, _ = json.Marshal(note)
					c.send <- *ms
					fmt.Printf("Sending device list to %s\n", c.device.Name)
				}

			}

		}
	}
}
