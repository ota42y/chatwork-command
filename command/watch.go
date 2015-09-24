package command

import (
    "fmt"
    "time"
    "strconv"

    "github.com/codegangsta/cli"
    "github.com/robfig/cron"
    chatwork "github.com/ota42y/go-chatwork"
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
        ch: chatwork.NewClient(apiToken),
        duration: time.Duration(duration) * time.Minute,
        verbose: time.Duration(verbose) * time.Minute,
        checkTime: time.Now(),
        verboseTime: time.Now(),
        unReads: make(map[int]int),
    }

    return watch, nil
}

type watchApi struct {
    ch *chatwork.Client
    verbose time.Duration
    duration time.Duration
    checkTime time.Time
    verboseTime time.Time
    unReads map[int]int
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
    fmt.Println("-----check------")
}

func (w *watchApi) checkApi() {
    var rooms Rooms = w.ch.Rooms()
    for _, room := range rooms {
        // if no unread message, not show (UnreadNum = 0 is skip)
        num := w.unReads[room.RoomId]
        if num < room.UnreadNum {
            w.showRoomMessage(room)
        }

        w.unReads[room.RoomId] = room.UnreadNum
    }
}

func (w *watchApi) showRoomMessage(room chatwork.Room) {
    messages := w.ch.RoomMessages(strconv.Itoa(room.RoomId))

    for _, message := range messages {
        fmt.Printf("%s %s\n", "room", room.Name)
        fmt.Printf("%s %s\n", "user", message.Account.Name)
        fmt.Printf("%s\n\n", message.Body)
    }
}