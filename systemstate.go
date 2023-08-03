package main

// This could be a lot better by having per-client reading/writing

import (
	"sync"
	"time"
)

type Service struct { // Is this not a protobuff
	Name           string
	IPAddress      string
	Status         string
	ParentService  string
	LastConnection time.Time
}

type systemState struct {
	Services        map[string]map[string]*Service // TODO shouldn't we be storing git here?
	ExpectedClients []string
	mutex           sync.RWMutex
}

func NewSystemState() *systemState {
	ss := &systemState{
		Services:        make(map[string]map[string]*Service),
		ExpectedClients: make([]string, 0),
	}
	return ss
}

func (ss *systemState) ClearClient(clientName string) {
	ss.mutex.Lock()
	defer ss.mutex.Lock()
	ss.Services[clientName] = make(map[string]*Service)
}
func (ss *systemState) UpdateService(clientName string, uS Service) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	if _, ok := ss.Services[clientName]; !ok {
		ss.Services[clientName] = make(map[string]*Service)
	}
	ss.Services[clientName][uS.Name] = &uS
	if ss.Services[clientName][uS.Name].LastConnection.IsZero() {
		ss.Services[clientName][uS.Name].LastConnection = time.Now()
	}
}

func (ss *systemState) UpdateTime(clientName string, serviceName string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.Services[clientName][serviceName].LastConnection = time.Now()
}

func (ss *systemState) UpdateClientList(names []string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.ExpectedClients = make([]string, len(names))
	copy(ss.ExpectedClients, names)
}
func (ss *systemState) ReadLock() {
	ss.mutex.RLock()
}
func (ss *systemState) ReadUnlock() {
	ss.mutex.RUnlock()
}

// TODO: maybe if we do a JSONMarshal function we don't need to call these manually
