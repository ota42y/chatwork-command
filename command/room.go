package command

import (
	"os"
	"io"
	"fmt"

	"github.com/codegangsta/cli"
	chatwork "github.com/yoppi/go-chatwork"
)

func CmdRoom(c *cli.Context) {
	roomID := c.String("r")
	room(roomID, os.Stdout)
}

func room(roomID string, writer io.Writer) {
	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		fmt.Fprintln(writer, err.Error())
		return
	}

	chatwork := chatwork.NewClient(apiToken)

	if roomID == "" {
		rooms := chatwork.Rooms()
		for _, room := range rooms {
			showRoomData(room, writer)
		}
	}else{
		messages := chatwork.RoomMessages(roomID)
		for _, message := range messages {
			showMessage(message, writer)
		}
	}
}

func showMessage(message chatwork.Message, writer io.Writer) {
	fmt.Fprintln(writer, message.Account.Name, message.Body, message.SendTime)
}

func showRoomData(room chatwork.Room, writer io.Writer) {
	fmt.Fprintln(writer, room.RoomId, room.Name, room.UnreadNum)
}