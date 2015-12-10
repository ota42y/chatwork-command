package command

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	chatwork "github.com/ota42y/go-chatwork-api"
)

func CmdSend(c *cli.Context) {
	if !sendMessage(int64(c.Int("r")), c.String("m")) {
		os.Exit(1)
	}
}

func sendMessage(roomID int64, message string) bool {
	if roomID == 0 {
		fmt.Println("please set roomID (-r option)")
		return false
	}
	if message == "" {
		fmt.Println("please set message (-m option)")
		return false
	}

	apiToken, err := getApiToken(ChatworkDomain)
	if err != nil {
		fmt.Println(err)
		return false
	}

	ch := chatwork.New(apiToken)
	if ch == nil {
		fmt.Println("create chatwork client error")
		return false
	}
	_, err = ch.PostMessage(roomID, message)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
