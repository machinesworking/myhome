package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
	"github.com/pions/webrtc/pkg/media"
	uuid "github.com/satori/go.uuid"
)

var m sync.RWMutex

// DataChanelTest sample data channel
var sampleChannels = make(map[uuid.UUID]chan<- media.RTCSample)
var cancelchan = make(chan bool)

func Rtc(sd string, id string) {

	var u1 uuid.UUID

	webrtc.RegisterDefaultCodecs()
	peerConnection, err := webrtc.New(webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState ice.ConnectionState) {

		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState.String() == "Disconnected" {
			delete(sampleChannels, u1)

			fmt.Printf("removed session %s\n", u1.String())
			return
		}

	})
	vp8Track, err := peerConnection.NewRTCTrack(webrtc.DefaultPayloadTypeH264, "video", "pion2")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = peerConnection.AddTrack(vp8Track)
	if err != nil {
		log.Println(err)
		return
	}

	offer := webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeOffer,
		Sdp:  sd,
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println(err)
		return
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println(err)
		return
	}

	//	fmt.Printf("answer: %v", answer)
	// send the answer back
	msg := new(message)
	msg.Type = "control"
	msg.Targets = append(msg.Targets, id)
	rwmsg := rawmsg{Type: "answer", Value: rawmsg{Type: *devicename, Value: answer}}
	msg.Msg, err = json.Marshal(rwmsg)
	if err != nil {
		fmt.Printf("Could not marshal answer: %s", err.Error())
		return
	}
	messageChan <- *msg

	/*
		mssg, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("mssg: %s", err.Error())
		}

		c.WriteMessage(1, []byte(mssg))
	*/
	u1 = uuid.Must(uuid.NewV4())
	m.RLock()
	sampleChannels[u1] = vp8Track.Samples
	m.RUnlock()
	fmt.Println("added samples channel " + u1.String())
	<-cancelChan
	fmt.Println("rtc completed")
}
