package gott

import (
	"testing"
	"log"
)

func TestStaticHeaderBytes(t *testing.T) {
	sh := StaticHeader{
		MessageType: SUBACK,
		Dup: true,
		Qos: 2,
		Remaining: 128,
	}
	log.Printf("BYTE: %+v\n", sh.Bytes())
}
