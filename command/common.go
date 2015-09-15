package command

import (
	"fmt"
	"github.com/typester/go-pit"
)

var ChatworkDomain string = "chatwork.com"

func getApiToken(pitDomain string) (string, error) {
	profile, err := pit.Get(pitDomain, pit.Requires{})
	if err != nil {
		return "", err
	}

	token := (*profile)["api_token"]
	if token == "" {
		return "", fmt.Errorf("api token not set in %s", pitDomain)
	}

	return token, nil
}
