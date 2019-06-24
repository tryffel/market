package modules

import "testing"

type testTask struct {
	Task
}

func (t *testTask) loop() {
	print("Running background loop\n")
	<-t.ChanStop
	print("Stopping background loop\n")
}

func TestBackgroundTask(t *testing.T) {
	task := testTask{}
	task.Name = "test"
	task.Init()
	task.Initialized = true
	task.Loop = task.loop

	err := task.Start()
	if err != nil {
		t.Error(err)
	}

	err = task.Start()
	if err == nil {
		t.Error("No error while starting same background task twice")
	}

	err = task.Stop()
	if err != nil {
		t.Error(err)
	}

	err = task.Stop()
	if err == nil {
		t.Error("No error while stopping same background task twice")
	}

}
