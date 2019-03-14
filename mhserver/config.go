package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"os"
)

//load the configuration
func loadnodes(nodesFile string) ([]byte, error) {

	//read the entire file into var body
	body, err := ioutil.ReadFile(msgfilepath + "/" + nodesFile)

	if err != nil {
		fmt.Printf("Couldn't load nodes file\r\n")
		return nil, err
	}
	//unmarshal the config JSON into the local config pointer
	//return the struct pointer
	/*
		err = json.Unmarshal(body, msg)
		if err != nil {
			fmt.Printf("could not unmarshal body: %s\n", err.Error())
		}

			msg.Msg, err = json.Marshal(ctrl)
			if err != nil {
				fmt.Printf("could not marshal control value: %s\n", err.Error())
			}
	*/
	return body, nil
}

func loadmsgnames() (*message, error) {
	files := make(map[string]json.RawMessage)
	//	names := []string{}
	//read the configs in the config folder
	fileinfos, err := ioutil.ReadDir(msgfilepath)
	if err != nil {
		fmt.Printf("Couldn't load config file names file\r\n")
		return nil, err
	}

	// list only the JSON files...no folders
	for _, info := range fileinfos {
		if info.IsDir() {
			continue
		}
		if strings.HasSuffix(info.Name(), ".json") {
			//names = append(names, info.Name())
			txt, _ := loadnodes(info.Name())

			files[info.Name()] = txt
		}
	}
	//create the map for the JSON message

	msg := new(message)
	msg.Type = "control"
	ctrl := rawmsg{Type: "msglist", Value: files}

	ctrlstr, err := json.Marshal(ctrl)
	if err != nil {
		fmt.Printf("could not marshal controls: %s", err.Error())
		return nil, err
	}

	msg.Targets = []string{"text", "hmi"}
	msg.Msg = ctrlstr
	/*
		send the map as JSON to the websocket connection
		err = websocket.JSON.Send(conn, msg)

		if err != nil {
			fmt.Printf("Could not send initial config list data %s", err.Error())
		}
	*/

	return msg, nil
}

//Save a changed Config file
func saveNodes(nodesFile string, msg interface{}) {

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Could not save file: %s", nodesFile)
	}

	//write the File to the correct path
	fmt.Printf("Saving file %s    with    %s\n", msgfilepath+"/"+nodesFile, data)
	err = ioutil.WriteFile(msgfilepath+"/"+nodesFile, data, os.FileMode.Perm(0666))
	if err != nil {
		fmt.Printf("failed to save message:  %s", err.Error())
	}

}
