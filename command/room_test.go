package command

import (
	"os"
	"testing"

	"github.com/codegangsta/cli"
)

func TestCmdRoom(t *testing.T) {
	room(0, os.Stdout)
}
