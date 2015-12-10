package command

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	if sendMessage(0, "message") {
		t.Errorf("when no set room id, sendMessage should return false")
	}
	if sendMessage(42, "") {
		t.Errorf("when no set message, sendMessage should return false")
	}
}
