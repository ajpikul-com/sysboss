package main

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ajpikul-com/wsssh/wsconn"
)

func ReadTexts(conn *wsconn.WSConn, name string) {
	defaultLogger.Debug("Starting to read texts")
	channel, _ := conn.SubscribeToTexts()
	buffer := bytes.NewBuffer([]byte{})
	for s := range channel {
		globalState.UpdateTime(name)
		buffer.WriteString(s) // we've received new input
		commandDecoder := json.NewDecoder(buffer)
		for {
			var command interface{}
			if err := commandDecoder.Decode(&command); err == nil || err == io.EOF {
				go processCommand(command)
				if err == io.EOF {
					break
				}
			} else if err != nil {
				io.Copy(buffer, commandDecoder.Buffered()) // TODO: okay to ignore error here?
				break
			}
		}
	}
	defaultLogger.Debug("ReadTexts Channel Closed")
}

func processCommand(command interface{}) {
	if commandService, ok := command.(Service); ok {
		globalState.UpdateService(commandService)
	} else {
		panic("Command was weird type")
	}
}