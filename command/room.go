package command

import (
	"os"
	"io"
	"fmt"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	chatwork "github.com/yoppi/go-chatwork"
)

// sort Sticky > UnreadNum > lastUpdateTime
type Rooms []chatwork.Room
func (r Rooms) Len() int {
	return len(r)
}

func (r Rooms) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rooms) Less(i, j int) bool {
	if r[i].Sticky != r[j].Sticky {
		return r[i].Sticky
	}

	if r[i].UnreadNum != 0 || r[j].UnreadNum != 0 {
		if r[i].UnreadNum != r[j].UnreadNum {
			// more big, more priority small
			return r[i].UnreadNum > r[j].UnreadNum
		}
	}

	return r[i].LastUpdateTime > r[j].LastUpdateTime
}

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
		var rooms Rooms = chatwork.Rooms()

		sort.Sort(rooms)
		showRoomsData(rooms, writer)
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

var IDLength = 9
var unReadLength = 7
var format = fmt.Sprintf("%%%dd %%%dd %%s\n", IDLength, unReadLength) // %9d %6d %s
func showRoomsData(rooms []chatwork.Room, writer io.Writer) {
	roomID := strings.Repeat(" ", IDLength - 2) + "ID"
	unReadNum := strings.Repeat(" ", unReadLength - 6) + "unRead"
	roomName := "RoomName"
	fmt.Fprintf(writer, "%s %s %s\n", roomID, unReadNum, roomName)

	for _, room := range rooms {
		fmt.Fprintf(writer, format, room.RoomId, room.UnreadNum, room.Name)
	}
}