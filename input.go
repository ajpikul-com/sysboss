package main

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ajpikul-com/gitstatus"
	"github.com/ajpikul-com/wsssh/wsconn"
)

func ReadTexts(conn *wsconn.WSConn, name string) {
	defaultLogger.Debug("Starting to read texts")
	channel, _ := conn.SubscribeToTexts()
	buffer := bytes.NewBuffer([]byte{})
	for s := range channel {
		globalState.UpdateTime(name, name)
		buffer.WriteString(s) // we've received new input
		commandDecoder := json.NewDecoder(buffer)
		for {

			command := make(map[string]json.RawMessage) // not very flexible, want to specify command
			if err := commandDecoder.Decode(&command); err == nil || err == io.EOF {

				go processCommand(command, name)

				if err == io.EOF {
					defaultLogger.Debug("Command EOF")
					break
				}
			} else if err != nil {
				defaultLogger.Debug("Command error, copying buffer")
				defaultLogger.Error(err.Error())
				io.Copy(buffer, commandDecoder.Buffered()) // TODO: okay to ignore error here?. Better with io.multireader.
				break
			}
		}
	}
	defaultLogger.Debug("ReadTexts Channel Closed")
}

// TODO: better to use a multireader instead of copy, will also help us if there is anything left over in the buffer after decode, instead of assumign decode takes all
func processCommand(command map[string]json.RawMessage, name string) {
	for k, v := range command {
		defaultLogger.Debug("Processing Command")
		defaultLogger.Debug(k)
		if k == "clearget" { // Hacky, time to refactor
			globalState.ClearClient(name)
		} else if k == "get" {
			service := new(Service)
			err := json.Unmarshal(v, service)
			if err == nil {
				globalState.UpdateService(name, *service)
			} else {
				defaultLogger.Error(err.Error())
			}
		} else if k == "git" {
			repostate := new(map[string]gitstatus.RepoState)

			err := json.Unmarshal(v, repostate)
			if err != nil {
				defaultLogger.Error(err.Error())
			} else {
				globalRepoState.ClearClient(name)
				for k2, v2 := range *repostate {
					defaultLogger.Debug("Dumping a repostate")
					globalRepoState.Update(name, k2, v2) // ohh all this needs to be cleared
				}
			}
		}
	}
}
