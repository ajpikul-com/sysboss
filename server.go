package main

import (
	"net/http"

	"golang.org/x/crypto/ssh"

	"github.com/ajpikul-com/wsssh/wsconn"
	"github.com/gorilla/websocket"
)

func ServeWSConn(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("Server: Incoming Req: " + r.Host + ", " + r.URL.Path)

	// Create websockets connection
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	gorrilaconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		defaultLogger.Error("Upgrade: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wsconn, err := wsconn.New(gorrilaconn)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		defaultLogger.Debug("Closing WSConn")
		// Doesn't warn client, just closes
		if err := wsconn.CloseAll(); err != nil {
			defaultLogger.Error("wsconn.CloseAll(): " + err.Error())
		}
	}()

	// Attach ssh to websockets
	sshconn, chans, reqs, err := GetServer(wsconn, globalConfig.PublicKeys, globalConfig.PrivateKey)
	if err != nil {
		defaultLogger.Error("GetServer(): " + err.Error())
		return
	}

	// User must have been legit
	defaultLogger.Info("Welcome, " + sshconn.Permissions.Extensions["comment"])

	// Record connection
	s := Service{ // Server never updates itself
		Name:      sshconn.Permissions.Extensions["comment"],
		IPAddress: r.Header.Get("x-forwarded-for"), // This is coming out as the proxy TODO maybe check to see if local
		Status:    "Online",
	}
	globalState.UpdateService(s)

	// Start Reading Input From User
	go ReadTexts(wsconn, s.Name)
	go ssh.DiscardRequests(reqs)
	for _ = range chans {
		// We're not accepting any session requests right now
		// I would like to see if server can act as client as well
	}

	defaultLogger.Info(sshconn.Permissions.Extensions["comment"] + " disconnected")
}
