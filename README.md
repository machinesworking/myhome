# myhome
A Home grown security system using Google Go, Raspberry Pi, and hacked Wizecams using Dafang hacks at https://github.com/Dafang-Hacks/rootfs

The system is intentially local only .....no cloud
Audible Notifications use picotts on wireless pi zeros
Lights are controlled using hacked sonoff switches
A thermostat is built from an esp8266
My Home server is written in go and uses websockets for device communications
Control is via web page using json messages.
Configuration is also a web page.
All devices are discovered automatically by the server connection and the control page is built dynamically.
 The Wizecams use onboard RtsptoWebrtc from  https://github.com/deepch/RTSPtoWebRTC Which uses https://github.com/pions/webrtc for the webrtc peer written entirely in Go.
