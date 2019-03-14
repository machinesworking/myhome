package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/pions/webrtc"
)

var msgfilepath = "view/config"
var processchan = make(chan *Client)
var answerchan = make(chan webrtc.RTCSessionDescription)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send   chan message
	msg    *message
	device *device
}

type device struct {
	Type    string
	Name    string
	Desc    string
	State   string
	Address string
}

type command struct {
	Name  string
	Value map[string]interface{}
}
type message struct {
	Type    string
	Targets []string
	//Conn    *websocket.Conn
	Msg json.RawMessage
}
type rawmsg struct {
	Type  string
	Value interface{}
}
