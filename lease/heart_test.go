package lease

import (
	"testing"
	"time"
)

func TestHeartCheck(t *testing.T) {
	heart := NewHeart(int64(time.Second))

	if !heart.Check() {
		t.Error("Heartbeat check failed first check too soon")
	}

	time.Sleep(time.Second + time.Nanosecond)

	if heart.Check() {
		t.Error("Heartbeat check succeeded after its ttl")
	}
}

func TestHeartValid(t *testing.T) {
	heart := NewHeart(int64(time.Second))

	if !heart.Valid() {
		t.Error("Heartbeat valid failed first check too soon")
	}

	time.Sleep(time.Second / 3)

	if !heart.Valid() {
		t.Error("Heartbeat valid failed second check too soon")
	}

	time.Sleep(time.Second / 2)

	if !heart.Valid() {
		t.Error("Hearbeat valid failed third check too soon")
	}

	time.Sleep(time.Second / 2)

	if !heart.Valid() {
		t.Error("Heartbeat valid failed fourth check too soon")
	}

	time.Sleep(time.Second)

	if heart.Valid() {
		t.Error("Heartbeat valid succeeded after its ttl")
	}
}
