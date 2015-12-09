package command

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/codegangsta/cli"
	chatwork "github.com/ota42y/go-chatwork-api"
	"github.com/robfig/cron"
)

// CmdWatch execute watch chatwork api
func CmdWatch(c *cli.Context) {
	watch, err := NewWatchAPI(c.Int("d"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cr := cron.New()
	cr.AddFunc("*/1 * * * * *", func() { watch.Check(time.Now()) })
	cr.Start()

	for {
		time.Sleep(10000000000000)
	}
}

func NewWatchAPI(duration int) (*watchAPI, error) {
	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		return nil, err
	}

	if duration <= 0 {
		duration = 1
	}

	watch := &watchAPI{
		ch:        chatwork.New(apiToken),
		duration:  time.Duration(duration) * time.Minute,
		checkTime: time.Now(),
		unReads:   make(map[int64]int64),
	}

	return watch, nil
}

type watchAPI struct {
	ch         *chatwork.Client
	duration   time.Duration
	checkTime  time.Time
	unReads    map[int64]int64
	jsonOutput bool
}

func (w *watchAPI) Check(now time.Time) {
	d := now.Sub(w.checkTime)
	if 0 < w.duration && w.duration < d {
		w.checkAPI()
		w.checkTime = now
	}
}

func (w *watchAPI) checkAPI() {
	rooms, err := w.ch.GetRooms()
	if err != nil {
		fmt.Println(err)
	} else {
		checkRoomIDs := w.checkUpdateRoom(rooms)
		w.showRoomMessages(checkRoomIDs)
	}
}

func (w *watchAPI) checkUpdateRoom(rooms []chatwork.Room) map[int64]string {
	checkRoomIDs := make(map[int64]string)

	for _, room := range rooms {
		num := w.unReads[room.RoomID]
		if num < room.UnreadNum {
			checkRoomIDs[room.RoomID] = room.Name
		}
		w.unReads[room.RoomID] = room.UnreadNum
	}
	return checkRoomIDs
}

type RoomMessages struct {
	Name     string
	Messages []chatwork.Message
}

func (w *watchAPI) showRoomMessages(checkRoomIDs map[int64]string) {
	for roomID, name := range checkRoomIDs {
		messages, _ := w.ch.GetMessage(roomID, false)
		msg := RoomMessages{
			Name:     name,
			Messages: messages,
		}
		b, _ := json.Marshal(msg)
		fmt.Printf("%s\n", string(b))
	}
}
