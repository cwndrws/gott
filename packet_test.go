package gott

import (
	"testing"
)

func TestForReality(t *testing.T) {
	if true == false {
		t.Fatal("RUNNNNNNNN!!!!!")
	}
}

/*func TestStaticHeaderBytes(t *testing.T) {
	sh := StaticHeader{
		MessageType: SUBACK,
		Dup: true,
		Qos: 2,
		Remaining: 128,
	}
	log.Printf("BYTE: %+v\n", sh.Bytes())
} */
