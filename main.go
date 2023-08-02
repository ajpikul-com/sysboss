package main

import (
	"encoding/json"
	"net/http"
	//	"strconv"
	"time"

	"github.com/ajpikul-com/gitstatus"
	"github.com/gorilla/mux"
)

func init() {
	initLogger()
	initConfig()
}

func main() {
	m := mux.NewRouter()
	// Add two gets here. Why are we using mux? I don't know
	m.HandleFunc("/", ServeWSConn)
	m.HandleFunc("/status", func(http.ResponseWriter, *http.Request) {})
	m.HandleFunc("/sysstatus", writeSystemStatus)
	s := &http.Server{
		Addr:           globalConfig.Hostname,
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	defaultLogger.Info("Initiating server")
	err := s.ListenAndServe()
	if err != nil {
		defaultLogger.Error("http.Server.ListenAndServe: " + err.Error())
	}
	defaultLogger.Info("Server exiting")
}

func writeSystemStatus(w http.ResponseWriter, r *http.Request) {
	// If r.Method = GET/POST
	jsonWrite := json.NewEncoder(w)
	jsonWrite.SetEscapeHTML(true)
	globalState.ReadLock() // TODO if global state has a jsonmarshall method it wouldn't need this
	jsonWrite.Encode(globalState)
	jsonWrite.Encode(map[string]*gitstatus.StateMap{"git": globalRepoState})
	//defaultLogger.Debug("Size of globalRepoState: " + strconv.Itoa(globalRepoState.Len()))
	globalState.ReadUnlock()
	w.(http.Flusher).Flush()
}
