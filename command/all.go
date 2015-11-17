package command

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	chatwork "github.com/ota42y/go-chatwork-api"
)

func CmdAll(c *cli.Context) {
	//flag := c.String("f")

	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		return
	}

	client := chatwork.New(apiToken)

	var rooms Rooms
	rooms, _ = client.GetRooms()
	for _, room := range rooms {
		// if no unread message, not show (UnreadNum = 0 is skip)
		if 0 < room.UnreadNum {
			showRoomMessage(client, room)
		}
	}
}

func showRoomMessage(client *chatwork.Client, room chatwork.Room) {
	messages, _ := client.GetMessage(room.RoomID, false)

	if len(messages) != 0 {
		fmt.Printf("%s %s\n", "room", room.Name)
		for _, message := range messages {
			fmt.Printf("  %s %s\n", "user", message.Account.Name)
			fmt.Printf("    %s\n\n", strings.Replace(message.Body, "\n", "\n    ", -1))
		}
	}
}
