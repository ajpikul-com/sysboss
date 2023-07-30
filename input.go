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
			command := new(Service) // not very flexible, want to specify command
			if err := commandDecoder.Decode(command); err == nil || err == io.EOF {
				// Does EOF get passed when there's data
				go processCommand(command)
				if err == io.EOF {
					defaultLogger.Debug("Command EOF")
					break
				}
			} else if err != nil {
				defaultLogger.Debug("Command error, copying buffer")
				io.Copy(buffer, commandDecoder.Buffered()) // TODO: okay to ignore error here?
				break
			}
		}
	}
	defaultLogger.Debug("ReadTexts Channel Closed")
}

// TODO: better to use a multireader instead of copy, will also help us if there is anything left over in the buffer after decode, instead of assumign decode takes all
func processCommand(command *Service) {
	if command.Name == "" {
		return
	}
	globalState.UpdateService(*command)
}
