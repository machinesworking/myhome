package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/base64"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

//var DataChanelTest chan<- media.RTCSample

//var rawrtpchan chan<- *rtp.Packet
var namechan = make(chan string)
var playctrlchan = make(chan string)
var vp8Track *webrtc.RTCTrack

func main() {
	dir, _ := os.Getwd()

	fmt.Printf("CWD: %s\n", dir)
	r := mux.NewRouter()
	r.HandleFunc("/receive", home)
	r.HandleFunc("/list", list)
	r.HandleFunc("/play", play)
	r.HandleFunc("/stop", stop)

	//	http.Handle("/static/", http.FileServer(http.Dir("/home/scarpenter/go/src/machinesworking.com/myhome/mhrecorder/static")))

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	go playgob(playctrlchan)
	go func() {
		err := http.ListenAndServe(":8181", r)
		if err != nil {
		}
	}()
	select {}
}

func list(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	list, err := IOReadDir(data)
	if err != nil {
		fmt.Println("list failed: " + err.Error())

	}
	filelist, err := json.Marshal(list)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	w.Write(filelist)
}

func play(w http.ResponseWriter, r *http.Request) {

	data := r.FormValue("data")
	namechan <- data
}
func stop(w http.ResponseWriter, r *http.Request) {

	playctrlchan <- "stop"
}
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data := r.FormValue("data")
	sd, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Println(err)
		return
	}
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
	})
	vp8Track, err = peerConnection.NewRTCSampleTrack(webrtc.DefaultPayloadTypeH264, "video", "pion2")

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
		Sdp:  string(sd),
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
	w.Write([]byte(base64.StdEncoding.EncodeToString([]byte(answer.Sdp))))
	//	DataChanelTest = vp8Track.Samples

}
