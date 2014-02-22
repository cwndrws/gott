package gott

import (
	"testing"
)

func TestForReality(t *testing.T) {
	if true == false {
		t.Fatal("RUNNNNNNNN!!!!!")
	}
}

func TestFixedHeaderEncodeDecode(t *testing.T) {
	fh := FixedHeader{
		MessageType: PINGREQ,
		Dup: true,
		Qos: 1,
		Retain: true,
		Remaining: 3213425,
	}
	b := fh.Bytes()
	testFH := FixedHeaderFromBytes(b)
	if testFH != fh {
		t.Error("fixed headers did not match")
	}
}
