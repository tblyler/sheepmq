package lease

import (
	"testing"
	"time"
)

func TestTimeoutValid(t *testing.T) {
	timeout := NewTimeout(int64(time.Second))

	if !timeout.Valid() {
		t.Error("timeout valid failed too soon")
	}

	time.Sleep(time.Nanosecond * 100)

	if !timeout.Valid() {
		t.Error("timeout valid failed too soon")
	}

	time.Sleep(time.Second - (time.Nanosecond * 100))

	if timeout.Valid() {
		t.Error("timeout valid should not be succeeding")
	}
}
