package modules

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

// Tasker can be run on background
type Tasker interface {
	Start() error
	Stop() error
}

// Common fields for task
type Task struct {
	Name        string
	lock        sync.RWMutex
	Initialized bool
	Running     bool
	ChanStop    chan bool
	Loop        func()
}

func (t *Task) Start() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.Running {
		return errors.New("background task already running")
	}

	if !t.Initialized {
		return errors.New("task not initialized properly")
	}

	if t.Loop == nil {
		return errors.New("no loop function defined")

	}

	logrus.Debug("Starting task ", t.Name)
	t.Running = true
	go t.Loop()
	return nil
}

func (t *Task) Stop() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.Running {
		return errors.New("background task not running")
	}

	logrus.Debug("Stopping task ", t.Name)
	t.ChanStop <- true
	t.Running = false
	return nil
}

func (t *Task) Init() {
	t.ChanStop = make(chan bool, 2)
}
