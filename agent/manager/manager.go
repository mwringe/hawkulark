package manager

import (
	"time"
	//"github.com/golang/glog"
)

type Manager interface {
	Start()
	Stop()
}

type Config struct {

}


type manager struct {

}

func NewManager() (Manager, error) {
	manager := manager{}
	return &manager, nil
}

func (m *manager) Start() {
	go m.timeLoop()
}

func (m *manager) Stop() {

}

func (m *manager) timeLoop() {
	for {
		//now := time.Now()

		select {
		case <- time.After(time.Second*15):
			//glog.Warning("HELLO WORLD")
		}
	}
}