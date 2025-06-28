package discord

import (
	"discord-bot/src/model"
	"os"
	"os/exec"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("should return a Discord instance", func(t *testing.T) {
		if New("fake-token") == (Discord{}) {
			t.Error("Expected a Discord instance, got nil")
		}
	})

	t.Run("should fatal when token is empty", func(t *testing.T) {
		if os.Getenv("BE_CRASHER") == "1" {
			New("")
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestNew")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	})
}

func TestSetCommands(t *testing.T) {
	d := New("fake-token")
	cmds := []model.Command{
		{Name: "command1", Description: "description1"},
		{Name: "command2", Description: "description2"},
	}
	d.SetCommands(cmds)

	if len(commands) != len(cmds) {
		t.Errorf("Expected %d commands, got %d", len(cmds), len(commands))
	}
}

func TestAddHandler(t *testing.T) {
	d := New("fake-token")
	d.SetCommands([]model.Command{{Name: "command1", Description: "description1"}})
	handler := model.Handler{Name: "handler1"}
	d.AddHandler(handler)

	if len(handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(handlers))
	}
}
