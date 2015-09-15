package command

import (
	"os"
	"io"
	"fmt"

	"github.com/codegangsta/cli"
	chatwork "github.com/yoppi/go-chatwork"
)

func CmdRoom(c *cli.Context) {
	room(c, os.Stdout)
}

func room(c *cli.Context, writer io.Writer) {
	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		fmt.Fprintln(writer, err.Error())
		return
	}

	chatwork := chatwork.NewClient(apiToken)
	rooms := chatwork.Rooms()

	for _, room := range rooms {
		showRoom(room, writer)
	}
}

func showRoom(room chatwork.Room, writer io.Writer) {
	fmt.Fprintln(writer, room.RoomId, room.Name, room.UnreadNum)
}