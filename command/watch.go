package command

import (
    "fmt"
    "time"

    "github.com/codegangsta/cli"
    "github.com/robfig/cron"
    chatwork "github.com/ota42y/go-chatwork-api"
    "strings"
)

func CmdWatch(c *cli.Context) {
    watch, err := NewWatchApi(c.Int("d"), c.Int("v"))
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    watch.Init()

    cr := cron.New()
    cr.AddFunc("*/1 * * * * *", func() { watch.Check(time.Now()) } )
    cr.Start()

    for {
        time.Sleep(10000000000000)
    }
}

func NewWatchApi(duration int, verbose int) (*watchApi, error) {
    apiToken, err := getApiToken(ChatworkDomain)
    if err != nil {
        return nil, err
    }

    if duration <= 0 {
        duration = 1
    }

    watch := &watchApi{
        ch: chatwork.New(apiToken),
        duration: time.Duration(duration) * time.Minute,
        verbose: time.Duration(verbose) * time.Minute,
        checkTime: time.Now(),
        verboseTime: time.Now(),
        unReads: make(map[int64]int64),
    }

    return watch, nil
}

type watchApi struct {
    ch *chatwork.Client
    verbose time.Duration
    duration time.Duration
    checkTime time.Time
    verboseTime time.Time
    unReads map[int64]int64
}

func (* watchApi) Init() {
    fmt.Println("start")
}

func (w *watchApi) Check(now time.Time) {
    d := now.Sub(w.checkTime)
    if 0 < w.duration && w.duration < d {
        w.checkApi()
        w.checkTime = now
    }

    d = now.Sub(w.verboseTime)
    if 0 < w.verbose && w.verbose < d {
        w.printInfo()
        w.verboseTime = now
    }
}

func (w *watchApi) printInfo() {
    fmt.Println("\n\n\n-----check------\n")
}

func (w *watchApi) checkApi() {
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

func (w *watchApi) showRoomMessage(room chatwork.Room) {
    messages, _ := w.ch.GetMessage(room.RoomID, false)

    fmt.Printf("%s %s\n", "room", room.Name)
    for _, message := range messages {
        fmt.Printf("  %s %s\n", "user", message.Account.Name)
        fmt.Printf("    %s\n\n", strings.Replace(message.Body, "\n", "\n    ", -1))
    }
}