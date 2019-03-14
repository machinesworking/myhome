package main

import (
	"fmt"
	"log"
	"time"

	rtsp "github.com/deepch/sample_rtsp"
	"github.com/pions/webrtc/pkg/media"
)

var (
	VideoWidth  int
	VideoHeight int
)

func rtsp2webrtc() {

	url := "rtsp://localhost:8554/unicast"
	sps := []byte{}
	pps := []byte{}
	fuBuffer := []byte{}
	count := 0
	Client := rtsp.RtspClientNew()
	Client.Debug = false
	syncCount := 0
	preTS := 0

	defer func() {
		//		cancelchan <- true
		Client.Close()
	}()

	writeNALU := func(sync bool, ts int, payload []byte) {

		var sample = media.RTCSample{Data: payload, Samples: uint32(ts - preTS)}
		for key := range sampleChannels {

			if sampleChannels[key] != nil && preTS != 0 {
				select {
				case sampleChannels[key] <- sample:
				default:

				}
			}
			preTS = ts
		}
	}
	handleNALU := func(nalType byte, payload []byte, ts int64) {
		if nalType == 7 {
			if len(sps) == 0 {
				sps = payload
			}
			//	writeNALU(true, int(ts), payload)
		} else if nalType == 8 {
			if len(pps) == 0 {
				pps = payload
			}
			//	writeNALU(true, int(ts), payload)
		} else if nalType == 5 {
			syncCount++
			lastkeys := append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), payload...)

			writeNALU(true, int(ts), lastkeys)
		} else {
			if syncCount > 0 {
				writeNALU(false, int(ts), payload)
			}
		}
	}
	if err := Client.Open(url); err != nil {
		fmt.Println("[RTSP] Error", err)
	} else {
		for {
			select {
			case <-Client.Signals:
				fmt.Println("Exit signals by rtsp")
				sps = []byte{}
				pps = []byte{}
				fuBuffer = []byte{}
				count = 0
				Client = rtsp.RtspClientNew()
				Client.Debug = false
				syncCount = 0
				preTS = 0
			RETRY1:
				if err := Client.Open(url); err != nil {
					fmt.Println("[RTSP] Error", err)
					time.Sleep(time.Second * 10)
					goto RETRY1
				}

			case <-cancelchan:
				Client.Close()
			case data := <-Client.Outgoing:
				count += len(data)
				//fmt.Println("recive  rtp packet size", len(data), "recive all packet size", count)
				if data[0] == 36 && data[1] == 0 {
					cc := data[4] & 0xF
					rtphdr := 12 + cc*4
					ts := (int64(data[8]) << 24) + (int64(data[9]) << 16) + (int64(data[10]) << 8) + (int64(data[11]))
					packno := (int64(data[6]) << 8) + int64(data[7])
					if false {
						log.Println("packet num", packno)
					}
					nalType := data[4+rtphdr] & 0x1F
					if nalType >= 1 && nalType <= 23 {
						handleNALU(nalType, data[4+rtphdr:], ts)
					} else if nalType == 28 {
						isStart := data[4+rtphdr+1]&0x80 != 0
						isEnd := data[4+rtphdr+1]&0x40 != 0
						nalType := data[4+rtphdr+1] & 0x1F
						nal := data[4+rtphdr]&0xE0 | data[4+rtphdr+1]&0x1F
						if isStart {
							fuBuffer = []byte{0}
						}
						fuBuffer = append(fuBuffer, data[4+rtphdr+2:]...)
						if isEnd {
							fuBuffer[0] = nal
							handleNALU(nalType, fuBuffer, ts)
						}
					}
				} else if data[0] == 36 && data[1] == 2 {
					//cc := data[4] & 0xF
					//rtphdr := 12 + cc*4
					//payload := data[4+rtphdr+4:]
				}
			}
		}
	}
	Client.Close()
}
