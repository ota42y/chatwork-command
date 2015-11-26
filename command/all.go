package command

import (
	"fmt"
	"strings"

	"encoding/json"
	"github.com/codegangsta/cli"
	chatwork "github.com/ota42y/go-chatwork-api"
)

type RoomMessages struct {
	Name     string
	Messages []chatwork.Message
}

func CmdAll(c *cli.Context) {
	flag := c.String("f")

	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		return
	}

	client := chatwork.New(apiToken)

	rooms, err := client.GetRooms()
	if err != nil {
		fmt.Println(err)
		return
	}

	var messageArray []RoomMessages
	for _, room := range rooms {
		// if no unread message, not show (UnreadNum = 0 is skip)
		if 0 < room.UnreadNum {
			messages, err := client.GetMessage(room.RoomID, false)
			if err != nil {
				// GetMessage api not return already received message even if it's unread
				// So, if room have unread message, we can't get all unread message
				if err.Error() != "unexpected end of JSON input" {
					fmt.Println(err)
				}
			} else {
				if len(messages) != 0 {
					msg := RoomMessages{
						Name:     room.Name,
						Messages: messages,
					}

					messageArray = append(messageArray, msg)
				}
			}
		}
	}

	if flag == "json" {
		b, _ := json.Marshal(messageArray)
		fmt.Printf("%s\n", string(b))
	} else {
		showRoomMessage(messageArray)
	}
}

func showRoomMessage(messageArray []RoomMessages) {
	for _, room := range messageArray {
		if len(room.Messages) != 0 {
			fmt.Printf("%s %s\n", "room", room.Name)
			for _, message := range room.Messages {
				fmt.Printf("  %s %s\n", "user", message.Account.Name)
				fmt.Printf("    %s\n\n", strings.Replace(message.Body, "\n", "\n    ", -1))
			}
		}
	}
}
