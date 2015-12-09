package command

import (
	"testing"

	chatwork "github.com/ota42y/go-chatwork-api"
)

func TestCheckUpdateRoom(t *testing.T) {
	watch := &watchAPI{
		unReads: make(map[int64]int64),
	}

	watch.unReads[0] = 1
	watch.unReads[1] = 10

	var rooms []chatwork.Room
	rooms = append(rooms, chatwork.Room{
		Name:      "room1",
		RoomID:    0,
		UnreadNum: 1,
	})
	rooms = append(rooms, chatwork.Room{
		Name:      "room2",
		RoomID:    1,
		UnreadNum: 11,
	})
	rooms = append(rooms, chatwork.Room{
		Name:      "room3",
		RoomID:    3,
		UnreadNum: 42,
	})

	checkRoomIDs := watch.checkUpdateRoom(rooms)
	if len(checkRoomIDs) != 2 {
		t.Errorf("checkUpdateRoom should return 1 array, but %d (%v)", len(checkRoomIDs), checkRoomIDs)
	}

	name, ok := checkRoomIDs[1]
	if !ok {
		t.Errorf("checkUpdateRoom should return ID:1 but not return (%v)", checkRoomIDs)
	}

	if name != "room2" {
		t.Errorf("checkUpdateRoom should return ID:1 with name %s but %s", "room2", name)
	}

	value, ok := watch.unReads[3]
	if !ok {
		t.Errorf("checkUpdateRoom should create ID:3 room data but not created")
	}
	if value != 42 {
		t.Errorf("checkUpdateRoom should create ID:3 room with 42 unread num but %d", value)
	}

	if watch.unReads[0] != 1 {
		t.Errorf("checkUpdateRoom shouldn't change ID:1 room's unread count but change %d", watch.unReads[0])
	}
	if watch.unReads[1] != 11 {
		t.Errorf("checkUpdateRoom should change ID:11 room's unread count to 11 but change %d", watch.unReads[1])
	}
}
