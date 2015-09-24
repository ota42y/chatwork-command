package command

import (
	"os"
	"io"
	"fmt"
	"sort"
	"time"
	"strings"

	"github.com/codegangsta/cli"
	chatwork "github.com/ota42y/go-chatwork"
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
		showMessage(messages, writer)
	}
}

var messageIDLength = 15
func showMessage(messages []chatwork.Message, writer io.Writer) {
	var maxNameLength = 0
	for _, message := range messages {
		length := len(message.Account.Name)
		if maxNameLength < length {
			maxNameLength = length
		}
	}

	// %15d %s %s
	messageID := "ID" + strings.Repeat(" ", messageIDLength- 2)
	unReadNum := strings.Repeat(" ", maxNameLength - 4) + "Name"
	roomName := "Message"
	fmt.Fprintf(writer, "%s %s %s\n", messageID, unReadNum, roomName)

	for _, message := range messages {
		fmt.Fprintln(writer, fmtMessage(maxNameLength, message))
	}
}

func fmtMessage(maxNameLength int, message chatwork.Message) string {
	timeStr := time.Unix(message.SendTime, 0).Format("01/02 15:04 JST")

	// %15d    %s %s
	messageFormat := fmt.Sprintf("%%-%dd %%s %s%%s %%s\n", messageIDLength, strings.Repeat(" ", maxNameLength - len(message.Account.Name)) )
	return fmt.Sprintf(messageFormat, message.MessageId, timeStr, message.Account.Name, message.Body)
}

var IDLength = 9
var unReadLength = 7
var roomNameFormat = fmt.Sprintf("%%%dd %%%dd %%s\n", IDLength, unReadLength) // %9d %6d %s
func showRoomsData(rooms []chatwork.Room, writer io.Writer) {
	roomID := strings.Repeat(" ", IDLength - 2) + "ID"
	unReadNum := strings.Repeat(" ", unReadLength - 6) + "unRead"
	roomName := "RoomName"
	fmt.Fprintf(writer, "%s %s %s\n", roomID, unReadNum, roomName)

	for _, room := range rooms {
		fmt.Fprintf(writer, roomNameFormat, room.RoomId, room.UnreadNum, room.Name)
	}
}