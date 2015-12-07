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
		ch:          chatwork.New(apiToken),
		duration:    time.Duration(duration) * time.Minute,
		checkTime:   time.Now(),
		unReads:     make(map[int64]int64),
	}

	return watch, nil
}

type watchAPI struct {
	ch          *chatwork.Client
	duration    time.Duration
	checkTime   time.Time
	unReads     map[int64]int64
	jsonOutput  bool
}

func (w *watchAPI) Check(now time.Time) {
	d := now.Sub(w.checkTime)
	if 0 < w.duration && w.duration < d {
		w.checkAPI()
		w.checkTime = now
	}
}

func (w *watchAPI) checkAPI() {
	var rooms Rooms
	rooms, _ = w.ch.GetRooms()
	for _, room := range rooms {
		// if no unread message, not show (UnreadNum = 0 is skip)
		num := w.unReads[room.RoomID]
		if num < room.UnreadNum {
			w.showRoomMessage(room)
		}

		w.unReads[room.RoomID] = room.UnreadNum
	}
}

func (w *watchAPI) showRoomMessage(room chatwork.Room) {
	messages, _ := w.ch.GetMessage(room.RoomID, false)
	
	if len(messages) != 0 {
		msg := RoomMessages{
			Name:     room.Name,
			Messages: messages,
		}
		b, _ := json.Marshal(msg)
		fmt.Printf("%s\n", string(b))
		
	}
}
