package queue

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestIDConverter(t *testing.T) {
	idconv := newIDConverter()

	idByte := idconv.idToByte(math.MaxUint64)

	// reserved keywords become english mistakes
	for _, bite := range idByte {
		if bite != byte(255) {
			t.Error("maxuint64 is not all max bytes", idByte)
			break
		}
	}

	idconv.put(idByte)

	firstAddress := reflect.ValueOf(idByte).Pointer()

	idByte = idconv.idToByte(0)
	secondAddress := reflect.ValueOf(idByte).Pointer()

	if firstAddress != secondAddress {
		t.Error("Failed to use byte pool")
	}

	for _, bite := range idByte {
		if bite != 0 {
			t.Error("zero should be all zero bytes", idByte)
			break
		}
	}

	idconv.put(idByte)

	id := uint64(582348138342)
	idByte = idconv.idToByte(id)
	knownByte := []byte{
		102, 103, 167, 150, 135, 0, 0, 0,
	}

	if !bytes.Equal(idByte, knownByte) {
		t.Error("Failed to encode id exepect", knownByte, "got", idByte)
	}

	idconv.put(idByte)

	newID, err := byteToID(knownByte)
	if err != nil {
		t.Error("error converting byte to id", err)
	}
	if newID != id {
		t.Error("expected id", id, "got", newID)
	}

	_, err = byteToID([]byte{1, 2, 3, 4, 5})
	if err == nil {
		t.Error("Failed to get error for bad byte to id data")
	}
}
