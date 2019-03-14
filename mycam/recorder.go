package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/pions/webrtc/pkg/media"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func record() error {
	var u1 uuid.UUID
	var address = "192.168.50.13:61000"
	var path = *devicename + "/" + *devicename + time.Now().Format(time.RFC3339) + ".gob\n"
	fmt.Printf("Openning: %s\n", time.Now().Format(time.RFC3339Nano))
	rw, err := open(address)
	if err != nil {
		return errors.Wrap(err, "Can't open recorder at "+address)

	}

	/*
		rw, err := openfile(path)
		if err != nil {
			fmt.Printf("couldn't open file")
			return errors.Wrap(err, "Can't open recorder at "+address)

		}
	*/
	defer func() {

		fmt.Println("Closing recorder")

	}()

	rw.WriteString("OPEN\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING command ")
	}
	rw.WriteString(path)
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING request ")
	}

	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}

	// Read the reply.
	log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'")
	}

	if response != "continue\n" {

		return errors.Wrap(nil, "STRING request: got a response: "+response)

	}
	log.Println("continuing")
	rw.WriteString("RECORD\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING command ")
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}

	u1 = uuid.Must(uuid.NewV4())
	mychan := make(chan media.RTCSample)
	m.RLock()
	sampleChannels[u1] = mychan
	m.RUnlock()
	fmt.Println("added samples recorder channel " + u1.String())

	enc := gob.NewEncoder(rw)
	for {
		select {
		case quit := <-cancelRecordChan:
			if quit == true {
				fmt.Println("Canceling record")
				err = rw.Flush()
				if err != nil {
					return errors.Wrap(err, "Flush failed.")
				}
				m.RLock()
				delete(sampleChannels, u1)
				m.RUnlock()
				fmt.Printf("deleted %s \n", u1.String())
				return errors.Wrap(nil, "recording canceled")
			}
		case sample := <-mychan:
			enc.Encode(sample)
		}
	}
}

//	return errors.Wrap(nil, "recording completed")

func open(addr string) (*bufio.ReadWriter, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}
func openfile(filename string) (*bufio.ReadWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Could not create file: %s", err.Error())
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file)), nil

}
