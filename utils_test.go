package main

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	want := make([]byte, RowSize)
	binary.LittleEndian.AppendUint32(want[0:4], 0)
	copy(want[4:36], "test")
	copy(want[36:291], "test@test.com")

	src := Row{
		Id:       0,
		Username: "test",
		Email:    "test@test.com",
	}
	got := make([]byte, RowSize)
	Serialize(src, got)

	if !bytes.Equal(got, want) {
		t.Errorf("Serialization failed. got: %v, want: %v", got, want)
	}
}

func TestDeserialize(t *testing.T) {
	t.Run("Deserialize row []byte into row Row", func(t *testing.T) {
		got := &Row{}
		want := &Row{
			Id:       0,
			Username: "test",
			Email:    "test@test.com",
		}

		src := make([]byte, RowSize)
		binary.LittleEndian.PutUint32(src[0:4], 0)
		copy(src[4:36], "test")
		copy(src[36:291], "test@test.com")

		Deserialize(src, got)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Test %s failed. got: %v, want: %v", t.Name(), got, want)
		}
	})

	t.Run("Fail on buffer smaller than row size", func(t *testing.T) {
		arbitraryNumber := 100
		src := make([]byte, RowSize-arbitraryNumber)
		r := &Row{}
		if err := Deserialize(src, r); err == nil {
			t.Errorf("Test %s failed. Safety check was bypassed", t.Name())
		}
	})
}
