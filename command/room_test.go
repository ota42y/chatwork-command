package command

import (
	"os"
	"testing"
)

func TestCmdRoom(t *testing.T) {
	room(0, os.Stdout)
}
